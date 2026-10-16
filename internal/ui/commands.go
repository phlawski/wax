package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"wax/internal/client"
)

// fetchCollection loads one page of the user's collection, optionally filtered by query.
func fetchCollection(c *client.Client, query string, page int) tea.Cmd {
	return func() tea.Msg {
		resp, err := c.GetCollectionSearch(query, page, 50)
		if err != nil {
			return fetchMsg{err: err}
		}
		return fetchMsg{
			releases: resp.Releases,
			total:    resp.Pagination.Items,
			pages:    resp.Pagination.Pages,
			query:    query,
		}
	}
}

// searchDB searches the Discogs database with structured params.
func searchDB(c *client.Client, p client.SearchParams, page int) tea.Cmd {
	return func() tea.Msg {
		resp, err := c.Search(p, page, 25)
		if err != nil {
			return searchResultMsg{err: err}
		}
		return searchResultMsg{
			results: resp.Results,
			total:   resp.Pagination.Items,
			pages:   resp.Pagination.Pages,
		}
	}
}

// fetchMarketplace loads marketplace listings for a release ID.
func fetchMarketplace(c *client.Client, releaseID, page int) tea.Cmd {
	return func() tea.Msg {
		resp, err := c.GetMarketplaceListings(releaseID, page, 25)
		if err != nil {
			return marketplaceMsg{err: err, releaseID: releaseID}
		}
		return marketplaceMsg{
			listings:  resp.Listings,
			total:     resp.Pagination.Items,
			pages:     resp.Pagination.Pages,
			releaseID: releaseID,
		}
	}
}

// fetchRelease loads full release details by ID.
func fetchRelease(c *client.Client, id int) tea.Cmd {
	return func() tea.Msg {
		r, err := c.GetRelease(id)
		return releaseDetailMsg{detail: r, err: err}
	}
}
