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
	local  = flag.String("local", "", "Local address (without Port)")
	remote = flag.String("remote", "", "Remote address (with Port)")
)

const (
	sizeMultiplier = 1000 * 1000 // MB
	payload        = sizeMultiplier * 5
	repetitions    = 10
	wait           = 5 * time.Second
)

type Result struct {
	duration time.Duration
	size     int
}

func main() {
	flag.Parse()
	if *local == "" {
		log.Fatal("Please specify the local address using -local")
	}

	if *remote == "" {
		log.Fatal("Please specify the remote address using -remote")
	}

	// warm-up
	selector := scion.DefaultPathSelector
	test(payload, selector)

	var duration time.Duration

	for i := 0; i < repetitions; i++ {
		time.Sleep(wait)

		duration += test(payload, selector)
	}

	fmt.Printf("Payload: %d\n", payload/sizeMultiplier)
	fmt.Printf("Duration (s): %f\n", duration.Seconds()/repetitions)

}

func test(payload int, selector scion.PathSelector) time.Duration {

	conn, err := scion.DialAddr(*local+":0", *remote, selector)
	if err != nil {
		log.Fatal("failed to connect", err)
	}
	defer conn.Close()

	header := header2.Header{
		Size:        payload,
		Repetitions: 1,
	}

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(&header); err != nil {
		log.Fatal("failed to read header", err)
	}

	start := time.Now()

	if n, err := io.CopyN(ioutil.Discard, conn, int64(header.Size)); err != nil && n != int64(header.Size) {
		log.Fatal("failed to read payload", err)
	}

	return time.Since(start)
}
