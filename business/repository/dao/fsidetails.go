package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

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

	// fsiPlaceholders := make([]string, len(fsis))
	// placeholderValues := make([]interface{}, len(fsis))

	// for i, fsi := range fsis {
	// 	fsiPlaceholders[i] = fmt.Sprintf("$%d", i+1)
	// 	placeholderValues[i] = fsi
	// }

	for _, fsi := range fsis {
		rows, err := d.db.QueryContext(ctx, FsiDetailsQueryTest, fsi)
		if err != nil && err != sql.ErrNoRows {
			return FsiStructs, fmt.Errorf("%s%w", "Error while fetching FSI Details: ", err)
		}
		var FsiStruct model.FsiStruct
		defer rows.Close()
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
			FsiStruct.Plans = append(FsiStruct.Plans, FsiDetail)

			err = json.Unmarshal(aboutData, &FsiStruct.About)
			if err != nil {
				return FsiStructs, err
			}

			err = json.Unmarshal(calculator, &FsiStruct.Calculator)
			if err != nil {
				return FsiStructs, err
			}

			var faq []model.FAQ
			tag := FsiDetail.Fsi
			rows, err := d.db.QueryContext(ctx, GetFAQsByTag, tag)
			if err != nil {
				return FsiStructs, err
			}
			defer rows.Close()
			for rows.Next() {
				err := rows.Scan(&faqData)
				if err != nil {
					return FsiStructs, err
				}
			}

			err = json.Unmarshal(faqData, &faq)
			if err != nil {
				return FsiStructs, err
			}

			FsiStruct.FAQs = faq
			FsiStructs = append(FsiStructs, FsiStruct)
		}
		//quotedPlaceholderString := strings.Join(fsiPlaceholders, ", ")
	}
	return FsiStructs, nil

}
