package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update is the Bubbletea update function.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case fetchMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			m.statusErr = true
			return m, nil
		}
		m.collectionReleases = msg.releases
		m.collectionTotal = msg.total
		m.collectionPages = msg.pages
		m.collectionCursor = 0
		return m, nil

	case searchResultMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			m.statusErr = true
			return m, nil
		}
		m.searchResults = msg.results
		m.searchTotal = msg.total
		m.searchPages = msg.pages
		m.searchCursor = 0
		m.mode = viewSearch
		return m, nil

	case marketplaceMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			m.statusErr = true
			return m, nil
		}
		m.marketplaceListings = msg.listings
		m.marketplaceTotal = msg.total
		m.marketplacePages = msg.pages
		m.marketplaceCursor = 0
		m.marketplaceReleaseID = msg.releaseID
		m.mode = viewMarketplace
		return m, nil

	case releaseDetailMsg:
		m.loading = false
		if msg.err != nil {
			m.status = "Error: " + msg.err.Error()
			m.statusErr = true
			return m, nil
		}
		m.detail = msg.detail
		m.detailScroll = 0
		m.mode = viewDetail
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	// search input active — route keys to text input
	if m.searching {
		switch msg.String() {
		case "enter":
			query := m.searchInput.Value()
			if query == "" {
				m.searching = false
				return m, nil
			}
			m.lastQuery = query
			m.searchPage = 1
			m.loading = true
			m.searching = false
			m.searchInput.Blur()
			return m, tea.Batch(m.spinner.Tick, searchDB(m.client, query, 1))
		case "esc":
			m.searching = false
			m.searchInput.Blur()
			return m, nil
		default:
			var cmd tea.Cmd
			m.searchInput, cmd = m.searchInput.Update(msg)
			return m, cmd
		}
	}

	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "k":
		if m.mode == viewDetail {
			if m.detailScroll > 0 {
				m.detailScroll--
			}
		} else {
			m = m.moveCursor(-1)
		}
	case "down", "j":
		if m.mode == viewDetail {
			m.detailScroll++
		} else {
			m = m.moveCursor(1)
		}

	case "n":
		return m.nextPage()
	case "p":
		return m.prevPage()

	case "/":
		m.searching = true
		m.searchInput.SetValue("")
		m.searchInput.Focus()
		return m, nil

	case "c":
		m.mode = viewCollection
		m.loading = true
		return m, tea.Batch(m.spinner.Tick, fetchCollection(m.client, m.collectionPage))

	case "enter":
		if m.mode == viewCollection && len(m.collectionReleases) > 0 {
			r := m.collectionReleases[m.collectionCursor]
			m.prevMode = viewCollection
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, fetchRelease(m.client, r.BasicInfo.ID))
		}
		// on search results: open marketplace for selected release
		if m.mode == viewSearch && len(m.searchResults) > 0 {
			r := m.searchResults[m.searchCursor]
			m.marketplacePage = 1
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, fetchMarketplace(m.client, r.ID, 1))
		}

	case "esc":
		if m.mode == viewDetail {
			m.mode = m.prevMode
		} else if m.mode == viewMarketplace {
			m.mode = viewSearch
		} else if m.mode == viewSearch {
			m.mode = viewCollection
		}
	}

	return m, nil
}

func (m Model) moveCursor(delta int) Model {
	switch m.mode {
	case viewCollection:
		n := len(m.collectionReleases)
		m.collectionCursor = clamp(m.collectionCursor+delta, 0, n-1)
	case viewSearch:
		n := len(m.searchResults)
		m.searchCursor = clamp(m.searchCursor+delta, 0, n-1)
	case viewMarketplace:
		n := len(m.marketplaceListings)
		m.marketplaceCursor = clamp(m.marketplaceCursor+delta, 0, n-1)
	}
	return m
}

func (m Model) nextPage() (Model, tea.Cmd) {
	switch m.mode {
	case viewCollection:
		if m.collectionPage < m.collectionPages {
			m.collectionPage++
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, fetchCollection(m.client, m.collectionPage))
		}
	case viewSearch:
		if m.searchPage < m.searchPages {
			m.searchPage++
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, searchDB(m.client, m.lastQuery, m.searchPage))
		}
	case viewMarketplace:
		if m.marketplacePage < m.marketplacePages {
			m.marketplacePage++
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, fetchMarketplace(m.client, m.marketplaceReleaseID, m.marketplacePage))
		}
	}
	return m, nil
}

func (m Model) prevPage() (Model, tea.Cmd) {
	switch m.mode {
	case viewCollection:
		if m.collectionPage > 1 {
			m.collectionPage--
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, fetchCollection(m.client, m.collectionPage))
		}
	case viewSearch:
		if m.searchPage > 1 {
			m.searchPage--
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, searchDB(m.client, m.lastQuery, m.searchPage))
		}
	case viewMarketplace:
		if m.marketplacePage > 1 {
			m.marketplacePage--
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, fetchMarketplace(m.client, m.marketplaceReleaseID, m.marketplacePage))
		}
	}
	return m, nil
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
