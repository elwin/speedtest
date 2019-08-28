package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/elwin/transmit2/scion"
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
	defer conn.Close()
	for {
		header := make([]byte, 4)
		if err := binary.Read(conn, binary.BigEndian, &header); err != nil {
			fmt.Println("failed to read header", err)
			return
		}

		length, err := binary.ReadUvarint(bytes.NewReader(header))
		if err != nil {
			fmt.Println("failed to read header", err)
			return
		}

		if n, err := io.CopyN(conn, rand.Reader, int64(length)); err != nil {
			fmt.Println("failed to send payload", err)
		} else {
			fmt.Printf("wrote %d bytes\n", n)
		}
	}
}
