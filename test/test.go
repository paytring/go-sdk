package main

import (
	"fmt"

	paytring "github.com/paytring/go-sdk"
)

var apiKey = "test_key"
var apiSecret = "test_secret"

func main() {
	var amount int64 = 1000
	receiptID := "TEST0001"
	callbackURL := "https://example.com/callback"
	customer := paytring.Customer{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Phone: "1234567890",
	}
	paymentConfig := paytring.PaymentConfig{
		Currency: "INR",
		Pg:       "razorpay",
	}

	billingAddress := paytring.BillingAddress{
		Firstname: "John",
		Lastname:  "Doe",
		Phone:     "1234567890",
		Line1:     "Address Line 1",
		Line2:     "Address Line 2",
		City:      "City",
		State:     "State",
		Country:   "Country",
		Zipcode:   "123456",
	}

	shippingAddress := paytring.ShippingAddress{
		Firstname: "John",
		Lastname:  "Doe",
		Phone:     "1234567890",
		Line1:     "Address Line 1",
		Line2:     "Address Line 2",
		City:      "City",
		State:     "State",
		Country:   "Country",
		Zipcode:   "123456",
	}

	tpv := []paytring.Tpv{
		// {
		// 	AccountNumber: "1234567890",
		// 	Name:          "John Doe",
		// 	Ifsc:          "IFSC1234",
		// },
		// {
		// 	AccountNumber: "9898989898",
		// 	Name:          "John Vick",
		// 	Ifsc:          "IFSC12345",
		// },
	}

	notes := paytring.Notes{
		Udf1: "udf1",
		Udf2: "udf2",
		Udf3: "udf3",
		Udf4: "udf4",
		Udf5: "udf5",
	}

	splitSettlement := paytring.SplitSettlement{
		SplitType: "percent",
		SplitRule: []paytring.SplitRule{
			// {
			// 	VendorId: "sub_merchant_id",
			// 	Amount:   50,
			// },
		},
	}

	orderID := "623842825435349418"
	paymentMethod := "card"
	paymentCode := "card"
	paymentData := paytring.PaymentData{
		Vpa:         "7027445661@paytm",
		CardNumber:  "4242424242424242",
		ExpiryMonth: "12",
		ExpiryYear:  "24",
		Cvv:         "123",
		HolderName:  "John Doe",
	}
	device := "android"

	paytringAPI := paytring.NewClient(apiKey, apiSecret)

	validateVPAResponse, errValidateVPA := paytringAPI.ValidateVPA("7027445661@paytm")
	if errValidateVPA != nil {
		fmt.Println(errValidateVPA)
		return
	}
	fmt.Println(validateVPAResponse)

	validateCardBinResponse, errValidateCardBin := paytringAPI.ValidateCard("424242")
	if errValidateCardBin != nil {
		fmt.Println(errValidateCardBin)
		return
	}
	fmt.Println(validateCardBinResponse)

	orderCreateResponse, errOrderCreate := paytringAPI.CreateOrder(amount, receiptID, callbackURL, customer, paymentConfig, billingAddress, shippingAddress, notes, tpv, splitSettlement)
	if errOrderCreate != nil {
		fmt.Println(errOrderCreate)
		return
	}
	fmt.Println(orderCreateResponse)

	orderID = orderCreateResponse["order_id"].(string)

	processOrderResponse, errProcessOrder := paytringAPI.ProcessOrder(orderID, paymentMethod, paymentCode, paymentData, device)
	if errProcessOrder != nil {
		fmt.Println(errProcessOrder)
		return
	}
	fmt.Println(processOrderResponse)

	fetchOrderResponse, errFetchOrder := paytringAPI.FetchOrder(orderCreateResponse["order_id"].(string))
	if errFetchOrder != nil {
		fmt.Println(errFetchOrder)
		return
	}
	fmt.Println(fetchOrderResponse)

}
