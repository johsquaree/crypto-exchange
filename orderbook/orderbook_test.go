package orderbook

import (
	"fmt"
	"reflect"
	"testing"
)

// assert checks if two values are deeply equal, and if not, logs an error with the provided testing.T instance.
func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

// TestLastMarketTrades tests the functionality of placing a market order and matching it against existing limit orders.
func TestLastMarketTrades(t *testing.T) {
	ob := NewOrderbook()
	price := 10000.0

	sellOrder := NewOrder(false, 10, 0)
	ob.PlaceLimitOrder(price, sellOrder)

	marketOrder := NewOrder(true, 10, 0)
	matches := ob.PlaceMarketOrder(marketOrder)
	assert(t, len(matches), 1)
	match := matches[0]

	assert(t, len(ob.Trades), 1)
	trade := ob.Trades[0]
	assert(t, trade.Price, price)
	assert(t, trade.Bid, marketOrder.Bid)
	assert(t, trade.Size, match.SizeFilled)
}

// TestLimit tests the basic functionality of the limit order data structure.
func TestLimit(t *testing.T) {
	l := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5, 0)
	buyOrderB := NewOrder(true, 8, 0)
	buyOrderC := NewOrder(true, 10, 0)

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)

	fmt.Println(l)
}

// TestPlaceLimitOrder tests the functionality of placing a limit order in the order book.
func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA := NewOrder(false, 10, 0)
	sellOrderB := NewOrder(false, 5, 0)
	ob.PlaceLimitOrder(10_000, sellOrderA)
	ob.PlaceLimitOrder(9_000, sellOrderB)

	assert(t, len(ob.Orders), 2)
	assert(t, ob.Orders[sellOrderA.ID], sellOrderA)
	assert(t, ob.Orders[sellOrderB.ID], sellOrderB)
	assert(t, len(ob.asks), 2)
}

// TestPlaceMarketOrder tests the functionality of placing a market order in the order book.
func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrder := NewOrder(false, 20, 0)
	ob.PlaceLimitOrder(10_000, sellOrder)

	buyOrder := NewOrder(true, 10, 0)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 10.0)
	assert(t, matches[0].Ask, sellOrder)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].SizeFilled, 10.0)
	assert(t, matches[0].Price, 10_000.0)
	assert(t, buyOrder.IsFilled(), true)
}

// TestPlaceMarketOrderMultiFill tests the functionality of placing a market order that fills multiple existing limit orders.
func TestPlaceMarketOrderMultiFill(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA := NewOrder(true, 5, 0) // filled fully
	buyOrderB := NewOrder(true, 8, 0) // partially filled
	buyOrderC := NewOrder(true, 1, 0)
	buyOrderD := NewOrder(true, 1, 0)

	ob.PlaceLimitOrder(5_000, buyOrderC)
	ob.PlaceLimitOrder(5_000, buyOrderD)
	ob.PlaceLimitOrder(9_000, buyOrderB)
	ob.PlaceLimitOrder(10_000, buyOrderA)

	assert(t, ob.BidTotalVolume(), 15.00)

	sellOrder := NewOrder(false, 10, 0)
	matches := ob.PlaceMarketOrder(sellOrder)

	assert(t, ob.BidTotalVolume(), 5.00)
	assert(t, len(ob.bids), 2)
	assert(t, len(matches), 2)
}

// TestCancelOrderAsk tests the functionality of cancelling an ask order in the order book.
func TestCancelOrderAsk(t *testing.T) {
	ob := NewOrderbook()
	sellOrder := NewOrder(false, 4, 0)
	price := 10_000.0
	ob.PlaceLimitOrder(price, sellOrder)

	assert(t, ob.AskTotalVolume(), 4.0)

	ob.CancelOrder(sellOrder)
	assert(t, ob.AskTotalVolume(), 0.0)

	_, ok := ob.Orders[sellOrder.ID]
	assert(t, ok, false)

	_, ok = ob.AskLimits[price]
	assert(t, ok, false)
}

// TestCancelOrderBid tests the functionality of cancelling a bid order in the order book.
func TestCancelOrderBid(t *testing.T) {
	ob := NewOrderbook()
	buyOrder := NewOrder(true, 4, 0)
	price := 10_000.0
	ob.PlaceLimitOrder(price, buyOrder)

	assert(t, ob.BidTotalVolume(), 4.0)

	ob.CancelOrder(buyOrder)
	assert(t, ob.BidTotalVolume(), 0.0)

	_, ok := ob.Orders[buyOrder.ID]
	assert(t, ok, false)

	_, ok = ob.BidLimits[price]
	assert(t, ok, false)
}
