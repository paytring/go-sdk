package paytring

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/go-resty/resty/v2"
)

func NewClient(apiKey string, apiSecret string) *Api {
	return &Api{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		ApiUrl:    "https://api.paytring.com/api/",
		UserAgent: "paytring-go-sdk/0",
		http:      resty.New(),
	}
}

type Api struct {
	ApiKey    string
	ApiSecret string
	ApiUrl    string
	UserAgent string
	http      *resty.Client
}

type Customer struct {
	Name  string
	Email string
	Phone string
}

type PaymentConfig struct {
	Currency    string
	Pg          string
	PgPoolId    string
	AutoCapture bool
}

type SplitRule struct {
	VendorId string
	Amount   int64
}

type SplitSettlement struct {
	SplitType string
	SplitRule []SplitRule
}

type BillingAddress struct {
	Firstname string
	Lastname  string
	Phone     string
	Line1     string
	Line2     string
	City      string
	State     string
	Country   string
	Zipcode   string
}

type ShippingAddress struct {
	Firstname string
	Lastname  string
	Phone     string
	Line1     string
	Line2     string
	City      string
	State     string
	Country   string
	Zipcode   string
}

type Notes struct {
	Udf1 string
	Udf2 string
	Udf3 string
	Udf4 string
	Udf5 string
}

type Tpv struct {
	AccountNumber string
	Name          string
	Ifsc          string
}

type PaymentData struct {
	Vpa         string
	CardNumber  string
	ExpiryMonth string
	ExpiryYear  string
	Cvv         string
	HolderName  string
}

func (c *Api) MakeAuthHeader() map[string]string {

	authString := c.ApiKey + ":" + c.ApiSecret
	authString = base64.StdEncoding.EncodeToString([]byte(authString))

	return map[string]string{
		"User-Agent":    c.UserAgent,
		"Content-Type":  "application/json",
		"Authorization": "Basic " + authString,
	}
}

func (c *Api) MakeHash(params map[string]interface{}) map[string]interface{} {

	valueString := ""

	keys := make([]string, 0, len(params))

	for key := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		if _, ok := params[key].(string); !ok {
			continue
		}
		valueString += string(params[key].(string)) + "|"
	}

	valueString += c.ApiSecret

	hash := sha512.Sum512([]byte(valueString))
	params["hash"] = fmt.Sprintf("%x", hash)
	return params
}

func addToMapIfNotBlank(m map[string]interface{}, key string, value interface{}) {
	if value != nil {
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Slice {
			if v.Len() > 0 {
				m[key] = value
			}
		} else if value != 0 && value != false && value != "" {
			m[key] = value
		}
	}
}

func (c *Api) HandleResponse(response map[string]interface{}) (map[string]interface{}, error) {
	if response["status"] == true {
		return response, nil
	}
	message := extractErrorMessage(response["error"])
	return nil, fmt.Errorf(message)
}

func extractErrorMessage(errors interface{}) string {
	if errors == nil {
		return "Something went wrong, invalid response received"
	}
	errorJSON, err := json.Marshal(errors)
	if err != nil {
		return fmt.Sprintf("Unexpected error response recieved: %v", err)
	}
	var errorMap map[string]interface{}
	if err := json.Unmarshal(errorJSON, &errorMap); err != nil {
		return fmt.Sprintf("Invalid error response recieved: %v", err)
	}

	message := errorMap["message"]

	if message == nil {
		return "Something went wrong, invalid data"
	}

	var messageMap map[string]interface{}
	if err := json.Unmarshal([]byte(message.(string)), &messageMap); err != nil {
		if s, ok := message.(string); ok {
			return s
		}
		return fmt.Sprintf("Error unmarshalling message: %v", err)
	}

	for _, v := range messageMap {
		if arr, ok := v.([]interface{}); ok && len(arr) > 0 {
			if s, ok := arr[0].(string); ok {
				return s
			}
		}
	}

	return "Something went wrong, please try again later."
}

func validatePaymentDataForCard(paymentData PaymentData) error {
	if paymentData.CardNumber == "" {
		return fmt.Errorf("card number is required")
	}
	if paymentData.ExpiryMonth == "" {
		return fmt.Errorf("expiry month is required")
	}
	if paymentData.ExpiryYear == "" {
		return fmt.Errorf("expiry year is required")
	}
	if paymentData.Cvv == "" {
		return fmt.Errorf("CVV is required")
	}
	return nil
}
