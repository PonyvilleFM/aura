package gophy

import (
	"net/http"
	"testing"

	"github.com/paddycarey/gophy/tests"
)

// Test that a client is constructed correctly with default options
func TestDefaultClient(t *testing.T) {
	co := &ClientOptions{}
	client := NewClient(co)
	if client.apiEndpoint != "https://api.giphy.com/v1" {
		t.Errorf("Default api endpoint should be \"https://api.giphy.com/v1\", not \"%s\"", client.apiEndpoint)
	}
	if client.apiKey != "dc6zaTOxFJmzC" {
		t.Errorf("Default api key should be \"dc6zaTOxFJmzC\", not \"%s\"", client.apiKey)
	}
	if client.httpClient == nil {
		t.Error("Default http client should not be nil")
	}
}

// Test that a custom HTTP client is properly added to the client
func TestCustomHttpClient(t *testing.T) {
	co := &ClientOptions{HttpClient: &http.Client{}}
	client := NewClient(co)
	if client.httpClient == nil {
		t.Error("Http client should not be nil")
	}
}

// Test GIF search functionality. This is a table-driven test, using the data
// at `tests.SearchTestData`.
func TestSearchGifs(t *testing.T) {
	server := tests.SetupTestServer()
	defer tests.TeardownTestServer(server)

	for _, td := range tests.SearchGifsTestData {

		// gophy client configured to use test server
		clientOptions := &ClientOptions{ApiKey: td.ApiKey, ApiEndpoint: server.URL}
		client := NewClient(clientOptions)

		searchResults, totalCount, err := client.SearchGifs(td.Q, td.Rating, td.Limit, td.Offset)
		if err != nil {
			if !td.ExpectedError {
				t.Errorf("Unexpected search error occurred (q=%s): %v", td.Q, err)
			}
			continue
		}
		if err == nil && td.ExpectedError {
			t.Errorf("Expected search error didn't happen (q=%s)", td.Q)
			continue
		}

		if numResultsReturned := len(searchResults); numResultsReturned != td.ExpectedNumReturned {
			t.Errorf("Expected %d search results, got %d: (q=%s)", td.ExpectedNumReturned, numResultsReturned, td.Q)
		}

		if totalCount != td.ExpectedTotalCount {
			t.Errorf("Expected %d total results, got %d: (q=%s)", td.ExpectedTotalCount, totalCount, td.Q)
		}

	}
}

// Test sticker search functionality. This is a table-driven test, using the
// data at `tests.SearchTestData`.
func TestSearchStickers(t *testing.T) {
	server := tests.SetupTestServer()
	defer tests.TeardownTestServer(server)

	for _, td := range tests.SearchStickersTestData {

		// gophy client configured to use test server
		clientOptions := &ClientOptions{ApiKey: td.ApiKey, ApiEndpoint: server.URL}
		client := NewClient(clientOptions)

		searchResults, totalCount, err := client.SearchStickers(td.Q, td.Rating, td.Limit, td.Offset)
		if err != nil {
			if !td.ExpectedError {
				t.Errorf("Unexpected search error occurred (q=%s): %v", td.Q, err)
			}
			continue
		}
		if err == nil && td.ExpectedError {
			t.Errorf("Expected search error didn't happen (q=%s)", td.Q)
			continue
		}

		if numResultsReturned := len(searchResults); numResultsReturned != td.ExpectedNumReturned {
			t.Errorf("Expected %d search results, got %d: (q=%s)", td.ExpectedNumReturned, numResultsReturned, td.Q)
		}

		if totalCount != td.ExpectedTotalCount {
			t.Errorf("Expected %d total results, got %d: (q=%s)", td.ExpectedTotalCount, totalCount, td.Q)
		}

	}
}

// Test get GIF by ID functionality. This is a table-driven test, using the
// data at `tests.GetGifByIdTestData`.
func TestGetGifById(t *testing.T) {
	server := tests.SetupTestServer()
	defer tests.TeardownTestServer(server)

	for _, td := range tests.GetGifByIdTestData {
		// gophy client configured to use test server
		clientOptions := &ClientOptions{ApiKey: td.ApiKey, ApiEndpoint: server.URL}
		client := NewClient(clientOptions)

		_, err := client.GetGifById(td.Id)
		if err != nil {
			if !td.ExpectedGetGifByIdError {
				t.Errorf("Unexpected GetGifById error occurred (id=%s): %v", td.Id, err)
			}
			continue
		}
		if err == nil && td.ExpectedGetGifByIdError {
			t.Errorf("Expected GetGifById error didn't happen (id=%s)", td.Id)
			continue
		}
	}
}

