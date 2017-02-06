package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/enkoder/library/server"
	"github.com/gorilla/mux"
)

var (
	Port   uint
	Host   string
	DBPath string
)

func main() {
	// collect user inputs
	flag.UintVar(&Port, "port", 8080, "Server port")
	flag.StringVar(&Host, "host", "0.0.0.0", "Server Host")
	flag.Parse()

	// Set up routers
	r := mux.NewRouter()
	r.HandleFunc("/api/${user:[a-zA-Z]}", server.UserHandler)
	r.HandleFunc("/api/${user:[a-zA-Z]}/book/${title}", server.BookHandler)
	r.HandleFunc("/api/${user:[a-zA-Z]}/book", server.BooksHandler)
	r.HandleFunc("/api/${user:[a-zA-Z]}/undo", server.UndoHandler)

	// Listen forever
	bind := fmt.Sprintf("%s:%d", Host, Port)
	log.Printf("Listening: %s", bind)
	log.Fatal(http.ListenAndServe(bind, r))
}
