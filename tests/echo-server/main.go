package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024*1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Printf("Read: %v", err)
		return
	}

	resp := append([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n"), buf[:n]...)
	conn.Write(resp)

	log.Printf("Connection from %s closed\n", conn.RemoteAddr())
}

func main() {
	const localAddr = "0.0.0.0"
	var (
		from string
		help bool
		h    bool
	)

	flag.StringVar(&from, "from", "1337", "from port")
	flag.BoolVar(&help, "help", false, "help")
	flag.BoolVar(&h, "h", false, "help")
	flag.Parse()

	if help || h {
		flag.PrintDefaults()
		return
	}

	envs := os.Environ()
	for _, e := range envs {
		fmt.Printf("Env: %s\n", e)
	}
	fmt.Println("")

	listenAddr, err := net.ResolveTCPAddr("tcp", localAddr+":"+from)
	if err != nil {
		log.Fatalf("ResolveTCPAddr: %v", err)
	}

	ln, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		log.Fatalf("ListenTCP: %v", err)
	}
	defer ln.Close()

	log.Printf("Listening on %s:%s\n", localAddr, from)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept: %v", err)
			continue
		}

		log.Printf("Connection from %s\n", conn.RemoteAddr())
		go handleConn(conn)
	}
}
