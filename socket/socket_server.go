package socket

import (
	"net"
	"os"
	"os/signal"
	"syscall"
	"github.com/devxfactor/quicklog/memstore"
	"time"
)

type Server interface {
	Run(f func(conn net.Conn, line string) error) error
}

type server struct {
	conn     net.Conn
	memstore memstore.Memstore
}

func NewServer(conn net.Conn, mstore memstore.Memstore) Server {
	server := &server{conn: conn, memstore: mstore}
	return server
}

func WaitForShutdownSignal() {
	sig_ch := make(chan os.Signal, 1)
	signal.Notify(sig_ch, os.Interrupt, syscall.SIGTERM)

	_ = <-sig_ch

	return
}

func (s *server) Run(f func(conn net.Conn, line string) error) error {
	return s.memstore.Tail(func(line string) error {
		err := f(s.conn, line)
		if err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		return nil
	})
}
