package downloader

type Downloader interface {
	Start()
	Stop()
	WaitCloser()
}

const (
	REALTIME = iota + 1
	REDIS
)

const (
	URL_KEY_PREFIX = "LAOSJ:URLS"
	WAITTING_KEY_PREFIX = URL_KEY_PREFIX + ":WAITTING"
)
