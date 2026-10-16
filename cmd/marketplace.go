package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"wax/pkg/output"
)

var (
	marketplacePage    int
	marketplacePerPage int
)

var marketplaceCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Browse the Discogs marketplace",
}

var marketplaceSearchCmd = &cobra.Command{
	Use:   "search QUERY",
	Short: "Search marketplace listings by release title or artist",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")
		c := newClient()

		// First find the release in the database, then show marketplace listings
		dbResp, err := c.Search(query, "release", 1, 5)
		if err != nil {
			return fmt.Errorf("database search failed: %w", err)
		}
		if len(dbResp.Results) == 0 {
			output.Warning("No releases found for %q", query)
			return nil
		}

		// Show top match and its listings
		top := dbResp.Results[0]
		output.Info("Top match: [%d] %s (%s)", top.ID, top.Title, top.Year)
		fmt.Println()

		listResp, err := c.GetMarketplaceListings(top.ID, marketplacePage, marketplacePerPage)
		if err != nil {
			return fmt.Errorf("marketplace fetch failed: %w", err)
		}

		output.Dim("Page %d of %d — %d listings", listResp.Pagination.Page, listResp.Pagination.Pages, listResp.Pagination.Items)
		fmt.Println()

		for _, l := range listResp.Listings {
			fmt.Printf("%-8s  %-20s  %.2f %s  @%s\n",
				l.Condition, l.SleeveCondition,
				l.Price.Value, l.Price.Currency,
				l.Seller.Username,
			)
			if l.Comments != "" {
				output.Dim("  %s", l.Comments)
			}
		}
		return nil
	},
}

var marketplaceListingsCmd = &cobra.Command{
	Use:   "listings RELEASE_ID",
	Short: "Show marketplace listings for a release ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var id int
		if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
			return fmt.Errorf("invalid release ID: %s", args[0])
		}

		c := newClient()
		resp, err := c.GetMarketplaceListings(id, marketplacePage, marketplacePerPage)
		if err != nil {
			return fmt.Errorf("marketplace fetch failed: %w", err)
		}

		output.Dim("Page %d of %d — %d listings", resp.Pagination.Page, resp.Pagination.Pages, resp.Pagination.Items)
		fmt.Println()

		for _, l := range resp.Listings {
			fmt.Printf("%-8s  %-20s  %.2f %s  @%s\n",
				l.Condition, l.SleeveCondition,
				l.Price.Value, l.Price.Currency,
				l.Seller.Username,
			)
			if l.Comments != "" {
				output.Dim("  %s", l.Comments)
			}
		}
		return nil
	},
}

func init() {
	marketplaceSearchCmd.Flags().IntVar(&marketplacePage, "page", 1, "Page number")
	marketplaceSearchCmd.Flags().IntVar(&marketplacePerPage, "per-page", 25, "Results per page")
	marketplaceListingsCmd.Flags().IntVar(&marketplacePage, "page", 1, "Page number")
	marketplaceListingsCmd.Flags().IntVar(&marketplacePerPage, "per-page", 25, "Results per page")

	marketplaceCmd.AddCommand(marketplaceSearchCmd, marketplaceListingsCmd)
}
