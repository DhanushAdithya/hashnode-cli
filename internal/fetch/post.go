package fetch

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
)

const like = `{
    "query": "mutation LikePost($input: LikePostInput!) { likePost(input: $input) { post { id } } }",
    "variables": {
        "input": {
            "postId": "%s"
        }
    }
}`

type Like struct {
	Data struct {
		Like struct {
			Post struct {
				Id string `json:"id"`
			} `json:"post"`
		} `json:"likePost"`
	} `json:"data"`
	Errors []utils.Error `json:"errors"`
}

func LikeResponse(postID string) Like {
	var response Like
	like, err := query(fmt.Sprintf(like, postID))
	if err != nil {
		utils.Exit(err)
	}
	defer like.Close()
	data, err := io.ReadAll(like)
	if err := json.Unmarshal(data, &response); err != nil {
		utils.Exit("Unable to parse response")
	}
	return response
}
