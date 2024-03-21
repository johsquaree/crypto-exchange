package client

import (
	"bytes" // Importing the bytes package for byte manipulation.
	"encoding/json" // Importing the json package for JSON encoding and decoding.
	"fmt" // Importing the fmt package for formatted I/O operations.
	"net/http" // Importing the http package for HTTP client and server implementations.

	"github.com/anthdm/crypto-exchange/orderbook" // Importing the orderbook package for order book operations.
	"github.com/anthdm/crypto-exchange/server" // Importing the server package for server operations.
)

// Endpoint defines the base URL for the client.
const Endpoint = "http://localhost:3000"

// PlaceOrderParams represents the parameters required for placing an order.
type PlaceOrderParams struct {
	UserID int64   // UserID holds the user identifier.
	Bid    bool    // Bid indicates whether the order is a bid or ask.
	Price  float64 // Price is required only for placing LIMIT orders.
	Size   float64 // Size represents the quantity of the order.
}

// Client represents an HTTP client for interacting with the exchange server.
type Client struct {
	*http.Client // Embedding the http client for underlying HTTP operations.
}

// NewClient creates a new instance of the client.
func NewClient() *Client {
	return &Client{
		Client: http.DefaultClient, // Using the default HTTP client.
	}
}

// GetTrades retrieves the trades for a given market.
func (c *Client) GetTrades(market string) ([]*orderbook.Trade, error) {
	// Constructing the URL for fetching trades.
	endpoint := fmt.Sprintf("%s/trades/%s", Endpoint, market)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	// Performing the HTTP request.
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decoding the JSON response into trades.
	trades := []*orderbook.Trade{}
	if err := json.NewDecoder(resp.Body).Decode(&trades); err != nil {
		return nil, err
	}

	return trades, nil
}

// GetOrders retrieves the orders for a specific user.
func (c *Client) GetOrders(userID int64) (*server.GetOrdersResponse, error) {
	endpoint := fmt.Sprintf("%s/order/%d", Endpoint, userID)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	orders := server.GetOrdersResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, err
	}

	return &orders, nil
}

// PlaceMarketOrder places a market order.
func (c *Client) PlaceMarketOrder(p *PlaceOrderParams) (*server.PlaceOrderResponse, error) {
	// Constructing the request parameters.
	params := &server.PlaceOrderRequest{
		UserID: p.UserID,
		Type:   server.MarketOrder,
		Bid:    p.Bid,
		Size:   p.Size,
		Market: server.MarketETH,
	}

	// Encoding the request body into JSON.
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// Constructing the URL for placing the order.
	endpoint := Endpoint + "/order"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decoding the response into a PlaceOrderResponse.
	placeOrderResponse := &server.PlaceOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(placeOrderResponse); err != nil {
		return nil, err
	}

	return placeOrderResponse, nil
}

// GetBestAsk retrieves the best ask order.
func (c *Client) GetBestAsk() (*server.Order, error) {
	endpoint := fmt.Sprintf("%s/book/ETH/ask", Endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	order := &server.Order{}
	if err := json.NewDecoder(resp.Body).Decode(order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetBestBid retrieves the best bid order.
func (c *Client) GetBestBid() (*server.Order, error) {
	endpoint := fmt.Sprintf("%s/book/ETH/bid", Endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	order := &server.Order{}
	if err := json.NewDecoder(resp.Body).Decode(order); err != nil {
		return nil, err
	}

	return order, nil
}

// CancelOrder cancels an existing order.
func (c *Client) CancelOrder(orderID int64) error {
	endpoint := fmt.Sprintf("%s/order/%d", Endpoint, orderID)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// PlaceLimitOrder places a limit order.
func (c *Client) PlaceLimitOrder(p *PlaceOrderParams) (*server.PlaceOrderResponse, error) {
	// Validation for the size of the order.
	if p.Size == 0.0 {
		return nil, fmt.Errorf("size cannot be 0 when placing a limit order")
	}

	// Constructing the request parameters.
	params := &server.PlaceOrderRequest{
		UserID: p.UserID,
		Type:   server.LimitOrder,
		Bid:    p.Bid,
		Size:   p.Size,
		Price:  p.Price,
		Market: server.MarketETH,
	}

	// Encoding the request body into JSON.
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// Constructing the URL for placing the order.
	endpoint := Endpoint + "/order"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decoding the response into a PlaceOrderResponse.
	placeOrderResponse := &server.PlaceOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(placeOrderResponse); err != nil {
		return nil, err
	}

	return placeOrderResponse, nil
}
