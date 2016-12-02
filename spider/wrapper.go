package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"sync"
	//"github.com/songtianyi/rrframework/connector/redis"
	//"io"
	//"net/http"
	//"os"
	"regexp"
	//"strconv"
	//"strings"
)

const (
	TEXT_LEAF = iota
	HTML_LEAF
)

type Spider struct {
	Rules    []string
	IndexUrl string
	LeafType int

	mu sync.Mutex
}

func (s *Spider) Run() ([]string, error) {
	if s.IndexUrl == "" || len(s.Rules) < 1 {
		return nil, fmt.Errorf("IndexUrl empty or Rules empty")
	}
	// start from level 0
	return s.do(s.IndexUrl, 0)
}

func (s *Spider) do(url string, level int) ([]string, error) {
	var (
		res = make([]string, 0) //for leaf = make
		err error
		wg  sync.WaitGroup
	)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, fmt.Errorf("Parent: url %s, level %d, error %s", url, level, err)
	}

	doc.Find(s.Rules[level]).Each(func(ix int, sl *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if len(s.Rules) > level+1 {
				// there a deep level
				// find herf url
				content, _ := sl.Html()
				m := regexp.MustCompile("href=\"(\\S+)\"").FindStringSubmatch(content)[1]
				t, err := s.do(m, level+1)
				if err != nil {
					return
				}
				s.mu.Lock()
				res = append(res, t...)
				s.mu.Unlock()

			} else {
				// last
				// text or html
				s.mu.Lock()
				if s.LeafType == TEXT_LEAF {
					res = append(res, sl.Text())
				} else if s.LeafType == HTML_LEAF {
					content, _ := sl.Html()
					res = append(res, content)
				}
				s.mu.Unlock()
			}
		}()
	})
	wg.Wait()
	return res, nil
}
