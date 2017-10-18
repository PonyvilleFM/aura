// Utilities for spinning up/down a HTTP server for use during tests.
package tests

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

type testResponse struct {
	Status  int
	Headers map[string]string
	Body    string
}

func testDataHandler(w http.ResponseWriter, r *http.Request) {

	u := strings.TrimLeft(r.URL.String(), "/")
	b64Filename := base64.URLEncoding.EncodeToString([]byte(u))

	testData, err := ioutil.ReadFile(fmt.Sprintf("tests/data/%s.json", b64Filename))
	if err != nil {
		fmt.Printf("File not found (add to urls.txt): %s\n", u)
		http.NotFound(w, r)
		return
	}

	tr := &testResponse{}
	if err := json.Unmarshal(testData, tr); err != nil {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(tr.Status)
	for k, v := range tr.Headers {
		w.Header().Set(k, v)
	}
	fmt.Fprint(w, tr.Body)
}

// SetupTestServer initialises a HTTP server that returns mock responses for
// testing.
func SetupTestServer() *httptest.Server {
	// test server
	mux := http.NewServeMux()
	mux.HandleFunc("/", testDataHandler)
	server := httptest.NewServer(mux)
	return server
}

// TeardownTestServer closes the test HTTP server.
func TeardownTestServer(s *httptest.Server) {
	s.Close()
}
