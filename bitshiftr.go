package main

import (
	"encoding/json"
	"fmt"
	"math"

	"code.google.com/p/portaudio-go/portaudio"

	"github.com/streadway/amqp"
)

var t int32 = 55000000
var playOn = false

const sampleRate = 44100

func bitshiftMain(ch *amqp.Channel) {

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

	portaudio.Initialize()
	defer portaudio.Terminate()
	//s := newStereoSine(256, 320, sampleRate)
	s := newStereoSine(255, 233, sampleRate)
	defer s.Close()
	chk(s.Start())
	defer s.Stop()

	for m := range msgs {
		var bm BpmMsg
		err := json.Unmarshal(m.Body, &bm)
		if err != nil {
			fmt.Println("blah", err)
		}
		if bm.TickCounter%8 == 1 && playOn == false {
			s.Start()
			playOn = true
		} else if bm.TickCounter%8 == 0 {
			s.Stop()
			playOn = false
		}

		t += 100000
	}
}

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
}

func newStereoSine(freqL, freqR, sampleRate float64) *stereoSine {
	s := &stereoSine{nil, freqL / sampleRate, 0, freqR / sampleRate, 0}
	var err error
	s.Stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, s.processAudio)
	chk(err)
	return s
}

func (g *stereoSine) processAudio(out [][]float32) {
	localt := t
	for i := range out[0] {
		//fmt.Println("processss... %d", t)
		inum := localt * ((localt>>uint(9) | localt>>uint(13)) & 25 & (localt >> uint(6)))
		num := scalr(inum)
		out[0][i] = num
		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
		// out[1][i] = float32(math.Sin(2 * math.Pi * g.phaseR))
		out[1][i] = num
		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
		localt++
	}
}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func scalr(x int32) float32 {
	low := float32(-2147483647) // 2147481575
	high := float32(2147483647)
	r1 := float32(high - low)
	lscal := float32(-1)
	hscal := float32(1)
	r2 := hscal - lscal
	return (r2 / r1) * (float32(x) + (-1))
}
