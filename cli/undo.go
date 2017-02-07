package cli

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var UndoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undoes the most recent command",
	Long:  `Undoes the most recent command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// send the post request
		endpoint := fmt.Sprintf("%s/%s/undo", Url, User)
		resp, err := http.Post(endpoint, "", nil)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// read body from the response to get at bytes
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// ensure the return code is expected
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%s (%d)\n", strings.Trim(string(body), "\n"), resp.StatusCode)
		}

		fmt.Println(string(body))
		return nil
	},
}
