package response

import "github.com/gin-gonic/gin"

// SuccessResponse is a standard success response for Swagger docs
// @Description Success response
// @example {"success":true,"message":"string","data":{}}
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data"`
}

// ErrorResponse is a standard error response for Swagger docs
// @Description Error response
// @example {"success":false,"message":"error message"}
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"error message"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}
