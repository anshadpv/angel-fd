package service

import (
	"context"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/repository/dao"
)

type PlansService interface {
	GetAllPlans(ctx context.Context) (model.Plans, error)
	GetFSIPlans(ctx context.Context, fsi string) (model.FsiPlans, error)
}

type PlansServiceImpl struct {
	plansDAO dao.PlansDAO
}

func DefaultPlansService() PlansService {
	return &PlansServiceImpl{plansDAO: dao.DefaultPlansDAO()}
}

func (service *PlansServiceImpl) GetAllPlans(ctx context.Context) (model.Plans, error) {
	response := model.Plans{}
	allPlans, err := service.plansDAO.FetchAllPlansDetails(ctx)
	if err != nil {
		return response, err
	}
	response.Plans = allPlans

	return response, nil
}

func (service *PlansServiceImpl) GetFSIPlans(ctx context.Context, fsi string) (model.FsiPlans, error) {
	fsiPlans, err := service.plansDAO.FetchFsiPlansDetails(ctx, fsi)

	maxInterestRate := 0.0
	for _, plan := range fsiPlans.Plans {
		if plan.InterestRate > maxInterestRate {
			maxInterestRate = plan.InterestRate
		}
	}
	fsiPlans.MaxInterestRate = maxInterestRate
	if err != nil {
		return fsiPlans, err
	}
	return fsiPlans, nil
}
