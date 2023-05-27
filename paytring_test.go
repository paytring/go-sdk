package paytring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var apiKey = "test_123"
var apiSecret = "secret_123"

func TestCreateOrder(t *testing.T) {

	var amount int64 = 1000
	receiptID := "TEST123"
	callbackURL := "https://example.com/callback"
	customer := Customer{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "1234567890",
	}
	options := Options{
		"currency": "USD",
	}

	paytring := NewClient(apiKey, apiSecret)
	resp, err := paytring.CreateOrder(amount, receiptID, callbackURL, customer, options)
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

	expectedHash := "10fb30b4375b834dea1153b9c501c175ad6cfd9944fa994341adefb0ddb46d9ac185abd58f7065734d7199410cd86b90eb3698aaec905bcefd7405c0d122261d"

	resp := client.MakeHash(params)

	assert.Equal(t, expectedHash, resp["hash"])
}
