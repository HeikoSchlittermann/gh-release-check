package gitea

import "fmt"

const url = `https://%s/api/v1/repos/%s/releases/latest`

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
