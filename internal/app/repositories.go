package app

import (
	"github.com/aberniando/xenshield/internal/usecases/transaction"
	"github.com/aberniando/xenshield/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type Repositories struct {
	TransactionRepository transaction.Repository
}

func InitRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		TransactionRepository: transaction.NewRepository(db, logger.GetLogger()),
	}
}
