package transaction

import (
	"context"
	"github.com/aberniando/xenshield/internal/entity"
	loggerPkg "github.com/aberniando/xenshield/pkg/logger"
)

type service struct {
	repo   Repository
	logger *loggerPkg.Logger
}

func NewService(repo Repository, logger *loggerPkg.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

//go:generate mockgen -source=service.go -destination=./mocks/service_mock.go -package=mock
type Service interface {
	InsertTransaction(ctx context.Context, request *entity.InsertTransactionRequest) (*entity.InsertTransactionResponse, error)
	GetIPAddressStolenCardHistory(ctx context.Context, ipAddress string) (*entity.GetIPAddressStolenCardHistoryResponse, error)
}

func (svc *service) InsertTransaction(ctx context.Context, request *entity.InsertTransactionRequest) (*entity.InsertTransactionResponse, error) {
	transaction := request.ToTransaction()
	err := svc.repo.InsertTransaction(ctx, transaction)
	if err != nil {
		svc.logger.Error("[Service] InsertTransaction: InsertTransaction repository function returned error")
		return nil, err
	}

	return transaction.ToInsertTransactionResponse(), nil
}

func (svc *service) GetIPAddressStolenCardHistory(ctx context.Context, ipAddress string) (*entity.GetIPAddressStolenCardHistoryResponse, error) {
	response := entity.GetIPAddressStolenCardHistoryResponse{
		HasTransaction: true,
	}

	hasTransaction, err := svc.repo.HasTransaction(ctx, ipAddress)
	if err != nil {
		svc.logger.Error("[Service] GetIPAddressStolenCardHistory: HasTransaction repository function returned error")
		return nil, err
	}

	if !hasTransaction {
		response.HasTransaction = false
		return &response, nil
	}

	isLinkedToStolenCard, err := svc.repo.GetIPAddressStolenCardHistory(ctx, ipAddress)
	if err != nil {
		svc.logger.Error("[Service] GetIPAddressStolenCardHistory: GetIPAddressStolenCardHistory repository function returned error")
		return nil, err
	}

	response.LinkedToStolenCard = isLinkedToStolenCard
	return &response, nil
}
