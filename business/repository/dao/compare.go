package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/commons/database"
)

type CompareDAO interface {
	FetchCompareList(ctx context.Context) ([]model.FsiDetails, error)
	FetchCompareFsiDetails(ctx context.Context, fsis []string) ([]model.CompareFSIDBDetails, error)
}

type compareDAOImpl struct {
	db *sql.DB
}

func DefaultCompareDAO() CompareDAO {
	return &compareDAOImpl{db: database.GetDBPool(true)}
}

func (d *compareDAOImpl) FetchCompareList(ctx context.Context) ([]model.FsiDetails, error) {
	var fsiList []model.FsiDetails
	rows, err := d.db.QueryContext(ctx, CompareLandingPageQuery)
	if err != nil && err != sql.ErrNoRows {
		return fsiList, fmt.Errorf("%s%w", "Error while fetching FSI list to compare: ", err)
	}

	defer rows.Close()
	for rows.Next() {
		var FsiDetail model.FsiDetails

		err := rows.Scan(
			&FsiDetail.Fsi,
			&FsiDetail.Name,
			&FsiDetail.ImageUrl,
			&FsiDetail.InterestRate,
		)
		if err != nil {
			return nil, fmt.Errorf("%s%w", "Error while fetching FSI list to compare: ", err)
		}
		fsiList = append(fsiList, FsiDetail)
	}

	return fsiList, nil
}

func (d *compareDAOImpl) FetchCompareFsiDetails(ctx context.Context, fsis []string) ([]model.CompareFSIDBDetails, error) {
	var compareFSIDBDetails []model.CompareFSIDBDetails

	fsiPlaceholders := make([]string, len(fsis))
	placeholderValues := make([]interface{}, len(fsis))

	for i, fsi := range fsis {
		fsiPlaceholders[i] = fmt.Sprintf("$%d", i+1)
		placeholderValues[i] = fsi
	}
	quotedPlaceholderString := strings.Join(fsiPlaceholders, ", ")

	query := fmt.Sprintf(CompareFsiQuery, quotedPlaceholderString)

	rows, err := d.db.QueryContext(ctx, query, placeholderValues...)
	if err != nil && err != sql.ErrNoRows {
		return compareFSIDBDetails, fmt.Errorf("%s%w", "Error while fetching compare FSI Details: ", err)
	}

	defer rows.Close()
	for rows.Next() {
		var compareFSIDBDetail model.CompareFSIDBDetails

		err := rows.Scan(
			&compareFSIDBDetail.FSI,
			&compareFSIDBDetail.Name,
			&compareFSIDBDetail.TenureYears,
			&compareFSIDBDetail.TenureMonths,
			&compareFSIDBDetail.TenureDays,
			&compareFSIDBDetail.InterestRate,
			&compareFSIDBDetail.MinDeposit,
			&compareFSIDBDetail.SeniorCitizenBenefit,
			&compareFSIDBDetail.BankAccount,
			&compareFSIDBDetail.InsuredAmount,
			&compareFSIDBDetail.ImageURL,
		)
		if err != nil {
			return nil, fmt.Errorf("%s%w", "Error while fetching Compare Details: ", err)
		}
		compareFSIDBDetails = append(compareFSIDBDetails, compareFSIDBDetail)
	}

	return compareFSIDBDetails, nil
}
