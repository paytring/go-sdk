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
	resp, err := paytring.FetchOrder("TEST123")
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
