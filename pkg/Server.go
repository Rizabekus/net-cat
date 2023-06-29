package pkg

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type Client struct {
	Name string
	addr string
	conn net.Conn
}

type Message struct {
	Name string
	Date string
	Text string
}

var (
	Clients  []Client
	gg       sync.Mutex
	join     chan string
	message  chan string
	leave    chan string
	numUsers int
	history  []Message
)

func Listener(port string) {
	numUsers = 0
	serv, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("Listening on the port: %s\n", "127.0.0.1:"+port)
	var m sync.Mutex
	go Broadcast(&m)
	for {
		if numUsers <= 10 {
			conn, err := serv.Accept()
			if err != nil {
				log.Println("Unable to connect")
				return
			}
			go ProcessClient(conn)
		}
	}
}

func ProcessClient(conn net.Conn) {
	conn.Close()
}

func Broadcast(m *sync.Mutex) {
	for {
		select {
		case msg := <-join:
			m.Lock()
			fmt.Println(msg)
			m.Unlock()
		case msg := <-message:
			m.Lock()
			fmt.Println(msg)
			m.Unlock()
		case msg := <-leave:
			m.Lock()
			fmt.Println(msg)
			m.Unlock()
		}
	}
}
