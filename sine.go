package main

import (
	"encoding/json"
	"fmt"
	"math"

	"code.google.com/p/portaudio-go/portaudio"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/streadway/amqp"
)

var tcp = 220
var udp = 220

func sineMain(ch *amqp.Channel) {

	go sniffy()
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
	s := newStereoSine(220, 220, sampleRate)
	defer s.Close()
	chk(s.Start())
	defer s.Stop()

	for m := range msgs {
		var bm BpmMsg
		err := json.Unmarshal(m.Body, &bm)
		if err != nil {
			fmt.Println("blah", err)
		}
		bpm = bm.Bpm
	}
}

func (g *stereoSine) processAudio(out [][]float32) {
	var t float64 = 0
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2*math.Pi*g.phaseL*(t/float64(bpm)))) / 2
		//out[0][i] = float32(math.Sin(2 * math.Pi * float64(tcp%400) * (t / float64(bpm))))
		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
		out[1][i] = float32(math.Sin(2*math.Pi*g.phaseR*(t/float64(bpm)))) / 2
		//out[1][i] = float32(math.Sin(2 * math.Pi * float64(tcp%400) * (t / float64(bpm))))
		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
		t++
	}
}

func sniffy() {
	LayerCount := make(map[string]int)

	if handle, err := pcap.OpenLive("en0", 1600, true, -1); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			// Iterate over all layers, printing out each layer type
			for _, layer := range packet.Layers() {
				//fmt.Println("PACKET LAYER:", layer.LayerType().String())
				LayerCount[layer.LayerType().String()] += 1
			}
			// fmt.Println("TCP!", LayerCount["TCP"])
			// fmt.Println("UDP!", LayerCount["UDP"])
			tcp = LayerCount["TCP"]
			udp = LayerCount["UDP"]
		}
	}
}
