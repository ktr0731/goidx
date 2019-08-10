package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/k0kubun/pp"
	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/ktr0731/goidx/index"
)

var (
	limit = flag.Int("limit", 2000, "entry limit")
	since = flag.String("since", "", "since")
)

func main() {
	flag.Parse()

	c, err := index.NewClient(nil)
	if err != nil {
		log.Fatalf("failed to instantiate a new index client")
	}
	queries := []index.Query{index.Limit(*limit)}
	if *since != "" {
		t, err := time.Parse(*since, time.RFC3339)
		if err != nil {
			log.Fatalf("failed to parse -since")
		}
		queries = append(queries, index.Since(t))
	}
	entries, err := c.Index(queries...)
	if err != nil {
		log.Fatalf("failed to get entries")
	}
	idx, err := fuzzyfinder.Find(
		entries,
		func(i int) string {
			return entries[i].Path
		},
		fuzzyfinder.WithPreviewWindow(func(i, _, _ int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintf("Version: %s\nTimestamp: %s\n", entries[i].Version, entries[i].Timestamp)
		}),
	)
	if err != nil {
		log.Fatalf("failed to find an entry")
	}
	pp.Println(entries[idx])
}
