package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func usage() {
	fmt.Println("Usage:", os.Args[0], "<client|server|bit>")
	os.Exit(1)
}

func main() {

	if len(os.Args) < 2 {
		usage()
	}

	done := make(chan bool)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	rabbitchan, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer rabbitchan.Close()

	exchanges := []string{"bpm", "updateMsgs"}
	for _, x := range exchanges {
		err = rabbitchan.ExchangeDeclare(
			x,        // name
			"fanout", // type
			false,    // durable
			true,     // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		failOnError(err, "Failed to declare exchange")
	}

	if os.Args[1] == "client" {
		fmt.Println("Starting client...")
		go clientMain(rabbitchan)
	} else if os.Args[1] == "server" {
		fmt.Println("Starting server...")
		go serverMain(rabbitchan)
	} else if os.Args[1] == "sine" {
		fmt.Println("Starting sine.....")
		go sineMain(rabbitchan)
	} else if os.Args[1] == "cmd" {
		fmt.Println("Starting cmd.....")
		go cmdMain(rabbitchan)
	} else {
		usage()
	}

	<-done // wait ferever...
}
