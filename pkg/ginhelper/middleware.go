package ginhelper

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	myerrors "github.com/homily707/ohmygo/pkg/errors"
	"github.com/sirupsen/logrus"
)

var commonLogger *logrus.Logger = logrus.StandardLogger()

func SetCommonLogger(logger *logrus.Logger) {
	commonLogger = logger
}

func ErrorHandler(c *gin.Context) {
	c.Next()
	code := http.StatusInternalServerError
	var errWithCode *myerrors.ErrorWithCode
	if len(c.Errors) > 0 {
		logger := GetLogger(c)
		// there better be only one error
		for _, err := range c.Errors {
			if errors.As(err, &errWithCode) {
				code = errWithCode.StatusCode
				logger.Infof("error with code %d", code)
			}
			logger.Error(err)
		}
		// print last non-500 error
		//TODO: json empty and code not working
		if code < 500 {
			c.JSON(code, errWithCode)
			return
		}
		// hide internal error message if code >= 500
		c.AbortWithStatus(code)
		return
	}
}

func LoggerSetter(c *gin.Context) {
	requestId := c.Request.Header.Get("Request-Id")
	if requestId == "" {
		requestId = c.Request.Header.Get("X-Request-ID")
	}
	if requestId == "" {
		requestId = fmt.Sprintf("gen-%d", time.Now().UnixNano())
	}
	logger := commonLogger.WithField("request_id", requestId)
	c.Set("LOGGER", logger)
	c.Next()
}

func GetLogger(c *gin.Context) *logrus.Entry {
	l, ok := c.MustGet("LOGGER").(*logrus.Entry)
	if !ok {
		l = commonLogger.WithField("request_id", "unknown")
	}
	return l
}
