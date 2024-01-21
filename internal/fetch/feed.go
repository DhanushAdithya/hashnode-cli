package fetch

import (
	"encoding/json"
	"io"

	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
)

const feed = `{
    "query": "query Feed($first: Int!) { feed(first: $first) { edges { node { title brief publishedAt content { markdown } url readTimeInMinutes author { name } } } pageInfo { hasNextPage endCursor } } }",
    "variables": {
        "first": 10,
    }
}`

// "after": "",
// "filter": {
//     "type": "%s",
//     "tags": %s,
//     "minReadTime": %d,
//     "maxReadTime": %d
// }

type Node struct {
	Title       string `json:"title"`
	Brief       string `json:"brief"`
	PublishedAt string `json:"publishedAt"`
	Author      struct {
		Name string `json:"name"`
	} `json:"author"`
	URL               string `json:"url"`
	ReadTimeInMinutes int    `json:"readTimeInMinutes"`
	Content           struct {
		Markdown string `json:"markdown"`
	} `json:"content"`
}

type Feed struct {
	Data struct {
		Feed struct {
			Edges []struct {
				Node Node `json:"node"`
			} `json:"edges"`
		} `json:"feed"`
	} `json:"data"`
	Errors []utils.Error `json:"errors"`
}

var FeedTypes = []string{
	"FOLLOWING",
	"PERSONALIZED",
	"RECENT",
	"RELEVANT",
	"FEATURED",
	"BOOKMARKS",
	"READING_HISTORY",
}

func FeedResponse() Feed {
	var response Feed
	feed, err := query(feed)
	if err != nil {
		utils.Exit(err)
	}
	defer feed.Close()
	data, err := io.ReadAll(feed)
	if err := json.Unmarshal(data, &response); err != nil {
		utils.Exit("Unable to parse response")
	}
	return response
}
