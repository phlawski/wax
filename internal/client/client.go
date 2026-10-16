package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.discogs.com"

// Client is an authenticated Discogs API client.
type Client struct {
	token      string
	username   string
	httpClient *http.Client
}

// New creates a Client with the given token and username.
func New(token, username string) *Client {
	return &Client{
		token:    token,
		username: username,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// get performs a GET request and decodes JSON into dest.
func (c *Client) get(path string, params url.Values, dest any) error {
	u := baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Discogs token="+c.token)
	req.Header.Set("User-Agent", "wax-cli/1.0 +https://github.com/you/wax")
	req.Header.Set("Accept", "application/vnd.discogs.v2.discogs+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(dest)
}

// Pagination holds paging metadata from API responses.
type Pagination struct {
	Pages   int `json:"pages"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Items   int `json:"items"`
}

// CollectionRelease represents one release in a collection folder.
type CollectionRelease struct {
	ID          int `json:"id"`
	InstanceID  int `json:"instance_id"`
	FolderID    int `json:"folder_id"`
	BasicInfo   ReleaseBasicInfo `json:"basic_information"`
}

// ReleaseBasicInfo holds the core fields of a release.
type ReleaseBasicInfo struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	Year     int      `json:"year"`
	Artists  []Artist `json:"artists"`
	Labels   []Label  `json:"labels"`
	Formats  []Format `json:"formats"`
	Genres   []string `json:"genres"`
	Styles   []string `json:"styles"`
}

// Artist holds artist name.
type Artist struct {
	Name string `json:"name"`
}

// Label holds label name and catalog number.
type Label struct {
	Name    string `json:"name"`
	CatNo   string `json:"catno"`
}

// Format holds format name and descriptions.
type Format struct {
	Name         string   `json:"name"`
	Qty          string   `json:"qty"`
	Descriptions []string `json:"descriptions"`
}

// CollectionResponse is the API response for collection listings.
type CollectionResponse struct {
	Pagination Pagination          `json:"pagination"`
	Releases   []CollectionRelease `json:"releases"`
}

// GetCollection returns paginated collection releases for folder 0 (all).
func (c *Client) GetCollection(page, perPage int) (*CollectionResponse, error) {
	path := fmt.Sprintf("/users/%s/collection/folders/0/releases", c.username)
	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))
	params.Set("sort", "artist")
	params.Set("sort_order", "asc")

	var result CollectionResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SearchResult is one item from the database search.
type SearchResult struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Year        string   `json:"year"`
	Country     string   `json:"country"`
	Format      []string `json:"format"`
	Label       []string `json:"label"`
	Genre       []string `json:"genre"`
	Style       []string `json:"style"`
	ResourceURL string   `json:"resource_url"`
	Type        string   `json:"type"`
}

// SearchResponse is the API response for database search.
type SearchResponse struct {
	Pagination Pagination     `json:"pagination"`
	Results    []SearchResult `json:"results"`
}

// Search queries the Discogs database.
func (c *Client) Search(query, searchType string, page, perPage int) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	if searchType != "" {
		params.Set("type", searchType)
	}
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))

	var result SearchResponse
	if err := c.get("/database/search", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// MarketplaceListing is one marketplace listing.
type MarketplaceListing struct {
	ID       int     `json:"id"`
	Status   string  `json:"status"`
	Price    Price   `json:"price"`
	Condition string `json:"condition"`
	SleeveCondition string `json:"sleeve_condition"`
	Comments string  `json:"comments"`
	Seller   Seller  `json:"seller"`
	Release  MarketplaceRelease `json:"release"`
}

// Price holds value and currency.
type Price struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

// Seller holds seller info.
type Seller struct {
	Username string `json:"username"`
	URL      string `json:"html_url"`
}

// MarketplaceRelease holds release info within a listing.
type MarketplaceRelease struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Format      string `json:"format"`
	ResourceURL string `json:"resource_url"`
}

// MarketplaceSearchResponse wraps marketplace listing search results.
type MarketplaceSearchResponse struct {
	Pagination Pagination           `json:"pagination"`
	Listings   []MarketplaceListing `json:"listings"`
}

// SearchCollection searches within a user's collection via the database search endpoint.
func (c *Client) SearchCollection(query string, page, perPage int, username string) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("type", "release")
	params.Set("collection", username)
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))

	var result SearchResponse
	if err := c.get("/database/search", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// Track holds a single tracklist entry.
type Track struct {
	Position string `json:"position"`
	Title    string `json:"title"`
	Duration string `json:"duration"`
}

// ReleaseStats holds community have/want counts.
type ReleaseStats struct {
	Have int `json:"in_collection"`
	Want int `json:"in_wantlist"`
}

// Community wraps community stats and lowest price.
type Community struct {
	Have       int          `json:"have"`
	Want       int          `json:"want"`
	LowestPrice *Price      `json:"lowest_price"`
	Rating      RatingStats `json:"rating"`
}

// RatingStats holds average rating and count.
type RatingStats struct {
	Count   int     `json:"count"`
	Average float64 `json:"average"`
}

// ReleaseDetail is the full release from /releases/{id}.
type ReleaseDetail struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Year        int       `json:"year"`
	Artists     []Artist  `json:"artists"`
	Labels      []Label   `json:"labels"`
	Formats     []Format  `json:"formats"`
	Genres      []string  `json:"genres"`
	Styles      []string  `json:"styles"`
	Country     string    `json:"country"`
	Notes       string    `json:"notes"`
	Tracklist   []Track   `json:"tracklist"`
	Community   Community `json:"community"`
}

// GetRelease fetches full release details by ID.
func (c *Client) GetRelease(id int) (*ReleaseDetail, error) {
	var result ReleaseDetail
	if err := c.get(fmt.Sprintf("/releases/%d", id), nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMarketplaceListings searches marketplace listings for a release ID.
func (c *Client) GetMarketplaceListings(releaseID, page, perPage int) (*MarketplaceSearchResponse, error) {
	path := fmt.Sprintf("/marketplace/search")
	params := url.Values{}
	params.Set("release_id", fmt.Sprintf("%d", releaseID))
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("per_page", fmt.Sprintf("%d", perPage))

	var result MarketplaceSearchResponse
	if err := c.get(path, params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
