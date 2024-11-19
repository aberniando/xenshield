package transaction_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aberniando/xenshield/internal/entity"
	"github.com/aberniando/xenshield/internal/enum"
	"github.com/aberniando/xenshield/internal/usecases/transaction"
	"github.com/aberniando/xenshield/pkg/logger"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"regexp"
	"time"
)

var _ = Describe("Repository", func() {
	var (
		db     *sql.DB
		mockDB sqlmock.Sqlmock
		repo   transaction.Repository
		err    error
	)

	BeforeEach(func() {
		db, mockDB, err = sqlmock.New()
		Expect(err).Should(BeNil())
		Expect(err).ShouldNot(HaveOccurred())

		repo = transaction.NewRepository(sqlx.NewDb(db, "sqlmock"), logger.GetLogger())
	})

	AfterEach(func() {
		err := mockDB.ExpectationsWereMet()
		Expect(err).Should(BeNil())
		_ = db.Close()
	})

	Context("InsertTransaction", func() {
		var (
			transaction   entity.Transaction
			transactionID string
			created       time.Time
			updated       time.Time
		)

		BeforeEach(func() {
			transactionID = uuid.NewString()
			created = time.Now()
			updated = time.Now()
			fr := enum.FailureReasonInsufficientBalance

			transaction = entity.Transaction{
				IPAddress:        "255.255.255.255",
				MaskedCardNumber: "XXXX-XXXX-XXXX-0856",
				Status:           enum.TransactionStatusFailed,
				Reason:           &fr,
			}
		})
		When("error inserting new transaction", func() {
			BeforeEach(func() {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO transactions (ip_address, masked_card_number, status, reason)
					VALUES ($1, $2, $3, $4)
					RETURNING id, created, updated`)).WillReturnError(sql.ErrConnDone)
			})

			It("should return error", func() {
				err := repo.InsertTransaction(context.Background(), &transaction)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("error when inserting new transaction: sql: connection is already closed"))
			})
		})

		When("success inserting new transaction", func() {
			BeforeEach(func() {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					INSERT INTO transactions (ip_address, masked_card_number, status, reason)
					VALUES ($1, $2, $3, $4)
					RETURNING id, created, updated`)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created", "updated"}).
						AddRow(transactionID, created, updated))
			})

			It("should not return error", func() {
				err := repo.InsertTransaction(context.Background(), &transaction)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(transaction.ID).Should(Equal(transactionID))
				Expect(transaction.Created).Should(Equal(created))
				Expect(transaction.Updated).Should(Equal(updated))
			})
		})
	})

	Context("HasTransaction", func() {
		var ipAddress string
		BeforeEach(func() {
			ipAddress = uuid.NewString()[0:15]
		})
		When("error checking transactions for given ip address", func() {
			BeforeEach(func() {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					SELECT EXISTS (
    					SELECT 1 
						FROM transactions 
						WHERE ip_address = $1
					) AS has_transaction`)).WillReturnError(sql.ErrConnDone)
			})

			It("should return error", func() {
				result, err := repo.HasTransaction(context.Background(), ipAddress)
				Expect(result).Should(Equal(false))
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(fmt.Sprintf("failed to check for transactions for IP address %s: sql: connection is already closed", ipAddress)))
			})
		})

		When("success checking transactions for given ip address (has no transaction)", func() {
			BeforeEach(func() {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					SELECT EXISTS (
    					SELECT 1 
						FROM transactions 
						WHERE ip_address = $1
					) AS has_transaction`)).
					WillReturnError(sql.ErrNoRows)
			})

			It("should return no error", func() {
				result, err := repo.HasTransaction(context.Background(), ipAddress)
				Expect(err).Should(BeNil())
				Expect(result).Should(Equal(false))
			})
		})

		When("success checking transactions for given ip address (has transaction)", func() {
			BeforeEach(func() {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					SELECT EXISTS (
    					SELECT 1 
						FROM transactions 
						WHERE ip_address = $1
					) AS has_transaction`)).
					WillReturnRows(sqlmock.NewRows([]string{"has_transaction"}).
						AddRow(true))
			})

			It("should return no error", func() {
				result, err := repo.HasTransaction(context.Background(), ipAddress)
				Expect(err).Should(BeNil())
				Expect(result).Should(Equal(true))
			})
		})
	})

	Context("GetIPAddressStolenCardHistory", func() {
		var ipAddress string
		BeforeEach(func() {
			ipAddress = uuid.NewString()[0:15]
		})
		When("error getting ip address' stolen card history", func() {
			BeforeEach(func() {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					SELECT EXISTS(
		    			SELECT 1
    					FROM transactions
    					WHERE ip_address = $1
      						AND status = 'FAILED'
      						AND reason = 'STOLEN_CARD'
    					) AS linked_to_stolen_card`)).
					WillReturnError(sql.ErrConnDone)
			})

			It("should return error", func() {
				result, err := repo.GetIPAddressStolenCardHistory(context.Background(), ipAddress)
				Expect(result).Should(Equal(false))
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(fmt.Sprintf("failed to check stolen card history for IP address %s: sql: connection is already closed", ipAddress)))
			})
		})

		When("success checking transactions for given ip address", func() {
			BeforeEach(func() {
				mockDB.ExpectQuery(regexp.QuoteMeta(`
					SELECT EXISTS(
		    			SELECT 1
    					FROM transactions
    					WHERE ip_address = $1
      						AND status = 'FAILED'
      						AND reason = 'STOLEN_CARD'
    					) AS linked_to_stolen_card`)).
					WillReturnRows(sqlmock.NewRows([]string{"linked_to_stolen_card"}).
						AddRow(true))
			})

			It("should return no error", func() {
				result, err := repo.GetIPAddressStolenCardHistory(context.Background(), ipAddress)
				Expect(err).Should(BeNil())
				Expect(result).Should(Equal(true))
			})
		})
	})
})
