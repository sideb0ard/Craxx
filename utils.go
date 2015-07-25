package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/mgutz/ansi"
)

const (
	playx = "/Users/sideboard/homebrew/bin/play"
)

var lime = ansi.ColorCode("green+h:black")
var reset = ansi.ColorCode("reset")

type Soxfilter struct {
	Effect string
	Val    string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func playrrr(sound string, args Soxfilter) {
	fmt.Println(lime, sound, reset)
	fmt.Println(args.Effect)
	if args.Effect != "" {
		fmt.Println("got an effect!", args.Effect, args.Val, "\n")
		cmd := exec.Command(playx, sound, args.Effect, args.Val)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Waiting for command to finish...")
		err = cmd.Wait()
		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	} else {
		cmd := exec.Command(playx, sound)
		err := cmd.Start()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Waiting for command to finish...")
		err = cmd.Wait()
		if err != nil {
			log.Printf("Command finished with error: %v", err)
		}
	}
}
