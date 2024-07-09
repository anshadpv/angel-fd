package dao

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/angel-one/fd-core/commons/database"
)

type FAQDAO interface {
	FetchFAQDetails(ctx context.Context, tag string) (json.RawMessage, error)
}

type faqDAOImpl struct {
	db *sql.DB
}

func DefaultFAQDAO() FAQDAO {
	return &faqDAOImpl{db: database.GetDBPool(true)}
}

func (d *faqDAOImpl) FetchFAQDetails(ctx context.Context, tag string) (json.RawMessage, error) {
	var faqData json.RawMessage
	rows, err := d.db.QueryContext(ctx, GetFAQsByTag, tag)
	if err != nil && err != sql.ErrNoRows {
		return faqData, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&faqData)
		if err != nil {
			return faqData, err
		}
	}

	return faqData, nil
}
