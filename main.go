package main

import (
	"otp-microservice/internal/otp"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.POST("/generate-otp", otp.GenerateOTPHandler)
	router.POST("/verify-otp", otp.VerifyOTPHandler)

	router.Run(":8080")
}
