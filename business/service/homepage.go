package service

import (
	"context"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/repository/dao"
)

type HomepageService interface {
	GetHomePageDetails(ctx context.Context, clientCode string, provider string) (model.Homepage, error)
}

type HomepageServiceImpl struct {
	plansDAO          dao.PlansDAO
	pendingJourneyDAO dao.PendingJourneyDAO
}

func DefaultHomepageService() HomepageService {
	return &HomepageServiceImpl{plansDAO: dao.DefaultPlansDAO(), pendingJourneyDAO: dao.DefaultPendingJourneyDAO()}
}

func (service *HomepageServiceImpl) GetHomePageDetails(ctx context.Context, clientCode string, provider string) (model.Homepage, error) {
	response := model.Homepage{}
	allFDs, err := service.plansDAO.FetchAllFDDetails(ctx)
	if err != nil {
		return response, err
	}
	mostBoughtPlans, err := service.plansDAO.FetchMostBoughtDetails(ctx)
	if err != nil {
		return response, err
	}

	pendingJourney, err := service.pendingJourneyDAO.FetchPendingJourneyDetails(ctx, clientCode, provider)
	if err != nil {
		return response, err
	}
	response.AllFDS = allFDs
	response.MostBought = mostBoughtPlans
	if pendingJourney != nil {
		response.Journey = model.Journey{Pending: pendingJourney.Pending, PendingState: model.PendingState{Payment: pendingJourney.Payment, KYC: pendingJourney.KYC}}
	}

	return response, nil
}
