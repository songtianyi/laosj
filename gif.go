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
	//"github.com/songtianyi/laosj/downloader"
	"flag"
	"github.com/songtianyi/laosj/spider"
	"github.com/songtianyi/rrframework/logs"
	"net/url"
)

var (
	content = flag.String("c", "", "search content")
)

func main() {
	flag.Parse()
	if *content == "" {
		logs.Error("please input the content you wanna search")
		return
	}
	key := *content
	uri := "http://www.gifmiao.com/search/" + url.QueryEscape(key) + "/3"
	s, err := spider.CreateSpiderFromUrl(uri)
	if err != nil {
		logs.Error(err)
		return
	}
	srcs, _ := s.GetAttr("div.wrap>div#main>ul#waterfall>li.item>div.img_block>a>img.gifImg", "xgif")
	logs.Debug(srcs)
}
