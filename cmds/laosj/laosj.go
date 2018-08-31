package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/rrframework/storage"

	"github.com/songtianyi/rrframework/connector/redis"

	"github.com/songtianyi/laosj/downloader"
	"github.com/songtianyi/laosj/sources"

	"github.com/urfave/cli"
)

type AppConfig struct {
	CClimit int
	DClimit int
	All     bool
	Mode    int
	Redis   string
	Dir     string
}

var appConfig *AppConfig

func init() {
	appConfig = &AppConfig{}
}
func startRealTimeDownloader(source sources.SourceWrapper) {
	d := &downloader.RealtimeDownloader{
		ConcurrencyLimit: appConfig.DClimit,
		UrlChannelFactor: 10,
		Store:            rrstorage.CreateLocalDiskStorage(strings.TrimSuffix(appConfig.Dir, "/") + "/" + source.Name() + "/"),
		Urls:             source.Receiver(),
	}
	d.Start()
}
func dealMode(source sources.SourceWrapper) error {
	switch appConfig.Mode {
	case downloader.REALTIME:
		startRealTimeDownloader(source)
		break
	case downloader.REDIS:
		// connect to redis
		err, rc := rrredis.GetRedisClient(appConfig.Redis)
		if err != nil {
			return err
		}
		for v := range source.Receiver() {
			if _, err := rc.RPush(source.Destination(), v.V); err != nil {
				logs.Error("push", v.V, "to", source.Destination(), "failed")
			}
		}
		break
	}
	return nil
}

func dealTestOrNot(source sources.SourceWrapper) sources.SourceWrapper {
	if appConfig.All {
		go func() {
			source.GetAll()
		}()
	} else {
		go func() {
			source.GetOne()
		}()
	}
	return source
}
func aissHandler(c *cli.Context) error {
	aissSource := sources.NewAiss(c.String("sub"), c.String("dq"), appConfig.CClimit)
	aissSource.SetReceiver(make(chan downloader.Url, 100))
	return dealMode(dealTestOrNot(aissSource))
}

func doubanAlbumHandler(c *cli.Context) error {
	doubanAlbumSource := sources.NewDoubanAlbum(
		c.String("sub"),
		c.String("id"),
		c.Int("ps"),
		c.Int("sp"),
		c.Int("lp"),
		c.String("dq"),
		appConfig.CClimit)
	logs.Debug(doubanAlbumSource)
	doubanAlbumSource.SetReceiver(make(chan downloader.Url, 100))
	return dealMode(dealTestOrNot(doubanAlbumSource))
}

func drainHandler(c *cli.Context) error {
	return nil
}

func main() {
	app := cli.NewApp()
	app.Usage = "A cli tool to crawl images"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "songtianyi",
			Email: "songtianyi630@163.com",
		},
	}
	app.Copyright = "Copyright (c) 2016-2018 songtianyi"
	app.Commands = []cli.Command{
		{
			Name:    "aiss",
			Aliases: []string{"aiss"},
			Usage:   "crawl aiss images",
			Action:  aissHandler,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "subdirectory, sub",
					Value: "aiss",
					Usage: "subdir, storage sub dir counting on it",
				},
				cli.StringFlag{
					Name:  "destination_queue, dq",
					Value: sources.AISS_DEFAULT_WAITING_QUEUE,
					Usage: "aiss default destination queue",
				},
			},
		},
		{
			Name:    "douban",
			Aliases: []string{"douban"},
			Usage:   "crawl douban album images",
			Action:  doubanAlbumHandler,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "subdirectory, sub",
					Value: "douban",
					Usage: "subdir, storage sub dir counting on it",
				},
				cli.IntFlag{
					Name:  "page_size, ps",
					Value: 18, // default
					Usage: "douban album page size",
				},
				cli.IntFlag{
					Name:  "last_page, lp",
					Value: 254,
					Usage: "douban album last page number, include itself",
				},
				cli.IntFlag{
					Name:  "start_page, sp",
					Value: 1, // from first page
					Usage: "set douban album start page number, include itself",
				},
				cli.StringFlag{
					Name:  "album_id, id",
					Value: "105181925",
					Usage: "douban album id",
				},
				cli.StringFlag{
					Name:  "destination_queue, dq",
					Value: sources.DOUBAN_ALBUM_WAITTING_QUEUE,
					Usage: "douban album default destination queue",
				},
			},
		},
		{
			Name:    "drain",
			Aliases: []string{"drain"},
			Usage:   "drain redis url queue",
			Action:  drainHandler,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "source_queue, sq",
					Value: downloader.URL_KEY_PREFIX,
					Usage: "key for url queue",
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "cclimit, ccl",
			Value:       1,
			Usage:       "concurrency limit for crawling, used when getting all images from source site",
			Destination: &appConfig.CClimit,
		},
		cli.IntFlag{
			Name:        "dclimit, dcl",
			Value:       3,
			Usage:       "concurrency limit for downloading",
			Destination: &appConfig.DClimit,
		},
		cli.BoolFlag{
			Name:        "all, a",
			Usage:       "false for get only on image from source, true for get all images",
			Destination: &appConfig.All,
		},
		cli.IntFlag{
			Name:        "mode, m",
			Value:       downloader.REALTIME,
			Usage:       "choose download mode, realtime downloading or put url into redis queue",
			Destination: &appConfig.Mode,
		},
		cli.StringFlag{
			Name:        "redis, r",
			Value:       "127.0.0.1:6379",
			Usage:       "redis ip:port",
			Destination: &appConfig.Redis,
		},
		cli.StringFlag{
			Name:        "directory, dir",
			Value:       "/Volumes/songtianyi/sexx",
			Usage:       "the local disk storage path prefix, no slash in the end",
			Destination: &appConfig.Dir,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	return
}
