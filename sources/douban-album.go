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
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/songtianyi/laosj/downloader"
	"github.com/songtianyi/laosj/spider"
	"github.com/songtianyi/rrframework/logs"
)

const (
	userAgent                   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36"
	albumPrefix                 = "https://www.douban.com/photos/album/"
	DOUBAN_ALBUM_WAITTING_QUEUE = downloader.WAITTING_KEY_PREFIX + ":DOUBAN"
)

type DoubanAlbum struct {
	name  string              // source name
	album string              // album id
	ps    int                 // page size
	sp    int                 // start page
	lp    int                 // last page
	dq    string              // destination queue
	urls  chan downloader.Url // url channel

	// internal use
	sema  chan struct{}
	refer string // header Refer
}

func NewDoubanAlbum(name string, album string, ps int, sp int, lp int, dq string, climit int) SourceWrapper {
	return &DoubanAlbum{
		name:  name,
		album: album,
		ps:    ps,
		sp:    sp,
		lp:    lp,
		dq:    dq,
		sema:  make(chan struct{}, climit),
		refer: "",
	}
}
func (s *DoubanAlbum) GetOne() {
	s.sema <- struct{}{}
	if err := s.doOnce(s.sp); err != nil {
		logs.Error(err)
	}
	s.waitCloser()
}

func (s *DoubanAlbum) GetAll() {
	page := s.sp
	ok := true

	for ok {
		select {
		case s.sema <- struct{}{}:
			go func(pg int) {
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

func (s *DoubanAlbum) doOnce(page int) error {
	defer func() { <-s.sema }() // release
	if page > s.lp {
		return EOS
	}
	startV := strconv.Itoa(page * s.ps)
	url := albumPrefix + "/" + s.album + "/?start=" + startV
	logs.Debug("url", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Referer", s.refer)
	// update refer
	s.refer = url
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("response code: %d", resp.StatusCode)
	}
	spi, err := spider.CreateSpiderFromResponse(resp)
	if err != nil {
		return err
	}
	rs, _ := spi.GetAttr("div.grid-16-8.clearfix>div.article>div.photolst.clearfix>div.photo_wrap>a>img", "src")
	for _, v := range rs {
		s.urls <- downloader.Url{
			V: v,
		}
	}
	return nil
}

func (s *DoubanAlbum) Destination() string {
	return s.dq
}

func (s *DoubanAlbum) SetReceiver(c chan downloader.Url) {
	s.urls = c
}

func (s *DoubanAlbum) Receiver() chan downloader.Url {
	return s.urls
}

func (s *DoubanAlbum) Name() string {
	return s.name
}

func (s *DoubanAlbum) waitCloser() {
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
