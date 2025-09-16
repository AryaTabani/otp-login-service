// AuthController.go

package controllers

import (
	"errors"
	"net/http"
	"otp-login-service/models"
	"otp-login-service/services"

	"github.com/gin-gonic/gin"
)

// @Summary      Request an OTP
// @Description  Sends a phone number to generate an OTP for login or registration.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body models.RequestOTPPayload true "Phone Number"
// @Success      200  {object}  models.APIResponse[any]
// @Failure      400  {object}  models.APIResponse[any]
// @Failure      429  {object}  models.APIResponse[any]
// @Router       /auth/request-otp [post]
func RequestOTPHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload models.RequestOTPPayload
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, models.APIResponse[any]{Success: false, Error: "Invalid request body: " + err.Error()})
			return
		}

		err := services.RequestOTP(payload.PhoneNumber)
		if err != nil {
			if errors.Is(err, services.ErrRateLimitExceeded) {
				ctx.JSON(http.StatusTooManyRequests, models.APIResponse[any]{Success: false, Error: err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, models.APIResponse[any]{Success: false, Error: "Failed to generate OTP"})
			return
		}
		ctx.JSON(http.StatusOK, models.APIResponse[any]{
			Success: true,
			Message: "OTP generated and printed to console. It will expire in 2 minutes.",
		})
	}
}

// @Summary      Verify OTP & Login/Register
// @Description  Verifies OTP. If valid, logs in an existing user or registers a new one, returning a JWT.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body models.VerifyOTPPayload true "Phone Number and OTP"
// @Success      200  {object}  models.APIResponse[models.LoginSuccessResponse]
// @Failure      400  {object}  models.APIResponse[any]
// @Failure      401  {object}  models.APIResponse[any]
// @Router       /auth/verify-otp [post]
func VerifyOTPHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var payload models.VerifyOTPPayload
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, models.APIResponse[any]{Success: false, Error: "Invalid request body: " + err.Error()})
			return
		}

		token, err := services.VerifyOTP(payload.PhoneNumber, payload.OTP)
		if err != nil {
			if errors.Is(err, services.ErrInvalidOTP) {
				ctx.JSON(http.StatusUnauthorized, models.APIResponse[any]{Success: false, Error: err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, models.APIResponse[any]{Success: false, Error: "Failed to verify OTP"})
			return
		}

		responseData := models.LoginSuccessResponse{Token: token}
		ctx.JSON(http.StatusOK, models.APIResponse[models.LoginSuccessResponse]{
			Success: true,
			Message: "User logged in successfully",
			Data:    responseData,
		})
	}
}
