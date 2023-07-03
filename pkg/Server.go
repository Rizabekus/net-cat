package pkg

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
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
	join    = make(chan Message)
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

	// buffer := make([]byte, 1024)
	// r, err := conn.Read(buffer)
	// if err != nil {
	// 	log.Fatal("Some errors with connection: conn.Read()")
	// }
	// for string(buffer[:r-1]) == "" {
	// 	conn.Write([]byte("[EMPTY NAME IS UNAVAILABLE, ENTER YOUR NAME]:"))
	// 	r, err = conn.Read(buffer)
	// 	if err != nil {
	// 		log.Fatal("Some errors with connection: conn.Read()")
	// 	}
	// }
	buffer := NameCheck(Clients, conn)

	Client := Client{
		Name: string(buffer),
		Addr: conn.LocalAddr().String(),
		Conn: conn,
	}

	Clients = append(Clients, Client)
	HelloText := "\r\n" + Client.Name + " has joined our chat..."
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	gg.Lock()
	join <- Message{
		Name: Client.Name,
		Date: currentTime,
		Text: HelloText,
	}
	gg.Unlock()

	if len(history) != 0 {
		for _, el := range history {
			conn.Write([]byte(el))
		}
	}

	input := bufio.NewScanner(conn)
	currentTime = time.Now().Format("2006-01-02 15:04:05")
	f := fmt.Sprintf("[%s][%s]:", currentTime, Client.Name)
	conn.Write([]byte(f))

	for input.Scan() {

		text := input.Text()
		if len(strings.Trim(text, " \r\n")) == 0 {
			currentTime = time.Now().Format("2006-01-02 15:04:05")
			f = fmt.Sprintf("[%s][%s]:", currentTime, Client.Name)
			conn.Write([]byte(f))

			continue
		}
		gg.Lock()
		message <- Message{
			Name: Client.Name,
			Date: currentTime,
			Text: text,
		}
		gg.Unlock()
		currentTime = time.Now().Format("2006-01-02 15:04:05")
		f = fmt.Sprintf("[%s][%s]:", currentTime, Client.Name)
		conn.Write([]byte(f))

	}

	LeavingText := Client.Name + " has left the chat..."
	for i, clients := range Clients {
		if Client.Name == clients.Name {
			Clients = append(Clients[:i], Clients[i+1:]...)
			break
		}
	}

	currentTime = time.Now().Format("2006-01-02 15:04:05")
	gg.Lock()
	leave <- Message{
		Name: Client.Name,
		Date: currentTime,
		Text: LeavingText,
	}
	defer conn.Close()
	gg.Unlock()
}

func Broadcast() {
	for {
		select {
		case msg := <-join:

			gg.Lock()

			for _, client := range Clients {
				welcome := fmt.Sprintf("[%s][%s]:", msg.Date, client.Name)
				if client.Name != msg.Name {
					client.Conn.Write([]byte(msg.Text + "\n"))
					client.Conn.Write([]byte(welcome))
				}
			}
			gg.Unlock()
		case msg := <-message:
			gg.Lock()
			// currentTime := time.Now().Format("2006-01-02 15:04:05")
			text := fmt.Sprintf("[%s][%s]:%s", msg.Date, msg.Name, msg.Text)

			for _, client := range Clients {
				w := fmt.Sprintf("[%s][%s]:", msg.Date, client.Name)
				if client.Name != msg.Name {
					client.Conn.Write([]byte("\r\n" + text + "\n"))
					client.Conn.Write([]byte(w))
				}
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
			gg.Unlock()
		}
	}
}

func NameCheck(Clients []Client, conn net.Conn) string {
	gg := bufio.NewScanner(conn)

	for gg.Scan() {
		flag := true
		text := gg.Text()
		if len(strings.Trim(text, " \r\n")) == 0 {
			conn.Write([]byte("[Empty name is unavailable.Write your name]:"))
			continue
		}
		if len(text) > 30 {
			conn.Write([]byte("[TOO LONG NAME. ENTER YOU NAME]:"))
			continue
		}
		for _, client := range Clients {
			if client.Name == text {
				conn.Write([]byte("[That name is already used. Choose another one]:"))
				flag = false
				break
			}
		}
		if flag == false {
			continue
		} else {
			return text
		}
	}
	return "AZALOH"
}
