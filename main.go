package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleReq(conn)
	}
}

func handleReq(conn net.Conn) {
	buf := make([]byte, 1024)

	for {
		resInt, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Recieved data (%s): %s\n", conn.LocalAddr(), string(buf[:resInt]))
		conn.Write([]byte("Data Send."))

	}

	// conn.Close()
}
