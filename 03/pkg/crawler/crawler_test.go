package crawler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawler(t *testing.T) {
	c := NewCrawler()

	c.SetTimeout(6)

	l := c.Parse([]string{"https://easymb.xyz"})

	assert.Equal(t, 1, len(l))
	assert.Equal(t, "EasyMB | ModBus - это просто!", l[0].Title)

	l = c.Parse([]string{"https://easymb.xyz"})

	assert.Equal(t, 1, len(l))
	assert.Equal(t, "EasyMB | ModBus - это просто!", l[0].Title)

	c.Clear()
	assert.Equal(t, 0, len(c.links))

	c.Clear()

	l = c.Parse([]string{"https://easymb.xyz\n/index.html"})
	assert.ErrorIs(t, l[0].Err, ErrInvalidURL)

	c.Clear()
	l = c.Parse([]string{"https://easymb.xyz/343534"})
	assert.ErrorIs(t, l[0].Err, ErrResponse)

	c.Clear()

	l = c.Parse([]string{"https://easymb.xyzsas"})
	assert.ErrorIs(t, l[0].Err, ErrRequest)

	c.Clear()

	l = c.Parse([]string{"https://easymb.xyz", "https://golang.org"})

	for _, v := range l {
		fmt.Printf("l: %+v\n", v)
	}
}
