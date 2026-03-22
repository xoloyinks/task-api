package main

import (
	"context"
	"log"
	"net/http"
	"task-tracker-api/config"
	"task-tracker-api/handlers"
	"task-tracker-api/repository"
	"task-tracker-api/routes"
	"task-tracker-api/services"
)

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
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, r))

}
