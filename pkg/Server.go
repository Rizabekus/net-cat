package pkg

import (
	"io"
	"log"
	"net"
	"os"
)

func Listener(port string) {
	serv, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("Listening on the port: %s\n", "127.0.0.1:"+port)
	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Println("Unable to connect")
		} else {
			go ProcessClient(conn)
		}

	}
}

func ProcessClient(conn net.Conn) {
	_, err := io.Copy(os.Stdout, conn)
	if err != nil {
		panic(err)
	}
	conn.Close()
}
