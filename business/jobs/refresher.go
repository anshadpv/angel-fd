package jobs

import (
	c "context"
	"encoding/json"
	"fmt"

	"github.com/angel-one/fd-core/business/repository/dao"
	"github.com/angel-one/fd-core/business/repository/entity"
	"github.com/angel-one/fd-core/commons/config"
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/angel-one/fd-core/external"
	"github.com/angel-one/fd-core/factory"
	"github.com/angel-one/goerr"
	"github.com/robfig/cron/v3"
)

const (
	PendingJourneyUpdateCronTest = "PendingJourneyUpdateCron"
)

type PendingJourneyUpdaterTest interface {
	DoJobTest(ctx c.Context, refresher string)
}

type pendingJourneyJobTest struct {
	upswing               external.UpSwing
	pendingJourneyDaoTest dao.PendingJourneyDAOTest
}

func DefaultpendingJourneyJobTest() cron.Job {
	return &pendingJourneyJobTest{upswing: factory.GetUpSwingExternalService(), pendingJourneyDaoTest: factory.GetPendingJourneyDAOTest()}
}

func DefaultPendingJourneyUpdaterTest() PendingJourneyUpdaterTest {
	return &pendingJourneyJobTest{upswing: factory.GetUpSwingExternalService(), pendingJourneyDaoTest: factory.GetPendingJourneyDAOTest()}
}

func (p *pendingJourneyJobTest) Run() {
	var ctx = context.Background(PendingJourneyUpdateCronTest)
	defer ctx.Done()

	enabled := isJobEnabled(ctx, PendingJourneyUpdateCronTest)
	if !enabled {
		log.Warn(ctx).Msg("pending journey job is marked as disabled in config, skipping its execution")
		return
	}
	p.DoJobTest(ctx, "")
}

func (p *pendingJourneyJobTest) DoJobTest(ctx c.Context, refresher string) {
	log.Info(ctx).Msg("starting pending journey job...")
	var instantRefresh bool
	if refresher == "instant" {
		instantRefresh = true
	}
	p.executeTest(ctx, instantRefresh)
	log.Info(ctx).Msg("stopping pending journey update job...")
}

func (p *pendingJourneyJobTest) executeTest(ctx c.Context, instantRefresh bool) {
	provider := getPendingJourneyUpdateProvider(ctx)

	clientList, err := p.pendingJourneyDaoTest.FetchClientListTest(ctx, provider, instantRefresh)
	if err != nil {
		log.Error(ctx).Err(err).Stack().Msg("fetching client list for pending journey job failed")
		return
	}

	var batchSize = config.Default().GetIntD(constants.ApplicationConfig, constants.PendingJourneyUpdateBatchSize, 50)
	var pendingJourneyEntities []entity.PendingJourneyEntity
	for _, clientCode := range clientList {
		var pendingJourneyEntity entity.PendingJourneyEntity
		var errRespMap map[string]interface{}
		var isError bool
		response, err := p.upswing.GetPendingJourneyData(ctx, clientCode)
		if err != nil {
			log.Error(ctx).Err(err).Stack().Msg("error from upswing API")
			errResp := goerr.ListStacks(err)[2]
			err = json.Unmarshal([]byte(errResp), &errRespMap)
			if err != nil {
				log.Error(ctx).Err(err).Stack().Msg("error unmarshalling upswing API error response JSON")
				return
			}
			isError = true
		}
		pendingJourneyEntity.ClientCode = clientCode
		pendingJourneyEntity.Provider = provider
		pendingJourneyEntity.CreatedBy = "pending_journey_update_job"
		pendingJourneyEntity.UpdatedBy = "pending_journey_update_job"
		if isError {
			key := constants.ErrorCode
			if _, exists := errRespMap[key]; exists {
				if errRespMap[key] == fmt.Sprintf(constants.ErrPciNotFound, clientCode) {
					pendingJourneyEntity.InvalidClient = true
				}
				pendingJourneyEntity.ApiError = errRespMap[key].(string)
			}
			pendingJourneyEntity.Pending = false
			pendingJourneyEntity.Payment = false
			pendingJourneyEntity.KYC = false
		} else {
			pendingJourneyEntity.Pending = response.JourneyPending
			pendingJourneyEntity.Payment = response.JourneyPendingOnPayment
			pendingJourneyEntity.KYC = response.JourneyPendingOnVkyc
		}
		pendingJourneyEntities = append(pendingJourneyEntities, pendingJourneyEntity)

		if len(pendingJourneyEntities) >= int(batchSize) {
			err = p.pendingJourneyDaoTest.BatchUpdatePendingJourneyTest(ctx, pendingJourneyEntities)
			if err != nil {
				log.Error(ctx).Err(err).Stack().Msg("batch updating pending journey failed")
				return
			}
			pendingJourneyEntities = pendingJourneyEntities[:0]
		}
	}

	if len(pendingJourneyEntities) > 0 {
		err = p.pendingJourneyDaoTest.BatchUpdatePendingJourneyTest(ctx, pendingJourneyEntities)
		if err != nil {
			log.Error(ctx).Err(err).Stack().Msg("batch updating pending journey failed")
			return
		}
	}

	if instantRefresh {
		err := p.pendingJourneyDaoTest.UpdateRefreshedPendingJourneyClientListTest(ctx, provider, clientList)
		if err != nil {
			log.Error(ctx).Err(err).Stack().Msg("error while update refreshed pending_journey client list")
			return
		}
	} else {
		err := p.pendingJourneyDaoTest.CleanStaleRecordsTest(ctx)
		if err != nil {
			log.Error(ctx).Err(err).Stack().Msg("error while cleaning up stale pending_journey records")
			return
		}
	}
}
