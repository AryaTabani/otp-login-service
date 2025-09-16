package main

import (
	"log"
	db "otp-login-service/DB"
	"otp-login-service/controllers"
	"otp-login-service/middleware"

	"github.com/gin-gonic/gin"

	_ "otp-login-service/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           OTP-Based Login API
// @version         1.0
// @description     A Go backend service for OTP login and registration.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	db.InitDB()

	router := gin.Default()

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/request-otp", controllers.RequestOTPHandler())
		authRoutes.POST("/verify-otp", controllers.VerifyOTPHandler())
	}

	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("", controllers.ListUsersHandler())
		userRoutes.GET("/:id", controllers.GetUserHandler())
		userRoutes.GET("/phone/:phone", controllers.GetUserByPhoneHandler())
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Server starting on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
