package cli

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
)

// Flags
var (
	Url     string
	User    string
	RootCmd = &cobra.Command{
		Use:   "library",
		Short: "Library cli tool",
		Long:  "Library cli tool",
	}
)

func init() {
	// Gets user informatuon to make the default config path
	u, err := user.Current()
	if err != nil {
		fmt.Println("Could not determine user information")
		os.Exit(-2)
	}

	RootCmd.PersistentFlags().StringVar(&Url, "url", "http://0.0.0.0:8080/api", "API host url")
	RootCmd.PersistentFlags().StringVar(&User, "user", u.Username, "Username to use with API")
	RootCmd.AddCommand(AddCmd)
	RootCmd.AddCommand(ReadCmd)
	RootCmd.AddCommand(ShowCmd)
}
