package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ArkjuniorK/go_micorservices/handlers"
	"github.com/gorilla/mux"
)

func main() {
	var err error

	// create new mux server and register the handlers
	mux := mux.NewRouter()

	// create the handler
	logger := log.New(os.Stdout, "user-api", log.LstdFlags)
	uh := handlers.UserPath(logger)

	// method GET
	get := mux.Methods(http.MethodGet).Subrouter()
	get.HandleFunc("/", uh.ListUsers)

	// method POST
	post := mux.Methods(http.MethodPost).Subrouter()
	post.HandleFunc("/create", uh.CreateUser)
	post.Use(uh.MwUserValidation)

	// method PUT
	put := mux.Methods(http.MethodPut).Subrouter()
	put.HandleFunc("/update/{_id}", uh.UpdateUser)
	put.Use(uh.MwUserValidation)

	// create server
	server := &http.Server{
		Addr:         ":2021",           // port
		Handler:      mux,               // handler to use
		ErrorLog:     logger,            // set error logger for server
		IdleTimeout:  120 * time.Second, // max time for connections
		ReadTimeout:  1 * time.Second,   // max time to read request
		WriteTimeout: 1 * time.Second,   // max time to write response
	}

	// listen and server on port 2021
	go func() {
		logger.Print("Start at port :2021")
		err = server.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	sig := <-channel
	log.Print("Process terminated, gracefully shutdown", sig)

	// set timeout before shutdown server
	// so it could resolve remaining request
	timeoutCtx, cancelCtx := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelCtx()
	server.Shutdown(timeoutCtx)
}
