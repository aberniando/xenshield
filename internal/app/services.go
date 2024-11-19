package app

import (
	"github.com/aberniando/xenshield/internal/usecases/transaction"
	"github.com/aberniando/xenshield/pkg/logger"
)

type Services struct {
	TransactionService transaction.Service
}

func InitServices(repositories *Repositories) *Services {
	return &Services{
		TransactionService: transaction.NewService(repositories.TransactionRepository, logger.GetLogger()),
	}
}
