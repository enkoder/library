package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
)

var (
	// Errors
	ErrUserCreate  = fmt.Errorf("Cannot create user")
	ErrInvalidBook = fmt.Errorf("Invalid book title")
)

type Book struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Read   bool   `json:"read"`
}

func SnakeCase(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "_", -1)
}

func PutBook(db *bolt.DB, user string, b *Book) error {
	return db.Update(func(tx *bolt.Tx) error {
		ubkt, err := tx.CreateBucketIfNotExists([]byte(user))
		if err != nil {
			return ErrUserCreate
		}

		// marshal the book into json
		buf, err := json.Marshal(b)
		if err != nil {
			return err
		}

		// put the book into the bucket with the key being the title
		return ubkt.Put([]byte(SnakeCase(b.Title)), buf)
	})
}

func GetBook(db *bolt.DB, user, title string) (*Book, error) {
	var b Book
	err := db.Update(func(tx *bolt.Tx) error {
		ubkt, err := tx.CreateBucketIfNotExists([]byte(user))
		if err != nil {
			return ErrUserCreate
		}

		v := ubkt.Get([]byte(title))
		if v == nil {
			return ErrInvalidBook
		}

		return json.Unmarshal(v, &b)
	})

	return &b, err
}

func GetBooks(db *bolt.DB, user string) ([]Book, error) {
	books := make([]Book, 0)
	err := db.Update(func(tx *bolt.Tx) error {
		ubkt, err := tx.CreateBucketIfNotExists([]byte(user))
		if err != nil {
			return ErrUserCreate
		}

		var b Book
		return ubkt.ForEach(func(k, v []byte) error {
			err := json.Unmarshal(v, &b)
			if err != nil {
				return err
			}

			books = append(books, b)
			return nil
		})
	})

	return books, err
}
func DeleteBook(db *bolt.DB, user, title string) error {
	return db.Update(func(tx *bolt.Tx) error {
		ubkt, err := tx.CreateBucketIfNotExists([]byte(user))
		if err != nil {
			return ErrUserCreate
		}

		err = ubkt.Delete([]byte(user))
		if err != nil {
			return ErrInvalidBook
		}
		return nil
	})
}
