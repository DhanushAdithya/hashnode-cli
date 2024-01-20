package fetch

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const baseURL = "https://gql.hashnode.com"

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func query(query string) (io.ReadCloser, error) {
	req, err := http.NewRequest(http.MethodPost, baseURL, strings.NewReader(query))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", viper.GetString("token"))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}
