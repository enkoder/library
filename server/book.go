package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

// api/${user}/book
func BooksHandler(db *bolt.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// Gets the id from the url
		vars := mux.Vars(r)
		user := vars["user"]

		switch r.Method {
		case "GET":
			books, err := GetBooks(db, user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Sets header, code, and marshal response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(books)

		case "POST":
			var b Book
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&b)
			if err != nil {
				http.Error(w, "Post body not formatted properly", http.StatusBadRequest)
				return
			}

			err = PutBook(db, user, &b)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Sets header, code, and marshal response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		}
	}
}

// api/${user}/book/${title}
func BookHandler(db *bolt.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Gets the id from the url
		vars := mux.Vars(r)
		user := vars["user"]
		title := vars["title"]

		// Need to get the book in all methods in this handler
		book, err := GetBook(db, user, title)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			// Sets header, code, and marshal response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(book)

		default:
			http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		}
	}
}
