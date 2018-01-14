package main

import (
	"github.com/devxfactor/quicklog/memstore"
	"time"
	"fmt"
)

func main() {
	mstore := memstore.NewMemstore()
	mstore.Info(time.Now(), "Start of log.")
	mstore.Each(func(e memstore.Entry) {
		fmt.Printf("%s %-5s %v\n", e.Time().Format("2006-01-02 03:04:05.999-07:00"), e.Level(), e.Line())
	})
}
