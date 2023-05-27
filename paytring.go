package paytring

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
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

type Options map[string]interface{}

func (c *Api) CreateOrder(
	amount int64,
	receiptId string,
	callbackUrl string,
	customer Customer,
	options Options,
) (map[string]interface{}, error) {

	client := request.New()

	requestBody := MergeMaps(Options{
		"amount":       strconv.FormatInt(amount, 10),
		"callback_url": "https://httpbin.org/post",
		"cname":        customer.Name,
		"email":        customer.Email,
		"phone":        customer.Phone,
		"key":          c.ApiKey,
		"receipt_id":   receiptId,
	},
		options,
	)

	body, err := json.Marshal(c.MakeHash(requestBody))
	if err != nil {
		print(err)
	}

	resp, err := client.WithHeaders(map[string]string{
		"Content-Type": "application/json",
	}).WithBody(body).Post(c.ApiUrl + "order/create")

	if err != nil {
		print(err)
	}

	bodyMap, err := resp.BodyMap()
	if err != nil {
		print(err)
	}

	return bodyMap, nil
}

func (c *Api) FetchOrder(orderId string) (map[string]interface{}, error) {

	client := request.New()

	requestBody := Options{
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
	}

	bodyMap, err := resp.BodyMap()
	if err != nil {
		print(err)
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
		valueString += string(params[key].(string)) + "|"
	}

	valueString += c.ApiSecret
	hash := sha512.Sum512([]byte(valueString))
	params["hash"] = fmt.Sprintf("%x", hash)
	return params
}
