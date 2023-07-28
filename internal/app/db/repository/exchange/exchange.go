package exchange

import (
	"akcom/internal/app/constants"
	"akcom/internal/app/db"
	"context"
	"fmt"

	exchange_DBModels "akcom/internal/app/db/dto/exchange"
)

type IExchangeRepository interface {
	CreateExchangeRate(ctx context.Context, exchange exchange_DBModels.Exchange) error
	GetExchangeRates(ctx context.Context, from, to string) ([]*exchange_DBModels.Exchange, error)
}

type ExchangeRepository struct {
	DBService *db.DBService
}

func NewExchangeRepository(dbService *db.DBService) IExchangeRepository {
	return &ExchangeRepository{
		DBService: dbService,
	}
}

func (u *ExchangeRepository) CreateExchangeRate(ctx context.Context, exchange exchange_DBModels.Exchange) error {
	tx := u.DBService.GetDB().Begin()                       // Start a database transaction
	defer tx.Rollback()                                     // Rollback the transaction if not committed
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE) // Set the database log mode

	conflict := fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s=excluded.%s, %s=excluded.%s, %s=excluded.%s",
		exchange_DBModels.COLUMN_EXCHANGE_DATE,

		exchange_DBModels.COLUMN_EUR,
		exchange_DBModels.COLUMN_EUR,

		exchange_DBModels.COLUMN_USD,
		exchange_DBModels.COLUMN_USD,

		exchange_DBModels.COLUMN_GBP,
		exchange_DBModels.COLUMN_GBP)

	err := tx.Table(exchange_DBModels.TABLE_NAME).Set("gorm:insert_option", conflict).Create(&exchange).Error // Create a new exchange record in the table
	if err != nil {
		return err // Return the error if the creation fails
	}
	tx.Commit() // Commit the transaction

	return nil // Return nil to indicate successful exchange creation
}

func (u *ExchangeRepository) GetExchangeRates(ctx context.Context, from, to string) ([]*exchange_DBModels.Exchange, error) {
	tx := u.DBService.GetDB()                               // Get the database instance
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE) // Set the database log mode

	var exchange []*exchange_DBModels.Exchange // Slice to store the retrieved user exchanges

	whr := fmt.Sprintf("%s BETWEEN '%s' AND '%s'", exchange_DBModels.COLUMN_EXCHANGE_DATE, from, to)
	// Retrieve exchanges with pagination
	if err := tx.Table(exchange_DBModels.TABLE_NAME).
		Where(whr).Scan(&exchange).Error; err != nil {
		return exchange, err // Return the retrieved exchanges and error, if any
	}

	return exchange, nil // Return the retrieved user exchanges and no error
}
