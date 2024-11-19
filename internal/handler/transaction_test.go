package handler_test

import (
	"encoding/json"
	"errors"
	"github.com/aberniando/xenshield/internal/entity"
	"github.com/aberniando/xenshield/internal/handler"
	mock "github.com/aberniando/xenshield/internal/usecases/transaction/mocks"
	l "github.com/aberniando/xenshield/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var _ = Describe("Handler", func() {
	var (
		ctrl               *gomock.Controller
		mockService        *mock.MockService
		transactionHandler *handler.TransactionHandler
		logger             *l.Logger
		ctx                *gin.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockService = mock.NewMockService(ctrl)
		logger = l.GetLogger()
		transactionHandler = handler.NewTransactionHandler(mockService, logger)
	})

	AfterEach(func() {
		ctrl.Finish()
		Expect(true).To(BeTrue())
	})

	Context("Insert", func() {
		When("invalid request body", func() {
			It("return http error 400 bad request", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(``))
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				transactionHandler.Insert(ctx)
				Expect(rec.Code).Should(Equal(http.StatusBadRequest))

				var resp handler.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)

				Expect(resp.Message).Should(Equal("invalid request body"))
			})
		})

		When("validation error request body 1", func() {
			It("return http error 400 bad request", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				transactionHandler.Insert(ctx)
				Expect(rec.Code).Should(Equal(http.StatusBadRequest))

				var resp handler.ValidationErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)

				for _, item := range resp.Message {
					Expect(item.Message).Should(Equal("This field is required"))
				}
			})
		})

		When("validation error request body 2", func() {
			It("return http error 400 bad request", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"ip_address":"n","masked_card_number":"XXXX-XXXX-XXXX-1234x","status":"SUCCESS","failure_reason":"STOLEN_CARD"}`))
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				transactionHandler.Insert(ctx)
				Expect(rec.Code).Should(Equal(http.StatusBadRequest))

				var resp handler.ValidationErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)

				for _, item := range resp.Message {
					switch item.Field {
					case "IPAddress":
						Expect(item.Message).Should(Equal("invalid value format"))
					case "MaskedCardNumber":
						Expect(item.Message).Should(Equal("invalid length"))
					case "FailureReason":
						Expect(item.Message).Should(Equal("failure reason must be null on success transaction"))
					}
				}
			})
		})

		When("validation error request body 3", func() {
			It("return http error 400 bad request", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"ip_address":"n","masked_card_number":"XXXX-XXXX-XXXX-1234x","status":"FAILED","failure_reason":""}`))
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				transactionHandler.Insert(ctx)
				Expect(rec.Code).Should(Equal(http.StatusBadRequest))

				var resp handler.ValidationErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)

				for _, item := range resp.Message {
					switch item.Field {
					case "IPAddress":
						Expect(item.Message).Should(Equal("invalid value format"))
					case "MaskedCardNumber":
						Expect(item.Message).Should(Equal("invalid length"))
					case "FailureReason":
						Expect(item.Message).Should(Equal("failure reason must not be null on failed transaction"))
					}
				}
			})
		})

		When("service return error", func() {
			BeforeEach(func() {
				req := entity.InsertTransactionRequest{}
				_ = json.NewDecoder(strings.NewReader(`{"ip_address":"198.51.100.33","masked_card_number":"XXXX-XXXX-XXXX-1234","status":"SUCCESS","failure_reason":""}`)).Decode(&req)
				mockService.EXPECT().InsertTransaction(gomock.Any(), &req).Return(nil, errors.New("internal server error"))
			})
			It("return http error 500 internal server error", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"ip_address":"198.51.100.33","masked_card_number":"XXXX-XXXX-XXXX-1234","status":"SUCCESS","failure_reason":""}`))
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				transactionHandler.Insert(ctx)
				Expect(rec.Code).Should(Equal(http.StatusInternalServerError))

				var resp handler.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)

				Expect(resp.Message).Should(Equal("Error inserting new transaction"))
			})
		})

		When("success", func() {
			var response entity.InsertTransactionResponse
			BeforeEach(func() {
				response = entity.InsertTransactionResponse{
					ID:               uuid.NewString(),
					IPAddress:        "198.51.100.33",
					MaskedCardNumber: "XXXX-XXXX-XXXX-1234",
					Status:           "SUCCESS",
					Reason:           "",
					Created:          time.Now().Format(time.RFC3339),
					Updated:          time.Now().Format(time.RFC3339),
				}

				req := entity.InsertTransactionRequest{}
				_ = json.NewDecoder(strings.NewReader(`{"ip_address":"198.51.100.33","masked_card_number":"XXXX-XXXX-XXXX-1234","status":"SUCCESS","failure_reason":""}`)).Decode(&req)
				mockService.EXPECT().InsertTransaction(gomock.Any(), &req).Return(&response, nil)
			})
			It("return http error 200 OK", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"ip_address":"198.51.100.33","masked_card_number":"XXXX-XXXX-XXXX-1234","status":"SUCCESS","failure_reason":""}`))
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				transactionHandler.Insert(ctx)
				Expect(rec.Code).Should(Equal(http.StatusOK))

				var actual entity.InsertTransactionResponse
				_ = json.NewDecoder(rec.Body).Decode(&actual)

				Expect(response.ID).Should(Equal(actual.ID))
				Expect(response.IPAddress).Should(Equal(actual.IPAddress))
				Expect(response.MaskedCardNumber).Should(Equal(actual.MaskedCardNumber))
				Expect(response.Status).Should(Equal(actual.Status))
				Expect(response.Reason).Should(Equal(actual.Reason))
			})
		})
	})

	Context("GetIPAddressStolenCardHistory", func() {
		When("empty ip address", func() {
			It("return http error 400 bad request", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				transactionHandler.GetIPAddressStolenCardHistory(ctx)
				Expect(rec.Code).Should(Equal(http.StatusBadRequest))

				var resp handler.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)

				Expect(resp.Message).Should(Equal("ip_address on path param is required"))
			})
		})

		When("invalid ip address", func() {
			It("return http error 400 bad request", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/?ip_address=xxx", nil)
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				ctx.Params = []gin.Param{
					{
						Key:   "ip_address",
						Value: "xx",
					},
				}
				transactionHandler.GetIPAddressStolenCardHistory(ctx)
				Expect(rec.Code).Should(Equal(http.StatusBadRequest))

				var resp handler.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)

				Expect(resp.Message).Should(Equal("invalid ip_address value"))
			})
		})

		When("service return error", func() {
			BeforeEach(func() {
				mockService.EXPECT().GetIPAddressStolenCardHistory(gomock.Any(), "192.0.0.1").Return(nil, errors.New("internal server error"))
			})
			It("return http error 500 internal server error", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				ctx.Params = []gin.Param{
					{
						Key:   "ip_address",
						Value: "192.0.0.1",
					},
				}
				transactionHandler.GetIPAddressStolenCardHistory(ctx)
				Expect(rec.Code).Should(Equal(http.StatusInternalServerError))

				var resp handler.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)

				Expect(resp.Message).Should(Equal("error getting ip address stolen card history"))
			})
		})

		When("ip address has not made any transaction", func() {
			BeforeEach(func() {
				mockService.EXPECT().GetIPAddressStolenCardHistory(gomock.Any(), "192.0.0.1").Return(&entity.GetIPAddressStolenCardHistoryResponse{}, nil)
			})
			It("return http error 404 not found", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/?ip_address=192.0.0.1", nil)
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				ctx.Params = []gin.Param{
					{
						Key:   "ip_address",
						Value: "192.0.0.1",
					},
				}
				transactionHandler.GetIPAddressStolenCardHistory(ctx)
				Expect(rec.Code).Should(Equal(http.StatusNotFound))

				var resp handler.ErrorResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)

				Expect(resp.Message).Should(Equal("given ip address has not made any transaction yet"))
			})
		})

		When("success", func() {
			BeforeEach(func() {
				mockService.EXPECT().GetIPAddressStolenCardHistory(gomock.Any(), "192.0.0.1").Return(&entity.GetIPAddressStolenCardHistoryResponse{
					HasTransaction:     true,
					LinkedToStolenCard: true,
				}, nil)
			})
			It("return http error 200 ok", func() {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodGet, "/?ip_address=192.0.0.1", nil)
				ctx, _ = gin.CreateTestContext(rec)
				ctx.Request = req
				ctx.Params = []gin.Param{
					{
						Key:   "ip_address",
						Value: "192.0.0.1",
					},
				}
				transactionHandler.GetIPAddressStolenCardHistory(ctx)
				Expect(rec.Code).Should(Equal(http.StatusOK))

				var resp entity.GetIPAddressStolenCardHistoryResponse
				_ = json.NewDecoder(rec.Body).Decode(&resp)

				Expect(resp.LinkedToStolenCard).Should(Equal(true))
			})
		})
	})
})
