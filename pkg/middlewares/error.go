package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GinErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Jika ada error
		err := c.Errors.Last()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
			return
		}
	}
}
