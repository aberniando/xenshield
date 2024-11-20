package integration_test

import (
	"fmt"
	. "github.com/Eun/go-hit"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"testing"
)

const (
	basePath = "http://app:8080"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func GenerateRandomIP() string {
	octets := make([]string, 4)

	for i := 0; i < 4; i++ {
		octets[i] = fmt.Sprintf("%d", rand.Intn(256))
	}

	return strings.Join(octets, ".")
}

func TestHTTPFailedNoAPIKeyProvided(t *testing.T) {
	ipAddress := GenerateRandomIP()

	Test(t,
		Description("FailedNoAPIKeyProvided"),
		Get(fmt.Sprintf(`%s/transactions/%s`, basePath, ipAddress)),
		Send().Headers("Content-Type").Add("application/json"),
		Expect().Status().Equal(http.StatusForbidden),
		Expect().Body().JSON().JQ(".message").Equal("Missing API key in request"),
	)
}

func TestHTTPFailedInvalidAPIKey(t *testing.T) {
	ipAddress := GenerateRandomIP()

	Test(t,
		Description("FailedNoAPIKeyProvided"),
		Get(fmt.Sprintf(`%s/transactions/%s`, basePath, ipAddress)),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("api-key").Add("xxxxx"),
		Expect().Status().Equal(http.StatusForbidden),
		Expect().Body().JSON().JQ(".message").Equal("Invalid API key"),
	)
}

func TestHTTPCreateTransactionWithoutMandatoryFields(t *testing.T) {
	failureReason := "STOLEN_CARD"

	body := fmt.Sprintf(`{
		"failure_reason": "%s"
	}`, failureReason)

	Test(t,
		Description("CreateTransactionWithoutMandatoryFields"),
		Post(basePath+"/transactions"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("api-key").Add("eGVuc2hpZWxk"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".message").Len().Equal(3),
	)
}

func TestHTTPCreateTransactionInvalidBody1(t *testing.T) {
	ipAddress := "xxxxx"
	maskedCardNumber := fmt.Sprintf(`****-****-*********-%d`, rand.Intn(9000)+1000)
	status := "SUCCESS"
	failureReason := "INSUFFICIENT_BALANCE"

	body := fmt.Sprintf(`{
		"ip_address": "%s",
		"masked_card_number": "%s",
		"status": "%s",
		"failure_reason": "%s"
	}`, ipAddress, maskedCardNumber, status, failureReason)

	Test(t,
		Description("CreateTransactionInvalidBody1"),
		Post(basePath+"/transactions"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("api-key").Add("eGVuc2hpZWxk"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".message").Len().Equal(3),
		Expect().Body().JSON().JQ(".message[0].message").Equal("invalid value format"),
		Expect().Body().JSON().JQ(".message[1].message").Equal("invalid length"),
		Expect().Body().JSON().JQ(".message[2].message").Equal("failure reason must be null on success transaction"),
	)
}

func TestHTTPCreateTransactionInvalidBody2(t *testing.T) {
	ipAddress := "xxxxx"
	maskedCardNumber := fmt.Sprintf(`****-****-*********-%d`, rand.Intn(9000)+1000)
	status := "FAILED"
	failureReason := ""

	body := fmt.Sprintf(`{
		"ip_address": "%s",
		"masked_card_number": "%s",
		"status": "%s",
		"failure_reason": "%s"
	}`, ipAddress, maskedCardNumber, status, failureReason)

	Test(t,
		Description("CreateTransactionInvalidBody1"),
		Post(basePath+"/transactions"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("api-key").Add("eGVuc2hpZWxk"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusBadRequest),
		Expect().Body().JSON().JQ(".message").Len().Equal(3),
		Expect().Body().JSON().JQ(".message[0].message").Equal("invalid value format"),
		Expect().Body().JSON().JQ(".message[1].message").Equal("invalid length"),
		Expect().Body().JSON().JQ(".message[2].message").Equal("failure reason must not be null on failed transaction"),
	)
}

func TestHTTPCreateTransactionNotLinkedToStolenCard(t *testing.T) {
	ipAddress := GenerateRandomIP()
	maskedCardNumber := fmt.Sprintf(`****-****-****-%d`, rand.Intn(9000)+1000)
	status := "FAILED"
	failureReason := "INSUFFICIENT_BALANCE"

	body := fmt.Sprintf(`{
		"ip_address": "%s",
		"masked_card_number": "%s",
		"status": "%s",
		"failure_reason": "%s"
	}`, ipAddress, maskedCardNumber, status, failureReason)

	Test(t,
		Description("CreateTransactionNotLinkedToStolenCard - POST"),
		Post(basePath+"/transactions"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("api-key").Add("eGVuc2hpZWxk"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".id").Len().Equal(36),
		Expect().Body().JSON().JQ(".ip_address").Equal(ipAddress),
		Expect().Body().JSON().JQ(".masked_card_number").Equal(maskedCardNumber),
		Expect().Body().JSON().JQ(".status").Equal(status),
		Expect().Body().JSON().JQ(".reason").Equal(failureReason),
	)

	Test(t,
		Description("CreateTransactionNotLinkedToStolenCard - GET"),
		Get(fmt.Sprintf(`%s/transactions/%s`, basePath, ipAddress)),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("api-key").Add("eGVuc2hpZWxk"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".linked_to_stolen_card").Equal(false),
	)
}

func TestHTTPCreateTransactionLinkedToStolenCard(t *testing.T) {
	ipAddress := GenerateRandomIP()
	maskedCardNumber := fmt.Sprintf(`****-****-****-%d`, rand.Intn(9000)+1000)
	status := "FAILED"
	failureReason := "STOLEN_CARD"

	body := fmt.Sprintf(`{
		"ip_address": "%s",
		"masked_card_number": "%s",
		"status": "%s",
		"failure_reason": "%s"
	}`, ipAddress, maskedCardNumber, status, failureReason)

	Test(t,
		Description("CreateTransactionLinkedToStolenCard - POST"),
		Post(basePath+"/transactions"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("api-key").Add("eGVuc2hpZWxk"),
		Send().Body().String(body),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".id").Len().Equal(36),
		Expect().Body().JSON().JQ(".ip_address").Equal(ipAddress),
		Expect().Body().JSON().JQ(".masked_card_number").Equal(maskedCardNumber),
		Expect().Body().JSON().JQ(".status").Equal(status),
		Expect().Body().JSON().JQ(".reason").Equal(failureReason),
	)

	Test(t,
		Description("CreateTransactionLinkedToStolenCard - GET"),
		Get(fmt.Sprintf(`%s/transactions/%s`, basePath, ipAddress)),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Headers("api-key").Add("eGVuc2hpZWxk"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".linked_to_stolen_card").Equal(true),
	)
}
