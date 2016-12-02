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
		Store:            storage.NewLocalDiskStorage("/data/sexx/mzituzp/"),
	}
	go func() {
		d.Start()
	}()

	// step1: find total pages
	s := &spider.Spider{
		IndexUrl: "http://www.mzitu.com/share",
		Rules: []string{
			"div.main>div.main-content>div.postlist>div>div.pagenavi-cm>a",
		},
		LeafType: spider.TEXT_LEAF,
	}
	rs, err := s.Run()
	if err != nil {
		logs.Error(err)
		return
	}
	max := spider.FindMaxFromSliceString(1, rs)

	// step2: for every page, find all img tags
	var wg sync.WaitGroup
	var mu sync.Mutex
	step2 := make([]string, 0)
	for i := 1; i <= max; i++ {
		wg.Add(1)
		go func(ix int) {
			defer wg.Done()
			ns := &spider.Spider{
				IndexUrl: s.IndexUrl + "/comment-page-" + strconv.Itoa(ix) + "#comments/",
				Rules: []string{
					"div.main>div.main-content>div.postlist>div>ul>li>div>p",
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
	err, rc := rrredis.GetRedisClient(d.RedisConnStr)
	if err != nil {
		logs.Error(err)
		return
	}
	// parse url
	for _, v := range step2 {
		re := regexp.MustCompile("src=\"(\\S+)\"")
		url := re.FindStringSubmatch(v)[1]
		key := d.SourceQueue
		if _, err := rc.RPush(key, url); err != nil {
			logs.Error(err)
			return
		}
	}
	d.WaitCloser()
}
