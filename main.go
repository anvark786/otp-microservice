package main

import (
	"log"
	"net/http"
	"os"
	"otp-microservice/internal/otp"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	expectedAPIKey := os.Getenv("API_KEY")

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		requestedAPIKey := c.GetHeader("Authorization")
		if requestedAPIKey != expectedAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	})

	router.POST("/generate-otp", otp.GenerateOTPHandler)
	router.POST("/verify-otp", otp.VerifyOTPHandler)

	router.Run(":8080")
}
