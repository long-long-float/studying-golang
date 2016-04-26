package main

import (
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal(err)
	}

	allCh := make(chan string)
	conns := make([]net.Conn, 0)

	go func() {
		for {
			msg := <-allCh
			for _, conn := range conns {
				conn.Write([]byte(msg))
			}
		}
	}()

	for {
		conn, err := ln.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatal(err)
			continue
		}
		conns = append(conns, conn)
		go func(conn net.Conn) {
			for {
				msgBuf := make([]byte, 1024)
				msgLen, err := conn.Read(msgBuf)
				if err != nil {
					log.Fatal(err)
				}

				msg := string(msgBuf[:msgLen])

				allCh <- msg
			}
		}(conn)
	}

}
