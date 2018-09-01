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

package sources

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/songtianyi/laosj/downloader"
	"github.com/songtianyi/laosj/spider"
	"github.com/songtianyi/rrframework/logs"
)

const (
	MZITU_WAITTING_QUEUE = downloader.WAITTING_KEY_PREFIX + ":MZITU"
	groupPrefix          = "http://www.mzitu.com/"
)

type Mzitu struct {
	name string              // source name
	sg   int                 // start group
	eg   int                 // end group
	dq   string              // destination queue
	urls chan downloader.Url // url channel

	// internal use
	sema  chan struct{}
	refer string // header Refer
}

func NewMzitu(name string, sg int, eg int, dq string, climit int) SourceWrapper {
	return &Mzitu{
		name:  name,
		sg:    sg,
		eg:    eg,
		dq:    dq,
		sema:  make(chan struct{}, climit),
		refer: "",
	}
}
func (s *Mzitu) GetOne() {
	s.sema <- struct{}{}
	if err := s.doOnce(s.sg); err != nil {
		logs.Error(err)
	}
	s.waitCloser()
}

func (s *Mzitu) GetAll() {
	page := s.sg
	ok := true

	for ok {
		select {
		case s.sema <- struct{}{}:
			go func(pg int) {
				logs.Debug("trying page", page)
				if err := s.doOnce(pg); err != nil {
					logs.Error(err)
					if err == EOS {
						// end
						ok = false
					} // continue
				}
			}(page)
			page++
		}
	}
	s.waitCloser()
}

func (s *Mzitu) doOnce(page int) error {
	defer func() { <-s.sema }() // release
	if page > s.eg {
		return EOS
	}
	url := groupPrefix + "/" + strconv.Itoa(page)
	ns1, err := spider.CreateSpiderFromUrl(url)
	if err != nil {
		return err
	}
	t1, err := ns1.GetText("div.main>div.content>div.pagenavi>a")
	if err != nil {
		return err
	}
	maxx := spider.FindMaxFromSliceString(1, t1)
	// step2: for every page
	for j := 1; j <= maxx; j++ {
		// step3: find img in this page
		subUrl := url + "/" + strconv.Itoa(j)
		ns2, err := spider.CreateSpiderFromUrl(subUrl)
		if err != nil {
			return err
		}
		t2, err := ns2.GetHtml("div.main>div.content>div.main-image>p>a")
		if err != nil {
			return err
		}
		if len(t2) < 1 {
			// ignore this page
			continue
		}
		sub := regexp.MustCompile("src=\"(\\S+)\"").FindStringSubmatch(t2[0])
		if len(sub) != 2 {
			// ignore this page
			continue
		}
		s.urls <- downloader.Url{
			V: sub[1],
			Header: http.Header{
				"Referer": []string{subUrl},
			},
		}
	}
	return nil
}

func (s *Mzitu) Destination() string {
	return s.dq
}

func (s *Mzitu) SetReceiver(c chan downloader.Url) {
	s.urls = c
}

func (s *Mzitu) Receiver() chan downloader.Url {
	return s.urls
}

func (s *Mzitu) Name() string {
	return s.name
}

func (s *Mzitu) waitCloser() {
	tick := time.Tick(6 * time.Second)
	logs.Alert("closing source url channel...")
loop:
	for {
		select {
		case <-tick:
			if len(s.urls) < 1 {
				close(s.urls)
				break loop
			}
			break
		}
	}
}
