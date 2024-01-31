package main

import (
	"otp-microservice/internal/otp"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.POST("/generate-otp", otp.GenerateOTPHandler)
	router.Run(":8080")
}
