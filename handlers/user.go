// Package classification of User API
//
// Documentation for user API
//
//  Schemes: http
//  BasePath: /user
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ArkjuniorK/go_micorservices/users"
	"github.com/gorilla/mux"
)

type Users struct {
	logger *log.Logger
}

func UserPath(logger *log.Logger) *Users {
	return &Users{logger}
}

// key for context
// it would be used as key when
// getting context from request
type KeyUser struct{}

// middleware for PUT and POST
func (u *Users) MwUserValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// create new User
		user := users.User{}

		// serialize d from JSON to Struct
		err := user.FromJSON(r.Body)
		if err != nil {
			u.logger.Println("[ERROR] deserializing user")
			http.Error(w, "Unable to deserialize user", http.StatusBadRequest)
			return
		}

		// validate user property that
		// defined inside `validate` tag
		err = user.Validate()
		if err != nil {
			u.logger.Println("[ERROR] validating user")
			http.Error(w, fmt.Sprintf("Unable to validate user %s", err), http.StatusBadRequest)
			return

		}

		// create new context for request
		// that hold user data so any
		// handler that using the middleware
		// could get the user data in struct type
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// A list of users that return in the response
// swagger:response usersResponse
type usersResponse struct {
	// All users in the system
	// in: body
	Body []users.User
}

// swagger:route GET /users users ListUsers
// Returns a list of users
// responses:
//  200: usersResponse

// GetUsers returns the users from data store
func (u *Users) ListUsers(w http.ResponseWriter, r *http.Request) {
	u.logger.Print("Handle GET Users")

	lu := users.ListUsers() // list users

	// convert to json
	err := lu.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to list user", http.StatusInternalServerError)
		return
	}
}

func (u *Users) CreateUser(w http.ResponseWriter, r *http.Request) {
	u.logger.Print("Handle POST Users")

	// create new User
	// get the user context from
	// *Request.Context()
	user := r.Context().Value(KeyUser{}).(users.User)

	// create new user using user variable
	usr := users.CreateUser(&user)

	// encode usr to JSON and send
	// it to client
	err := usr.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshall user", http.StatusInternalServerError)
		return
	}
}

func (u Users) UpdateUser(w http.ResponseWriter, r *http.Request) {
	u.logger.Print("Handle POST Users")

	// get the _id from URL
	// then convert to xid.ID
	vars := mux.Vars(r)
	_id := vars["_id"]

	// create new User
	user := r.Context().Value(KeyUser{}).(users.User)
	usr, err := users.UpdateUser(_id, &user)
	if err == users.ErrUserNotFound {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Unbale to update user", http.StatusInternalServerError)
		return
	}

	err = usr.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal user", http.StatusInternalServerError)
		return
	}

}
