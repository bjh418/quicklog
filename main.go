package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
	"github.com/devxfactor/quicklog/memstore"
	"github.com/devxfactor/quicklog/socket"
)

func ServeLoggers(socketName string, mstore memstore.Memstore) {
	sock, err := net.Listen("unix", socketName)
	if err != nil {
		log.Fatalf("Error listening on unix:%s: %v", socketName, err)
	}

	for {
		conn, err := sock.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection on unix:%s: %v", socketName, err)
		}

		if unixConn, ok := conn.(*net.UnixConn); ok {
			if err := unixConn.CloseWrite(); err != nil {
				fmt.Errorf("Error on CloseWrite of unix:%v: %v", socketName, err)
			}
		}

		// start reader (for logger)
		go func() {
			reader := bufio.NewReader(conn)
			for {
				line, err := reader.ReadString('\n')
				if len(line) > 0 {
					mstore.Log(time.Now(), memstore.Level("INFO"), line[0:len(line) - 1])
				}
				if err != nil {
					if err != io.EOF {
						fmt.Errorf("Error reading unix:%s: %v", socketName, err)
					}
					conn.Close()
					return
				}
			}
		}()
	}
}

func ServeTailers(socketName string, mstore memstore.Memstore) {
	sock, err := net.Listen("unix", socketName)
	if err != nil {
		log.Fatalf("Error listening on unix:%s: %v", socketName, err)
	}

	for {
		conn, err := sock.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection on unix:%s: %v", socketName, err)
		}

		if unixConn, ok := conn.(*net.UnixConn); ok {
			if err := unixConn.CloseRead(); err != nil {
				fmt.Errorf("Error on CloseRead of unix:%v: %v", socketName, err)
			}
		}

		// start writer (for tailer)
		go func() {
			writer := bufio.NewWriter(conn)
			err := mstore.Tail(func(e memstore.Entry) error {
				line := fmt.Sprintf("%s %-5s %v\n", e.Time().Format("2006-01-02 03:04:05.000-07:00"), e.Level(), e.Line())
				nWritten := 0
				for nWritten < len(line) {
					nWrote, err := writer.Write([]byte(line[nWritten:]))
					if err != nil {
						return err
					}
					nWritten += nWrote
				}
				writer.Flush()
				return nil
			})
			if err != nil {
				if err != io.EOF {
					fmt.Errorf("Error writing unix:%s: %v", socketName, err)
				}
				conn.Close()
				return
			}
		}()
	}
}

func main() {
	mstore := memstore.NewMemstore()

	mstore.Errorf(time.Now(), "Log level %s is enabled.", "ERROR")
	mstore.Warnf(time.Now(), "Log level %s is enabled.", "WARN")
	mstore.Notef(time.Now(), "Log level %s is enabled.", "NOTE")
	mstore.Infof(time.Now(), "Log level %s is enabled.", "INFO")
	mstore.Debugf(time.Now(), "Log level %s is enabled.", "DEBUG")
	mstore.Tracef(time.Now(), "Log level %s is enabled.", "TRACE")

	go ServeLoggers("./quicklogger.sock", mstore)
	go ServeTailers("./quicktailer.sock", mstore)

	socket.WaitForShutdownSignal()
	os.Exit(0)
}
