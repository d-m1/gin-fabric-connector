package main

import (
	"fmt"
	"gin-fabric-connector/api"
	"gin-fabric-connector/blockchain"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	client, err := blockchain.GetClient()
	if err != nil {
		fmt.Println("Error initializing fabric client:", err.Error())
		panic(err)
	}

	authMiddleware := api.AuthMiddleware()
	api.RegisterRoutes(r, authMiddleware, *client)

	fmt.Println("Server running on :8080")
	r.Run(":8080")
}
