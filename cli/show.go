package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays books in the library",
	Long:  `Displays books in the library. Args "all", "unread" and "read" are supported`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Must add an argument of 'all', 'unread', or 'reag'")
		}

		if args[0] != "all" &&
			args[0] != "unread" &&
			args[0] != "read" {
			return fmt.Errorf("Must add an argument of 'all', 'unread', or 'reag'")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		var endpoint string
		if args[0] == "all" {
			endpoint = fmt.Sprintf("%s/%s/book", Url, User)
		} else if args[0] == "read" {
			endpoint = fmt.Sprintf("%s/%s/book?read=true", Url, User)
		} else if args[0] == "unread" {
			endpoint = fmt.Sprintf("%s/%s/book?read=false", Url, User)
		}

		// send the get request
		resp, err := http.Get(endpoint)
		if err != nil {
			return err
		}

		// read body from the response to get at bytes
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// ensure the return code is expected
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%s (%d)\n", string(body), resp.StatusCode)
		}

		var books []Book
		err = json.Unmarshal(body, &books)
		if err != nil {
			return err
		}

		for _, book := range books {
			var read string
			if book.Read {
				read = "read"
			} else {
				read = "unread"
			}
			fmt.Printf("\"%s\" by \"%s\" (%s)\n", book.Title, book.Author, read)
		}

		return nil
	},
}
