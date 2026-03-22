package routes

import (
	"net/http"
	"task-tracker-api/handlers"
	"task-tracker-api/utils"
)

func SetupRoutes(taskHandler *handlers.TaskHandler, authHandler *handlers.AuthHandler) http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("POST /tasks", taskHandler.CreateTask)
	r.HandleFunc("GET /tasks", taskHandler.GetAllTasks)
	r.HandleFunc("GET /task/{id}", taskHandler.GetTask)
	r.HandleFunc("GET /view/{id}", taskHandler.CheckId)
	r.HandleFunc("GET /", taskHandler.HomeRoute)
	r.HandleFunc("PATCH /task/{id}", taskHandler.UpdateTask)
	r.HandleFunc("PATCH /status/{id}", taskHandler.CompleteTask)
	r.HandleFunc("DELETE /task/{id}", taskHandler.DeleteTask)
	r.HandleFunc("POST /createAccount", authHandler.CreateAccount)
	r.HandleFunc("POST /login", utils.Make(authHandler.Login))

	return r
}
