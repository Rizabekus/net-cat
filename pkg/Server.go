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
	arr     []net.Conn
	Clients []Client
	gg      sync.Mutex
	join    = make(chan string)
	message = make(chan string)
	leave   = make(chan string)

	history []Message
)

func Listener(port string) {
	tmpWelcome, err := os.ReadFile("welcome.txt")
	if err != nil {
		log.Println("File reading error!")
		return
	}
	welcome := string(tmpWelcome)

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
	conn.Write([]byte(welcome))
	conn.Write([]byte("\n[ENTER YOUR NAME]:"))
	buffer := make([]byte, 1024)
	r, err := conn.Read(buffer)
	if err != nil {
		log.Fatal("Some errors with connection: conn.Read()")
	}

	// conn.Write([]byte(string(buffer[:r-1]) + " has joined the chat..."))

	Client := Client{
		Name: string(buffer[:r-1]),
		Addr: conn.LocalAddr().String(),
		Conn: conn,
	}

	Clients = append(Clients, Client)
	arr = append(arr, conn)

	join <- Client.Name + " has joined our chat..."
	input := bufio.NewScanner(conn)
	for input.Scan() {
		message <- input.Text()
	}

	defer conn.Close()
}

func Broadcast() {
	for {
		select {
		case msg := <-join:

			gg.Lock()

			for _, client := range arr {
				client.Write([]byte(msg))
				fmt.Println(msg)

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
