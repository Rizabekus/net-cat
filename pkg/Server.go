package pkg

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
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
	Clients []Client
	gg      sync.Mutex
	join    = make(chan string)
	message = make(chan Message)
	leave   = make(chan Message)

	history []string
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
	gg.Lock()
	buffer := make([]byte, 1024)
	r, err := conn.Read(buffer)
	if err != nil {
		log.Fatal("Some errors with connection: conn.Read()")
	}

	Client := Client{
		Name: string(buffer[:r-1]),
		Addr: conn.LocalAddr().String(),
		Conn: conn,
	}

	Clients = append(Clients, Client)
	HelloText := Client.Name + " has joined our chat..."

	join <- HelloText

	if len(history) != 0 {
		for _, el := range history {
			conn.Write([]byte(el))
		}
	}
	gg.Unlock()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		gg.Lock()
		text := input.Text()
		currentTime := time.Now().Format("2006-01-02 15:04:05")

		message <- Message{
			Name: Client.Name,
			Date: currentTime,
			Text: text,
		}
		gg.Unlock()

	}
	gg.Lock()
	LeavingText := Client.Name + " has left the chat..."

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	leave <- Message{
		Name: Client.Name,
		Date: currentTime,
		Text: LeavingText,
	}
	gg.Unlock()
	defer conn.Close()
}

func Broadcast() {
	for {
		select {
		case str := <-join:

			gg.Lock()

			for _, client := range Clients {
				client.Conn.Write([]byte(str + "\n"))
			}
			history = append(history, str+"\n")
			gg.Unlock()
		case msg := <-message:
			gg.Lock()
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			text := fmt.Sprintf("[%s][%s]:%s", currentTime, msg.Name, msg.Text)
			for _, client := range Clients {
				client.Conn.Write([]byte(text + "\n"))
			}
			history = append(history, text+"\n")
			gg.Unlock()
		case msg := <-leave:
			gg.Lock()

			for _, client := range Clients {
				if client.Name != msg.Name {
					client.Conn.Write([]byte(msg.Text + "\n"))
				}
			}
			history = append(history, msg.Text+"\n")
			gg.Unlock()
		}
	}
}
