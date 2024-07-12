package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/repository/dao"
)

type HomeInfoService interface {
	GetHomeInfoDetails(ctx context.Context, clientCode string, provider string) (model.HomeInfo, error)
}

type HomeInfoServiceImpl struct {
	homeInfoDAO dao.HomeInfoDAO
}

func DefaultHomeInfoService() HomeInfoService {
	return &HomeInfoServiceImpl{homeInfoDAO: dao.DefaultHomeInfoDAO()}
}

func (service *HomeInfoServiceImpl) GetHomeInfoDetails(ctx context.Context, clientCode string, provider string) (model.HomeInfo, error) {
	response := model.HomeInfo{}

	allFDs, err := service.homeInfoDAO.FetchAllFDDetails(ctx)
	if err != nil {
		return response, err
	}
	mostBoughtPlans, err := service.homeInfoDAO.FetchMostBoughtDetails(ctx)
	if err != nil {
		return response, err
	}

	pendingJourney, err := service.homeInfoDAO.FetchPendingJourneyDetails(ctx, clientCode, provider)
	if err != nil {
		return response, err
	}
	response.AllFDS = allFDs
	response.MostBought = mostBoughtPlans
	if pendingJourney != nil {
		response.Journey = model.Journey{Pending: pendingJourney.Pending, PendingState: model.PendingState{Payment: pendingJourney.Payment, KYC: pendingJourney.KYC}}
	}

	var faq []model.FAQ
	tag := "home"
	faqData, err := service.homeInfoDAO.FetchFAQDetails(ctx, tag)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(faqData, &faq)
	fmt.Print(faq)
	if err != nil {
		return response, err
	}
	response.FAQs = faq

	return response, nil
}
