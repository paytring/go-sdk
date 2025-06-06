package paytring

import (
	"encoding/json"
	"fmt"
)

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
		"key":      c.ApiKey,
		"bin_code": bin,
	}

	body, err := json.Marshal(c.MakeHash(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for ValidateCard: %w", err)
	}

	resp, err := c.http.R().SetHeaders(c.MakeAuthHeader()).SetBody(body).Post(c.ApiUrl + "v1/health/bin")

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
func (c *Api) CurrencyConversion(from string, to string) (map[string]interface{}, error) {

	requestBody := map[string]interface{}{
		"key":  c.ApiKey,
		"from": from,
		"to":   to,
	}

	body, err := json.Marshal(c.MakeHash(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body for CurrencyConversion: %w", err)
	}

	resp, err := c.http.R().
		SetHeaders(c.MakeAuthHeader()).
		SetBody(body).
		Post(c.ApiUrl + "v1/currency/get")

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
