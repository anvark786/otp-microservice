package otp

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func init() {

	redisClient = redis.NewClient(&redis.Options{

		Addr: "localhost:6379", // Redis server address,
		DB:   0,
	})

}

type OTPResponse struct {
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Otp        string `json:"otp,omitempty"`
}

func generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func GenerateOTPHandler(c *gin.Context) {
	otp := generateOTP()
	var requestData struct {
		OtpType string `json:"otpType"`
		UserID  int64  `json:"userID"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, OTPResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Invalid Request body.",
		})
		return
	}

	otpType := requestData.OtpType
	userId := requestData.UserID

	key := fmt.Sprintf("%s_user_%d", otpType, userId)
	println(key, "Key")

	err := redisClient.Set(context.Background(), key, otp, 5*time.Minute).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, OTPResponse{
			Status:     "error",
			Message:    "Internal Server Error.",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, OTPResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "OTP generated successfully.",
		Otp:        otp,
	})
}

func VerifyOTPHandler(c *gin.Context) {

	var requestData map[string]string

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, OTPResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Invalid Request body.",
		})
		return
	}
	key := requestData["key"]
	storedOTP, err := redisClient.Get(context.Background(), key).Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, OTPResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Otp expired or not found.",
		})
		return
	}
	receivedOTP := requestData["otp"]
	if receivedOTP == storedOTP {
		err := redisClient.Del(context.Background(), key).Err()
		if err != nil {
			log.Println("Error clearing OTP from Redis:", err)
		}
		c.JSON(http.StatusOK, OTPResponse{
			StatusCode: http.StatusOK,
			Status:     "success",
			Message:    "Otp verified successfully.",
		})
	} else {
		c.JSON(http.StatusUnauthorized, OTPResponse{

			StatusCode: http.StatusUnauthorized,
			Status:     "error",
			Message:    "Invalid OTP",
		})
		return
	}

}
