package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/commons/database"
)

type FsiDetailsDAO interface {
	FetchFsiDetailsList(ctx context.Context, fsis []string) ([]model.FsiStruct, error)
}

type fsiDetailsDAOImpl struct {
	db *sql.DB
}

func DefaultFsiDetailsDAO() FsiDetailsDAO {
	return &fsiDetailsDAOImpl{db: database.GetDBPool(true)}
}

func (d *fsiDetailsDAOImpl) FetchFsiDetailsList(ctx context.Context, fsis []string) ([]model.FsiStruct, error) {
	var FsiStructs []model.FsiStruct
	var aboutData, calculator []byte
	var faqData json.RawMessage

	fsiPlaceholders := make([]string, len(fsis))
	placeholderValues := make([]interface{}, len(fsis))

	for i, fsi := range fsis {
		fsiPlaceholders[i] = fmt.Sprintf("$%d", i+1)
		placeholderValues[i] = fsi
	}
	quotedPlaceholderString := strings.Join(fsiPlaceholders, ", ")
	query := fmt.Sprintf(FsiDetailsQueryTest, quotedPlaceholderString)

	rows, err := d.db.QueryContext(ctx, query, placeholderValues...)
	if err != nil && err != sql.ErrNoRows {
		return FsiStructs, fmt.Errorf("%s%w", "Error while fetching FSI Details: ", err)
	}

	defer rows.Close()

	fsiDetailsMap := make(map[string]*model.FsiStruct)

	for rows.Next() {
		var FsiDetail model.FsiDetailPlans
		err := rows.Scan(
			&FsiDetail.Fsi,
			&FsiDetail.Name,
			&FsiDetail.Type,
			&FsiDetail.InterestRate,
			&FsiDetail.LockinMonths,
			&FsiDetail.WomenBenefit,
			&FsiDetail.SeniorCitizen,
			&FsiDetail.ImageURL,
			&FsiDetail.Description,
			&FsiDetail.InsuredAmount,
			&aboutData,
			&calculator,
		)
		if err != nil {
			return FsiStructs, err
		}

		fsiStruct, exists := fsiDetailsMap[FsiDetail.Fsi]
		if !exists {
			fsiStruct = &model.FsiStruct{}
			err = json.Unmarshal(aboutData, &fsiStruct.About)
			if err != nil {
				return FsiStructs, err
			}

			err = json.Unmarshal(calculator, &fsiStruct.Calculator)
			if err != nil {
				return FsiStructs, err
			}

			fsiDetailsMap[FsiDetail.Fsi] = fsiStruct
		}
		fsiStruct.Plans = append(fsiStruct.Plans, FsiDetail)
	}

	// Fetching FAQs
	for fsi, fsiStruct := range fsiDetailsMap {
		rows, err := d.db.QueryContext(ctx, GetFAQsByTag, fsi)
		if err != nil {
			return FsiStructs, err
		}
		defer rows.Close()

		var faq []model.FAQ
		for rows.Next() {
			err := rows.Scan(&faqData)
			if err != nil {
				return FsiStructs, err
			}
			err = json.Unmarshal(faqData, &faq)
			if err != nil {
				return FsiStructs, err
			}
		}
		fsiStruct.FAQs = faq
	}

	// Converting map to slice
	for _, fsiStruct := range fsiDetailsMap {
		FsiStructs = append(FsiStructs, *fsiStruct)
	}

	return FsiStructs, nil

}
