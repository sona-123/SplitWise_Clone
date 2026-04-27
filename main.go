package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sona-123/splitwise_clone/api/handlers"
	"github.com/sona-123/splitwise_clone/business"
	_ "github.com/sona-123/splitwise_clone/docs"
	"github.com/sona-123/splitwise_clone/infra"
	"github.com/sona-123/splitwise_clone/middleware"
	"github.com/sona-123/splitwise_clone/repository"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	h := &handlers.Handler{Service: svc}

	// Initialize Gin router
	r := gin.Default()

	//Public routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/users", h.UserHandler)  //SignUp
	r.POST("/api/login", h.LoginHandler) //Login
	// Protected routes
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.POST("/expenses", h.ExpenseHandler)
		authorized.POST("/groups", h.CreateGroupHandler)
		authorized.GET("/groups/:id/balances", h.BalancesHandler)
		authorized.POST("/groups/:id/members", h.AddMemberHandler)
		authorized.GET("/user/summary", h.UserSummaryHandler)
	}

	// Run server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
