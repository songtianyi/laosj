package sources

import (
	"github.com/songtianyi/laosj/downloader"
)

// images sources

type SourceWrapper interface {
	GetOne()
	GetAll()
	SetReceiver(chan downloader.Url)
	Receiver() chan downloader.Url
	Destination() string
	Name() string
}
