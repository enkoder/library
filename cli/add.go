package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a book to the library",
	Long:  `Adds a book to the library`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("Must add title and author arguments")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		book := Book{
			Title:  args[0],
			Author: args[1],
		}

		// Marshal the book into bytes
		buf, err := json.Marshal(book)
		if err != nil {
			return err
		}

		// send the post request
		endpoint := fmt.Sprintf("%s/%s/book", Url, User)
		resp, err := http.Post(endpoint, "application/json", bytes.NewReader(buf))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// ensure the return code is expected
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("couldn't add \"%s\" by \"%s\": %d\n", book.Title, book.Author, resp.StatusCode)
		}
		fmt.Printf("Added \"%s\" by \"%s\"\n", book.Title, book.Author)
		return nil
	},
}
