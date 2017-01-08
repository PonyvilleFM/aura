// gophy is a simple library designed to give easy access to the Giphy API.
// gophy aims to have 100% API coverage with a full test suite.
package gophy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ClientOptions is used when initialising a new `Client` instance via the
// `NewClient` function. All values are optional.
type ClientOptions struct {
	ApiKey      string
	ApiEndpoint string
	HttpClient  *http.Client
}

type Client struct {
	apiKey      string
	apiEndpoint string
	httpClient  *http.Client
}

func NewClient(co *ClientOptions) *Client {

	client := &Client{}

	// set default api key if not set
	if co.ApiKey == "" {
		client.apiKey = "dc6zaTOxFJmzC"
	} else {
		client.apiKey = co.ApiKey
	}

	// set default endpoint if not set (mostly used for overriding the server
	// url during test runs)
	if co.ApiEndpoint == "" {
		client.apiEndpoint = "https://api.giphy.com/v1"
	} else {
		client.apiEndpoint = strings.TrimRight(co.ApiEndpoint, "/")
	}

	// set default http client if not set. Useful in situations where you need
	// special behaviour or aren't able to use a standard `http.Client`
	// instance (like on appengine).
	if co.HttpClient == nil {
		client.httpClient = &http.Client{}
	} else {
		client.httpClient = co.HttpClient
	}

	return client
}

func (c *Client) makeRequest(suffix string, qs *url.Values, ds interface{}) error {

	// inject configured API key into url
	qs.Set("api_key", c.apiKey)

	// execute HTTP request
	u := fmt.Sprintf("%s/%s?%s", c.apiEndpoint, suffix, qs.Encode())
	resp, err := c.httpClient.Get(u)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	// unmarshal HTTP response as JSON into the provided data structure
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(ds); err != nil {
		return err
	}

	return nil
}

// GetGifById returns a single GIF from the Giphy API.
func (c *Client) GetGifById(id string) (*Gif, error) {

	// construct and execute the HTTP request
	u := fmt.Sprintf("gifs/%s", id)
	sr := &singleResult{}
	if err := c.makeRequest(u, &url.Values{}, sr); err != nil {
		return nil, err
	}

	return sr.Data, nil
}

// GetGifsById returns a number of GIFs from the Giphy API. This method accepts
// a maximum of 100 ids.
func (c *Client) GetGifsById(ids ...string) ([]*Gif, error) {

	// check that a query string was actually passed in
	if len(ids) > 100 {
		err := errors.New("Cannot fetch more than 100 GIFs at a time.")
		return nil, err
	}

	// build query string that will be appended to url
	qs := &url.Values{}
	qs.Set("ids", strings.Join(ids, ","))

	// construct and execute the HTTP request
	pr := &paginatedResults{}
	if err := c.makeRequest("gifs", qs, pr); err != nil {
		return nil, err
	}

	return pr.Data, nil
}

// searchCommon provides common search functionality for both GIF and sticker
// search endpoints.
func (c *Client) searchCommon(q string, rating string, limit int, offset int, urlFragment string) ([]*Gif, int, error) {

	// ensure the value for `limit` is valid
	if err := validateLimit(limit, 1, 100); err != nil {
		return nil, 0, err
	}

	// ensure the value for `offset` is valid
	if offset < 0 {
		err := fmt.Errorf("%d is not a valid value for `offset`, must not be negative.", offset)
		return nil, 0, err
	}

	// check that a query string was actually passed in
	if len(q) < 1 {
		err := errors.New("`q` parameter must not be empty.")
		return nil, 0, err
	}

	// check that the given rating is valid
	if err := validateRating(rating); err != nil {
		return nil, 0, err
	}

	// build query string that will be appended to url
	qs := &url.Values{}
	qs.Set("q", q)
	if rating != "" {
		qs.Set("rating", rating)
	}
	qs.Set("limit", strconv.Itoa(limit))
	qs.Set("offset", strconv.Itoa(offset))

	// construct and execute the HTTP request
	sr := &paginatedResults{}
	if err := c.makeRequest(urlFragment, qs, sr); err != nil {
		return nil, 0, err
	}

	return sr.Data, sr.Pagination.TotalCount, nil
}

