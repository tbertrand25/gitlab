package gitlab

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type glClient struct {
	baseURL    string
	auth       string
	httpClient http.Client
}

func MakeGlClient() glClient {
	return glClient{
		"https://www.gitlab.com/api/v4",
		os.Getenv("GITLAB_API_KEY"),
		http.Client{Timeout: 10 * time.Second},
	}
}

func (c glClient) Do(rType string, URL string, params map[string]string) (*http.Response, error) {

	req, err := http.NewRequest(rType, URL, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	q := req.URL.Query()

	for name, value := range params {
		q.Add(name, value)
	}

	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c glClient) parseLinkHeader(h string) map[string]string {
	linkMap := make(map[string]string)

	links := strings.Split(h, ", ")

	for _, elem := range links {
		pair := strings.Split(elem, "; rel=")
		ind := strings.Trim(pair[1], "\"")
		link := strings.Trim(pair[0], "<")
		link = strings.Trim(link, ">")
		linkMap[ind] = link
	}
	return linkMap
}
