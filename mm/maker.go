package mm // Package declaration

import (
	"time" // Importing the time package for time-related functionalities

	"github.com/anthdm/crypto-exchange/client" // Importing client package from anthdm/crypto-exchange
	"github.com/sirupsen/logrus" // Importing logrus package for logging
)

// Config struct to hold configuration parameters
type Config struct {
	UserID         int64         // User ID for the market maker
	OrderSize      float64       // Size of each order
	MinSpread      float64       // Minimum spread allowed
	SeedOffset     float64       // Offset for seeding the market
	ExchangeClient *client.Client // Client for interacting with the exchange
	MakeInterval   time.Duration // Interval for making orders
	PriceOffset    float64       // Offset for adjusting price
}

// MarketMaker struct to represent a market maker
type MarketMaker struct {
	userID         int64           // User ID for the market maker
	orderSize      float64         // Size of each order
	minSpread      float64         // Minimum spread allowed
	seedOffset     float64         // Offset for seeding the market
	priceOffset    float64         // Offset for adjusting price
	exchangeClient *client.Client  // Client for interacting with the exchange
	makeInterval   time.Duration   // Interval for making orders
}

// NewMakerMaker creates a new MarketMaker instance with provided config
func NewMakerMaker(cfg Config) *MarketMaker {
	return &MarketMaker{
		userID:         cfg.UserID,
		orderSize:      cfg.OrderSize,
		minSpread:      cfg.MinSpread,
		seedOffset:     cfg.SeedOffset,
		exchangeClient: cfg.ExchangeClient,
		makeInterval:   cfg.MakeInterval,
		priceOffset:    cfg.PriceOffset,
	}
}

// Start starts the market maker
func (mm *MarketMaker) Start() {
	// Logging the start of the market maker with relevant information
	logrus.WithFields(logrus.Fields{
		"id":           mm.userID,
		"orderSize":    mm.orderSize,
		"makeInterval": mm.makeInterval,
		"minSpread":    mm.minSpread,
		"priceOffset":  mm.priceOffset,
	}).Info("starting market maker")

	// Start the maker loop in a separate goroutine
	go mm.makerLoop()
}

// makerLoop is the main loop for the market maker
func (mm *MarketMaker) makerLoop() {
	// Creating a ticker for the make interval
	ticker := time.NewTicker(mm.makeInterval)

	for {
		// Getting the best bid from the exchange
		bestBid, err := mm.exchangeClient.GetBestBid()
		if err != nil {
			logrus.Error(err)
			break
		}

		// Getting the best ask from the exchange
		bestAsk, err := mm.exchangeClient.GetBestAsk()
		if err != nil {
			logrus.Error(err)
			break
		}

		// If both bid and ask prices are zero, seed the market
		if bestAsk.Price == 0 && bestBid.Price == 0 {
			if err := mm.seedMarket(); err != nil {
				logrus.Error(err)
				break
			}
			continue
		}

		// Adjusting bid price if necessary
		if bestBid.Price == 0 {
			bestBid.Price = bestAsk.Price - mm.priceOffset*2
		}

		// Adjusting ask price if necessary
		if bestAsk.Price == 0 {
			bestAsk.Price = bestBid.Price + mm.priceOffset*2
		}

		// Calculating spread
		spread := bestAsk.Price - bestBid.Price

		// If spread is less than or equal to minSpread, continue to next iteration
		if spread <= mm.minSpread {
			continue
		}

		// Placing bid order
		if err := mm.placeOrder(true, bestBid.Price+mm.priceOffset); err != nil {
			logrus.Error(err)
			break
		}

		// Placing ask order
		if err := mm.placeOrder(false, bestAsk.Price-mm.priceOffset); err != nil {
			logrus.Error(err)
			break
		}

		// Waiting for the next tick
		<-ticker.C
	}
}

// placeOrder places an order on the exchange
func (mm *MarketMaker) placeOrder(bid bool, price float64) error {
	// Creating order parameters
	bidOrder := &client.PlaceOrderParams{
		UserID: mm.userID,
		Size:   mm.orderSize,
		Bid:    bid,
		Price:  price,
	}
	// Placing the order on the exchange
	_, err := mm.exchangeClient.PlaceLimitOrder(bidOrder)
	return err
}

// seedMarket seeds the market by placing initial bid and ask orders
func (mm *MarketMaker) seedMarket() error {
	// Simulating fetching current ETH price from another exchange
	currentPrice := simulateFetchCurrentETHPrice()

	// Logging seeding of the market with relevant information
	logrus.WithFields(logrus.Fields{
		"currentETHPrice": currentPrice,
		"seedOffset":      mm.seedOffset,
	}).Info("orderbooks empty => seeding market!")

	// Placing bid order to seed the market
	bidOrder := &client.PlaceOrderParams{
		UserID: mm.userID,
		Size:   mm.orderSize,
		Bid:    true,
		Price:  currentPrice - mm.seedOffset,
	}
	_, err := mm.exchangeClient.PlaceLimitOrder(bidOrder)
	if err != nil {
		return err
	}

	// Placing ask order to seed the market
	askOrder := &client.PlaceOrderParams{
		UserID: mm.userID,
		Size:   mm.orderSize,
		Bid:    false,
		Price:  currentPrice + mm.seedOffset,
	}
	_, err = mm.exchangeClient.PlaceLimitOrder(askOrder)

	return err
}

// simulateFetchCurrentETHPrice simulates fetching current ETH price from another exchange
func simulateFetchCurrentETHPrice() float64 {
	// Simulating delay in fetching the price
	time.Sleep(80 * time.Millisecond)

	// Returning a simulated current ETH price
	return 1000.0
}
