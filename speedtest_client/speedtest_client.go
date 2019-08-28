package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	header2 "github.com/elwin/speedtest/header"

	"github.com/elwin/transmit2/scion"
)

var (
	local   = flag.String("local", "", "Local address (with Port)")
	remote  = flag.String("remote", "", "Remote address (with Port)")
	size    = flag.Int("size", 1024, "bytes to be sent")
	packets = flag.Int("packets", 1, "number of packets to be sent")
)

const (
	sizeMuliplier = 1 // KB
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

	header := header2.Header{
		Size:        *size,
		Repetitions: *packets,
	}

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(&header); err != nil {
		log.Fatal("failed to read header", err)
	}
	start := time.Now()

	for i := 0; i < *packets; i++ {

		if n, err := io.CopyN(ioutil.Discard, conn, int64(*size)*sizeMuliplier); err != nil && n != int64(*size)*sizeMuliplier {
			log.Fatal("failed to read payload", err)
		}

		if i%100 == 0 {
			fmt.Print(".")
		}

	}

	fmt.Println()

	fmt.Println(float64(header.Size*header.Repetitions)/1024/time.Since(start).Seconds(), " KB/s")

}
