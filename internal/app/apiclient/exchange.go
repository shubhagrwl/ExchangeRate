package apiclient

import (
	"akcom/internal/app/constants"
	"akcom/internal/app/service/dto/response"
	"akcom/internal/app/service/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// IExchange is an interface that defines the methods for a exchange.
type IExchange interface {
	GetHistoricalExchangeRates(ctx context.Context, date string) (response.ExchangeRatesResponse, error)
}

// Exchange is a struct that implements the IExchange interface.
type Exchange struct {
	ApiClient IApiClient
}

// NewExchange is a constructor function that creates a new Exchange.
func NewExchange(apiClient IApiClient,
) IExchange {
	return &Exchange{
		ApiClient: apiClient,
	}
}

func (e *Exchange) GetHistoricalExchangeRates(ctx context.Context, date string) (response.ExchangeRatesResponse, error) {
	log := logger.Logger(ctx)
	var rspr response.ExchangeRatesResponse

	exchangeUrl := fmt.Sprintf(constants.Config.Exchange.EXCHANGE_URL, date,
		constants.Config.Exchange.EXCHANGE_KEY)

	// * Creates a request body
	req, requestError := e.ApiClient.CreateJSONRequest(ctx, http.MethodGet, exchangeUrl, nil, nil, nil)
	if requestError != nil {
		log.Errorf("Error for calling exchange api, Exchange: ", exchangeUrl, ", err: ", requestError)
		return rspr, requestError
	}

	// * Makes an api call to the OTP service
	status, rsp, err := e.ApiClient.RestExecute(ctx, req)
	if err != nil || status != http.StatusOK {
		log.Errorf("Error for calling exchange api, Exchange: ", exchangeUrl, ", err: ", requestError)
		return rspr, err
	}

	err = json.Unmarshal([]byte(rsp), &rspr)
	if err != nil {
		return rspr, err
	}
	return rspr, nil
}
