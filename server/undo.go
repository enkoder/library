package server

import (
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

// api/${user}/undo
func UndoHandler(db *bolt.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		}

		// Gets the id from the url
		vars := mux.Vars(r)
		user := vars["user"]

		// execute the undo on the db
		text, err := ExecuteUndo(db, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Sets header, code, and marshal response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(text))
	}
}
