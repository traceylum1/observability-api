package handlers

import (
	"net/http"
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/traceylum1/observability-api/internal/models"
)

var users = []models.User{
	{UserID: "1", Name: "Jimbo", Email: "jimbo@gmail.com"},
	{UserID: "2", Name: "James", Email: "james@gmail.com"},
	{UserID: "3", Name: "Marcus", Email: "marcus@gmail.com"},
}


func Hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("user found: %s", users[1].Name)))
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")

	for _, u := range users {
		if u.UserID == user_id {
			w.Write([]byte(fmt.Sprintf("user found: %s, email: %s", u.Name, u.Email)))
			return
		}
	}
	w.WriteHeader(404)
	w.Write([]byte(fmt.Sprintf("user with ID %s not found", user_id)))
}