// SearchGifs searches the Giphy API for GIFs with the specified options.
// Returns a slice containing the returned gifs, the total number of images
// available for the specified query (so that you can paginate your requests as
// required), and an error if one occured.
func (c *Client) SearchGifs(q string, rating string, limit int, offset int) ([]*Gif, int, error) {
	return c.searchCommon(q, rating, limit, offset, "gifs/search")
}

// SearchStickers replicates the functionality and requirements of the classic
// Giphy search, but returns animated stickers rather than gifs.
func (c *Client) SearchStickers(q string, rating string, limit int, offset int) ([]*Gif, int, error) {
	return c.searchCommon(q, rating, limit, offset, "stickers/search")
}

// translateCommon provides common search functionality for both GIF and
// sticker translate endpoints.
func (c *Client) translateCommon(q string, rating string, urlFragment string) (*Gif, error) {

	// check that a query string was actually passed in
	if len(q) < 1 {
		err := errors.New("`q` parameter must not be empty.")
		return nil, err
	}

	// check that the given rating is valid
	if err := validateRating(rating); err != nil {
		return nil, err
	}

	// build query string that will be appended to url
	qs := &url.Values{}
	qs.Set("s", q)
	if rating != "" {
		qs.Set("rating", rating)
	}

	// construct and execute the HTTP request
	sr := &singleResult{}
	if err := c.makeRequest(urlFragment, qs, sr); err != nil {
		return nil, err
	}

	return sr.Data, nil
}

// TranslateGif is prototype endpoint for using Giphy as a translation engine
// for a GIF dialect. The translate API draws on search, but uses the Giphy
// "special sauce" to handle translating from one vocabulary to another. In
// this case, words and phrases to GIFs. Returns a single GIF from the Giphy
// API.
func (c *Client) TranslateGif(q string, rating string) (*Gif, error) {
	return c.translateCommon(q, rating, "gifs/translate")
}

// TranslateSticker replicates the functionality and requirements of the
// classic Giphy translate endpoint, but returns animated stickers rather than
// gifs.
func (c *Client) TranslateSticker(q string, rating string) (*Gif, error) {
	return c.translateCommon(q, rating, "stickers/translate")
}

// trendingCommon provides functionality common to bothe the GIF and sticker
// trending endpoints.
func (c *Client) trendingCommon(rating string, limit int, urlFragment string) ([]*Gif, error) {

	// ensure the value for `limit` is valid
	if err := validateLimit(limit, 1, 100); err != nil {
		return nil, err
	}

	// check that the given rating is valid
	if err := validateRating(rating); err != nil {
		return nil, err
	}

	// build query string that will be appended to url
	qs := &url.Values{}
	if rating != "" {
		qs.Set("rating", rating)
	}
	qs.Set("limit", strconv.Itoa(limit))

	// construct and execute the HTTP request
	sr := &paginatedResults{}
	if err := c.makeRequest(urlFragment, qs, sr); err != nil {
		return nil, err
	}

	return sr.Data, nil
}

// TrendingGifs fetches GIFs currently trending online. The data returned
// mirrors that used to create The Hot 100 list of GIFs on Giphy.
func (c *Client) TrendingGifs(rating string, limit int) ([]*Gif, error) {
	return c.trendingCommon(rating, limit, "gifs/trending")
}

// TrendingStickers replicates the functionality and requirements of the
// classic Giphy trending endpoint, but returns animated stickers rather than
// gifs.
func (c *Client) TrendingStickers(rating string, limit int) ([]*Gif, error) {
	return c.trendingCommon(rating, limit, "stickers/trending")
}

// validateRating checks if the given string matches an allowed value for the
// giphy API's `rating` parameter.
func validateRating(r string) error {

	r = strings.ToLower(r)
	switch r {
	case "y", "g", "pg", "pg-13", "r", "":
		return nil
	}

	fmtString := "\"%s\" is not a valid value for `rating`, must be one of \"y\", \"g\", \"pg\", \"pg-13\", \"r\" or \"\""
	return fmt.Errorf(fmtString, r)
}

// validateLimit checks if the given integer falls within an allowed range for
// the giphy API's `limit` parameter.
func validateLimit(v, min, max int) error {

	// ensure the value for `limit` is valid
	if v < min || v > max {
		err := fmt.Errorf("%d is not a valid value for `limit`, must be between %d-%d (inclusive).", v, min, max)
		return err
	}

	return nil
}
