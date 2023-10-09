package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"net/http"

	"order-service/utils"
)

func HandlePanic(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.SetLevel(log.ErrorLevel)
			log.Errorf("error occured: %v", r)
			c.JSON(http.StatusBadRequest, utils.ResponseError(r.(error))) //nolint:forcetypeassert
			return
		}
	}()
	c.Next()
}
