package jobs

import (
	"github.com/angel-one/fd-core/commons/context"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/fd-core/external"
	"github.com/angel-one/fd-core/factory"
	"github.com/robfig/cron/v3"
)

const (
	TokenRenewalCron = "tokenRenewalCron"
)

type tokenRenewalJob struct {
	upswing external.UpSwing
}

func DefaultTokenRenewalJob() cron.Job {
	return &tokenRenewalJob{upswing: factory.GetUpSwingExternalService()}
}

func (t *tokenRenewalJob) Run() {
	var ctx = context.Background(TokenRenewalCron)
	defer ctx.Done()

	log.Info(ctx).Msg("starting token renewal job...")
	enabled := isJobEnabled(ctx, TokenRenewalCron)
	if !enabled {
		log.Warn(ctx).Msg("renew token job is marked as disabled in config, skipping its execution")
		return
	}

	t.upswing.ValidateToken(ctx)
	log.Info(ctx).Msg("stopping token renewal job...")
}
