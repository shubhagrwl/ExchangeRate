package response

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

// ErrorResponseData -
type ErrorResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorResponse -
type ErrorResponse struct {
	Success bool              `json:"success" default:"false"`
	Error   ErrorResponseData `json:"data"`
}

type Response struct {
	Success bool `json:"success"`
	Data    Data `json:"data"`
}

type ResponseV2 struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Request interface{} `json:"request,omitempty"`
	// List    interface{} `json:"list,omitempty"`
}
type Data struct {
	Message string `json:"message"`
}

type Rates struct {
	USD float64 `json:"USD"`
	EUR float64 `json:"EUR"`
	GBP float64 `json:"GBP"`
}
type ExchangeRatesResponse struct {
	Success   bool   `json:"success"`
	Timestamp int64  `json:"timestamp"`
	Base      string `json:"base"`
	Date      string `json:"date"`
	Rates     Rates  `json:"rates"`
}
