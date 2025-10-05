// Package http
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Jrdc12/habit-tracker-go-api/internal/user"
	"github.com/Jrdc12/habit-tracker-go-api/internal/validation"
)

type userService interface {
	GetUser(id int) (user.User, error)
	CreateUser(r user.CreateRequest) (user.User, error)
	DeleteUser(id int) error
}

func NewRouter(svc userService) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /user", createUserHandler(svc))
	mux.HandleFunc("GET /user/{id}", getUserHandler(svc))
	mux.HandleFunc("DELETE /user/{id}", deleteUserHandler(svc))

	return mux
}

func createUserHandler(svc userService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req user.CreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := validation.ValidateCreateUser(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u, err := svc.CreateUser(req)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrEmailExists):
				http.Error(w, "email already exists", http.StatusConflict)
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/user/"+strconv.Itoa(u.ID))
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(u)
	}
}

func getUserHandler(svc userService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		u, err := svc.GetUser(id)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrNotFound):
				http.Error(w, "User not found", http.StatusNotFound)
			default:
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(u)
	}
}

func deleteUserHandler(svc userService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if err := svc.DeleteUser(id); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
