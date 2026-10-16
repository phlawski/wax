package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"wax/internal/client"
	"wax/internal/config"
	"wax/pkg/table"
)

var (
	collectionPage    int
	collectionPerPage int
)

var collectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "Browse your Discogs collection",
}

var collectionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List releases in your collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := newClient()
		resp, err := c.GetCollection(collectionPage, collectionPerPage)
		if err != nil {
			return fmt.Errorf("failed to fetch collection: %w", err)
		}

		table.PrintMeta("Page %d of %d — %d total releases", resp.Pagination.Page, resp.Pagination.Pages, resp.Pagination.Items)

		t := table.New(collectionColumns())
		rows := make([][]string, len(resp.Releases))
		for i, r := range resp.Releases {
			rows[i] = releaseRow(r.BasicInfo)
		}
		t.Render(rows)
		return nil
	},
}

var collectionSearchCmd = &cobra.Command{
	Use:   "search QUERY",
	Short: "Search releases in your collection",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		c := newClient()

		resp, err := c.SearchCollection(query, collectionPage, collectionPerPage, config.AppConfig.Discogs.Username)
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

func newClient() *client.Client {
	cfg := config.AppConfig.Discogs
	if cfg.Token == "" || cfg.Token == "<your_discogs_personal_access_token>" {
		log.Fatal("Discogs token not set. Edit ~/.config/wax/config.yaml")
	}
	if cfg.Username == "" || cfg.Username == "<your_discogs_username>" {
		log.Fatal("Discogs username not set. Edit ~/.config/wax/config.yaml")
	}
	return client.New(cfg.Token, cfg.Username)
}

func collectionColumns() []table.Column {
	return []table.Column{
		{Title: "Artist", Width: 32},
		{Title: "Title", Width: 56},
		{Title: "Year", Width: 6},
		{Title: "Format", Width: 14},
		{Title: "Label", Width: 38},
	}
}

func searchColumns() []table.Column {
	return []table.Column{
		{Title: "ID", Width: 10},
		{Title: "Title", Width: 50},
		{Title: "Year", Width: 6},
		{Title: "Format", Width: 20},
		{Title: "Label", Width: 24},
	}
}

func releaseRow(r client.ReleaseBasicInfo) []string {
	year := fmt.Sprintf("%d", r.Year)
	if r.Year == 0 {
		year = "-"
	}
	return []string{
		artistNames(r.Artists),
		r.Title,
		year,
		formatNames(r.Formats),
		labelNames(r.Labels),
	}
}

func searchRow(r client.SearchResult) []string {
	return []string{
		fmt.Sprintf("%d", r.ID),
		r.Title,
		r.Year,
		strings.Join(r.Format, ", "),
		strings.Join(r.Label, ", "),
	}
}

func artistNames(artists []client.Artist) string {
	names := make([]string, len(artists))
	for i, a := range artists {
		names[i] = a.Name
	}
	return strings.Join(names, ", ")
}

func labelNames(labels []client.Label) string {
	seen := make(map[string]bool)
	var names []string
	for _, l := range labels {
		if !seen[l.Name] {
			seen[l.Name] = true
			names = append(names, l.Name)
		}
	}
	return strings.Join(names, ", ")
}

func formatNames(formats []client.Format) string {
	names := make([]string, len(formats))
	for i, f := range formats {
		names[i] = f.Name
	}
	return strings.Join(names, ", ")
}

func init() {
	collectionListCmd.Flags().IntVar(&collectionPage, "page", 1, "Page number")
	collectionListCmd.Flags().IntVar(&collectionPerPage, "per-page", 50, "Results per page")
	collectionSearchCmd.Flags().IntVar(&collectionPage, "page", 1, "Page number")
	collectionSearchCmd.Flags().IntVar(&collectionPerPage, "per-page", 25, "Results per page")

	collectionCmd.AddCommand(collectionListCmd, collectionSearchCmd)
}
