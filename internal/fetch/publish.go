package fetch

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
)

const publish = `{
    "query": "mutation PublishPost($input: PublishPostInput!) { publishPost(input: $input) { post { id } } }",
    "variables": {
        "input": {
            "title": "%s",
            "contentMarkdown": "%s",
            "tags": [],
            "coverImageOptions": {
                "coverImageURL": "%s"
            },
            "publicationId": "%s"
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

func PublishResponse(title, content, coverImg, publicationId string) Publish {
	var response Publish
	publish, err := query(fmt.Sprintf(publish, title, content, coverImg, publicationId))
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
