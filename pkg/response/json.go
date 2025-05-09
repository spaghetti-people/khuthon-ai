package response

import (
	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"status": "error", "message": msg})
}
