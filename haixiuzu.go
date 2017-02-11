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
	"fmt"
	//"github.com/songtianyi/laosj/downloader"
	"github.com/songtianyi/laosj/spider"
	//"github.com/songtianyi/rrframework/config"
	//"github.com/songtianyi/rrframework/connector/redis"
)

func main() {
	s, err := spider.CreateSpiderFromUrl("http://www.douban.com/group/haixiuzu/discussion")
	if err != nil {
		fmt.Println(err)
	}
	rs, _ := s.GetAttr("div.grid-16-8.clearfix>div.article>div>table.olt>tbody>tr>td.title>a", "href")
	fmt.Println(rs)
}
