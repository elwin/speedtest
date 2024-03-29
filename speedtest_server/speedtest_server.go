package main

import (
	"crypto/rand"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"

	header2 "github.com/elwin/speedtest/header"

	"github.com/elwin/scionFTP/scion"
)

var (
	local = flag.String("local", "", "Local address (with Port)")
)

func main() {
	flag.Parse()
	if *local == "" {
		log.Fatal("Please specify the local address using -local")
	}

	listener, err := scion.Listen(*local)
	if err != nil {
		log.Fatal("failed to listen", err)
	}

	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed to accept connection", err)
		}

		go handleConnection(conn)

	}
}

func handleConnection(conn *scion.Connection) {
	defer func() {
		err := conn.Close()
		if err != nil {
			fmt.Println("failed to close conn", err)
		}
	}()

	decoder := gob.NewDecoder(conn)
	header := header2.Header{}
	if err := decoder.Decode(&header); err != nil {
		fmt.Println("failed to decode header", err)
		return
	}

	for i := 0; i < header.Repetitions; i++ {

		if _, err := io.CopyN(conn, rand.Reader, int64(header.Size)); err != nil {
			fmt.Println("failed to send payload", err)
			return
		}
	}

}

var _ io.Reader = EmptyReader{}

type EmptyReader struct{}

func (EmptyReader) Read(p []byte) (n int, err error) {
	return len(p), nil
}
