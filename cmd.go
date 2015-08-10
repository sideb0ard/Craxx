package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/streadway/amqp"
)

func cmdUsage() {
	fmt.Println("Publish one off change to message bus - Usage:\n", os.Args[0], " <string-name> <int-value>")
	os.Exit(1)
}

func cmdMain(ch *amqp.Channel) {
	fmt.Println(len(os.Args))
	if len(os.Args) != 4 {
		cmdUsage()
	}
	name := os.Args[2]
	val, _ := strconv.Atoi(os.Args[3])
	fmt.Println("ARGS", os.Args, name, val)

	m := UpdateMsg{name, val}

	msg, err := json.Marshal(m)
	if err != nil {
		log.Fatal("Barf!", err)
	}
	err = ch.Publish(
		"updateMsgs", // exchange
		"",           // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        msg,
		})
	failOnError(err, "Failed to publish a message")
	log.Printf("Sending msg -- %s", msg)
	os.Exit(0)

}
