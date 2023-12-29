package api

import (
	"gin-fabric-connector/blockchain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes defines the REST API routes protecting them with authMiddleware
func RegisterRoutes(r *gin.Engine, authMiddleware gin.HandlerFunc, client blockchain.FabricClient) {

	r.GET("/transaction", authMiddleware, func(c *gin.Context) {
		var transactionDto blockchain.Transaction

		err := c.ShouldBindJSON(&transactionDto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
			return
		}

		response := client.Evaluate(transactionDto)
		if response.Status == "KO" {
			c.JSON(http.StatusInternalServerError, gin.H{"data": response})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	})

	r.POST("/transaction", authMiddleware, func(c *gin.Context) {
		var transactionDto blockchain.Transaction

		err := c.ShouldBindJSON(&transactionDto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
			return
		}

		response := client.Submit(transactionDto)
		if response.Status == "KO" {
			c.JSON(http.StatusInternalServerError, gin.H{"data": response})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	})
}
