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

	clientConfig "order-service/clients/config"

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
			c.JSON(http.StatusUnauthorized, response.Response{
				Status:  constantError.Error,
				Message: constantError.ErrUnauthorized,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthenticateRBAC() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(constant.Authorization)
		if token == "" {
			c.JSON(http.StatusUnauthorized, response.Response{
				Status:  constantError.Error,
				Message: constantError.ErrUnauthorized.Error(),
			})
			c.Abort()
			return
		}

		client := clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.RBAC.Host),
			clientConfig.WithSecretKey(config.Config.InternalService.RBAC.SecretKey))

		rbac := NewRBACMiddleware(client)
		user, err := rbac.GetUserLogin(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Response{
				Status:  constantError.Error,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set(constant.Token, token)
		c.Set(constant.UserLogin, user)
		c.Next()
	}
}

func CheckPermission(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := c.Get(constant.Token)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.Response{
				Status:  constantError.Error,
				Message: constantError.ErrUnauthorized.Error(),
			})
			c.Abort()
			return
		}

		client := clientConfig.NewClientConfig(
			clientConfig.WithBaseURL(config.Config.InternalService.RBAC.Host),
			clientConfig.WithSecretKey(config.Config.InternalService.RBAC.SecretKey))

		rbac := NewRBACMiddleware(client)
		user, err := rbac.CheckPermission(token.(string), permissions)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Response{
				Status:  constantError.Error,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		if !user.Allowed {
			c.JSON(http.StatusForbidden, response.Response{
				Status:  constantError.Error,
				Message: constantError.ErrForbidden.Error(),
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
