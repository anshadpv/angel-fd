package service

import (
	"context"

	"github.com/angel-one/fd-core/business/jobs"
)

type JobsService interface {
	InvokePortfolioJob(ctx context.Context, refresher string)
	InvokePendingJourneyJob(ctx context.Context, refresher string)
}

type jobsServiceImpl struct {
	portfolioJob      jobs.PortfolioUpdater
	pendingJourneyJob jobs.PendingJourneyUpdater
}

func DefaultJobsService() JobsService {
	return &jobsServiceImpl{portfolioJob: jobs.DefaultPortfolioUpdater(), pendingJourneyJob: jobs.DefaultPendingJourneyUpdater()}
}

func (service *jobsServiceImpl) InvokePortfolioJob(ctx context.Context, refresher string) {
	service.portfolioJob.DoJob(ctx, refresher)
}

func (service *jobsServiceImpl) InvokePendingJourneyJob(ctx context.Context, refresher string) {
	service.pendingJourneyJob.DoJob(ctx, refresher)
}
