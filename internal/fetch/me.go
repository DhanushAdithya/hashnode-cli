package fetch

import (
	"encoding/json"
	"io"

	"github.com/DhanushAdithya/hashnode-cli/internal/utils"
)

const me = `{ "query": "query Me { me { id username name followersCount followingsCount newPost: posts(pageSize: 1, page: 1, sortBy: DATE_PUBLISHED_ASC) { nodes { title url } } oldPost: posts(pageSize: 1, page: 1, sortBy: DATE_PUBLISHED_DESC) { nodes { title url } } publications(first: 10) { edges { node { id title } } } } }" }`

type Me struct {
	Data struct {
		Me struct {
			ID              string `json:"id"`
			Username        string `json:"username"`
			Name            string `json:"name"`
			FollowersCount  int    `json:"followersCount"`
			FollowingsCount int    `json:"followingsCount"`
			NewPost         struct {
				Nodes []struct {
					Title string `json:"title"`
					URL   string `json:"url"`
				} `json:"nodes"`
			} `json:"newPost"`
			OldPost struct {
				Nodes []struct {
					Title string `json:"title"`
					URL   string `json:"url"`
				} `json:"nodes"`
			} `json:"oldPost"`
			Publications struct {
				Edges []struct {
					Node struct {
						ID    string `json:"id"`
						Title string `json:"title"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"publications"`
		} `json:"me"`
	} `json:"data"`
	Errors []utils.Error `json:"errors"`
}

func MeResponse() Me {
	var response Me
	me, err := query(me)
	if err != nil {
		utils.Exit(err)
	}
	defer me.Close()
	data, err := io.ReadAll(me)
	if err := json.Unmarshal(data, &response); err != nil {
		utils.Exit("Unable to parse response")
	}
	return response
}
