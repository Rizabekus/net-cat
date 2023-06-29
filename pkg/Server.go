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
	if err != nil {
		panic(err)
	}
	log.Printf("Listening on the port: %s\n", "127.0.0.1:"+port)

	go Broadcast()
	for {

		conn, err := serv.Accept()
		if err != nil {
			log.Println("Unable to connect")
			return
		}

		go ProcessClient(conn, welcome)

	}
}

func ProcessClient(conn net.Conn, welcome string) {
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
			fmt.Println("Azaloh")
			gg.Lock()

			for _, client := range Clients {
				fmt.Fprintf(client.Conn, msg)
			}
			gg.Unlock()
		case msg := <-message:
			gg.Lock()

			for _, client := range Clients {
				fmt.Fprintf(client.Conn, msg)
			}
			gg.Unlock()
		case msg := <-leave:
			gg.Lock()

			for _, client := range Clients {
				fmt.Fprintf(client.Conn, msg)
			}
			gg.Unlock()
		}
	}
}
