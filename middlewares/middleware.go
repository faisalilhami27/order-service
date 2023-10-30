package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"order-service/utils/response"

	"net/http"
)

func HandlePanic(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.SetLevel(log.ErrorLevel)
			log.Errorf("error occured: %v", r)
			c.JSON(http.StatusBadRequest, response.ResponsePanicError(r))
			return
		}
	}()
	c.Next()
}
