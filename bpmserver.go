package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var bpm = 90

type BpmMsg struct {
	Bpm         int
	MicroTick   int
	TickLength  int
	Beat        int
	TickCounter int
}

func bpm_server(ch *amqp.Channel, done chan bool) {

	tickCounter := 1

	for {
		tickLength := (60000 / bpm) / 4 // 4 ticks per BPM measure ( one minute in milliseconds divided by Beats Per Minute )
		timer := time.NewTimer(time.Duration(tickLength) * time.Millisecond)

		beatTick := tickCounter % 32
		beat := (beatTick + 3) / 4
		microTick := tickCounter % 4
		if microTick == 0 {
			microTick = 4
		}

		err := ch.ExchangeDeclare(
			"bpm",    // name
			"fanout", // type
			false,    // durable
			true,     // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		failOnError(err, "Failed to declare an exchange")
		m := BpmMsg{bpm, microTick, tickLength, beat, tickCounter}

		msg, err := json.Marshal(m)
		if err != nil {
			log.Fatal("Barf!", err)
		}
		err = ch.Publish(
			"bpm", // exchange
			"",    // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType: "text/json",
				Body:        msg,
			})
		failOnError(err, "Failed to publish a message")
		log.Printf("Sending msg -- %s", msg)
		tickCounter++

		<-timer.C // pause before next iteration - length of bpm
	}

	//done <- true

}

func serverMain() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	done := make(chan bool)

	go bpm_server(ch, done)

	<-done

}
