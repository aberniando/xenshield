package transaction_test

import (
	"context"
	"errors"
	"github.com/aberniando/xenshield/internal/entity"
	"github.com/aberniando/xenshield/internal/enum"
	"github.com/aberniando/xenshield/internal/usecases/transaction"
	mock "github.com/aberniando/xenshield/internal/usecases/transaction/mocks"
	"github.com/aberniando/xenshield/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Service", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mock.MockRepository
		svc      transaction.Service
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mock.NewMockRepository(ctrl)
		svc = transaction.NewService(mockRepo, logger.GetLogger())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("InsertTransaction", func() {
		var (
			request       entity.InsertTransactionRequest
			transactionID string
			created       time.Time
			updated       time.Time
		)

		BeforeEach(func() {
			transactionID = uuid.NewString()
			created = time.Now()
			updated = time.Now()

			request = entity.InsertTransactionRequest{
				IPAddress:        "255.255.255.255",
				MaskedCardNumber: "XXXX-XXXX-XXXX-0856",
				Status:           enum.TransactionStatusFailed,
				FailureReason:    enum.FailureReasonInsufficientBalance,
			}
		})
		When("error inserting new transaction", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().InsertTransaction(gomock.Any(), request.ToTransaction()).
					Return(errors.New("error when inserting new transaction: sql: connection is already closed"))
			})

			It("should return error", func() {
				resp, err := svc.InsertTransaction(context.Background(), &request)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(Equal("error when inserting new transaction: sql: connection is already closed"))
				Expect(resp).Should(BeNil())
			})
		})

		When("success inserting new transaction", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().InsertTransaction(gomock.Any(), request.ToTransaction()).
					DoAndReturn(func(_ context.Context, transaction *entity.Transaction) error {
						transaction.ID = transactionID
						transaction.Created = created
						transaction.Updated = updated
						return nil
					})
			})

			It("should return no error", func() {
				resp, err := svc.InsertTransaction(context.Background(), &request)
				Expect(err).Should(BeNil())
				Expect(resp).ShouldNot(BeNil())
				Expect(resp.ID).Should(Equal(transactionID))
			})
		})
	})

	Context("GetIPAddressStolenCardHistory", func() {
		ipAddress := "255.255.255.255"

		When("error getting IP address transaction activity", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().HasTransaction(gomock.Any(), ipAddress).
					Return(false, errors.New("internal server error"))
			})
			It("should return error", func() {
				resp, err := svc.GetIPAddressStolenCardHistory(context.Background(), ipAddress)
				Expect(resp).Should(BeNil())
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(Equal("internal server error"))
			})
		})

		When("IP address does not have transaction activity", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().HasTransaction(gomock.Any(), ipAddress).
					Return(false, nil)
			})
			It("should return no error", func() {
				resp, err := svc.GetIPAddressStolenCardHistory(context.Background(), ipAddress)
				Expect(resp.HasTransaction).Should(Equal(false))
				Expect(resp.LinkedToStolenCard).Should(Equal(false))
				Expect(err).Should(BeNil())
			})
		})

		When("error getting IP address link to stolen card", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().HasTransaction(gomock.Any(), ipAddress).
					Return(true, nil)
				mockRepo.EXPECT().GetIPAddressStolenCardHistory(gomock.Any(), ipAddress).
					Return(false, errors.New("internal server error"))
			})
			It("should return error", func() {
				resp, err := svc.GetIPAddressStolenCardHistory(context.Background(), ipAddress)
				Expect(resp).Should(BeNil())
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(Equal("internal server error"))
			})
		})

		When("success getting IP address link to stolen card", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().HasTransaction(gomock.Any(), ipAddress).
					Return(true, nil)
				mockRepo.EXPECT().GetIPAddressStolenCardHistory(gomock.Any(), ipAddress).
					Return(true, nil)
			})
			It("should return error", func() {
				resp, err := svc.GetIPAddressStolenCardHistory(context.Background(), ipAddress)
				Expect(err).Should(BeNil())
				Expect(resp.HasTransaction).Should(Equal(true))
				Expect(resp.LinkedToStolenCard).Should(Equal(true))
			})
		})
	})
})
