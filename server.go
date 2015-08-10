package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var bpm = 145

func msgListener(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		"updateMsgs", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		var m UpdateMsg
		err := json.Unmarshal(d.Body, &m)
		if err != nil {
			fmt.Println("blah", err)
		}
		if m.Name == "bpm" {
			fmt.Println("Setting BPM to", m.Value)
			bpm = m.Value
		}
	}
}

func serverMain(ch *amqp.Channel) {

	go msgListener(ch)
	tickCounter := 1
	tickLength := (60000 / bpm) / 4 // 1min divided by bpm divided by 4 microticks

	for {

		if bpm != 0 {
			tickLength = (60000 / bpm) / 4
		}
		timer := time.NewTimer(time.Duration(tickLength) * time.Millisecond)
		beatTick := tickCounter % 32
		beat := (beatTick + 3) / 4
		microTick := tickCounter % 4
		if microTick == 0 {
			microTick = 4
		}

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
}
