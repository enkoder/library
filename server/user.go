package server

import (
	"net/http"

	"github.com/boltdb/bolt"
)

func UserHandler(db *bolt.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
