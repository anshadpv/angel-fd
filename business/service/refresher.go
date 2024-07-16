package service

import (
	"context"

	"github.com/angel-one/fd-core/business/jobs"
)

type RefresherService interface {
	InvokePendingJourneyRefresher(ctx context.Context, refresher string)
}

type refresherServiceImpl struct {
	pendingJourneyRefresher jobs.PendingJourneyUpdaterTest
}

func DefaultRefresherService() RefresherService {
	return &refresherServiceImpl{pendingJourneyRefresher: jobs.DefaultPendingJourneyUpdaterTest()}
}

func (service *refresherServiceImpl) InvokePendingJourneyRefresher(ctx context.Context, refresher string) {
	service.pendingJourneyRefresher.DoJobTest(ctx, refresher)
}
