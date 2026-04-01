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
	"task-tracker-api/sse"
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

	hub := sse.NewHub()

	client := config.ConnectMongo(cfg)
	defer client.Disconnect(context.TODO())
	db := client.Database(cfg.MongoDBName)

	taskCollection := db.Collection(cfg.MongoCollection)
	userCollection := db.Collection(cfg.UserCollection)
	boardCollection := db.Collection(cfg.BoardCollection)
	columnCollection := db.Collection(cfg.ColumnCollection)
	teamCollection := db.Collection(cfg.TeamCollection)
	teamMemberCollection := db.Collection(cfg.TeamMemberCollection)

	if err := config.CreateIndexes(db); err != nil {
		log.Fatal(err)
	}

	sseHandler := sse.NewSSEHandler(hub)

	authRepo := repository.NewAuthRepository(taskCollection, userCollection)
	authService := services.NewAuthServices(authRepo)
	authHandler := handlers.NewAuthHandler(authService)

	taskRepo := repository.NewTaskRepository(taskCollection, columnCollection, boardCollection)
	taskService := services.NewTaskServices(taskRepo)
	taskHandler := handlers.NewTaskHandler(taskService, hub)

	boardRepo := repository.NewBoardReposity(boardCollection, columnCollection, taskCollection)
	boardService := services.NewBoardServices(boardRepo)
	boardHandler := handlers.NewBoardHandler(boardService)

	teamRepo := repository.NewTeamRepository(teamCollection, userCollection, teamMemberCollection, taskCollection, boardCollection, columnCollection)
	teamService := services.NewTeamServices(teamRepo)
	teamHandler := handlers.NewTeamHandler(teamService)

	r := routes.SetupRoutes(taskHandler, authHandler, teamHandler, boardHandler, sseHandler)

	loggedRouter := middleware.LoggerMiddleware(middleware.RateLimiterMiddleware(r))
	log.Fatal(http.ListenAndServe(":"+cfg.AppPort, loggedRouter))

}
