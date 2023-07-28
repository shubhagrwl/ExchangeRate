package user

import (
	"context"
	"fmt"

	"akcom/internal/app/constants"
	db "akcom/internal/app/db"
	users_DBModels "akcom/internal/app/db/dto/users"
	"akcom/internal/app/service/dto/request"
)

type IUserRepository interface {
	GetUser(ctx context.Context, whr string) (users_DBModels.Users, error)
	GetAllUsers(ctx context.Context, pagination request.Pagination) ([]*users_DBModels.Users, error)

	CreateUser(ctx context.Context, User users_DBModels.Users) error
	UpdateUser(ctx context.Context, whr string, patch map[string]interface{}) error
}

type UserRepository struct {
	DBService *db.DBService
}

func NewUserRepository(dbService *db.DBService) IUserRepository {
	return &UserRepository{
		DBService: dbService,
	}
}

func (u *UserRepository) CreateUser(ctx context.Context, User users_DBModels.Users) error {
	tx := u.DBService.GetDB().Begin()                       // Start a database transaction
	defer tx.Rollback()                                     // Rollback the transaction if not committed
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE) // Set the database log mode

	if err := tx.Table(users_DBModels.TABLE_NAME).Create(&User).Error; err != nil {
		return err // Return the error if user creation fails
	}
	tx.Commit() // Commit the transaction

	return nil // Return nil to indicate successful user creation
}

func (u *UserRepository) GetUser(ctx context.Context, whr string) (users_DBModels.Users, error) {
	tx := u.DBService.GetDB()                               // Get the database instance
	var user users_DBModels.Users                           // Variable to store the retrieved user
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE) // Set the database log mode

	if err := tx.Table(users_DBModels.TABLE_NAME).Where(whr).Scan(&user).Error; err != nil {
		return user, err // Return the retrieved user and error, if any
	}

	return user, nil // Return the retrieved user and no error
}

func (u *UserRepository) UpdateUser(ctx context.Context, whr string, patch map[string]interface{}) error {
	tx := u.DBService.GetDB().Begin()                       // Start a database transaction
	defer tx.Rollback()                                     // Rollback the transaction if not committed
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE) // Set the database log mode

	if err := tx.Table(users_DBModels.TABLE_NAME).Where(whr).Update(patch).Error; err != nil {
		return err // Return the error if user update fails
	}
	tx.Commit() // Commit the transaction

	return nil // Return nil to indicate successful user update
}

func (u *UserRepository) GetAllUsers(ctx context.Context, pagination request.Pagination) ([]*users_DBModels.Users, error) {
	tx := u.DBService.GetDB()                               // Get the database instance
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE) // Set the database log mode

	var users []*users_DBModels.Users // Slice to store the retrieved users
	var whr string                    // Variable to hold the WHERE condition

	if pagination.Query != "" {
		// Create a WHERE condition to filter users based on the search query
		whr = fmt.Sprintf(`
		(coalesce(%s) || coalesce(%s) || coalesce(%s) || coalesce(%s))  LIKE '%%%s%%'`,
			users_DBModels.COLUMN_EMAIL,
			users_DBModels.COLUMN_FIRST_NAME,
			users_DBModels.COLUMN_LAST_NAME,
			users_DBModels.COLUMN_PHONE,
			pagination.Query,
		)
	}

	if pagination.GetAllData {
		// Retrieve all users without pagination
		if err := tx.Table(users_DBModels.TABLE_NAME).
			Order(fmt.Sprintf("%s %s", users_DBModels.COLUMN_CREATED_AT, pagination.Sort)).
			Where(whr).Scan(&users).Error; err != nil {
			return users, err // Return the retrieved users and error, if any
		}
	} else {
		// Retrieve users with pagination
		if err := tx.Table(users_DBModels.TABLE_NAME).Limit(*pagination.Limit).Offset(pagination.Offset).
			Order(fmt.Sprintf("%s %s", users_DBModels.COLUMN_CREATED_AT, pagination.Sort)).
			Where(whr).Scan(&users).Error; err != nil {
			return users, err // Return the retrieved users and error, if any
		}
	}

	return users, nil // Return the retrieved users and no error
}
