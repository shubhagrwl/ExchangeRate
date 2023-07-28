package user

import (
	"akcom/internal/app/constants"
	"akcom/internal/app/controller"
	users_DBModels "akcom/internal/app/db/dto/users"
	userDB "akcom/internal/app/db/repository/user"

	"encoding/json"
	"fmt"

	"akcom/internal/app/service/correlation"
	"akcom/internal/app/service/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// IUserController is an interface that defines the methods for a user controller.
type IUserController interface {
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
}

// UserController is a struct that implements the IUserController interface.
type UserController struct {
	UserDBClient userDB.IUserRepository // UserDBClient represents the database client for user-related operations.
}

// NewUserController is a constructor function that creates a new UserController.
func NewUserController(UserClient userDB.IUserRepository) IUserController {
	return &UserController{
		UserDBClient: UserClient,
	}
}

func (u UserController) GetProfile(c *gin.Context) {
	ctx := correlation.WithReqContext(c) // Get the request context
	log := logger.Logger(ctx)            // Get the logger

	context, exist := c.Get(constants.CTK_CLAIM_KEY.String()) // Retrieve the user context from the request
	if !exist {
		log.Errorf("context not found")
		controller.RespondWithError(c, http.StatusUnauthorized, "Context Not Found")
		return
	}
	usr := context.(*users_DBModels.Users) // Type assertion to retrieve the user information

	user, err := u.UserDBClient.GetUser(ctx, fmt.Sprintf("%s='%s' AND %s=true",
		users_DBModels.COLUMN_EMAIL, usr.Email,
		users_DBModels.COLUMN_STATUS))
	if err != nil {
		log.Errorf("Error getting user", err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}
	user.Password = "" // Clear the password field for security reasons

	// Respond with success message and the user profile (with password field cleared)
	controller.RespondWithSuccess(c, http.StatusAccepted, "User Profile", user)
}

func (u UserController) UpdateProfile(c *gin.Context) {
	ctx := correlation.WithReqContext(c) // Get the request context
	log := logger.Logger(ctx)            // Get the logger

	context, exist := c.Get(constants.CTK_CLAIM_KEY.String()) // Retrieve the user context from the request
	if !exist {
		log.Errorf("context not found")
		controller.RespondWithError(c, http.StatusUnauthorized, "Context Not Found")
		return
	}
	usr := context.(*users_DBModels.Users) // Type assertion to retrieve the user information

	dataFromBody := users_DBModels.Users{}                       // Create an empty Users struct
	err := json.NewDecoder(c.Request.Body).Decode(&dataFromBody) // Decode the request body into the Users struct
	if err != nil {
		log.Errorf(constants.BadRequest, err)
		controller.RespondWithError(c, http.StatusBadRequest, constants.BadRequest)
		return
	}

	var patcher = make(map[string]interface{}) // Create a patcher map to hold the fields to be updated

	// Check if each field is present in the request body and add it to the patcher map if not empty
	if dataFromBody.FirstName != "" {
		patcher["first_name"] = dataFromBody.FirstName
	}
	if dataFromBody.LastName != "" {
		patcher["last_name"] = dataFromBody.LastName
	}
	if dataFromBody.CountryCode != "" {
		patcher["country_code"] = dataFromBody.CountryCode
	}
	if dataFromBody.Phone != "" {
		patcher["phone"] = dataFromBody.Phone
	}
	if dataFromBody.Address != "" {
		patcher["address"] = dataFromBody.Address
	}

	filter := fmt.Sprintf("%s=%d", users_DBModels.COLUMN_ID, usr.Id) // Create a filter string to match the user ID

	if err := u.UserDBClient.UpdateUser(ctx, filter, patcher); err != nil {
		// Update the user using the UserDBClient and the filter and patcher
		log.Error("error while updating user", err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	r, _ := u.UserDBClient.GetUser(ctx, fmt.Sprintf("%s=%d AND %s=true",
		users_DBModels.COLUMN_ID, usr.Id,
		users_DBModels.COLUMN_STATUS))
	r.Password = ""

	// Respond with success message and the updated user profile (with password field cleared)
	controller.RespondWithSuccess(c, http.StatusAccepted, "User Profile Updated Successfully", r)
}
