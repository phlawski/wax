package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"wax/internal/client"
	"wax/pkg/table"
)

var (
	searchType    string
	searchPage    int
	searchPerPage int
)

var searchCmd = &cobra.Command{
	Use:   "search QUERY",
	Short: "Search the Discogs database",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		c := newClient()

		resp, err := c.Search(client.SearchParams{Query: query, SearchType: searchType}, searchPage, searchPerPage)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		table.PrintMeta("Page %d of %d — %d results", resp.Pagination.Page, resp.Pagination.Pages, resp.Pagination.Items)

		t := table.New(searchColumns())
		rows := make([][]string, len(resp.Results))
		for i, r := range resp.Results {
			rows[i] = searchRow(r)
		}
		t.Render(rows)
		return nil
	},
}

func init() {
	searchCmd.Flags().StringVar(&searchType, "type", "release", "Search type: release, master, artist, label")
	searchCmd.Flags().IntVar(&searchPage, "page", 1, "Page number")
	searchCmd.Flags().IntVar(&searchPerPage, "per-page", 25, "Results per page")
}
