package main

import (
	"github.com/devxfactor/quicklog/memstore"
	"time"
	"fmt"
	"net"
	"log"
	"github.com/devxfactor/quicklog/socket"
	"os"
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

	conn1, err := listener.Accept()
	if err != nil {
		log.Fatalf("Error accepting connection on unix:%s: %v", socketName, err)
	}

	conn2, err := listener.Accept()
	if err != nil {
		log.Fatalf("Error accepting connection on unix:%s: %v", socketName, err)
	}

	server1 := socket.NewServer(conn1, mstore)
	server2 := socket.NewServer(conn2, mstore)

	f := func(conn net.Conn, e memstore.Entry) {
		fmt.Fprintf(conn, "%s %-5s %v\n", e.Time().Format("2006-01-02 03:04:05.999-07:00"), e.Level(), e.Line())
	}

	go server1.Run(f)
	go server2.Run(f)

	time.Sleep(1 * time.Second)
	mstore.Errorf(time.Now(), "This isn't shown by the server runs because they only snapshot and not tail")

	socket.WaitForShutdownSignal()
	listener.Close()
	os.Exit(0)
}
