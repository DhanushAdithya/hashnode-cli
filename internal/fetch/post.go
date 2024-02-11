package fetch

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
)

const postQuery = `{
    "query": "query Post($id: ID!) { post(id: $id) { id title brief publishedAt content { markdown } url readTimeInMinutes author { name } } }",
    "variables": {
        "id": "%s"
    }
}`

type Post struct {
	Data struct {
		Post utils.Post `json:"post"`
	} `json:"data"`
	Errors []utils.Error `json:"errors"`
}

func PostResponse(postID string) Post {
	var response Post
	post, err := query(fmt.Sprintf(postQuery, postID))
	if err != nil {
		utils.Exit(err)
	}
	defer post.Close()
	data, err := io.ReadAll(post)
	if err := json.Unmarshal(data, &response); err != nil {
		utils.Exit("Unable to parse response")
	}
	return response
}

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
