package paytring

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

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
		requestBody["pg"] = paymentConfig.Pg
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

	requestBody["hash"] = "none"

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for CreateOrder: %w", err)
	}

	resp, err := c.http.R().
		SetHeaders(c.MakeAuthHeader()).
		SetBody(body).
		Post(c.ApiUrl + "v2/order/create")

	if err != nil {
		print(err)
		return nil, err
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(resp.Body(), &bodyMap)
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

	requestBody := map[string]interface{}{
		"key":  c.ApiKey,
		"id":   orderId,
		"hash": "none",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for FetchOrder: %w", err)
	}

	resp, err := c.http.R().
		SetHeaders(c.MakeAuthHeader()).
		SetBody(body).
		Post(c.ApiUrl + "v2/order/fetch")

	if err != nil {
		print(err)
		return nil, err
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(resp.Body(), &bodyMap)
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

func (c *Api) FetchOrderByReceipt(receiptId string) (map[string]interface{}, error) {

	requestBody := map[string]interface{}{
		"key": c.ApiKey,
		"id":  receiptId,
	}

	body, err := json.Marshal(c.MakeHash(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for FetchOrderByReceipt: %w", err)
	}

	resp, err := c.http.R().
		SetHeaders(c.MakeAuthHeader()).
		SetBody(body).
		Post(c.ApiUrl + "v2/order/fetch/receipt")

	if err != nil {
		print(err)
		return nil, err
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(resp.Body(), &bodyMap)
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

func (c *Api) MakeAuthHeader() map[string]string {

	authString := c.ApiKey + ":" + c.ApiSecret
	authString = base64.StdEncoding.EncodeToString([]byte(authString))

	return map[string]string{
		"User-Agent":    c.UserAgent,
		"Content-Type":  "application/json",
		"Authorization": "Basic " + authString,
	}
}

func (c *Api) ValidateVPA(vpa string) (map[string]interface{}, error) {

	requestBody := map[string]interface{}{
		"key": c.ApiKey,
		"vpa": vpa,
	}

	body, err := json.Marshal(c.MakeHash(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for ValidateVPA: %w", err)
	}

	resp, err := c.http.R().SetHeaders(c.MakeAuthHeader()).SetBody(body).Post(c.ApiUrl + "v1/info/vpa")

	if err != nil {
		print(err)
		return nil, err
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(resp.Body(), &bodyMap)
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

func (c *Api) ValidateCard(bin string) (map[string]interface{}, error) {

	requestBody := map[string]interface{}{
		"key": c.ApiKey,
		"bin": bin,
	}

	body, err := json.Marshal(c.MakeHash(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for ValidateCard: %w", err)
	}

	resp, err := c.http.R().SetHeaders(c.MakeAuthHeader()).SetBody(body).Post(c.ApiUrl + "v1/info/bin")

	if err != nil {
		print(err)
		return nil, err
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(resp.Body(), &bodyMap)
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

func (c *Api) ProcessOrder(orderId string, paymentMethod string, paymentCode string, paymentData PaymentData, device string) (map[string]interface{}, error) {

	requestBody := map[string]interface{}{
		"key":      c.ApiKey,
		"order_id": orderId,
		"method":   paymentMethod,
		"code":     paymentCode,
	}

	if device != "" {
		requestBody["device"] = device
	}

	if paymentMethod == "upi" {
		if paymentCode == "collect" {
			if paymentData.Vpa == "" {
				return nil, fmt.Errorf("VPA is required for UPI collect")
			}
			requestBody["vpa"] = paymentData.Vpa
		}
	}

	if paymentMethod == "card" {
		errInValidatingPaymentDataForCard := validatePaymentDataForCard(paymentData)
		if errInValidatingPaymentDataForCard != nil {
			return nil, errInValidatingPaymentDataForCard
		}
		requestBody["card"] = map[string]interface{}{
			"holder_name":  paymentData.HolderName,
			"number":       paymentData.CardNumber,
			"expiry_month": paymentData.ExpiryMonth,
			"expiry_year":  paymentData.ExpiryYear,
			"cvv":          paymentData.Cvv,
		}
	}

	body, err := json.Marshal(c.MakeHash(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for ProcessOrder: %w", err)
	}

	resp, err := c.http.R().SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   c.UserAgent, // Added User-Agent consistent with MakeAuthHeader
	}).SetBody(body).Post(c.ApiUrl + "v1/order/process")

	if err != nil {
		print(err)
		return nil, err
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(resp.Body(), &bodyMap)
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

func (c *Api) RefundOrder(orderID string) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{
		"key":  c.ApiKey,
		"id":   orderID,
		"hash": "null",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for RefundOrder: %w", err)
	}

	resp, err := c.http.R().
		SetHeaders(c.MakeAuthHeader()).
		SetBody(body).
		Post(c.ApiUrl + "v2/order/refund")

	if err != nil {
		return nil, fmt.Errorf("RefundOrder request failed: %w", err)
	}

	var bodyMap map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &bodyMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body for RefundOrder: %w", err)
	}

	return c.HandleResponse(bodyMap)
}

func (c *Api) FetchRefundStatus(refundID string) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{
		"key":  c.ApiKey,
		"id":   refundID,
		"hash": "null",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for FetchRefundStatus: %w", err)
	}

	resp, err := c.http.R().
		SetHeaders(c.MakeAuthHeader()).
		SetBody(body).
		Post(c.ApiUrl + "v2/order/refund/fetch")

	if err != nil {
		return nil, fmt.Errorf("FetchRefundStatus request failed: %w", err)
	}

	var bodyMap map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &bodyMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body for FetchRefundStatus: %w", err)
	}

	return c.HandleResponse(bodyMap)
}

func (c *Api) PartialRefund(orderID string, amount int64) (map[string]interface{}, error) {
	requestPayload := map[string]interface{}{
		"key":    c.ApiKey,
		"id":     orderID,
		"amount": strconv.FormatInt(amount, 10), // Amount as string for hashing consistency
	}

	hashedPayload := c.MakeHash(requestPayload) // MakeHash will add the "hash" field

	body, err := json.Marshal(hashedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for PartialRefund: %w", err)
	}

	resp, err := c.http.R().
		SetHeaders(c.MakeAuthHeader()).
		SetBody(body).
		Post(c.ApiUrl + "v2/order/refund/partial")

	if err != nil {
		return nil, fmt.Errorf("PartialRefund request failed: %w", err)
	}

	var bodyMap map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &bodyMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body for PartialRefund: %w", err)
	}

	return c.HandleResponse(bodyMap)
}

func (c *Api) FetchRefundAttempts(orderID string) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{
		"key":      c.ApiKey,
		"order_id": orderID,
		"hash":     "null",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for FetchRefundAttempts: %w", err)
	}

	resp, err := c.http.R().
		SetHeaders(c.MakeAuthHeader()).
		SetBody(body).
		Post(c.ApiUrl + "v2/order/refund/attempts")

	if err != nil {
		return nil, fmt.Errorf("FetchRefundAttempts request failed: %w", err)
	}

	var bodyMap map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &bodyMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body for FetchRefundAttempts: %w", err)
	}

	return c.HandleResponse(bodyMap)
}

func (c *Api) FetchRefund(refundID string) (map[string]interface{}, error) {
	requestBody := map[string]interface{}{
		"key":  c.ApiKey,
		"id":   refundID,
		"hash": "null",
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for FetchRefund: %w", err)
	}

	resp, err := c.http.R().
		SetHeaders(c.MakeAuthHeader()).
		SetBody(body).
		Post(c.ApiUrl + "v2/order/refund/fetch")

	if err != nil {
		return nil, fmt.Errorf("FetchRefund request failed: %w", err)
	}

	var bodyMap map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &bodyMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body for FetchRefund: %w", err)
	}

	return c.HandleResponse(bodyMap)
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
