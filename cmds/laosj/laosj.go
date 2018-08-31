package main

import (
	"fmt"
	"log"
	"os"
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
}

var appConfig *AppConfig

func init() {
	appConfig = &AppConfig{}
}
func startRealTimeDownloader(source sources.SourceWrapper) {
	d := &downloader.RealtimeDownloader{
		ConcurrencyLimit: appConfig.DClimit,
		UrlChannelFactor: 10,
		Store:            rrstorage.CreateLocalDiskStorage("/data/sexx/" + source.Name()),
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
	fmt.Println(appConfig)
	aissSource := sources.NewAiss("aiss", c.String("dq"), appConfig.CClimit)
	aissSource.SetReceiver(make(chan downloader.Url, 100))
	return dealMode(dealTestOrNot(aissSource))
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
	app.Copyright = "(c) 2018 songtianyi"
	app.Commands = []cli.Command{
		{
			Name:    "aiss",
			Aliases: []string{"aiss"},
			Usage:   "crawl aiss images",
			Action:  aissHandler,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "destination_queue, dq",
					Value: sources.AISS_DEFAULT_WAITING_QUEUE,
					Usage: "aiss default destination queue",
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
			Value:       10,
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
			Usage:       "true for get only on image from source, false for get all images",
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
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	return
}
