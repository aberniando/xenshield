package enum

import "encoding/json"

type TransactionStatus string

const TransactionStatusSuccess TransactionStatus = "SUCCESS"
const TransactionStatusFailed TransactionStatus = "FAILED"

func (ts *TransactionStatus) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch TransactionStatus(value) {
	case TransactionStatusSuccess, TransactionStatusFailed:
		*ts = TransactionStatus(value)
	default:
		*ts = ""
	}

	return nil
}

type FailureReason string

const FailureReasonStolenCard FailureReason = "STOLEN_CARD"
const FailureReasonCardDeclined FailureReason = "CARD_DECLINED"
const FailureReasonInsufficientBalance FailureReason = "INSUFFICIENT_BALANCE"

func (fr *FailureReason) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch FailureReason(value) {
	case FailureReasonStolenCard, FailureReasonCardDeclined, FailureReasonInsufficientBalance:
		*fr = FailureReason(value)
	default:
		*fr = ""
	}

	return nil
}
