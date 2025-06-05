package paytring

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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
