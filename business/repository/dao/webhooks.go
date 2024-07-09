package dao

import (
	"context"
	"database/sql"

	"github.com/angel-one/fd-core/business/repository/entity"
	"github.com/angel-one/fd-core/commons/database"
	"github.com/angel-one/goerr"
)

type WebhooksEventsDAO interface {
	SaveNewEvent(ctx context.Context, entity entity.WebhookEvent) error
}

type webhooksDAOImpl struct {
	db *sql.DB
}

func DefaultWebhookEventsDAO() WebhooksEventsDAO {
	return &webhooksDAOImpl{db: database.GetDBPool(true)}
}

func (d *webhooksDAOImpl) SaveNewEvent(ctx context.Context, entity entity.WebhookEvent) error {
	_, err := d.db.ExecContext(ctx, InsertWebookEvent, entity.ClientCode, entity.Vendor, entity.TrackingId, entity.EventType, entity.Institution, entity.Type, entity.Amount, entity.TenureMonths, entity.TenureDays, entity.FailureReason, entity.CreatedBy, entity.UpdatedBy)
	if err != nil {
		return goerr.New(err, "dao failed: new webhook save failed")
	}
	return nil
}
