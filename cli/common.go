package cli

import (
	"strings"
)

type Book struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Read   bool   `json:"read,omitempty"`
}

func SnakeCase(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "_", -1)
}
