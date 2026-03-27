package routes

import (
	"net/http"
	"task-tracker-api/handlers"
	"task-tracker-api/middleware"
	"task-tracker-api/utils"

	_ "task-tracker-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(taskHandler *handlers.TaskHandler, authHandler *handlers.AuthHandler, teamHandler *handlers.TeamHandler, boardHandler *handlers.BoardHandler) http.Handler {
	r := http.NewServeMux()

	r.Handle("/swagger/", httpSwagger.WrapHandler)

	r.HandleFunc("POST /tasks", utils.Make(middleware.AuthMiddleware(taskHandler.CreateTask)))
	r.HandleFunc("GET /tasks", utils.Make(middleware.AuthMiddleware(taskHandler.GetTasks)))
	r.HandleFunc("GET /task/{id}", utils.Make(middleware.AuthMiddleware(taskHandler.GetTask)))
	r.HandleFunc("PATCH /tasks/{id}", utils.Make(middleware.AuthMiddleware(taskHandler.UpdateTask)))
	r.HandleFunc("PATCH /tasks/{id}/complete", utils.Make(middleware.AuthMiddleware(taskHandler.CompleteTask)))
	r.HandleFunc("DELETE /tasks/{id}", utils.Make(middleware.AuthMiddleware(taskHandler.DeleteTask)))

	r.HandleFunc("POST /createAccount", utils.Make(authHandler.CreateAccount))
	r.HandleFunc("POST /login", utils.Make(authHandler.Login))

	r.HandleFunc("POST /team", utils.Make(middleware.AuthMiddleware(teamHandler.CreateTeam)))
	r.HandleFunc("POST /addMember/{id}", utils.Make(middleware.AuthMiddleware(teamHandler.AddMember)))
	r.HandleFunc("GET /teams", utils.Make(middleware.AuthMiddleware(teamHandler.GetTeams)))
	r.HandleFunc("GET /teams/{id}", utils.Make(middleware.AuthMiddleware(teamHandler.GetTeam)))
	r.HandleFunc("PATCH /teams/{id}", utils.Make(middleware.AuthMiddleware(teamHandler.UpdateTeam)))
	r.HandleFunc("DELETE /teams/{id}", utils.Make(middleware.AuthMiddleware(teamHandler.DeleteTeam)))

	r.HandleFunc("POST /board/{id}", utils.Make(middleware.AuthMiddleware(boardHandler.CreateBoard)))
	r.HandleFunc("GET /boards", utils.Make(middleware.AuthMiddleware(boardHandler.GetBoardsByDestination)))
	r.HandleFunc("GET /boards/{id}", utils.Make(middleware.AuthMiddleware(boardHandler.GetBoard)))
	r.HandleFunc("PATCH /boards/{id}", utils.Make(middleware.AuthMiddleware(boardHandler.UpdateBoard)))
	r.HandleFunc("DELETE /boards/{id}", utils.Make(middleware.AuthMiddleware(boardHandler.DeleteBoard)))

	r.HandleFunc("PATCH /columns/{id}", utils.Make(middleware.AuthMiddleware(taskHandler.UpdateColumn)))

	return r

}
