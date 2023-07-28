package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"akcom/internal/app/api/middleware/jwt"
	"akcom/internal/app/constants"
	"akcom/internal/app/controller"
	users_DBModels "akcom/internal/app/db/dto/users"
	userDB "akcom/internal/app/db/repository/user"

	"akcom/internal/app/service/correlation"
	"akcom/internal/app/service/dto/request"
	"akcom/internal/app/service/logger"
	"akcom/internal/app/service/util"

	"github.com/gin-gonic/gin"
)

// VerifyUser represents a user to be verified
type VerifyUser struct {
	Email       string
	OTP         string
	CreatedTime time.Time
}

// NewUserMapper is a map to store new user verification data
var NewUserMapper = make(map[string]VerifyUser)

// IOAuthController represents the interface for OAuthController
type IOAuthController interface {
	SignUp(c *gin.Context)
	Login(c *gin.Context)

	GetAllUser(c *gin.Context)
	DisableUser(c *gin.Context)
}

// OAuthController is the implementation of the IOAuthController interface
type OAuthController struct {
	UserDBClient userDB.IUserRepository
	JWT          jwt.IJwtService
}

// NewOAuthController creates a new instance of OAuthController
func NewOAuthController(
	UserClient userDB.IUserRepository,
	jwt jwt.IJwtService,
) IOAuthController {
	return &OAuthController{
		UserDBClient: UserClient,
		JWT:          jwt,
	}
}

// SignUp handles the sign-up functionality
func (u OAuthController) SignUp(c *gin.Context) {
	ctx := correlation.WithReqContext(c) // Get the request context
	log := logger.Logger(ctx)            // Get the logger

	dataFromBody := users_DBModels.Users{}                       // Create an empty Users struct
	err := json.NewDecoder(c.Request.Body).Decode(&dataFromBody) // Decode the request body into the Users struct
	if err != nil {
		log.Errorf(constants.BadRequest, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BadRequest)
		return
	}

	if err := dataFromBody.ValidateSignUpDetails(); err != nil {
		// Validate the sign-up details
		controller.RespondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := util.GenerateHash(dataFromBody.Password) // Generate a hashed password
	if err != nil {
		log.Errorf(constants.InternalServerError, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	dataFromBody.Password = hashedPassword     // Set the hashed password in the Users struct
	dataFromBody.Status = true                 // Set the status to true
	dataFromBody.CreatedAt = time.Now()        // Set the current time as the creation time
	dataFromBody.Role = string(constants.USER) // Set the role as "USER"

	if err = u.UserDBClient.CreateUser(ctx, dataFromBody); err != nil {
		// Create a new user using the UserDBClient
		log.Errorf(constants.InternalServerError, err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	// Respond with success message and nil data
	controller.RespondWithSuccess(c, http.StatusAccepted, "User Created Successfully", nil)
}

// Login handles the login functionality
func (u OAuthController) Login(c *gin.Context) {
	ctx := correlation.WithReqContext(c) // Get the request context
	log := logger.Logger(ctx)            // Get the logger

	dataFromBody := users_DBModels.Users{}                       // Create an empty Users struct
	err := json.NewDecoder(c.Request.Body).Decode(&dataFromBody) // Decode the request body into the Users struct
	if err != nil {
		log.Errorf(constants.BadRequest, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BadRequest)
		return
	}

	user, err := u.UserDBClient.GetUser(ctx, fmt.Sprintf("%s='%s' AND %s=true",
		users_DBModels.COLUMN_EMAIL, dataFromBody.Email,
		users_DBModels.COLUMN_STATUS))
	if err != nil {
		log.Errorf("Error getting user", err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	if !util.ValidatePassword(dataFromBody.Password, user.Password) {
		// Validate the provided password with the stored password
		log.Errorf("wrong credentials")
		controller.RespondWithError(c, http.StatusUnauthorized, "invalid email/password")
		return
	}

	token, err := u.JWT.CreateNewTokens(ctx, user) // Create new access tokens for the user
	if err != nil {
		log.Errorf("error while creating access token", err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	// Respond with success message and the generated token
	controller.RespondWithSuccess(c, http.StatusAccepted, "Login Successfully", token)
}

// GetAllUser retrieves all users with pagination
func (u OAuthController) GetAllUser(c *gin.Context) {
	ctx := correlation.WithReqContext(c) // Get the request context
	log := logger.Logger(ctx)            // Get the logger

	var pagination request.Pagination // Create a Pagination struct

	// Bind the query parameters to the Pagination struct
	if err := c.ShouldBindQuery(&pagination); err != nil {
		log.Error("error while fetching pagination", err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BadRequest)
		return
	}

	pagination.Validate() // Validate the pagination parameters

	users, err := u.UserDBClient.GetAllUsers(ctx, pagination) // Get all users using the UserDBClient and pagination
	if err != nil {
		log.Errorf("Error getting user", err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	// Respond with success and the list of users
	controller.RespondWithSuccess(c, http.StatusAccepted, "Users List", users)
}

// DisableUser disables a user based on their ID
func (u OAuthController) DisableUser(c *gin.Context) {
	ctx := correlation.WithReqContext(c) // Get the request context
	log := logger.Logger(ctx)            // Get the logger

	id, err := strconv.Atoi(c.Query("id")) // Get the user ID from the query parameter
	if err != nil {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BadRequest)
		return
	}

	whr := fmt.Sprintf("%s=%d", users_DBModels.COLUMN_ID, id) // Create a WHERE clause to match the user ID
	var patcher = make(map[string]interface{})
	patcher[users_DBModels.COLUMN_STATUS] = false // Set the status to false for disabling the user

	if err := u.UserDBClient.UpdateUser(ctx, whr, patcher); err != nil {
		// Update the user using the UserDBClient
		log.Error("error while updating user", err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	// Respond with success message and nil data
	controller.RespondWithSuccess(c, http.StatusAccepted, "User Disable Successfully", nil)
}
