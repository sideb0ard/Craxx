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

func serverMain(ch *amqp.Channel) {

	tickCounter := 1
	tickLength := (60000 / bpm) / 4 // 1min divided by bpm divided by 4 microticks

	for {

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
