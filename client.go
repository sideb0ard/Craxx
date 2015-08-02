package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/streadway/amqp"
)

const (
	kick  = "WuTangDrumz/WTC_kyKX/W1_K_40_.wav"
	snare = "WuTangDrumz/Cynerz/36ChamberSnarEZ/GVD_snr_47_.wav"
	hat   = "WuTangDrumz/Perkussin/WU_HH_074.wav"
)

func clientMain(ch *amqp.Channel) {

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
		q.Name, // queue name
		"",     // routing key
		"bpm",  // exchange
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
		var bm BpmMsg
		err := json.Unmarshal(d.Body, &bm)
		if err != nil {
			fmt.Println("blah", err)
		}
		if bm.TickCounter%3 == 0 {
			fmt.Println("KICK Got one Modulo 4!", bm.TickCounter)
			//go playrrr(kick, Soxfilter{})
			if rand.Intn(10) > 2 {
				go playrrr(kick, Soxfilter{Effect: "pitch", Val: strconv.Itoa(int(math.Pow(float64(bm.TickCounter%1000), 3.0)) % 1000)})
			} else {
				fmt.Println("RANDY MISSS!")
			}
		}
		if bm.TickCounter%4 == 0 {
			fmt.Println("SNARE Got one Modulo 3!", bm.TickCounter)
			//go playrrr(snare, Soxfilter{})
			if rand.Intn(10) > 1 {
				//go playrrr(snare, Soxfilter{Effect: "pitch", Val: strconv.Itoa(int(math.Pow(float64(bm.TickCounter%1000), 3.0)) % 1000)})
				go playrrr(snare, Soxfilter{})
			}
		}
		if bm.MicroTick%1 == 0 {
			fmt.Println("HAT Got one Modulo 1!", bm.TickCounter)
			// go playrrr(hat + " pitch " + strconv.Itoa(bm.TickCounter%1000))
			if rand.Intn(10) > 4 {
				go playrrr(hat, Soxfilter{Effect: "pitch", Val: strconv.Itoa(int(math.Pow(float64(bm.TickCounter%1000), 3.0)) % 1000)})
			}
		}
	}
}
