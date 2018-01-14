package main

import (
	"github.com/devxfactor/quicklog/memstore"
	"time"
	"fmt"
)

func main() {
	mstore := memstore.NewMemstore()

	mstore.Errorf(time.Now(), "Log level %s is enabled.", "ERROR")
	mstore.Warnf(time.Now(), "Log level %s is enabled.", "WARN")
	mstore.Notef(time.Now(), "Log level %s is enabled.", "NOTE")
	mstore.Infof(time.Now(), "Log level %s is enabled.", "INFO")
	mstore.Debugf(time.Now(), "Log level %s is enabled.", "DEBUG")
	mstore.Tracef(time.Now(), "Log level %s is enabled.", "TRACE")

	mstore.Each(func(e memstore.Entry) {
		fmt.Printf("%s %-5s %v\n", e.Time().Format("2006-01-02 03:04:05.999-07:00"), e.Level(), e.Line())
	})
}
