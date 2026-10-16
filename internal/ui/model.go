package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"wax/internal/client"
)

// viewMode represents which screen is active.
type viewMode int

const (
	viewCollection viewMode = iota
	viewSearch
	viewMarketplace
	viewDetail
)

// fetchMsg carries a fresh page of collection releases.
type fetchMsg struct {
	releases []client.CollectionRelease
	total    int
	pages    int
	query    string
	err      error
}

// searchResultMsg carries database search results.
type searchResultMsg struct {
	results []client.SearchResult
	total   int
	pages   int
	err     error
}

// marketplaceMsg carries marketplace listings.
type marketplaceMsg struct {
	listings  []client.MarketplaceListing
	total     int
	pages     int
	releaseID int
	err       error
}

// statusMsg carries a transient status line message.
type statusMsg struct {
	text  string
	isErr bool
}

// releaseDetailMsg carries a full release detail response.
type releaseDetailMsg struct {
	detail *client.ReleaseDetail
	err    error
}

// Model is the root Bubbletea model.
type Model struct {
	client    *client.Client
	mode      viewMode
	width     int
	height    int
	loading   bool
	spinner   spinner.Model
	status    string
	statusErr bool

	// collection state
	collectionReleases []client.CollectionRelease
	collectionPage     int
	collectionPages    int
	collectionTotal    int
	collectionCursor   int
	collectionQuery    string

	// search state
	searching        bool
	searchInput      textinput.Model
	searchResults    []client.SearchResult
	searchPage       int
	searchPages      int
	searchTotal      int
	searchCursor     int
	lastQuery        string
	lastSearchParams client.SearchParams

	// format filter (k9s-style)
	filtering      bool
	filterInput    textinput.Model
	formatFilter   string
	filterMatches  []string
	filterMatchIdx int

	// marketplace state
	marketplaceListings  []client.MarketplaceListing
	marketplacePage      int
	marketplacePages     int
	marketplaceTotal     int
	marketplaceCursor    int
	marketplaceReleaseID int

	// detail state
	detail         *client.ReleaseDetail
	detailScroll   int
	prevMode       viewMode
}

// New creates the initial model. Init() fires the first fetch.
func New(c *client.Client) Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = styleFilter

	si := textinput.New()
	si.Placeholder = "search Discogs…"
	si.CharLimit = 100

	fi := textinput.New()
	fi.Placeholder = "Vinyl, CD, Box Set…"
	fi.CharLimit = 40

	return Model{
		client:         c,
		mode:           viewCollection,
		loading:        true,
		spinner:        sp,
		searchInput:    si,
		filterInput:    fi,
		collectionPage: 1,
	}
}

// Init fires the initial collection fetch and spinner tick.
func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchCollection(m.client, "", 1))
}
