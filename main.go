package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sona-123/splitwise_clone/api"
	"github.com/sona-123/splitwise_clone/business"
	"github.com/sona-123/splitwise_clone/infra"
	"github.com/sona-123/splitwise_clone/repository"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	db := infra.InitDB()

	repo := &repository.Repo{DB: db}
	svc := &business.Service{Repo: repo}
	h := &api.Handler{Service: svc}

	// Initialize Gin router
	r := gin.Default()

	// Group routes
	v1 := r.Group("/api")
	{
		v1.POST("/users", h.UserHandler)
		v1.POST("/expenses", h.ExpenseHandler)
		v1.POST("/groups", h.CreateGroupHandler)
		v1.GET("/groups/:id/balances", h.BalancesHandler)
		v1.POST("/groups/:id/members", h.AddMemberHandler)
	}

	// Run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
