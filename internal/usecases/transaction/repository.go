package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/aberniando/xenshield/pkg/logger"
	"github.com/jmoiron/sqlx"

	"github.com/aberniando/xenshield/internal/entity"
)

type repository struct {
	db     *sqlx.DB
	logger *logger.Logger
}

func NewRepository(db *sqlx.DB, logger *logger.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

//go:generate mockgen -source=repository.go -destination=./mocks/repository_mock.go -package=mock
type Repository interface {
	InsertTransaction(ctx context.Context, transaction *entity.Transaction) error
	HasTransaction(ctx context.Context, ipAddress string) (bool, error)
	GetIPAddressStolenCardHistory(ctx context.Context, ipAddress string) (bool, error)
}

func (r *repository) InsertTransaction(ctx context.Context, transaction *entity.Transaction) error {

	err := r.db.GetContext(ctx, transaction, `
		INSERT INTO transactions (ip_address, masked_card_number, status, reason)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created, updated
	`, transaction.IPAddress, transaction.MaskedCardNumber, transaction.Status, transaction.Reason)

	if err != nil {
		r.logger.Error(err.Error())
		return fmt.Errorf("error when inserting new transaction: %w", err)
	}

	return nil
}

func (r *repository) HasTransaction(ctx context.Context, ipAddress string) (bool, error) {

	var hasTransaction bool

	err := r.db.GetContext(ctx, &hasTransaction, `
		SELECT EXISTS (
    		SELECT 1 
    		FROM transactions 
    		WHERE ip_address = $1
		) AS has_transaction`, ipAddress)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		r.logger.Error(err.Error())
		return false, fmt.Errorf("failed to check for transactions for IP address %s: %w", ipAddress, err)
	}

	return hasTransaction, nil
}

func (r *repository) GetIPAddressStolenCardHistory(ctx context.Context, ipAddress string) (bool, error) {
	var result bool
	err := r.db.GetContext(ctx, &result, `
		SELECT EXISTS(
		    SELECT 1
    		FROM transactions
    		WHERE ip_address = $1
      			AND status = 'FAILED'
      			AND reason = 'STOLEN_CARD'
		) AS linked_to_stolen_card`, ipAddress)

	if err != nil {
		r.logger.Error(err.Error())
		return false, fmt.Errorf("failed to check stolen card history for IP address %s: %w", ipAddress, err)
	}

	return result, nil
}
