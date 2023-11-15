package middleware

import (
	"net/http"

	"github.com/adityatresnobudi/job-portal/shared"
	"github.com/gin-gonic/gin"
)

func GlobalErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			switch e := err.Err.(type) {
			case *shared.CustomError:
				c.AbortWithStatusJSON(e.StatusCode, e.ToErrorDTO())
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": err.Error(),
				})
			}
			c.Abort()
		}
	}
}
