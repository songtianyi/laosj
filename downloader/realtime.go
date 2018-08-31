package downloader

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/songtianyi/rrframework/logs"
	"github.com/songtianyi/rrframework/storage"
)

type RealtimeDownloader struct {
	// exported
	ConcurrencyLimit int                      // max number of goroutines to download
	SourceQueue      string                   // url queue
	Store            rrstorage.StorageWrapper // for saving downloaded binary
	UrlChannelFactor int
	Urls             chan Url

	// inner use
	sema chan struct{} // for concurrency-limiting
	flag chan struct{} // stop flag
}

func (s *RealtimeDownloader) Start() {
	// create channel
	s.sema = make(chan struct{}, s.ConcurrencyLimit)
	s.flag = make(chan struct{})
	tick := time.Tick(2 * time.Second)
	logs.Info("realtime downloader started.")

loop2:
	for {
		select {
		case <-s.flag:
			// be stopped
			for url := range s.Urls {
				// log undownloaded urls
				logs.Info("url", url.V, "droped")
			}
			// end RealtimeDownloader
			break loop2
		case s.sema <- struct{}{}:
			// s.sema not full
			url, ok := <-s.Urls
			if !ok {
				// channel closed
				logs.Alert("Channel s.Urls closed")
				// TODO what's the right way to deal this situation?
				break loop2
			}
			go func() {
				if err := s.download(url.V); err != nil {
					// download fail
					// push back to redis
					logs.Error("Download %s fail, %s", url.V, err)
				} else {
					// download success
				}
			}()
		case <-tick:
			// print this every 2 seconds
			logs.Info("In queue: %d, doing: %d", len(s.Urls), len(s.sema))
		}
	}

}

// Stop RealtimeDownloader
func (s *RealtimeDownloader) Stop() {
	close(s.flag)
}

// Wait all urls in redis queue be downloaded
func (s *RealtimeDownloader) WaitCloser() {
loop:
	for {
		select {
		case <-time.After(1 * time.Second):
			// len
			if len(s.Urls) > 0 || len(s.sema) > 1 {
				// TODO there is a chance that last url downloading process be interupted
				continue
			}
			break loop
		}
	}
}

func (s *RealtimeDownloader) download(url string) error {

	defer func() { <-s.sema }() // release

	logs.Info("Downloading %s", url)
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) { return net.DialTimeout(network, addr, 3*time.Second) },
		},
	}
	response, err := client.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return fmt.Errorf("StatusCode %d", response.StatusCode)
	}

	// read binary from body
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	urlv := strings.Split(url, "/")
	if len(urlv) < 1 {
		return fmt.Errorf("invalid url %s", url)
	}
	filename := urlv[len(urlv)-1]
	// save binary to storage
	if err := s.Store.Save(b, filename); err != nil {
		return err
	}
	return nil
}
