package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewServer(listnerAddr string) *Server {
	return &Server{
		listenAddr: listnerAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		rep, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Error :", err)
			continue
		}

		fmt.Println("new connection to the server", rep.RemoteAddr())

		go s.getData(rep)
	}
}

func (s *Server) getData(rep net.Conn) {
	defer rep.Close()
	reader := bufio.NewReader(rep)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error:", err)
			}
			break
		}
		s.msgch <- Message{
			from:    rep.RemoteAddr().String(),
			payload: []byte(line),
		}
	}
}

func main() {
	server := NewServer(":3000")

	go func() {
		for msg := range server.msgch {
			fmt.Printf("message recived(%s):%s\n", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())
}

// func handleReq(conn net.Conn) {
// 	buf := make([]byte, 1024)

// 	for {
// 		resInt, err := conn.Read(buf)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("Recieved data (%s): %s\n", conn.LocalAddr(), string(buf[:resInt]))
// 		conn.Write([]byte("Data Send."))

// 	}

// 	// conn.Close()
// }
