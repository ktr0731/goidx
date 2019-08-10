package index

import (
	"net/url"
	"strconv"
	"time"
)

type query struct {
	since time.Time
	limit int
}

func (q *query) encode() string {
	v := make(url.Values)
	if q.since != (time.Time{}) {
		v["since"] = append(v["since"], q.since.String())
	}
	if q.limit != 0 {
		v["limit"] = append(v["limit"], strconv.Itoa(q.limit))
	}
	return v.Encode()
}

type Query func(*query)

func Since(t time.Time) Query {
	return func(o *query) {
		o.since = t
	}
}

func Limit(n int) Query {
	return func(o *query) {
		o.limit = n
	}
}
