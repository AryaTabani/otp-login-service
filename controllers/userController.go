package controllers

import (
	"net/http"
	"otp-login-service/models"
	"otp-login-service/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary      Get a single user by ID
// @Description  Retrieves the details of a single user by their unique ID.
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  models.APIResponse[models.User]
// @Failure      401  {object}  models.APIResponse[any]
// @Failure      404  {object}  models.APIResponse[any]
// @Security     BearerAuth
// @Router       /users/{id} [get]
func GetUserHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.APIResponse[any]{Success: false, Error: "Invalid user ID"})
			return
		}

		user, err := services.GetUserByID(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, models.APIResponse[any]{Success: false, Error: "User not found"})
			return
		}

		ctx.JSON(http.StatusOK, models.APIResponse[*models.User]{Success: true, Data: user})
	}
}

// @Summary      List all users
// @Description  Retrieves a paginated list of all registered users. Supports searching by phone number.
// @Tags         Users
// @Produce      json
// @Param        page   query     int  false  "Page number"
// @Param        limit  query     int  false  "Users per page"
// @Param        search query     string false "Search by phone number"
// @Success      200    {object}  models.APIResponse[models.UserListResponse]
// // @Failure      401    {object}  models.APIResponse[any]
// @Security     BearerAuth
// @Router       /users [get]
func ListUsersHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
		search := ctx.Query("search")

		users, total := services.ListUsers(page, limit, search)

		responseData := models.UserListResponse{
			Users: users,
			Total: total,
			Page:  page,
			Limit: limit,
		}
		ctx.JSON(http.StatusOK, models.APIResponse[models.UserListResponse]{
			Success: true,
			Data:    responseData,
		})
	}
}

// @Summary      Get a single user by phone number
// @Description  Retrieves the details of a single user by their phone number.
// @Tags         Users
// @Produce      json
// @Param        phone   path      string  true  "Phone Number"
// @Success      200  {object}  models.APIResponse[models.User]
// @Failure      401  {object}  models.APIResponse[any]
// @Failure      404  {object}  models.APIResponse[any]
// @Security     BearerAuth
// @Router       /users/phone/{phone} [get]
func GetUserByPhoneHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		phone := ctx.Param("phone")
		if phone == "" {
			ctx.JSON(http.StatusBadRequest, models.APIResponse[any]{Success: false, Error: "Phone number is required"})
			return
		}

		user, err := services.GetUserByPhoneNumber(phone)
		if err != nil {
			ctx.JSON(http.StatusNotFound, models.APIResponse[any]{Success: false, Error: "User not found"})
			return
		}

		ctx.JSON(http.StatusOK, models.APIResponse[*models.User]{Success: true, Data: user})
	}
}
