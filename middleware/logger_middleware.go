package middleware

import (
	"time"

	"github.com/adityatresnobudi/job-portal/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

func Logger(l logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		finish := time.Now()
		latency := finish.Sub(start)
		method := c.Request.Method
		uri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		ip := c.ClientIP()

		param := map[string]interface{}{
			"id":          requestid.Get(c),
			"method":      method,
			"latency":     latency,
			"uri":         uri,
			"status_code": statusCode,
			"ip":          ip,
		}

		log := logger.NewLogger()

		if len(c.Errors) == 0 {
			log.Info(param)
		} else {
			errList := []error{}
			for _, err := range c.Errors {
				errList = append(errList, err)
			}

			if len(errList) > 0 {
				param["errors"] = errList
				log.Errorf("%v", param)
			}
		}

		c.Next()
	}
}
