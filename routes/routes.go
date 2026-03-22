package routes

import (
	"net/http"
	"task-tracker-api/handlers"
	"task-tracker-api/middleware"
	"task-tracker-api/utils"
)

func SetupRoutes(taskHandler *handlers.TaskHandler, authHandler *handlers.AuthHandler) http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("POST /tasks", utils.Make(middleware.AuthMiddleware(taskHandler.CreateTask)))
	r.HandleFunc("GET /tasks/{userId}", utils.Make(middleware.AuthMiddleware(taskHandler.GetAllTasks)))
	r.HandleFunc("GET /task/{id}", utils.Make(middleware.AuthMiddleware(taskHandler.GetTask)))
	r.HandleFunc("PATCH /tasks/{id}", utils.Make(middleware.AuthMiddleware(taskHandler.UpdateTask)))
	r.HandleFunc("PATCH /tasks/{id}/complete", utils.Make(middleware.AuthMiddleware(taskHandler.CompleteTask)))
	r.HandleFunc("DELETE /tasks/{id}", utils.Make(middleware.AuthMiddleware(taskHandler.DeleteTask)))

	r.HandleFunc("POST /createAccount", utils.Make(authHandler.CreateAccount))
	r.HandleFunc("POST /login", utils.Make(authHandler.Login))

	return r
}
