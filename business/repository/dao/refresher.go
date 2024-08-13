package dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/angel-one/fd-core/business/repository/entity"
	"github.com/angel-one/fd-core/commons/database"
	"github.com/angel-one/goerr"
)

type PendingJourneyDAOTest interface {
	FetchClientListTest(ctx context.Context, provider string, instantRefresh bool) ([]string, error)
	FetchPendingJourneyDetailsTest(ctx context.Context, clientCode string, provider string) (*entity.PendingJourneyEntity, error)
	BatchUpdatePendingJourneyTest(ctx context.Context, pendingJourneyEntities []entity.PendingJourneyEntity) error
	UpdateRefreshedPendingJourneyClientListTest(ctx context.Context, provider string, clientList []string) error
	CleanStaleRecordsTest(ctx context.Context) error
}

type pendingJourneyDAOImplTest struct {
	db *sql.DB
}

func DefaultPendingJourneyDAOTest() PendingJourneyDAOTest {
	return &pendingJourneyDAOImplTest{db: database.GetDBPool(true)}
}

func (p *pendingJourneyDAOImplTest) FetchPendingJourneyDetailsTest(ctx context.Context, clientCode string, provider string) (*entity.PendingJourneyEntity, error) {
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

func (p *pendingJourneyDAOImplTest) FetchClientListTest(ctx context.Context, provider string, instantRefresh bool) ([]string, error) {
	var clientList []string
	var rows *sql.Rows
	var err error

	if instantRefresh {
		rows, err = p.db.QueryContext(ctx, FetchRefreshPendingJourneyClientListByProviderTest, provider, false, true)
	} else {
		rows, err = p.db.QueryContext(ctx, FetchPendingJourneyClientListByProviderTest, provider, false)
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

	// if !instantRefresh {
	// 	fsiPlaceholders := make([]string, len(clientList))

	// 	for i, _ := range clientList {
	// 		fsiPlaceholders[i] = fmt.Sprintf("$%d", i+1)
	// 	}

	// 	quotedPlaceholderString := strings.Join(fsiPlaceholders, ", ")

	// 	query := fmt.Sprintf(FetchPendingJourneyWithPending, quotedPlaceholderString)

	// 	rows, err = p.db.QueryContext(ctx, query, clientList)
	// }

	return clientList, nil
}

func (p *pendingJourneyDAOImplTest) BatchUpdatePendingJourneyTest(ctx context.Context, pendingJourneyEntities []entity.PendingJourneyEntity) error {
	if len(pendingJourneyEntities) == 0 {
		return nil
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString(InsertPendingJourneyDetailsTest)

	values := []interface{}{}
	valueStrings := []string{}
	paramIndex := 1

	for _, pendingJourneyEntity := range pendingJourneyEntities {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			paramIndex, paramIndex+1, paramIndex+2, paramIndex+3, paramIndex+4, paramIndex+5, paramIndex+6, paramIndex+7, paramIndex+8))
		values = append(values, pendingJourneyEntity.ClientCode, pendingJourneyEntity.Provider, pendingJourneyEntity.Pending, pendingJourneyEntity.Payment, pendingJourneyEntity.KYC, pendingJourneyEntity.CreatedBy, pendingJourneyEntity.UpdatedBy, pendingJourneyEntity.InvalidClient, pendingJourneyEntity.ApiError)
		paramIndex += 9
	}

	queryBuilder.WriteString(strings.Join(valueStrings, ", "))
	queryBuilder.WriteString(UpdatePendingJourneyDetails)
	query := queryBuilder.String()

	// Execute the batch update
	_, err := p.db.ExecContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (p *pendingJourneyDAOImplTest) UpdateRefreshedPendingJourneyClientListTest(ctx context.Context, provider string, clientList []string) error {
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

	query := fmt.Sprintf(UpdateRefreshPendingJourneyClientListTest, placeholderString)

	_, err := p.db.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (p *pendingJourneyDAOImplTest) CleanStaleRecordsTest(ctx context.Context) error {

	_, err := p.db.ExecContext(ctx, CleanStalePendingJourneyRecordsTest)

	if err != nil {
		return err
	}

	return nil
}
