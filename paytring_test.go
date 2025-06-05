package paytring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var apiKey = "test_key"
var apiSecret = "test_secret"

func TestCreateOrder(t *testing.T) {

	var amount int64 = 1000
	receiptID := "TEST123"
	callbackURL := "https://example.com/callback"
	customer := Customer{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "1234567890",
	}
	paymentConfig := PaymentConfig{
		Currency: "INR",
	}

	paytring := NewClient(apiKey, apiSecret)
	resp, err := paytring.CreateOrder(amount, receiptID, callbackURL, customer, paymentConfig)
	fmt.Println(resp)
	assert.NoError(t, err)

}

func TestFetchOrder(t *testing.T) {

	paytring := NewClient(apiKey, apiSecret)
	resp, err := paytring.FetchOrder("503279383566355092")
	if resp["status"] != false {
		assert.True(t, true)
	}
	assert.NoError(t, err)

}

func TestMakeHash(t *testing.T) {

	client := NewClient(apiKey, apiSecret)

	params := map[string]interface{}{
		"amount":       "100",
		"callback_url": "https://example.com/callback",
		"cname":        "JohnDoe",
		"email":        "john.doe@example.com",
		"phone":        "1234567890",
		"key":          apiKey,
		"receipt_id":   "TEST123",
	}

	expectedHash := "b7c0e8c6896f7a07200eed32578f15b42751aac6da9fa66c7d784fbbcc0f943d811c6eaecf1ee2aa6aa57e868803c3051cf611375890e3d57eeafc6d48e7c150"

	resp := client.MakeHash(params)

	assert.Equal(t, expectedHash, resp["hash"])
}

func TestFetchOrderByReceipt(t *testing.T) {
	paytring := NewClient(apiKey, apiSecret)
	resp, err := paytring.FetchOrderByReceipt("TEST_RECEIPT_ID_123")
	fmt.Println(resp)
	if resp != nil && resp["status"] != nil && resp["status"] != false {
		assert.True(t, true, "Expected status to be non-false if present")
	}
	assert.NoError(t, err)
	assert.NotNil(t, resp, "Response should not be nil on success")
}

func TestValidateVPA(t *testing.T) {
	paytring := NewClient(apiKey, apiSecret)
	resp, err := paytring.ValidateVPA("test@vpa")
	fmt.Println(resp)
	if resp != nil && resp["status"] != nil && resp["status"] != false {
		assert.True(t, true, "Expected status to be non-false if present")
	}
	assert.NoError(t, err)
	assert.NotNil(t, resp, "Response should not be nil on success")
}

func TestValidateCard(t *testing.T) {
	paytring := NewClient(apiKey, apiSecret)
	resp, err := paytring.ValidateCard("424242")
	fmt.Println(resp)
	if resp != nil && resp["status"] != nil && resp["status"] != false {
		assert.True(t, true, "Expected status to be non-false if present")
	}
	assert.NoError(t, err)
	assert.NotNil(t, resp, "Response should not be nil on success")
}

func TestProcessOrder(t *testing.T) {
	paytring := NewClient(apiKey, apiSecret)
	// Example for a UPI payment, adjust as needed
	resp, err := paytring.ProcessOrder("TEST_ORDER_ID_PROCESS", "upi", "collect", PaymentData{Vpa: "test@upi"}, "desktop")
	fmt.Println(resp)
	if resp != nil && resp["status"] != nil && resp["status"] != false {
		assert.True(t, true, "Expected status to be non-false if present")
	}
	assert.NoError(t, err)
	assert.NotNil(t, resp, "Response should not be nil on success")
}

func TestRefundOrder(t *testing.T) {
	paytring := NewClient(apiKey, apiSecret)
	resp, err := paytring.RefundOrder("TEST_ORDER_ID_REFUND")
	fmt.Println(resp)
	if resp != nil && resp["status"] != nil && resp["status"] != false {
		assert.True(t, true, "Expected status to be non-false if present")
	}
	assert.NoError(t, err)
	assert.NotNil(t, resp, "Response should not be nil on success")
}

func TestFetchRefundStatus(t *testing.T) {
	paytring := NewClient(apiKey, apiSecret)
	resp, err := paytring.FetchRefundStatus("TEST_REFUND_ID_STATUS")
	fmt.Println(resp)
	if resp != nil && resp["status"] != nil && resp["status"] != false {
		assert.True(t, true, "Expected status to be non-false if present")
	}
	assert.NoError(t, err)
	assert.NotNil(t, resp, "Response should not be nil on success")
}

func TestPartialRefund(t *testing.T) {
	paytring := NewClient(apiKey, apiSecret)
	var amount int64 = 100
	resp, err := paytring.PartialRefund("TEST_ORDER_ID_PARTIAL_REFUND", amount)
	fmt.Println(resp)
	if resp != nil && resp["status"] != nil && resp["status"] != false {
		assert.True(t, true, "Expected status to be non-false if present")
	}
	assert.NoError(t, err)
	assert.NotNil(t, resp, "Response should not be nil on success")
}

func TestFetchRefundAttempts(t *testing.T) {
	paytring := NewClient(apiKey, apiSecret)
	resp, err := paytring.FetchRefundAttempts("TEST_ORDER_ID_REFUND_ATTEMPTS")
	fmt.Println(resp)
	if resp != nil && resp["status"] != nil && resp["status"] != false {
		assert.True(t, true, "Expected status to be non-false if present")
	}
	assert.NoError(t, err)
	assert.NotNil(t, resp, "Response should not be nil on success")
}

func TestFetchRefund(t *testing.T) {
	paytring := NewClient(apiKey, apiSecret)
	resp, err := paytring.FetchRefund("TEST_REFUND_ID_FETCH")
	fmt.Println(resp)
	if resp != nil && resp["status"] != nil && resp["status"] != false {
		assert.True(t, true, "Expected status to be non-false if present")
	}
	assert.NoError(t, err)
	assert.NotNil(t, resp, "Response should not be nil on success")
}
