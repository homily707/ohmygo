package ginhelper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	myerrors "github.com/homily707/ohmygo/pkg/errors"
	"github.com/sirupsen/logrus"
)

var fallbackLogger logrus.FieldLogger = logrus.NewEntry(logrus.StandardLogger())

var getLogger func(c *gin.Context) logrus.FieldLogger

func SetGetLogger(f func(c *gin.Context) logrus.FieldLogger) {
	getLogger = f
}

func ErrorHandler(c *gin.Context) {
	logger := getLogger(c)
	if logger == nil {
		logger = fallbackLogger
	}
	c.Next()
	code := http.StatusInternalServerError
	var errWithCode *myerrors.ErrorWithCode
	if len(c.Errors) > 0 {
		// there better be only one error
		for _, e := range c.Errors {
			if httpError, ok := e.Err.(myerrors.HttpError); ok {
				code = httpError.StatusCode()
				logger.Infof("error with code %d", code)
			}
			logger.Error(e)
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
	logger := fallbackLogger.WithField("request_id", requestId)
	c.Set("LOGGER", logger)
	c.Next()
}

func GetLogger(c *gin.Context) *logrus.Entry {
	l, ok := c.MustGet("LOGGER").(*logrus.Entry)
	if !ok {
		l = fallbackLogger.WithField("request_id", "unknown")
	}
	return l
}
