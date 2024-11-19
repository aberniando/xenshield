package app

import (
	"github.com/aberniando/xenshield/internal/handler"
	"github.com/aberniando/xenshield/pkg/logger"
)

type Handlers struct {
	TransactionHandler *handler.TransactionHandler
}

func InitHandlers(services *Services, logger *logger.Logger) *Handlers {
	return &Handlers{
		TransactionHandler: handler.NewTransactionHandler(services.TransactionService, logger),
	}
}
