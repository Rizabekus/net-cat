package main

import (
	"fmt"
	"git/rzhampeis/net-cat/pkg"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		pkg.Listener(os.Args[1])

	} else {
		pkg.Listener("8989")
	}
}
