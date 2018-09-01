package downloader

import (
	"net/http"
)

type Downloader interface {
	Start()
	Stop()
	WaitCloser()
}

type Url struct {
	V      string
	Header http.Header
}

// RedisDownloader get urls from redis SourceQueue

const (
	REALTIME = iota + 1
	REDIS
)

const (
	URL_KEY_PREFIX      = "LAOSJ:URLS"
	WAITTING_KEY_PREFIX = URL_KEY_PREFIX + ":WAITTING"
)
