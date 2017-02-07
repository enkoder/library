package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var ReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Marks a book as being read",
	Long:  `Marks a book as being read`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Must add a title argument")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// normalize the title client side
		title := SnakeCase(args[0])

		// Hand writing json like a boss
		body := []byte("{\"read\":true}")

		// send the post request
		endpoint := fmt.Sprintf("%s/%s/book/%s", Url, User, title)
		resp, err := http.Post(endpoint, "application/json", bytes.NewReader(body))
		if err != nil {
			return err
		}

		// ensure the return code is expected
		if resp.StatusCode != http.StatusOK {
			// read body from the response to get at bytes
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			return fmt.Errorf("%s (%d)\n", strings.Trim(string(body), "\n"), resp.StatusCode)
		}

		// send the get request
		resp, err = http.Get(endpoint)
		if err != nil {
			return err
		}

		// read body from the response to get at bytes
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// unmarshal the book
		var book Book
		err = json.Unmarshal(body, &book)
		if err != nil {
			return err
		}

		fmt.Printf("\"%s\" by \"%s\" marked as read\n", book.Title, book.Author)
		return nil
	},
}
