package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/commons/database"
)

type PlansDAO interface {
	FetchAllPlansDetails(ctx context.Context) ([]model.Plan, error)
	FetchFsiPlansDetails(ctx context.Context, fsi string) (model.FsiPlans, error)
	FetchAllFDDetails(ctx context.Context) ([]model.Plan, error)
	FetchMostBoughtDetails(ctx context.Context) ([]model.Plan, error)
}

type plansDAOImpl struct {
	db *sql.DB
}

func DefaultPlansDAO() PlansDAO {
	return &plansDAOImpl{db: database.GetDBPool(true)}
}

func (d *plansDAOImpl) FetchAllFDDetails(ctx context.Context) ([]model.Plan, error) {
	var allFDs []model.Plan
	rows, err := d.db.QueryContext(ctx, FetchAllFDDetails)
	if err != nil && err != sql.ErrNoRows {
		return allFDs, fmt.Errorf("%s%w", "Error while fetching All FD details: ", err)
	}

	defer rows.Close()
	for rows.Next() {
		var PlanDetail model.Plan
		var isMostBought bool
		err := rows.Scan(
			&PlanDetail.Fsi,
			&PlanDetail.Name,
			&PlanDetail.Type,
			&PlanDetail.TenureYears,
			&PlanDetail.TenureMonths,
			&PlanDetail.TenureDays,
			&PlanDetail.InterestRate,
			&PlanDetail.LockinMonths,
			&PlanDetail.WomenBenefit,
			&PlanDetail.SeniorCitizen,
			&PlanDetail.ImageURL,
			&isMostBought,
			&PlanDetail.Description,
			&PlanDetail.InsuredAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("%s%w", "Error while fetching All FD details: ", err)
		}
		allFDs = append(allFDs, PlanDetail)
	}

	return allFDs, nil
}
func (d *plansDAOImpl) FetchMostBoughtDetails(ctx context.Context) ([]model.Plan, error) {
	var mostBoughtPlans []model.Plan
	rows, err := d.db.QueryContext(ctx, FetchMostBoughtPlanDetails)
	if err != nil && err != sql.ErrNoRows {
		return mostBoughtPlans, fmt.Errorf("%s%w", "Error while fetching Most bought FD details: ", err)
	}

	defer rows.Close()
	for rows.Next() {
		var PlanDetail model.Plan

		err := rows.Scan(
			&PlanDetail.Fsi,
			&PlanDetail.Name,
			&PlanDetail.Type,
			&PlanDetail.TenureYears,
			&PlanDetail.TenureMonths,
			&PlanDetail.TenureDays,
			&PlanDetail.InterestRate,
			&PlanDetail.LockinMonths,
			&PlanDetail.WomenBenefit,
			&PlanDetail.SeniorCitizen,
			&PlanDetail.ImageURL,
			&PlanDetail.Description,
			&PlanDetail.InsuredAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("%s%w", "Error while fetching Most bought FD details: ", err)
		}
		mostBoughtPlans = append(mostBoughtPlans, PlanDetail)
	}

	return mostBoughtPlans, nil
}

func (d *plansDAOImpl) FetchAllPlansDetails(ctx context.Context) ([]model.Plan, error) {
	var allPlans []model.Plan
	rows, err := d.db.QueryContext(ctx, FetchAllPlansDetails)
	if err != nil && err != sql.ErrNoRows {
		return allPlans, err
	}
	defer rows.Close()
	for rows.Next() {
		var PlanDetail model.Plan

		err := rows.Scan(
			&PlanDetail.Fsi,
			&PlanDetail.Name,
			&PlanDetail.Type,
			&PlanDetail.TenureYears,
			&PlanDetail.TenureMonths,
			&PlanDetail.TenureDays,
			&PlanDetail.InterestRate,
			&PlanDetail.LockinMonths,
			&PlanDetail.WomenBenefit,
			&PlanDetail.SeniorCitizen,
			&PlanDetail.ImageURL,
			&PlanDetail.Description,
			&PlanDetail.InsuredAmount,
		)
		if err != nil {
			return nil, err
		}
		allPlans = append(allPlans, PlanDetail)
	}

	return allPlans, nil
}

func (d *plansDAOImpl) FetchFsiPlansDetails(ctx context.Context, fsi string) (model.FsiPlans, error) {
	var fsiPlans model.FsiPlans
	var insuredAmount, minInvestmentAmount int
	var compareFsi, compareFsiName, compareFsiImageUrl string
	var compareFsiInterestRate float64
	var aboutData, calculator []byte

	rows, err := d.db.QueryContext(ctx, FetchFsiPlansDetails, fsi)
	if err != nil && err != sql.ErrNoRows {
		return fsiPlans, err
	}
	defer rows.Close()
	for rows.Next() {
		var PlanDetail model.Plan

		err := rows.Scan(
			&PlanDetail.Fsi,
			&PlanDetail.Name,
			&PlanDetail.Type,
			&PlanDetail.TenureYears,
			&PlanDetail.TenureMonths,
			&PlanDetail.TenureDays,
			&PlanDetail.InterestRate,
			&PlanDetail.LockinMonths,
			&PlanDetail.WomenBenefit,
			&PlanDetail.SeniorCitizen,
			&PlanDetail.ImageURL,
			&aboutData,
			&calculator,
			&PlanDetail.Description,
			&insuredAmount,
			&minInvestmentAmount,
			&compareFsi,
			&compareFsiName,
			&compareFsiInterestRate,
			&compareFsiImageUrl,
		)
		if err != nil {
			return fsiPlans, err
		}
		PlanDetail.InsuredAmount = insuredAmount
		fsiPlans.Plans = append(fsiPlans.Plans, PlanDetail)
	}
	fsiPlans.InsuredAmount = insuredAmount
	fsiPlans.MinInvestment = minInvestmentAmount
	fsiPlans.CompareFsi = compareFsi
	fsiPlans.CompareFsiName = compareFsiName
	fsiPlans.CompareFsiInterestRate = compareFsiInterestRate
	fsiPlans.CompareFsiImageUrl = compareFsiImageUrl

	err = json.Unmarshal(aboutData, &fsiPlans.About)
	if err != nil {
		return fsiPlans, err
	}

	err = json.Unmarshal(calculator, &fsiPlans.Calculator)
	if err != nil {
		return fsiPlans, err
	}
	return fsiPlans, nil
}
