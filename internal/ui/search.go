package ui

import (
	"strings"

	"wax/internal/client"
)

// parseQuery parses a search string into SearchParams.
// Supports "artist: 999", "artist:999"
func parseQuery(input string) client.SearchParams {
	input = strings.TrimSpace(input)

	lower := strings.ToLower(input)
	for _, prefix := range []string{"artist:", "artist: "} {
		if strings.HasPrefix(lower, prefix) {
			rest := strings.TrimSpace(input[len(prefix):])
			if rest != "" {
				return client.SearchParams{Artist: rest, SearchType: "release"}
			}
		}
	}

	return client.SearchParams{Query: input, SearchType: "release"}
}
