package handlers

import (
	"log"
	"net/http"
)

type Greeting struct {
	logger *log.Logger
}

func NewGreeting(logger *log.Logger) *Greeting {
	return &Greeting{logger}
}

func (g *Greeting) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	g.logger.Print("Success")
	wr.WriteHeader(http.StatusOK)
	wr.Write([]byte("Hello World!"))
}
