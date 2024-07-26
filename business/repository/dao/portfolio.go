package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/angel-one/fd-core/business/model"
	"github.com/angel-one/fd-core/business/repository/entity"
	"github.com/angel-one/fd-core/commons/cache"
	"github.com/angel-one/fd-core/commons/database"
	"github.com/angel-one/fd-core/constants"
	"github.com/angel-one/goerr"
)

type PortfolioDAO interface {
	FindByClient(ctx context.Context, clientCode string, provider string) (*entity.PortfolioEntity, error)
	FetchPortfolioFromRedis(ctx context.Context, clientCode string, provider string) (model.Portfolio, error)
	FetchClientList(ctx context.Context, provider string, instantRefresh bool) ([]string, error)
	BatchUpdatePortfolio(ctx context.Context, portfolioUpdateEntities []entity.PortfolioEntity) error
	UpdateRefreshedPortfolioClientList(ctx context.Context, provider string, clientList []string) error
	CleanStaleRecords(ctx context.Context) error
}

type portfolioDAOImpl struct {
	db *sql.DB
}

func DefaultPortfolioDAO() PortfolioDAO {
	return &portfolioDAOImpl{db: database.GetDBPool(true)}
}

func (p *portfolioDAOImpl) FindByClient(ctx context.Context, clientCode string, provider string) (*entity.PortfolioEntity, error) {
	var entity entity.PortfolioEntity
	err := p.db.QueryRowContext(ctx, PortfolioByClientCode, clientCode, provider).Scan(&entity.ClientCode, &entity.TotalActiveDeposits, &entity.Provider, &entity.InvestedValue, &entity.CurrentValue, &entity.InterestEarned, &entity.ReturnsValue, &entity.ReturnsPercentage)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, goerr.New(err, fmt.Sprintf("dao failed: find by portfolio failed for clientCode: %s", clientCode))
		}
	}
	return &entity, nil
}

func (p *portfolioDAOImpl) FetchClientList(ctx context.Context, provider string, instantRefresh bool) ([]string, error) {
	var clientList []string
	var rows *sql.Rows
	var err error
	if instantRefresh {
		rows, err = p.db.QueryContext(ctx, FetchRefreshPortfolioClientListByProvider, provider, false, true)
	} else {
		rows, err = p.db.QueryContext(ctx, FetchPortfolioClientListByProvider, provider, false)
	}

	if err != nil && err != sql.ErrNoRows {
		return clientList, err
	}

	defer rows.Close()
	for rows.Next() {
		var clientCode string

		err := rows.Scan(
			&clientCode,
		)
		if err != nil {
			return clientList, err
		}
		clientList = append(clientList, clientCode)
	}
	return clientList, nil
}

func (p *portfolioDAOImpl) BatchUpdatePortfolio(ctx context.Context, portfolioUpdateEntities []entity.PortfolioEntity) error {
	if len(portfolioUpdateEntities) == 0 {
		return nil
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(InsertClientPortfolio)

	values := []interface{}{}
	valueStrings := []string{}
	paramIndex := 1

	for _, portfolioUpdateEntity := range portfolioUpdateEntities {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			paramIndex, paramIndex+1, paramIndex+2, paramIndex+3, paramIndex+4, paramIndex+5, paramIndex+6, paramIndex+7, paramIndex+8, paramIndex+9, paramIndex+10, paramIndex+11))
		values = append(values, portfolioUpdateEntity.ClientCode, portfolioUpdateEntity.Provider, portfolioUpdateEntity.TotalActiveDeposits, portfolioUpdateEntity.InvestedValue, portfolioUpdateEntity.CurrentValue, portfolioUpdateEntity.InterestEarned, portfolioUpdateEntity.InterestEarned, portfolioUpdateEntity.ReturnsPercentage, portfolioUpdateEntity.CreatedBy, portfolioUpdateEntity.UpdatedBy, portfolioUpdateEntity.InvalidClient, portfolioUpdateEntity.ApiError)
		paramIndex += 12
	}

	queryBuilder.WriteString(strings.Join(valueStrings, ", "))
	queryBuilder.WriteString(UpdateClientPortfolio)
	query := queryBuilder.String()

	// Execute the batch update
	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (p *portfolioDAOImpl) UpdateRefreshedPortfolioClientList(ctx context.Context, provider string, clientList []string) error {
	if len(clientList) == 0 {
		return nil
	}

	placeholders := make([]string, len(clientList))
	args := make([]interface{}, len(clientList))
	for i, client := range clientList {
		placeholders[i] = "$" + fmt.Sprintf("%d", i+1)
		args[i] = client
	}

	placeholderString := strings.Join(placeholders, ", ")

	query := fmt.Sprintf(UpdateRefreshPortfolioClientList, placeholderString)

	_, err := p.db.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (p *portfolioDAOImpl) CleanStaleRecords(ctx context.Context) error {

	_, err := p.db.ExecContext(ctx, CleanStalePortfolioRecords)

	if err != nil {
		return err
	}

	return nil
}

func (p *portfolioDAOImpl) FetchPortfolioFromRedis(ctx context.Context, clientCode string, provider string) (model.Portfolio, error) {
	var portfolioData model.Portfolio
	var err error

	response, err := cache.GetRedisClient().GetMultipleHashStringMap(ctx, constants.PortfolioRedisKey)
	if err != nil {
		return portfolioData, goerr.New(err, fmt.Sprintf("dao failed: fetch portfolio details from redis failed for clientCode: %s", clientCode))
	} else {
		clientPortfolio := response[constants.PortfolioRedisKey][clientCode]

		err = json.Unmarshal([]byte(clientPortfolio), &portfolioData)
		if err != nil {
			return portfolioData, fmt.Errorf("error while unmarshalling portfolio data fetched from redis for key: %s and clientcode: %s", constants.PortfolioRedisKey, clientCode)
		}
	}

	return portfolioData, nil
}
