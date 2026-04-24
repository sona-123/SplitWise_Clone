package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sona-123/splitwise_clone/api"
	"github.com/sona-123/splitwise_clone/business"
	"github.com/sona-123/splitwise_clone/infra"
	"github.com/sona-123/splitwise_clone/repository"
)

func main() {
	godotenv.Load()
	db := infra.InitDB()

	repo := &repository.Repo{DB: db}
	svc := &business.Service{Repo: repo}
	h := &api.Handler{Service: svc}

	//Initialize Gin router
	r := gin.Default()

	//Group routes for better organization
	v1 := r.Group("/api")
	{
		v1.POST("/users", h.UserHandler)
		v1.POST("/expenses", h.ExpenseHandler)
	}
	r.Run(":8080") // Defaults to listening on 0.0.0.0:8080
}
