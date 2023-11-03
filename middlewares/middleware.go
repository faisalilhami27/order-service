package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"order-service/config"
	"order-service/constant"
	constantError "order-service/constant/error"
	"order-service/utils/response"
)

func HandlePanic(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.SetLevel(log.ErrorLevel)
			log.Errorf("error occured: %v", r)
			c.JSON(http.StatusBadRequest, response.Response{
				Status:  constantError.Error,
				Message: r.(error), //nolint:forcetypeassert
			})
			return
		}
	}()
	c.Next()
}

func ValidateAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		signatureKey := config.Config.SignatureKey
		apiKey := c.GetHeader(constant.XApiKey)
		requestAt := c.GetHeader(constant.XRequestAt)
		serviceName := c.GetHeader(constant.XServiceName)
		validateAPIKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)
		hash := sha256.New()
		hash.Write([]byte(validateAPIKey))
		apiKeyHash := hex.EncodeToString(hash.Sum(nil))

		if apiKey != apiKeyHash {
			newError := fmt.Sprintf("Unauthorized") //nolint:gosimple
			c.JSON(http.StatusUnauthorized, response.Response{
				Status:  constantError.Error,
				Message: newError,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, response.Response{
				Status:  constantError.Error,
				Message: constantError.ErrTooManyRequest.Error(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
