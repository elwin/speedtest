package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
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

	results := make([]Result, 0)

	payload := 1000 * 1000 * 5

	// warm-up
	test(payload/10000, 10000)

	for size := 100; size <= 100000; size *= 10 {
		packets := payload / size
		duration := test(packets, size)

		results = append(results, Result{duration, size})
	}

	fmt.Println("packet size (bytes),bandwidth (MB/s)")
	for _, result := range results {
		fmt.Print(strconv.Itoa(result.size) + ",")
		fmt.Println(float64(payload) / (sizeMultiplier * result.duration.Seconds()))
	}

}

func test(packets, size int) time.Duration {
	conn, err := scion.DialAddr(*local+":0", *remote, scion.DefaultPathSelector)
	if err != nil {
		log.Fatal("failed to connect", err)
	}
	defer conn.Close()

	header := header2.Header{
		Size:        size,
		Repetitions: packets,
	}

	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(&header); err != nil {
		log.Fatal("failed to read header", err)
	}

	start := time.Now()

	for i := 0; i < packets; i++ {

		if n, err := io.CopyN(ioutil.Discard, conn, int64(header.Size)); err != nil && n != int64(header.Size) {
			log.Fatal("failed to read payload", err)
		}

	}

	return time.Since(start)
}
