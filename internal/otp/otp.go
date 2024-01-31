package otp

import (
	"fmt"
	"math/rand"
	"net/http"

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
}

func generateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func GenerateOTPHandler(context *gin.Context) {
	otp := generateOTP()
	fmt.Println("OTP", otp)

	var requestData struct {
		OtpType string `json:"otpType"`
		UserID  int64  `json:"userID"`
	}
	if err := context.ShouldBindJSON(&requestData); err != nil {
		context.JSON(http.StatusBadRequest, OTPResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "error",
			Message:    "Invalid Request body",
		})
		return
	}

	otpType := requestData.OtpType
	userId := requestData.UserID

	key := fmt.Sprintf("%s_user_%d", otpType, userId)

	fmt.Println("key:", key)

	context.JSON(http.StatusOK, OTPResponse{
		StatusCode: http.StatusOK,
		Status:     "success",
		Message:    "OTP generated successfully",
	})
}
