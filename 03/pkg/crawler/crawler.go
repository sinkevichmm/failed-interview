package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Crawler struct {
	links   []Link
	timeout time.Duration
}

type Link struct {
	URL   string
	Valid bool
	Title string
	Err   error
}

//TODO: сделать проверку на len(links)==0
// NewCrawler создает новый Crawler с таймаутом по умолчанию 5 сек.
func NewCrawler() *Crawler {
	c := &Crawler{timeout: 5 * time.Second, links: []Link{}}
	return c
}

// AddLink добавляет link
func (c *Crawler) addLink(link string) {
	for _, l := range c.links {
		if l.URL == link {
			return
		}
	}

	_, err := url.Parse(link)
	if err != nil {
		c.links = append(c.links, Link{URL: link, Err: ErrInvalidURL})

		return
	}

	c.links = append(c.links, Link{URL: link, Valid: true})
}

// Clear отчищает links
func (c *Crawler) Clear() {
	c.links = []Link{}
}

// SetTimeout устанавливает таймаут в секундах
func (c *Crawler) SetTimeout(t uint) {
	c.timeout = time.Duration(t) * time.Second
}

// Parse заполняет Links
func (c *Crawler) Parse(links []string) []Link {
	m := make(map[string]struct{})

	for _, l := range links {
		m[l] = struct{}{}

		c.addLink(l)
	}

	var wg sync.WaitGroup

	mux := &sync.Mutex{}

	for i := 0; i < len(c.links); i++ {
		if c.links[i].Valid && c.links[i].Title == "" {
			if _, ok := m[c.links[i].URL]; !ok {
				continue
			}

			wg.Add(1)

			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
				defer cancel()

				title, err := c.getPage(ctx, c.links[i].URL)

				mux.Lock()
				defer mux.Unlock()

				c.links[i].Title = title
				c.links[i].Err = err
			}(i, &wg)
		}
	}

	wg.Wait()

	return c.links
}

// getPage запрашивает страницу по ссылке
func (c *Crawler) getPage(ctx context.Context, url string) (title string, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		err = fmt.Errorf("%w: %s", ErrInternal, err)
		return
	}

	cc := &http.Client{}
	resp, err := cc.Do(req)

	if err != nil {
		err = fmt.Errorf("%w: %s", ErrRequest, err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		err = fmt.Errorf("%w: %s", ErrResponse, resp.Status)

		return
	}

	return c.parsePage(html.NewTokenizer(resp.Body))
}

// parsePage ищет первый тег title
func (c *Crawler) parsePage(z *html.Tokenizer) (string, error) {
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return "", fmt.Errorf("%w: %s", ErrParse, z.Err())
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "title" {
				z.Next()
				t = z.Token()

				return t.Data, nil
			}
		}
	}
}
