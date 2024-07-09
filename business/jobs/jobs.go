package jobs

import (
	"context"

	"github.com/angel-one/fd-core/commons/config"
	fdctx "github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/constants"
	"github.com/robfig/cron/v3"
)

func StartJobs() {
	crons := map[string]cron.Job{
		TokenRenewalCron:         DefaultTokenRenewalJob(),
		PortfolioUpdateCron:      DefaultPortfolioUpdateJob(),
		PendingJourneyUpdateCron: DefaultPendingJourneyJob(),
	}

	c := cron.New()
	for k, v := range crons {
		schedule := GetConfig(k)
		_, _ = c.AddJob(schedule, v)
	}
	c.Start()
	log.Info(fdctx.Background("jobs")).Msgf("inited crons: %+v\n", c.Entries())
}

func GetConfig(key string) string {
	return config.Default().GetStringD(constants.ApplicationConfig, key, "")
}

func isJobEnabled(ctx context.Context, name string) bool {
	//todo: read from config/db for dynamism
	enabled := false
	if name == "tokenRenewalCron" {
		enabled = true
	}
	return enabled
}

func getPortfolioUpdateProvider(ctx context.Context) string {
	return config.Default().GetStringD(constants.ApplicationConfig, constants.PortfolioProvider, "")
}

func getPendingJourneyUpdateProvider(ctx context.Context) string {
	return config.Default().GetStringD(constants.ApplicationConfig, constants.PendingJourneyProvider, "")
}
