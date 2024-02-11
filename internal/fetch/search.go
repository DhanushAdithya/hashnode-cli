package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
	"github.com/spf13/viper"
)

const searchURL = "https://amerdmzm12-dsn.algolia.net/1/indexes/posts_prod/query"

type Search struct {
	Hits []struct {
		Id          string `json:"objectID"`
		Title       string `json:"title"`
		Brief       string `json:"brief"`
		PublishedAt string `json:"dateAdded"`
		Author      string `json:"authorName"`
	} `json:"hits"`
}

const searchQuery = `{
    "query": "%s",
    "filters": "(type:question OR type:story) AND (isActive=1 AND isDelisted=0 AND sB=0) AND totalReactions >= 25",
    "hitsPerPage": 10,
    "page": 0
}`

func SearchResponse(r chan struct{}, query string) Search {
	var response Search
	searchQuery := fmt.Sprintf(searchQuery, query)
	req, err := http.NewRequest(http.MethodPost, searchURL, strings.NewReader(searchQuery))
	if err != nil {
		close(r)
		utils.Exit("Unable to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Algolia-API-Key", viper.GetString("search-token"))
	req.Header.Set("X-Algolia-Application-Id", "AMERDMZM12")
	req.Header.Set("Authorization", viper.GetString("token"))
	res, err := client.Do(req)
	if err != nil {
		close(r)
		utils.Exit("Unable to make request")
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err := json.Unmarshal(data, &response); err != nil {
		close(r)
		utils.Exit("Unable to parse response")
	}
	return response
}
