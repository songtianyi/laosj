## laosj(老司机)

[![Build Status](https://travis-ci.org/songtianyi/laosj.svg?branch=master)](https://travis-ci.org/songtianyi/laosj)
[![Go Report Card](https://goreportcard.com/badge/github.com/songtianyi/laosj)](https://goreportcard.com/report/github.com/songtianyi/laosj)
[![codebeat badge](https://codebeat.co/badges/c05ec05d-e902-4091-b5e0-c1656f88ae3c)](https://codebeat.co/projects/github-com-songtianyi-laosj)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

[![logo](https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTaiDDQDv9P90h7lu9jARb1O8i6hmVMpgEuK9qY57l0CZjRVue2)](https://github.com/songtianyi/laosj)


基于goquery的轻量级爬虫, 支持分布式爬取和下载。
### 展示
![laosj-demo](http://owm6k6w0y.bkt.clouddn.com/laosj-demo.gif)

### CLI

```shell
go install github.com/songtianyi/laosj
```

### 图片源

* aiss(已不可用)

  ```
  ./laosj help aiss
  ./laosj aiss 
  ```

* douban相册

  ```shell
  ./laosj help douban
  ./laosj douban --sp 1
  ```

* [妹子图](http://meizitu.com/)(待重构)

* [javlibrary](http://www.javlibrary.com/cn/)(待重构)
> 可以直接下载Release的二进制文件使用

### 代码上手

###### 下载
```shell
go get -u -v github.com/songtianyi/laosj
```

###### 安装redis
	略，使用redis作为下载队列需安装

###### golang.org/x依赖安装
```shell
mkdir $GOPATH/src/golang.org/x
cd $GOPATH/src/golang.org/x
git clone https://github.com/golang/net.git
```

###### 编译并运行

```shell
cd cmds/laosj/ && go build .
./laosj douban --sp 1
```

### 微信交流群

### <img src="http://owm6k6w0y.bkt.clouddn.com/17-9-21/70665214.jpg" width="480" height="480"/>