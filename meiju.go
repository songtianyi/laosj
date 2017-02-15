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
	"github.com/songtianyi/laosj/spider"
	"github.com/songtianyi/rrframework/logs"
)

func main() {
	uri := "http://bbs.ncar.cc/thread-28825-1-1.html"
	s, err := spider.CreateSpiderFromUrl(uri)
	if err != nil {
		panic(err)
	}
	srcs, _ := s.GetText("div.wp>div.wp.cl>div.pl.bm>table>tbody>tr>td.plc.ptm.pbn.vwthd>h1.ts>span")
	logs.Debug(srcs)
}
