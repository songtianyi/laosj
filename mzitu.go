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

package main

import (
	"github.com/songtianyi/laosj/downloader"
	"github.com/songtianyi/laosj/spider"
	"github.com/songtianyi/laosj/storage"
	"github.com/songtianyi/rrframework/connector/redis"
	"github.com/songtianyi/rrframework/logs"
	"regexp"
	"strconv"
	"sync"
)

func main() {
	d := &downloader.Downloader{
		ConcurrencyLimit: 10,
		UrlChannelFactor: 10,
		RedisConnStr:     "127.0.0.1:6379",
		SourceQueue:      "DATA:IMAGE:MZITU:XINGGAN",
		Store:            storage.NewLocalDiskStorage("/data/sexx/taiwan/"),
	}
	go func() {
		d.Start()
	}()

	// step1: find total index pages
	s := &spider.Spider{
		IndexUrl: "http://www.mzitu.com/taiwan",
		Rules: []string{
			"div.main>div.main-content>div.postlist>nav.navigation.pagination>div.nav-links>a.page-numbers",
		},
		LeafType: spider.TEXT_LEAF,
	}
	rs, err := s.Run()
	if err != nil {
		logs.Error(err)
		return
	}
	max := spider.FindMaxFromSliceString(1, rs)

	// step2: for every index page, find every post entrance
	var wg sync.WaitGroup
	var mu sync.Mutex
	step2 := make([]string, 0)
	for i := 1; i <= max; i++ {
		wg.Add(1)
		go func(ix int) {
			defer wg.Done()
			ns := &spider.Spider{
				IndexUrl: s.IndexUrl + "/page/" + strconv.Itoa(ix),
				Rules: []string{
					"div.main>div.main-content>div.postlist>ul>li",
				},
				LeafType: spider.HTML_LEAF,
			}
			t, err := ns.Run()
			if err != nil {
				logs.Error(err)
				return
			}
			mu.Lock()
			step2 = append(step2, t...)
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	// parse url
	for i, v := range step2 {
		re := regexp.MustCompile("href=\"(\\S+)\"")
		url := re.FindStringSubmatch(v)[1]
		step2[i] = url
	}

	for _, v := range step2 {
		// step3: step in entrance, find max pagenum
		ns1 := &spider.Spider{
			IndexUrl: v,
			Rules: []string{
				"div.main>div.content>div.pagenavi>a",
			},
			LeafType: spider.TEXT_LEAF,
		}
		t1, err := ns1.Run()
		if err != nil {
			logs.Error(err)
			return
		}
		maxx := spider.FindMaxFromSliceString(1, t1)
		// step4: for every page
		for j := 1; j <= maxx; j++ {

			// step5: find img in this page
			ns2 := &spider.Spider{
				IndexUrl: v + "/" + strconv.Itoa(j),
				Rules: []string{
					"div.main>div.content>div.main-image>p>a",
				},
				LeafType: spider.HTML_LEAF,
			}
			t2, err := ns2.Run()
			if err != nil {
				logs.Error(err)
				return
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
			err, rc := rrredis.GetRedisClient(d.RedisConnStr)
			if err != nil {
				logs.Error(err)
				return
			}
			key := d.SourceQueue
			if _, err := rc.RPush(key, sub[1]); err != nil {
				logs.Error(err)
				return
			}
		}
	}
	d.WaitCloser()
}
