package repository

import (
	"failed-interview/03/entity"
	c "failed-interview/03/pkg/crawler"
)

type Crawler struct {
	c *c.Crawler
}

// NewcrawlerPG create new repository
func NewCrawler() *Crawler {
	return &Crawler{c: c.NewCrawler()}
}

// List links
func (c *Crawler) List(link []string, timeout uint) (e []*entity.Links) {
	if timeout > 0 {
		c.c.SetTimeout(timeout)
	}

	l := c.c.Parse(link)
	for _, v := range l {
		e = append(e, crawlerLinksToEntityLinks(v))
	}

	return e
}

func crawlerLinksToEntityLinks(cl c.Link) (e *entity.Links) {
	e = &entity.Links{}
	e.URL = cl.URL
	e.Title = cl.Title
	err := ""

	if cl.Err != nil {
		err = cl.Err.Error()
	}

	e.Error = err

	return e
}
