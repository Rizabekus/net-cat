package main

import (
	"fmt"
	"git/rzhampeis/net-cat/pkg"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		pkg.Listener(os.Args[1])
	} else if len(os.Args) == 1 {
		pkg.Listener("8989")
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
	}
}
