package paytring

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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
