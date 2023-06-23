package main

import (
	"git/rzhampeis/net-cat/pkg"
	"os"
)

func main() {
	if len(os.Args) == 2 {
		pkg.Listener(os.Args[1])
	} else {
		pkg.Listener("8989")
	}
}
