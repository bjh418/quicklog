package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"github.com/devxfactor/quicklog/memstore"
	"github.com/devxfactor/quicklog/socket"
)

func main() {
	mstore := memstore.NewMemstore()

	mstore.Errorf(time.Now(), "Log level %s is enabled.", "ERROR")
	mstore.Warnf(time.Now(), "Log level %s is enabled.", "WARN")
	mstore.Notef(time.Now(), "Log level %s is enabled.", "NOTE")
	mstore.Infof(time.Now(), "Log level %s is enabled.", "INFO")
	mstore.Debugf(time.Now(), "Log level %s is enabled.", "DEBUG")
	mstore.Tracef(time.Now(), "Log level %s is enabled.", "TRACE")

	socketName := "./quicklog.sock"

	listener, err := net.Listen("unix", socketName)
	if err != nil {
		log.Fatalf("Error listening on unix:%s: %v", socketName, err)
	}

	f := func(conn net.Conn, e memstore.Entry) error {
		fmt.Fprintf(conn, "%s %-5s %v\n", e.Time().Format("2006-01-02 03:04:05.999-07:00"), e.Level(), e.Line())
		return nil
	}

	serve := func() {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection on unix:%s: %v", socketName, err)
		}

		server := socket.NewServer(conn, mstore)
		server.Run(f)
	}

	go serve()
	go serve()

	go func() {
		for {
			mstore.Log(time.Now(), memstore.Level("INFO"), "another line")
			time.Sleep(500 * time.Millisecond)
		}
	}()

	socket.WaitForShutdownSignal()
	listener.Close()
	os.Exit(0)
}
