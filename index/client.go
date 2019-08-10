package index

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/morikuni/failure"
)

type Entry struct {
	Path      string    `json:"path"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

type Client struct {
	endpoint *url.URL
	c        *http.Client
}

func NewClient(c *http.Client) (*Client, error) {
	if c == nil {
		c = http.DefaultClient
	}
	u, err := url.Parse("https://index.golang.org")
	if err != nil {
		return nil, failure.Wrap(err)
	}
	return &Client{endpoint: u, c: c}, nil
}

func (c *Client) Index(queries ...Query) ([]*Entry, error) {
	u := *c.endpoint
	u.Path = "index"

	if queries != nil {
		var q query
		for _, query := range queries {
			query(&q)
		}
		u.RawQuery = q.encode()
	}

	res, err := c.c.Get(u.String())
	if err != nil {
		return nil, failure.Wrap(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, failure.Wrap(
			errors.New("non-ok status returned"),
			failure.Context{
				"status": res.Status,
				"url":    u.String(),
			})
	}

	var entries []*Entry
	dec := json.NewDecoder(res.Body)
	for {
		var entry Entry
		if err := dec.Decode(&entry); err == io.EOF {
			break
		} else if err != nil {
			return nil, failure.Wrap(err)
		}
		entries = append(entries, &entry)
	}
	return entries, nil
}
