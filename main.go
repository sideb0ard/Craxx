package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Println("Usage:", os.Args[0], "<client|server>")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	if os.Args[1] == "client" {
		fmt.Println("Starting client...")
		clientMain()
	} else if os.Args[1] == "server" {
		fmt.Println("Starting server...")
		serverMain()
	} else {
		usage()
	}
}
