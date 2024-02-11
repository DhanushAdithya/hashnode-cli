package fetch

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
)

const feed = `{
    "query": "query Feed($first: Int!, $after: String, $filter: FeedFilter) { feed(first: $first, after: $after, filter: $filter) { edges { node { id title brief publishedAt content { markdown } url readTimeInMinutes author { name } } } pageInfo { hasNextPage endCursor } } }",
    "variables": {
        "first": 10,
        "filter": {
            "type": "%s",
            "minReadTime": %d,
            "maxReadTime": %d
        },
        "after": "%s"
    }
}`

type Feed struct {
	Data struct {
		Feed struct {
			Edges []struct {
				Node utils.Post `json:"node"`
			} `json:"edges"`
			PageInfo struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
		} `json:"feed"`
	} `json:"data"`
	Errors []utils.Error `json:"errors"`
}

var FeedTypes = []string{
	"FEATURED",
	"PERSONALIZED",
	"FOLLOWING",
	"RECENT",
	"RELEVANT",
	"BOOKMARKS",
	"READING_HISTORY",
}

func FeedResponse(
	r chan struct{},
	feedType string,
	minRead int,
	maxRead int,
	after string,
) Feed {
	var response Feed
	feed, err := query(
		fmt.Sprintf(feed, feedType, minRead, maxRead, after),
	)
	if err != nil {
		close(r)
		utils.Exit(err)
	}
	defer feed.Close()
	data, err := io.ReadAll(feed)
	if err := json.Unmarshal(data, &response); err != nil {
		close(r)
		utils.Exit("Unable to parse response")
	}
	return response
}
