// Copyright 2016 laosj Author @songtianyi. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/songtianyi/rrframework/logs"
	"regexp"
	"sync"
)

const (
	TEXT_LEAF = iota // content type goquery.Selection.Text()
	HTML_LEAF        // content type goquery.Selection.Html()
)

// Spider
type Spider struct {
	Rules    []string // goquery rules
	IndexUrl string   // first page that spider would deal with
	LeafType int      // return Text() or Html()

	mu sync.Mutex
}

// Start spider
func (s *Spider) Run() ([]string, error) {
	if s.IndexUrl == "" || len(s.Rules) < 1 {
		return nil, fmt.Errorf("IndexUrl empty or Rules empty")
	}
	// start from level 0
	return s.do(s.IndexUrl, 0)
}

func (s *Spider) do(url string, level int) ([]string, error) {
	var (
		res = make([]string, 0) //for leaf
		err error
		wg  sync.WaitGroup
	)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, fmt.Errorf("url %s, level %d, error %s", url, level, err)
	}

	doc.Find(s.Rules[level]).Each(func(ix int, sl *goquery.Selection) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if len(s.Rules) > level+1 {
				// there a deeper level page
				// find herf url
				content, _ := sl.Html()
				m := regexp.MustCompile("href=\"(\\S+)\"").FindStringSubmatch(content)
				if len(m) < 2 {
					logs.Error("Find href error, %s", "len(m) < 2")
					return
				}
				href := m[1]
				t, err := s.do(href, level+1)
				if err != nil {
					logs.Error(err)
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
