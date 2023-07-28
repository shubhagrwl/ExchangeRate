package user

import (
	"akcom/internal/app/apiclient"
	"akcom/internal/app/constants"
	"akcom/internal/app/controller"
	"akcom/internal/app/service/correlation"
	"akcom/internal/app/service/logger"
	"context"
	"net/http"
	"time"

	exchange_DBModels "akcom/internal/app/db/dto/exchange"
	exchangeDB "akcom/internal/app/db/repository/exchange"

	"github.com/gin-gonic/gin"
)

// IExchangeController is an interface that defines the methods for a exchange controller.
type IExchangeController interface {
	GetExchangeRates(c *gin.Context)
	GetPreviousTenDaysExchangeRates(ctx context.Context)
}

// ExchangeController is a struct that implements the IExchangeController interface.
type ExchangeController struct {
	Exchange         apiclient.IExchange
	ExchangeDBClient exchangeDB.IExchangeRepository
}

// NewExchangeController is a constructor function that creates a new ExchangeController.
func NewExchangeController(exchange apiclient.IExchange, exchangeDBClient exchangeDB.IExchangeRepository) IExchangeController {
	return &ExchangeController{
		Exchange:         exchange,
		ExchangeDBClient: exchangeDBClient,
	}
}

// exchange handles creating a new exchange
func (u ExchangeController) GetExchangeRates(c *gin.Context) {
	ctx := correlation.WithReqContext(c) // Get the request context
	log := logger.Logger(ctx)            // Get the logger

	// Retrieve the user context from the request
	_, exist := c.Get(constants.CTK_CLAIM_KEY.String())
	if !exist {
		log.Errorf("context not found")
		controller.RespondWithError(c, http.StatusUnauthorized, "Context Not Found")
		return
	}

	var from, to string
	currentDate := time.Now()

	date := c.Query("date")
	if date == "" {
		from = currentDate.AddDate(0, 0, -10).Format("2006-01-02")
		to = currentDate.Format("2006-01-02")
	} else {
		from = date
		to = date
	}

	// fetch the current date exchange rates and update in DB
	if date == currentDate.Format("2006-01-02") || from == currentDate.Format("2006-01-02") || to == currentDate.Format("2006-01-02") {
		exchangeRate, err := u.Exchange.GetHistoricalExchangeRates(ctx, currentDate.Format("2006-01-02"))
		if err != nil {
			log.Error("error while getting exchange rates for this date", currentDate.Format("2006-01-02"))

		}

		u.ExchangeDBClient.CreateExchangeRate(ctx, exchange_DBModels.Exchange{
			ExchangeDate: currentDate.Format("2006-01-02"),
			Base:         exchangeRate.Base,
			Usd:          exchangeRate.Rates.USD,
			Eur:          exchangeRate.Rates.EUR,
			Gbp:          exchangeRate.Rates.GBP,
		})
	}

	exchangeRates, err := u.ExchangeDBClient.GetExchangeRates(ctx, from, to)
	if err != nil {
		log.Errorf("Error getting exchange rates", err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	// If the exchange is created successfully, respond with success and the created exchange
	controller.RespondWithSuccess(c, http.StatusOK, "Exchange Rates", exchangeRates)
}

// exchange handles creating a new exchange
func (u ExchangeController) GetPreviousTenDaysExchangeRates(ctx context.Context) {

	log := logger.Logger(ctx) // Get the logger
	log.Info("started collecting historical exchange rates")
	// Get the current date
	currentDate := time.Now()

	// Loop through the previous 10 days
	for i := 0; i < 10; i++ {
		// Calculate the date for the current iteration
		previousDate := currentDate.AddDate(0, 0, -i).Format("2006-01-02")

		exchangeRate, err := u.Exchange.GetHistoricalExchangeRates(ctx, previousDate)
		if err != nil {
			log.Error("error while getting exchange rates for this date", previousDate)
			continue
		}

		u.ExchangeDBClient.CreateExchangeRate(ctx, exchange_DBModels.Exchange{
			ExchangeDate: previousDate,
			Base:         exchangeRate.Base,
			Usd:          exchangeRate.Rates.USD,
			Eur:          exchangeRate.Rates.EUR,
			Gbp:          exchangeRate.Rates.GBP,
		})
	}
}
