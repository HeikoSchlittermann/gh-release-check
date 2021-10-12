package github

import "fmt"

const url = `https://api.%s/repos/%s/releases/latest`

type client struct {
	url string
}

func NewClient(baseurl, project string) *client {
	return &client{
		url: fmt.Sprintf(url, baseurl, project),
	}
}

func (c client) URL() string {
	return c.url
}
