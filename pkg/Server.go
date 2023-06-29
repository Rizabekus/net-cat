package pkg

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type Client struct {
	Name string
	Addr string
	Conn net.Conn
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
	num      int
)

func Listener(port string) {
	tmpWelcome, err := os.ReadFile("welcome.txt")
	if err != nil {
		log.Println("File reading error!")
		return
	}
	welcome := string(tmpWelcome)
	numUsers = 0
	serv, err := net.Listen("tcp", "localhost:"+port)
	var wg sync.WaitGroup
	if err != nil {
		panic(err)
	}
	log.Printf("Listening on the port: %s\n", "127.0.0.1:"+port)

	go Broadcast()
	for {
		if numUsers <= 10 {
			wg.Add(1)
			conn, err := serv.Accept()
			if err != nil {
				log.Println("Unable to connect")
				return
			}
			num++
			go ProcessClient(conn, welcome, &wg, &num)

		}
	}
}

func ProcessClient(conn net.Conn, welcome string, wg *sync.WaitGroup, num *int) {
	defer wg.Done()
	reader := bufio.NewReader(os.Stdin)

	conn.Write([]byte(welcome))
	conn.Write([]byte("\n[ENTER YOUR NAME]:"))
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Println("Some error with reading from terminal with reader.ReadString")
		return

	}

	Client := Client{
		Name: name,
		Addr: conn.LocalAddr().String(),
		Conn: conn,
	}
	Clients = append(Clients, Client)
	join <- Client.Name + " has joined our chat..."
	conn.Close()
	// *num = *num - 1
}

func Broadcast() {
	for {
		select {
		case msg := <-join:
			gg.Lock()

			fmt.Fprintf(Clients[0].Conn, msg)
			fmt.Fprintf(Clients[1].Conn, msg)
			fmt.Fprintf(Clients[2].Conn, msg)
			gg.Unlock()
		case msg := <-message:
			gg.Lock()

			fmt.Fprintf(Clients[0].Conn, msg)
			fmt.Fprintf(Clients[1].Conn, msg)
			fmt.Fprintf(Clients[2].Conn, msg)
			gg.Unlock()
		case msg := <-leave:
			gg.Lock()

			fmt.Fprintf(Clients[0].Conn, msg)
			fmt.Fprintf(Clients[1].Conn, msg)
			fmt.Fprintf(Clients[2].Conn, msg)
			gg.Unlock()
		}
	}
}
