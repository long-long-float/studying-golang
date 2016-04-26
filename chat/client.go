package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:12345")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			msgBuf := make([]byte, 1024)
			msgLen, err := conn.Read(msgBuf)
			if err != nil {
				log.Fatal(err)
			}

			msg := string(msgBuf[:msgLen])

			fmt.Println(msg)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		conn.Write([]byte(scanner.Text()))
	}

}
