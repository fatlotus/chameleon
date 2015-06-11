package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var server = flag.Bool("server", false, "true to listen, false to connect.")
var port = flag.String("port", ":8080", "addr:port to bind on")
var track = flag.String("track", "", "which track file to use")

func main() {
	var conn net.Conn
	var err error

	flag.Parse()

	if *server {
		cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
		if err != nil {
			log.Fatal(err)
			return
		}

		listener, err := tls.Listen("tcp", *port, &tls.Config{
			Certificates: []tls.Certificate{
				cert,
			},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			},
		})

		if err != nil {
			log.Fatal(err)
			return
		}

		if conn, err = listener.Accept(); err != nil {
			log.Fatal(err)
			return
		}

	} else {
		conn, err = tls.Dial("tcp", "ninetales"+*port, &tls.Config{
			InsecureSkipVerify: true,
		})

		if err != nil {
			log.Fatal(err)
			return
		}
	}

	log.Println("Connected!")

	go func() {
		buf := make([]byte, 10*1024)

		for {
			_, err := conn.Read(buf)
			if err != nil {
				log.Fatal(err)
				break
			}
		}
	}()

	blob := make([]byte, 5000)
	for i := 0; i < len(blob); i++ {
		blob[i] = 'A'
	}

	file, err := os.Open(*track)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)
	times := make([]int64, 0)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		frags := strings.Split(line[:len(line) - 1], " ")
		lng, err := strconv.ParseInt(frags[0], 10, 32)
		if err != nil {
			panic(err)
		}

		tim, err := strconv.ParseFloat(frags[1], 32)
		if err != nil {
			panic(err)
		}

		lng += 1375 - 1452 // measured empirically :p
		// lng = 1375
		if lng < 0 {
			lng = 0
		}

		tm := int(tim * 500)

		for len(times) <= tm {
			times = append(times, 0)
		}
		
		for times[tm] != 0 {
			tm += 1
			
			for len(times) <= tm {
				times = append(times, 0)
			}
		}
		
		times[tm] = lng
	}

	delay := time.Tick(2 * time.Millisecond)
	
	for _, lng := range times {
		<-delay
		
		if lng != 0 {
			_, err := conn.Write(blob[:lng])
			if err != nil {
				panic(err)
			}
		}
	}
	
	time.Sleep(2)
}
