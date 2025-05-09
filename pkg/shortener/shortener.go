package shortener

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Shortener is used to manage url shortener.
type Shortener struct {
	// Api key usually obtained from <api_homepage>/member/tools/api.
	ApiKey string `json:"api_key,omitempty" bson:"api_key,omitempty"`
	// Url of the homepage of the shortener for ex: https://gplinks.com.
	RootURL string `json:"root_url,omitempty" bson:"root_url,omitempty"`
	// If set the shortener will return the raw api url in the format <shortener_url>/api?api=<apikey>&url=<url>
	RawURL bool `json:"raw_url,omitempty" bson:"raw_url,omitempty"`
}

// ShortenURLOpts are optional parameter for ShortenURL.
type ShortenURLOpts struct {
	// If set api call will be made returning the result url.
}

// shortenResult is the response from api request.
type shortenResult struct {
	// Shortened url.
	Url string `json:"url,omitempty"`
	// Status returned by most shorteners either Success or Error.
	Status string `json:"status"`
	// Error message returned by most shorteners.
	Message interface{} `json:"message,omitempty"` // type can vary from list to string so we'll just stringify it
}

func (c *Shortener) ShortenURL(inputURL string) (string, error) {
	// protocol is already added when saving
	requestURL := fmt.Sprintf("%v/api?api=%v&url=%v", c.RootURL, c.ApiKey, inputURL)

	if c.RawURL {
		return requestURL, nil
	}

	req, err := http.NewRequest("GET", requestURL, http.NoBody)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res shortenResult

	err = json.Unmarshal(bytes, &res)
	if err != nil {
		return "", nil
	}

	if strings.EqualFold(res.Status, "error") {
		return "", fmt.Errorf("shortener error response: %v", res.Message)
	}

	return res.Url, nil
}
