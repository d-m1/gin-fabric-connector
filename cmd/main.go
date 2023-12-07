package main

import (
	"fmt"
	"gin-fabric-connector/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	authMiddleware := api.AuthMiddleware()
	api.RegisterRoutes(r, authMiddleware)

	fmt.Println("Server running on :8080")
	r.Run(":8080")
}
