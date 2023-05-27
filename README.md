# Paytring Go Package

The paytring package is a Go client library for interacting with the Paytring API. It provides convenient methods for creating and fetching orders using the Paytring service.

## Installation
To use the paytring package in your Go project, you can install it using the go get command:
```go
go get "github.com/paytring/go-sdk"
```

## Usage
Import the `paytring` package in your Go code:

```go
import "github.com/paytring/go-sdk"
```

Create a new Paytring client by calling the `NewClient` function:
```go
client := paytring.NewClient(apiKey, apiSecret)
```

### Create an Order

To create an order, use the `CreateOrder` method:
```go
amount := int64(100)
receiptID := "your-receipt-id"
callbackURL := "https://your-callback-url.com"
customer := paytring.Customer{
	Name:  "John Doe",
	Email: "john.doe@example.com",
	Phone: "1234567890",
}
options := paytring.Options{
	"currency": "USD",
}

response, err := client.CreateOrder(amount, receiptID, callbackURL, customer, options)
if err != nil {
	// Handle error
}

```

### Fetch an Order
To fetch an existing order, use the `FetchOrder` method:

```go
orderID := "your-order-id"

response, err := client.FetchOrder(orderID)
if err != nil {
	// Handle error
}
```

## API Documentation

### type Api
```go
type Api struct {
	ApiKey    string
	ApiSecret string
	ApiUrl    string
}
```
The `Api` struct represents the Paytring API client. It contains the API key, API secret, and API URL.
