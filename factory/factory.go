package factory

import (
	"context"

	"github.com/angel-one/fd-core/business/repository/dao"
	v1 "github.com/angel-one/fd-core/business/service/v1"
	"github.com/angel-one/fd-core/external"
)

var upSwingService external.UpSwing
var portfolioService v1.PortfolioService
var portfolioDAO dao.PortfolioDAO
var pendingJourneyDAO dao.PendingJourneyDAO

func Init(ctx context.Context) {
	upSwingService = external.DefaultUpSwing(ctx)

	// services
	portfolioService = v1.DefaultPortfolioService()

	//dao
	portfolioDAO = dao.DefaultPortfolioDAO()
	pendingJourneyDAO = dao.DefaultPendingJourneyDAO()
}

func GetUpSwingExternalService() external.UpSwing {
	return upSwingService
}

func GetPortfolioService() v1.PortfolioService {
	return portfolioService
}

func GetPortfolioDAO() dao.PortfolioDAO {
	return portfolioDAO
}

func GetPendingJourneyDAO() dao.PendingJourneyDAO {
	return pendingJourneyDAO
}
