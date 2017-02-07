package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

var ShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays books in the library",
	Long:  `Displays books in the library. Args "all", "unread" and "read" are supported`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Args format: show [all|read|unread] [by AUTHOR]")
		}

		if args[0] != "all" &&
			args[0] != "unread" &&
			args[0] != "read" {
			return fmt.Errorf("Args format: show [all|read|unread] [by AUTHOR]")
		}

		if len(args) > 1 {
			if args[1] != "by" {
				return fmt.Errorf("Args format: show [all|read|unread] [by AUTHOR]")
			}
			if len(args[2]) != 3 {
				return fmt.Errorf("Args format: show [all|read|unread] [by AUTHOR]")
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Use url package to make it easier to build the url
		url, err := url.Parse(fmt.Sprintf("%s/%s/book", Url, User))
		if err != nil {
			return err
		}

		q := url.Query()
		if args[0] == "read" {
			q.Set("read", "true")
		} else if args[0] == "unread" {
			q.Set("read", "false")
		}

		if len(args) > 1 {
			q.Set("by", args[2])
		}
		url.RawQuery = q.Encode()

		// send the get request
		resp, err := http.Get(url.String())
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
