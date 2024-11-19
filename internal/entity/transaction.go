package entity

import (
	"github.com/aberniando/xenshield/internal/enum"
	"time"
)

type InsertTransactionRequest struct {
	IPAddress        string                 `json:"ip_address" binding:"required"`
	MaskedCardNumber string                 `json:"masked_card_number" binding:"required"`
	Status           enum.TransactionStatus `json:"status" binding:"required"`
	FailureReason    enum.FailureReason     `json:"failure_reason"`
}

type InsertTransactionResponse struct {
	ID               string                 `json:"id"`
	IPAddress        string                 `json:"ip_address"`
	MaskedCardNumber string                 `json:"masked_card_number"`
	Status           enum.TransactionStatus `json:"status"`
	Reason           enum.FailureReason     `json:"reason"`
	Created          string                 `json:"created"`
	Updated          string                 `json:"updated"`
}

type GetIPAddressStolenCardHistoryResponse struct {
	HasTransaction     bool `json:"-"`
	LinkedToStolenCard bool `json:"linked_to_stolen_card"`
}

type Transaction struct {
	ID               string                 `db:"id"`
	IPAddress        string                 `db:"ip_address"`
	MaskedCardNumber string                 `db:"masked_card_number"`
	Status           enum.TransactionStatus `db:"status"`
	Reason           *enum.FailureReason    `db:"reason"`
	Created          time.Time              `db:"created"`
	Updated          time.Time              `db:"updated"`
}

func (itr *InsertTransactionRequest) ToTransaction() *Transaction {
	return &Transaction{
		IPAddress:        itr.IPAddress,
		MaskedCardNumber: itr.MaskedCardNumber,
		Status:           itr.Status,
		Reason: func(fr enum.FailureReason) *enum.FailureReason {
			if itr.FailureReason == "" {
				return nil
			}
			return &fr
		}(itr.FailureReason),
	}
}

func (t *Transaction) ToInsertTransactionResponse() *InsertTransactionResponse {
	return &InsertTransactionResponse{
		ID:               t.ID,
		IPAddress:        t.IPAddress,
		MaskedCardNumber: t.MaskedCardNumber,
		Status:           t.Status,
		Reason: func(r *enum.FailureReason) enum.FailureReason {
			if r == nil {
				return ""
			}
			return *r
		}(t.Reason),
		Created: t.Created.Format(time.RFC3339),
		Updated: t.Updated.Format(time.RFC3339),
	}
}
