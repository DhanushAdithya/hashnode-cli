package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const me = `{ "query": "query Me { me { id username name followersCount followingsCount newPost: posts(pageSize: 1, page: 1, sortBy: DATE_PUBLISHED_ASC) { nodes { title url } } oldPost: posts(pageSize: 1, page: 1, sortBy: DATE_PUBLISHED_DESC) { nodes { title url } } } }" }`

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
		} `json:"me"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func MeResponse() Me {
	var response Me
	me, err := query(me)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer me.Close()
	data, err := io.ReadAll(me)
	if err := json.Unmarshal(data, &response); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return response
}
