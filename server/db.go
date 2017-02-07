package server

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
)

const (
	UndoBucketName = "undos"
	UndoTypeUnread = "unread"
	UndoTypeDelete = "delete"
)

var (
	// Errors
	ErrUserCreate  = fmt.Errorf("Cannot create user")
	ErrInvalidBook = fmt.Errorf("Invalid book title")
	ErrInvalidUser = fmt.Errorf("Invalid user")
	ErrUndoCreate  = fmt.Errorf("Cannot create undo bucket")
)

type Book struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Read   bool   `json:"read"`
}

type Undo struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Title string `json:"title"`
}

func SnakeCase(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "_", -1)
}

func StageUndoDelete(db *bolt.DB, user string, b *Book) error {
	return db.Update(func(tx *bolt.Tx) error {
		ubkt, err := tx.CreateBucketIfNotExists([]byte(UndoBucketName))
		if err != nil {
			return ErrUndoCreate
		}

		// Create and stage the undo for the user
		undo := Undo{
			Type:  UndoTypeDelete,
			Text:  fmt.Sprintf("Removed \"%s\" by \"%s\"", b.Title, b.Author),
			Title: b.Title,
		}

		// marshal the undo into json
		buf, err := json.Marshal(undo)
		if err != nil {
			return err
		}

		// put the book into the bucket with the key being the title
		return ubkt.Put([]byte(user), buf)
	})
}

func StageUndoUnread(db *bolt.DB, user string, b *Book) error {
	return db.Update(func(tx *bolt.Tx) error {
		ubkt, err := tx.CreateBucketIfNotExists([]byte(UndoBucketName))
		if err != nil {
			return ErrUndoCreate
		}

		// Create and stage the undo for the user
		undo := Undo{
			Type:  UndoTypeUnread,
			Text:  fmt.Sprintf("\"%s\" by \"%s\" was marked as unread", b.Title, b.Author),
			Title: b.Title,
		}

		// marshal the undo into json
		buf, err := json.Marshal(undo)
		if err != nil {
			return err
		}

		// put the book into the bucket with the key being the title
		return ubkt.Put([]byte(user), buf)
	})
}

func ExecuteUndo(db *bolt.DB, user string) (string, error) {
	var text string
	err := db.Update(func(tx *bolt.Tx) error {
		ubkt, err := tx.CreateBucketIfNotExists([]byte(UndoBucketName))
		if err != nil {
			return ErrUndoCreate
		}

		v := ubkt.Get([]byte(user))
		if v == nil {
			return ErrInvalidUser
		}

		// get the latest undo for the user
		var undo Undo
		err = json.Unmarshal(v, &undo)
		if err != nil {
			return err
		}

		// Be sure to set the value of the text to the caller
		text = undo.Text

		// Execute the two different types of undos
		switch undo.Type {
		case UndoTypeDelete:
			err := deleteBook(tx, user, undo.Title)
			if err != nil {
				return err
			}

		case UndoTypeUnread:
			// Get the book which should now be in a read state
			b, err := getBook(tx, user, undo.Title)
			if err != nil {
				return err
			}

			// Update the books
			b.Read = false

			// Jam it back in the db
			return putBook(tx, user, b)

		default:
			return fmt.Errorf("Unknows undo type: %s", undo.Type)
		}

		// deletes the undo since its executed
		return ubkt.Delete([]byte(user))
	})

	return text, err
}

func putBook(tx *bolt.Tx, user string, b *Book) error {
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
}

func PutBook(db *bolt.DB, user string, b *Book) error {
	return db.Update(func(tx *bolt.Tx) error {
		return putBook(tx, user, b)
	})
}

func getBook(tx *bolt.Tx, user, title string) (*Book, error) {
	var book Book
	ubkt, err := tx.CreateBucketIfNotExists([]byte(user))
	if err != nil {
		return nil, ErrUserCreate
	}

	v := ubkt.Get([]byte(title))
	if v == nil {
		return nil, ErrInvalidBook
	}

	err = json.Unmarshal(v, &book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func GetBook(db *bolt.DB, user, title string) (*Book, error) {
	var book *Book
	var err error
	err = db.Update(func(tx *bolt.Tx) error {
		book, err = getBook(tx, user, title)
		if err != nil {
			return err
		}
		return nil
	})

	return book, err
}

func GetBooks(db *bolt.DB, user string, isread *bool) ([]Book, error) {
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

			// Used to filter books by read or unread
			if isread == nil || b.Read == *isread {
				books = append(books, b)
			}
			return nil
		})
	})

	return books, err
}

func deleteBook(tx *bolt.Tx, user, title string) error {
	ubkt, err := tx.CreateBucketIfNotExists([]byte(user))
	if err != nil {
		return ErrUserCreate
	}

	err = ubkt.Delete([]byte(user))
	if err != nil {
		return ErrInvalidBook
	}
	return nil
}

func DeleteBook(db *bolt.DB, user, title string) error {
	return db.Update(func(tx *bolt.Tx) error {
		return deleteBook(tx, user, title)
	})
}
