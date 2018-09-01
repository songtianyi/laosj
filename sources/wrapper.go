package sources

import (
	"errors"

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

// end of source error
var EOS = errors.New("End of source")
