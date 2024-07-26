package v1

import (
	"context"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/repository/dao"
	"github.com/angel-one/fd-core/commons/log"
	"github.com/angel-one/goerr"
)

type PortfolioService interface {
	GetPortfolio(ctx context.Context, provider string, clientCode string) (*model.Portfolio, error)
	GetPortfolioFromRedis(ctx context.Context, provider string, clientCode string) (model.PortfolioDetails, error)
}

type portfolioServiceImpl struct {
	portfolioDAO dao.PortfolioDAO
}

func DefaultPortfolioService() PortfolioService {
	return &portfolioServiceImpl{portfolioDAO: dao.DefaultPortfolioDAO()}
}

func (p *portfolioServiceImpl) GetPortfolio(ctx context.Context, clientCode string, provider string) (*model.Portfolio, error) {
	entity, err := p.portfolioDAO.FindByClient(ctx, clientCode, provider)
	if err != nil {
		return nil, goerr.New(err, "service: GetPortfolio by client failed")
	}
	if entity == nil {
		log.Debug(ctx).Msgf("no portfolio exists for clientCode: %s", clientCode)
		return &model.EmptyPortfolio, nil
	}
	portfolio := model.Portfolio{TotalActiveDeposits: entity.TotalActiveDeposits, InvestedValue: entity.InvestedValue, CurrentValue: entity.CurrentValue, InterestEarned: entity.InterestEarned, ReturnsValue: entity.ReturnsValue, ReturnsPercentage: entity.ReturnsPercentage}
	return &portfolio, nil
}

func (p *portfolioServiceImpl) GetPortfolioFromRedis(ctx context.Context, clientCode string, provider string) (model.PortfolioDetails, error) {
	var portfolioDetails model.PortfolioDetails
	portfolio, err := p.portfolioDAO.FetchPortfolioFromRedis(ctx, clientCode, provider)
	if err != nil {
		return portfolioDetails, goerr.New(err, "service: GetPortfolioFromRedis by client failed")
	}

	portfolioDetails.Portfolio = portfolio
	return portfolioDetails, nil
}
