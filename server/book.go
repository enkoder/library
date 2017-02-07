package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

// api/${user}/book?read=true&by=kodie
func BooksHandler(db *bolt.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Gets the id from the url
		vars := mux.Vars(r)
		user := vars["user"]

		// extract the query params
		var isread *bool
		var author *string
		q := r.URL.Query()
		a := q.Get("by")
		read := q.Get("read")
		if read != "" {
			b, err := strconv.ParseBool(q.Get("read"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			isread = &b
		}
		if a != "" {
			author = &a
		}

		switch r.Method {
		case "GET":
			books, err := GetBooks(db, user, isread, author)
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

			err = StageUndoDelete(db, user, &b)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to stage undo: %v", err), http.StatusInternalServerError)
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

		case "POST":
			pbody := struct {
				Read *bool `json:"read"`
			}{}

			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&pbody)
			if err != nil {
				http.Error(w, "Post body not formatted properly", http.StatusBadRequest)
				return
			}

			// Ensure they send us the read field
			if pbody.Read == nil {
				http.Error(w, "Post body didn't contain read json field", http.StatusBadRequest)
				return
			}

			// Set the read state of the book
			book.Read = *pbody.Read

			// insert the new book data into the db
			err = PutBook(db, user, book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = StageUndoUnread(db, user, book)
			if err != nil {
				http.Error(w, fmt.Sprintf("Unable to stage undo: %v", err), http.StatusInternalServerError)
				return
			}

			// Sets header, code, and marshal response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

		default:
			http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		}
	}
}
