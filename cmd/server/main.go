package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
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
	flag.StringVar(&DBPath, "db", "library.db", "Path to the library database")
	flag.Parse()

	// open boltdb
	db, err := bolt.Open(DBPath, 0600, &bolt.Options{Timeout: time.Second * 10})
	if err != nil {
		log.Fatalf("bolt.Open(%s): %v", DBPath, err)
	}

	// Set up routers
	r := mux.NewRouter()
	r.HandleFunc("/api/${user:[a-zA-Z]}", server.UserHandler(db))
	r.HandleFunc("/api/${user:[a-zA-Z]}/book/${title}", server.BookHandler(db))
	r.HandleFunc("/api/${user:[a-zA-Z]}/book", server.BooksHandler(db))
	r.HandleFunc("/api/${user:[a-zA-Z]}/undo", server.UndoHandler(db))

	// Listen forever
	bind := fmt.Sprintf("%s:%d", Host, Port)
	log.Printf("Listening: %s", bind)
	log.Fatal(http.ListenAndServe(bind, r))
}
