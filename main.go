package main

import (
	"context"
	"log"
	"net/http"
	"task-tracker-api/config"
	_ "task-tracker-api/docs"
	"task-tracker-api/handlers"
	"task-tracker-api/middleware"
	"task-tracker-api/repository"
	"task-tracker-api/routes"
	"task-tracker-api/services"
)

// main.go
// @title           Task Tracker API
// @version         1.0
// @description     A simple task tracker API built with Go and MongoDB
// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := config.Load()

	client := config.ConnectMongo(cfg)
	defer client.Disconnect(context.TODO())

	collection := client.Database(cfg.MongoDBName).Collection(cfg.MongoCollection)
	userCollection := client.Database(cfg.MongoDBName).Collection(cfg.UserCollection)

	authRepo := repository.NewAuthRepository(collection, userCollection)
	authService := services.NewAuthServices(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	taskRepo := repository.NewTaskRepository(collection, userCollection)
	taskService := services.NewTaskServices(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService)

	r := routes.SetupRoutes(taskHandler, authHandler)

	loggedRouter := middleware.LoggerMiddleware(middleware.RateLimiterMiddleware(r))
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, loggedRouter))

}
