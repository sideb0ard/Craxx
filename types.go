package main

type BpmMsg struct {
	Bpm         int
	MicroTick   int
	TickLength  int
	Beat        int
	TickCounter int
}

type UpdateMsg struct {
	Name  string
	Value int
}
