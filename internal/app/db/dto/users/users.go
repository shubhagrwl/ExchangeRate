package users

import (
	"fmt"
	"time"
)

const (
	TABLE_NAME          = "users"
	COLUMN_ID           = "id"
	COLUMN_FIRST_NAME   = "first_name"
	COLUMN_LAST_NAME    = "last_name"
	COLUMN_EMAIL        = "email"
	COLUMN_PASSWORD     = "password"
	COLUMN_COUNTRY_CODE = "country_code"
	COLUMN_PHONE        = "phone"
	COLUMN_ADDRESS      = "address"
	COLUMN_STATUS       = "status"
	COLUMN_CREATED_AT   = "created_at"
)

type Users struct {
	Id          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Password    string    `json:"password,omitempty"`
	CountryCode string    `json:"country_code"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	Status      bool      `json:"status"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

func (u Users) ValidateSignUpDetails() error {
	if u.FirstName == "" {
		return fmt.Errorf("first_name can't be empty")
	}
	if u.Email == "" {
		return fmt.Errorf("email can't be empty")
	}
	if u.Password == "" {
		return fmt.Errorf("password can't be empty")
	}
	return nil
}
