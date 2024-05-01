package paytring

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	request "github.com/dagar-in/http-client"
)

func NewClient(apiKey string, apiSecret string) *Api {
	return &Api{
		ApiKey:    apiKey,
		ApiSecret: apiSecret,
		ApiUrl:    "https://api.paytring.com/api/v1/",
	}
}

type Api struct {
	ApiKey    string
	ApiSecret string
	ApiUrl    string
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

func (c *Api) CreateOrder(
	amount int64,
	receiptId string,
	callbackUrl string,
	customer Customer,
	opts ...interface{},
) (map[string]interface{}, error) {

	var paymentConfig PaymentConfig
	var billingAddress BillingAddress
	var shippingAddress ShippingAddress
	var notes Notes
	var tpv []Tpv
	var splitSettlement SplitSettlement

	for _, opt := range opts {
		switch v := opt.(type) {
		case PaymentConfig:
			paymentConfig = v
		case BillingAddress:
			billingAddress = v
		case ShippingAddress:
			shippingAddress = v
		case Notes:
			notes = v
		case []Tpv:
			tpv = v
		case SplitSettlement:
			splitSettlement = v
		}
	}

	client := request.New()

	requestBody := map[string]interface{}{
		"key":          c.ApiKey,
		"receipt_id":   receiptId,
		"amount":       strconv.FormatInt(amount, 10),
		"cname":        customer.Name,
		"phone":        customer.Phone,
		"email":        customer.Email,
		"callback_url": callbackUrl,
	}

	if paymentConfig.Currency != "" {
		requestBody["currency"] = paymentConfig.Currency
	} else {
		requestBody["currency"] = "INR"
	}

	if paymentConfig.Pg != "" {
		requestBody["pg_pool_id"] = paymentConfig.Pg
	}

	if !paymentConfig.AutoCapture {
		requestBody["auto_capture"] = "false"
	} else {
		requestBody["auto_capture"] = "true"
	}

	billingAddressMap := make(map[string]interface{})
	addToMapIfNotBlank(billingAddressMap, "firstname", billingAddress.Firstname)
	addToMapIfNotBlank(billingAddressMap, "lastname", billingAddress.Lastname)
	addToMapIfNotBlank(billingAddressMap, "phone", billingAddress.Phone)
	addToMapIfNotBlank(billingAddressMap, "line1", billingAddress.Line1)
	addToMapIfNotBlank(billingAddressMap, "line2", billingAddress.Line2)
	addToMapIfNotBlank(billingAddressMap, "city", billingAddress.City)
	addToMapIfNotBlank(billingAddressMap, "state", billingAddress.State)
	addToMapIfNotBlank(billingAddressMap, "country", billingAddress.Country)
	addToMapIfNotBlank(billingAddressMap, "zipcode", billingAddress.Zipcode)

	shippingAddressMap := make(map[string]interface{})
	addToMapIfNotBlank(shippingAddressMap, "firstname", shippingAddress.Firstname)
	addToMapIfNotBlank(shippingAddressMap, "lastname", shippingAddress.Lastname)
	addToMapIfNotBlank(shippingAddressMap, "phone", shippingAddress.Phone)
	addToMapIfNotBlank(shippingAddressMap, "line1", shippingAddress.Line1)
	addToMapIfNotBlank(shippingAddressMap, "line2", shippingAddress.Line2)
	addToMapIfNotBlank(shippingAddressMap, "city", shippingAddress.City)
	addToMapIfNotBlank(shippingAddressMap, "state", shippingAddress.State)
	addToMapIfNotBlank(shippingAddressMap, "country", shippingAddress.Country)
	addToMapIfNotBlank(shippingAddressMap, "zipcode", shippingAddress.Zipcode)

	notesMap := make(map[string]interface{})
	addToMapIfNotBlank(notesMap, "udf1", notes.Udf1)
	addToMapIfNotBlank(notesMap, "udf2", notes.Udf2)
	addToMapIfNotBlank(notesMap, "udf3", notes.Udf3)
	addToMapIfNotBlank(notesMap, "udf4", notes.Udf4)
	addToMapIfNotBlank(notesMap, "udf5", notes.Udf5)

	var tpvMapMap []map[string]interface{}
	for _, tpvAccount := range tpv {
		tpvMap := make(map[string]interface{})
		addToMapIfNotBlank(tpvMap, "account_number", tpvAccount.AccountNumber)
		addToMapIfNotBlank(tpvMap, "name", tpvAccount.Name)
		addToMapIfNotBlank(tpvMap, "ifsc", tpvAccount.Ifsc)
		tpvMapMap = append(tpvMapMap, tpvMap)
	}

	if splitSettlement.SplitType != "" {
		requestBody["split_type"] = splitSettlement.SplitType
	}

	var splitSettlementMap []map[string]interface{}
	for _, splitRule := range splitSettlement.SplitRule {
		var splitSettlementRuleMap = make(map[string]interface{})
		addToMapIfNotBlank(splitSettlementRuleMap, "vendor_id", splitRule.VendorId)
		addToMapIfNotBlank(splitSettlementRuleMap, "amount", splitRule.Amount)
		splitSettlementMap = append(splitSettlementMap, splitSettlementRuleMap)
	}

	if len(splitSettlementMap) > 0 {
		requestBody["split_settlement"] = splitSettlementMap
	}
	if len(billingAddressMap) > 0 {
		requestBody["billing_address"] = billingAddressMap
	}
	if len(shippingAddressMap) > 0 {
		requestBody["shipping_address"] = shippingAddressMap
	}
	if len(notesMap) > 0 {
		requestBody["notes"] = notesMap
	}
	if len(tpvMapMap) > 0 {
		requestBody["tpv"] = tpvMapMap
	}

	body, err := json.Marshal(c.MakeHash(requestBody))
	if err != nil {
		print(err)
		return nil, err
	}

	resp, err := client.WithHeaders(map[string]string{
		"Content-Type": "application/json",
	}).WithBody(body).Post(c.ApiUrl + "order/create")

	if err != nil {
		print(err)
		return nil, err
	}

	bodyMap, err := resp.BodyMap()
	if err != nil {
		print(err)
		return nil, err
	}

	response, err := c.HandleResponse(bodyMap)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Api) FetchOrder(orderId string) (map[string]interface{}, error) {

	client := request.New()

	requestBody := map[string]interface{}{
		"key": c.ApiKey,
		"id":  orderId,
	}

	body, err := json.Marshal(c.MakeHash(requestBody))
	if err != nil {
		print(err)
	}

	resp, err := client.WithHeaders(map[string]string{
		"Content-Type": "application/json",
	}).WithBody(body).Post(c.ApiUrl + "order/fetch")

	if err != nil {
		print(err)
		return nil, err
	}

	bodyMap, err := resp.BodyMap()
	if err != nil {
		print(err)
		return nil, err
	}

	return bodyMap, nil
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
