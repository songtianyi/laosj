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
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/songtianyi/laosj/downloader"
	"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/logs"
)

var (
	oss = "http://com-pmkoo-img.oss-cn-beijing.aliyuncs.com/picture/"
)

const (
	AISS_DEFAULT_WAITING_QUEUE = downloader.WAITTING_KEY_PREFIX + ":AISS"
)

type ReqBody struct {
	page   int
	userId int
}

type Aiss struct {
	// set when create
	name string // source name
	urls chan downloader.Url
	dq   string // destination queue

	// internal use
	sema chan struct{}
	max  int // max images to get every req
}

func NewAiss(name string, dq string, limit int) SourceWrapper {
	return &Aiss{
		dq:   dq,
		sema: make(chan struct{}, limit),
		max:  10,
		name: name,
	}
}

func (s *Aiss) getSuiteList(page int) ([]byte, error) {
	uri := "http://api.pmkoo.cn/aiss/suite/suiteList.do"
	para := "page=" + strconv.FormatInt(int64(page), s.max) + "&userId=153044"
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, strings.NewReader(para))
	req.Header.Add("Host", "api.pmkoo.cn")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Language", "zh-Haq=1, en-CN;q=0.9")
	req.Header.Add("User-Agent", "aiss/1.0 (iPhone; iOS 10.2; Scale/2.00)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func (s *Aiss) doOnce(pg int) error {
	defer func() { <-s.sema }() // release
	b, err := s.getSuiteList(pg)
	if err != nil {
		return err
	}
	jc, _ := rrconfig.LoadJsonConfigFromBytes(b)
	// du, _ := jc.Dump()
	// fmt.Println(du)
	ics, err := jc.GetInterfaceSlice("data.list")
	if err != nil {
		return &SourceEOF{}
	}
	for _, v := range ics {
		vm := v.(map[string]interface{})
		vsource := vm["source"].(map[string]interface{})
		catlog := vsource["catalog"].(string)
		pictureCount := int(vm["pictureCount"].(float64))
		issue := int(vm["issue"].(float64))
		for j := 0; j < pictureCount; j++ {
			uri := oss + catlog + "/"
			uri += strconv.FormatInt(int64(issue), 10) + "/"
			uri += strconv.FormatInt(int64(j), 10) + ".jpg"
			s.urls <- downloader.Url{
				V: uri,
			}
		}
	}
	return nil
}
func (s *Aiss) waitCloser() {
	tick := time.Tick(30 * time.Second)
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
func (s *Aiss) GetOne() {
	s.sema <- struct{}{}
	s.doOnce(1) // 1-based
	s.waitCloser()
}
func (s *Aiss) GetAll() {
	page := 1 // 1-based
	ok := true

	for ok {
		select {
		case s.sema <- struct{}{}:
			go func(pg int) {
				if err := s.doOnce(pg); err != nil {
					ok = false
				}
			}(page)
			page++
		}
	}
	s.waitCloser()
}

func (s *Aiss) SetReceiver(c chan downloader.Url) {
	s.urls = c
}
func (s *Aiss) Receiver() chan downloader.Url {
	return s.urls
}

func (s *Aiss) Destination() string {
	return s.dq
}

func (s *Aiss) Name() string {
	return s.name
}
