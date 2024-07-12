package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/repository/entity"
	"github.com/angel-one/fd-core/commons/database"
	"github.com/angel-one/goerr"
)

type HomeInfoDAO interface {
	FetchAllFDDetails(ctx context.Context) ([]model.Plan, error)
	FetchMostBoughtDetails(ctx context.Context) ([]model.Plan, error)
	FetchPendingJourneyDetails(ctx context.Context, clientCode string, provider string) (*entity.PendingJourneyEntity, error)
	FetchFAQDetails(ctx context.Context, tag string) (json.RawMessage, error)
}

type homeInfoDAOImpl struct {
	db *sql.DB
}

func DefaultHomeInfoDAO() HomeInfoDAO {
	return &homeInfoDAOImpl{db: database.GetDBPool(true)}
}

func (d *homeInfoDAOImpl) FetchAllFDDetails(ctx context.Context) ([]model.Plan, error) {
	var allFDs []model.Plan

	rows, err := d.db.QueryContext(ctx, FetchFDHighIR)
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

	rows, err = d.db.QueryContext(ctx, FetchFDMinIR)
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

func (d *homeInfoDAOImpl) FetchMostBoughtDetails(ctx context.Context) ([]model.Plan, error) {
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

func (p *homeInfoDAOImpl) FetchPendingJourneyDetails(ctx context.Context, clientCode string, provider string) (*entity.PendingJourneyEntity, error) {
	var entity entity.PendingJourneyEntity
	err := p.db.QueryRowContext(ctx, FetchPendingForClient, clientCode, provider).Scan(&entity.Pending, &entity.Payment, &entity.KYC)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, goerr.New(err, fmt.Sprintf("dao failed: fetch pending journey failed for clientCode: %s", clientCode))
		}
	}
	return &entity, nil
}

func (d *homeInfoDAOImpl) FetchFAQDetails(ctx context.Context, tag string) (json.RawMessage, error) {
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
