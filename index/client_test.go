package index

import (
	"testing"

	"github.com/k0kubun/pp"
)

func TestClient(t *testing.T) {
	c, err := NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := c.Index(Limit(3))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range entries {
		pp.Println(e)
	}
}
