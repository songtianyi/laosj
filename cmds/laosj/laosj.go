package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/songtianyi/laosj/downloader"
	"github.com/songtianyi/laosj/sources"

	"github.com/urfave/cli"
)

type AppConfig struct {
	Climit int
	All    bool
	Mode   int
	Redis  string
}

var appConfig *AppConfig

func init() {
	appConfig = &AppConfig{}
}

func dealMode(urls []string) {
	switch appConfig.Mode {
	case downloader.REALTIME:
		fmt.Println(urls)
		break
	case downloader.REDIS:
		fmt.Println(urls)
		break
	}
}

func dealTestOrNot(source sources.SourceWrapper) []string {
	if appConfig.All {
		return source.GetAll()
	} else {
		return source.GetOne()
	}
}
func aissHandler(c *cli.Context) error {
	fmt.Println(appConfig)
	aissSource := sources.NewAiss(c.Int("c"))
	dealMode(dealTestOrNot(aissSource))
	return nil
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
	app.Copyright = "(c) 1999 Serious Enterprise"
	app.Commands = []cli.Command{
		{
			Name:    "aiss",
			Aliases: []string{"aiss"},
			Usage:   "crawl aiss images",
			Action:  aissHandler,
		},
		{
			Name:    "drain",
			Aliases: []string{"drain"},
			Usage:   "drain redis url queue",
			Action:  drainHandler,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "queue, q",
					Value: downloader.URL_KEY_PREFIX,
					Usage: "key for url queue",
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "climit, c",
			Value:       10,
			Usage:       "concurrency limit for crawling, used when getting all images from source site",
			Destination: &appConfig.Climit,
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
