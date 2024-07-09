package jobs

import (
	c "context"
	"encoding/json"

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
	PortfolioUpdateCron = "portfolioUpdateCron"
)

type PortfolioUpdater interface {
	DoJob(ctx c.Context, refresher string)
}

type portfolioUpdateJob struct {
	upswing      external.UpSwing
	portfolioDao dao.PortfolioDAO
}

func DefaultPortfolioUpdateJob() cron.Job {
	return &portfolioUpdateJob{upswing: factory.GetUpSwingExternalService(), portfolioDao: factory.GetPortfolioDAO()}
}

func DefaultPortfolioUpdater() PortfolioUpdater {
	return &portfolioUpdateJob{upswing: factory.GetUpSwingExternalService(), portfolioDao: factory.GetPortfolioDAO()}
}

func (p *portfolioUpdateJob) Run() {
	var ctx = context.Background(PortfolioUpdateCron)
	defer ctx.Done()
	enabled := isJobEnabled(ctx, PortfolioUpdateCron)
	if !enabled {
		log.Warn(ctx).Msg("portfolio update job is marked as disabled in config, skipping its execution")
		return
	}

	p.DoJob(ctx, "")
}

func (p *portfolioUpdateJob) DoJob(ctx c.Context, refresher string) {
	log.Info(ctx).Msg("starting portfolio update job...")
	var instantRefresh bool
	if refresher == "instant" {
		instantRefresh = true
	}
	p.execute(ctx, instantRefresh)
	log.Info(ctx).Msg("stopping portfolio update job...")
}

func (p *portfolioUpdateJob) execute(ctx c.Context, instantRefresh bool) {
	provider := getPortfolioUpdateProvider(ctx)
	clientList, err := p.portfolioDao.FetchClientList(ctx, provider, instantRefresh)
	if err != nil {
		log.Error(ctx).Err(err).Stack().Msg("fetching client list for portfolio update job failed")
		return
	}

	var batchSize = config.Default().GetIntD(constants.ApplicationConfig, constants.PortfolioUpdateBatchSize, 50)
	var portfolioUpdateEntities []entity.PortfolioEntity

	for _, clientCode := range clientList {
		var totalInterestPercentage float64
		var portfolioUpdateEntity entity.PortfolioEntity
		var errRespMap map[string]interface{}
		var isError bool
		response, err := p.upswing.GetNetWorthData(ctx, clientCode)
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
		portfolioUpdateEntity.ClientCode = clientCode
		portfolioUpdateEntity.Provider = provider
		portfolioUpdateEntity.CreatedBy = "portfolio_update_job"
		portfolioUpdateEntity.UpdatedBy = "portfolio_update_job"

		if isError {
			key := constants.ErrorCode
			if _, exists := errRespMap[key]; exists {
				if errRespMap[key] == constants.ErrClientNotFound {
					portfolioUpdateEntity.InvalidClient = true
				}
				portfolioUpdateEntity.ApiError = errRespMap[key].(string)
			}
			portfolioUpdateEntity.TotalActiveDeposits = 0
			portfolioUpdateEntity.InvestedValue = 0.0
			portfolioUpdateEntity.CurrentValue = 0.0
			portfolioUpdateEntity.InterestEarned = 0.0
			portfolioUpdateEntity.ReturnsValue = 0.0
			portfolioUpdateEntity.ReturnsPercentage = 0.0
		} else {
			if response.TotalInvestedAmount.Amount == 0.00 {
				totalInterestPercentage = 0.00
			} else {
				totalInterestPercentage = (response.TotalInterestEarned.Amount / response.TotalInvestedAmount.Amount) * 100
			}
			portfolioUpdateEntity.TotalActiveDeposits = response.ActiveTermDepositCount
			portfolioUpdateEntity.InvestedValue = response.TotalInvestedAmount.Amount
			portfolioUpdateEntity.CurrentValue = response.CurrentAmount.Amount
			portfolioUpdateEntity.InterestEarned = response.TotalInterestEarned.Amount
			portfolioUpdateEntity.ReturnsValue = response.TotalInterestEarned.Amount
			portfolioUpdateEntity.ReturnsPercentage = totalInterestPercentage
		}

		portfolioUpdateEntities = append(portfolioUpdateEntities, portfolioUpdateEntity)

		if len(portfolioUpdateEntities) >= int(batchSize) {
			err = p.portfolioDao.BatchUpdatePortfolio(ctx, portfolioUpdateEntities)
			if err != nil {
				log.Error(ctx).Err(err).Stack().Msg("batch updating client portfolios failed")
				return
			}
			portfolioUpdateEntities = portfolioUpdateEntities[:0]
		}
	}

	if len(portfolioUpdateEntities) > 0 {
		err = p.portfolioDao.BatchUpdatePortfolio(ctx, portfolioUpdateEntities)
		if err != nil {
			log.Error(ctx).Err(err).Stack().Msg("batch updating client portfolios failed")
			return
		}
	}

	if instantRefresh {
		err := p.portfolioDao.UpdateRefreshedPortfolioClientList(ctx, provider, clientList)
		if err != nil {
			log.Error(ctx).Err(err).Stack().Msg("error while update refreshed portfolio client list")
			return
		}
	} else {
		err := p.portfolioDao.CleanStaleRecords(ctx)
		if err != nil {
			log.Error(ctx).Err(err).Stack().Msg("error while cleaning up stale portfolio records")
			return
		}
	}
}
