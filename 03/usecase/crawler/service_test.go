package crawler

import (
	"failed-interview/03/infrastructure/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_List(t *testing.T) {
	repo := repository.NewCrawler()
	m := NewService(repo)

	links := []string{"https://easymb.xyz", "https://golang.org"}

	l := m.GetList(links, 6)

	assert.Equal(t, 2, len(l))
}