// Test get GIFs by ID functionality. This is a table-driven test, using the
// data at `tests.GetGifsByIdTestData`.
func TestGetGifsById(t *testing.T) {
	server := tests.SetupTestServer()
	defer tests.TeardownTestServer(server)

	for _, td := range tests.GetGifsByIdTestData {
		// gophy client configured to use test server
		clientOptions := &ClientOptions{ApiKey: td.ApiKey, ApiEndpoint: server.URL}
		client := NewClient(clientOptions)

		gifs, err := client.GetGifsById(td.Ids...)
		if err != nil {
			if !td.ExpectedGetGifsByIdError {
				t.Errorf("Unexpected GetGifsById error occurred (ids=%v): %v", td.Ids, err)
			}
			continue
		}
		if err == nil && td.ExpectedGetGifsByIdError {
			t.Errorf("Expected GetGifsById error didn't happen (ids=%v)", td.Ids)
			continue
		}
		if numGifs := len(gifs); td.ExpectedNumReturned != numGifs {
			t.Errorf("Expected %d results, got %d: (ids=%v)", td.ExpectedNumReturned, numGifs, td.Ids)
		}

	}
}

// Test GIF translate functionality. This is a table-driven test, using the
// data at `tests.TranslateTestData`.
func TestTranslateGif(t *testing.T) {
	server := tests.SetupTestServer()
	defer tests.TeardownTestServer(server)

	for _, td := range tests.TranslateGifTestData {
		// gophy client configured to use test server
		clientOptions := &ClientOptions{ApiKey: td.ApiKey, ApiEndpoint: server.URL}
		client := NewClient(clientOptions)

		_, err := client.TranslateGif(td.Q, td.Rating)
		if err != nil {
			if !td.ExpectedTranslateError {
				t.Errorf("Unexpected translate error occurred (q=%s): %v", td.Q, err)
			}
			continue
		}
		if err == nil && td.ExpectedTranslateError {
			t.Errorf("Expected translate error didn't happen (q=%s)", td.Q)
			continue
		}

	}
}

// Test sticker translate functionality. This is a table-driven test, using the
// data at `tests.TranslateTestData`.
func TestTranslateSticker(t *testing.T) {
	server := tests.SetupTestServer()
	defer tests.TeardownTestServer(server)

	for _, td := range tests.TranslateStickerTestData {
		// gophy client configured to use test server
		clientOptions := &ClientOptions{ApiKey: td.ApiKey, ApiEndpoint: server.URL}
		client := NewClient(clientOptions)

		_, err := client.TranslateSticker(td.Q, td.Rating)
		if err != nil {
			if !td.ExpectedTranslateError {
				t.Errorf("Unexpected translate error occurred (q=%s): %v", td.Q, err)
			}
			continue
		}
		if err == nil && td.ExpectedTranslateError {
			t.Errorf("Expected translate error didn't happen (q=%s)", td.Q)
			continue
		}

	}
}

// Test trending GIF functionality. This is a table-driven test, using the
// data at `tests.TrendingTestData`.
func TestTrendingGifs(t *testing.T) {
	server := tests.SetupTestServer()
	defer tests.TeardownTestServer(server)

	for _, td := range tests.TrendingGifsTestData {
		// gophy client configured to use test server
		clientOptions := &ClientOptions{ApiKey: td.ApiKey, ApiEndpoint: server.URL}
		client := NewClient(clientOptions)

		trendingResults, err := client.TrendingGifs(td.Rating, td.Limit)
		if err != nil {
			if !td.ExpectedError {
				t.Errorf("Unexpected error occurred: %v", err)
			}
			continue
		}
		if err == nil && td.ExpectedError {
			t.Error("Expected error didn't happen")
			continue
		}

		if numResultsReturned := len(trendingResults); numResultsReturned != td.ExpectedNumReturned {
			t.Errorf("Expected %d trending results, got %d", td.ExpectedNumReturned, numResultsReturned)
		}

	}
}

// Test trending functionality. This is a table-driven test, using the
// data at `tests.TrendingTestData`.
func TestTrendingStickers(t *testing.T) {
	server := tests.SetupTestServer()
	defer tests.TeardownTestServer(server)

	for _, td := range tests.TrendingStickersTestData {
		// gophy client configured to use test server
		clientOptions := &ClientOptions{ApiKey: td.ApiKey, ApiEndpoint: server.URL}
		client := NewClient(clientOptions)

		trendingResults, err := client.TrendingStickers(td.Rating, td.Limit)
		if err != nil {
			if !td.ExpectedError {
				t.Errorf("Unexpected error occurred: %v", err)
			}
			continue
		}
		if err == nil && td.ExpectedError {
			t.Error("Expected error didn't happen")
			continue
		}

		if numResultsReturned := len(trendingResults); numResultsReturned != td.ExpectedNumReturned {
			t.Errorf("Expected %d trending results, got %d", td.ExpectedNumReturned, numResultsReturned)
		}

	}
}
