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
	//"github.com/songtianyi/rrframework/config"
	"github.com/songtianyi/rrframework/storage"
	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/rrframework/connector/redis"
)

func main() {
	url := "http://www.douban.com/group/haixiuzu/discussion"
	d := &downloader.Downloader{
		ConcurrencyLimit: 10,
		UrlChannelFactor: 10,
		RedisConnStr: "10.19.147.75:6379",
		SourceQueue: "DATA:IMAGE:HAIXIUZU",
		Store: rrstorage.CreateLocalDiskStorage("/data/sexx/haixiuzu/"),
	}
	err, rc := rrredis.GetRedisClient("10.19.147.75:6379")
	if err != nil {
		panic(err)
	}
	go func() {
		d.Start()
	}()

	for {
		s, err := spider.CreateSpiderFromUrl(url)
		if err != nil {
			logs.Debug(err)
			continue
		}
		rs, _ := s.GetAttr("div.grid-16-8.clearfix>div.article>div>table.olt>tbody>tr>td.title>a", "href")
		for _, v := range rs {
			s01, err := spider.CreateSpiderFromUrl(v)
			if err != nil {
				logs.Error(err)
				continue
			}
			rs01, _ := s01.GetAttr("div.grid-16-8.clearfix>div.article>div.topic-content.clearfix>div.topic-doc>div#link-report>div.topic-content>div.topic-figure.cc>img", "src")
			for _, vv := range rs01 {
				if _, err := rc.RPush("DATA:IMAGE:HAIXIUZU", vv); err != nil {
					logs.Error(err)
				}
			}
		}
		rs1, _ := s.GetAttr("div.grid-16-8.clearfix>div.article>div.paginator>span.next>a", "href")
		if len(rs1) != 1 {
			break
		}
		url = rs1[0]
	}
}
