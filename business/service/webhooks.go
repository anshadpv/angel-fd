package service

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/repository/dao"
	"github.com/angel-one/fd-core/business/repository/entity"
	"github.com/angel-one/goerr"
)

var avgDays float64 = 30.417

type WebhookService interface {
	RegisterNewEvent(ctx context.Context, vendor string, event model.UpSwingWebhookEvent) error

	// private
	extractTenure(tenure string) (int, int)
}

type webhookServiceImpl struct {
	webhookDAO dao.WebhooksEventsDAO
}

func DefaultWebhookService() WebhookService {
	return &webhookServiceImpl{webhookDAO: dao.DefaultWebhookEventsDAO()}
}

func (w *webhookServiceImpl) RegisterNewEvent(ctx context.Context, vendor string, event model.UpSwingWebhookEvent) error {
	entity := entity.WebhookEvent{ClientCode: event.Pci, Vendor: vendor, TrackingId: event.JourneyID, EventType: event.EventType, Institution: event.Fsi, Type: event.TermDepositType, Amount: event.Amount, FailureReason: event.Reason, CreatedBy: "webhook-api", UpdatedBy: "webhook-api"}
	entity.TenureMonths, entity.TenureDays = w.extractTenure(event.Tenure)
	err := w.webhookDAO.SaveNewEvent(ctx, entity)
	if err != nil {
		return goerr.New(err, "service: webhook event registration failed")
	}
	return nil
}

// sample formats: 22m0d 0m1826d 0m582d 11m11d
func (w *webhookServiceImpl) extractTenure(tenure string) (int, int) {
	if tenure != "" && len(tenure) > 0 {
		monthsIndex := strings.Index(tenure, "m")
		daysIndex := strings.Index(tenure, "d")
		months, _ := strconv.Atoi(tenure[0:monthsIndex])
		days, _ := strconv.Atoi(tenure[monthsIndex+1 : daysIndex])
		if months == 0 && days > 31 {
			months = int(math.Round(float64(days) / avgDays))
			days = 0
		}
		return months, days
	}
	return 0, 0
}
