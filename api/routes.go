package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, authMiddleware gin.HandlerFunc) {
	r.GET("/public", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Public endpoint"})
	})

	r.GET("/secured", authMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Secured endpoint"})
	})
}
