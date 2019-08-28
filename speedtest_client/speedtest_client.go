package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/elwin/transmit2/scion"
)

var (
	local  = flag.String("local", "", "Local address (with Port)")
	remote = flag.String("remote", "", "Remote address (with Port)")
	size   = flag.Int("size", 1024, "KB to be sent")
)

const (
	sizeMuliplier = 1024 // KB
)

func main() {
	flag.Parse()
	if *local == "" {
		log.Fatal("Please specify the local address using -local")
	}

	if *remote == "" {
		log.Fatal("Please specify the remote address using -remote")
	}

	conn, err := scion.DialAddr(*local, *remote, scion.DefaultPathSelector)
	if err != nil {
		log.Fatal("failed to connect", err)
	}
	defer conn.Close()

	header := make([]byte, 4)
	binary.PutUvarint(header, uint64(*size*sizeMuliplier))

	start := time.Now()

	if err := binary.Write(conn, binary.BigEndian, header); err != nil {
		log.Fatal("failed to write size", err)
	}

	if n, err := io.CopyN(ioutil.Discard, conn, int64(*size)*sizeMuliplier); err != nil && n != int64(*size)*sizeMuliplier {
		log.Fatal("failed to read payload", err)
	} else {
		fmt.Printf("read %d KB\n", *size)
	}

	fmt.Println(float64(*size)/time.Since(start).Seconds(), " KB/s")

}
