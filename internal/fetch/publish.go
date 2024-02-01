package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
)

const publish = `{
    "query": "mutation PublishPost($input: PublishPostInput!) { publishPost(input: $input) { post { id } } }",
    "variables": {
        "input": {
            "title": "%s",
            "contentMarkdown": "%s",
            "coverImageOptions": {
                "coverImageURL": "%s"
            },
            "publicationId": "%s",
            "tags": %s
        }
    }
}`

type Publish struct {
	Data struct {
		PublishPost struct {
			Post struct {
				Id string `json:"id"`
			} `json:"post"`
		} `json:"publishPost"`
	} `json:"data"`
	Errors []utils.Error `json:"errors"`
}

func PublishResponse(title, content, coverImg, publicationId string, tags []string) Publish {
	var response Publish

	tagsMap := []map[string]string{}
	for _, tag := range tags {
		tagsMap = append(tagsMap, map[string]string{"name": tag, "slug": strings.ReplaceAll(tag, " ", "-")})
	}

	publish, err := query(fmt.Sprintf(publish, title, strings.Trim(content, "\n"), coverImg, publicationId, utils.Listify(tagsMap)))
	if err != nil {
		utils.Exit(err)
	}
	defer publish.Close()
	data, err := io.ReadAll(publish)
	if err := json.Unmarshal(data, &response); err != nil {
		utils.Exit("Unable to parse response")
	}
	return response
}
