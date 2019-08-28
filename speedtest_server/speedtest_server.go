package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/elwin/transmit2/scion"
)

var (
	local = flag.String("local", "", "Local address (with Port)")
)

func registerContext() context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()
	return ctx
}

func main() {
	flag.Parse()
	if *local == "" {
		log.Fatal("Please specify the local address using -local")
	}

	ctx := registerContext()

	listener, err := scion.Listen(*local)
	if err != nil {
		log.Fatal("failed to listen", err)
	}
	defer listener.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("failed to accept connection", err)
			}

			go handleConnection(ctx, conn)
		}
	}
}

func handleConnection(ctx context.Context, conn *scion.Connection) {
	defer conn.Close()
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return
		default:
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
}
