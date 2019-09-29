package main

import "gobot.io/x/gobot/platforms/keyboard"

const (
	intensity int = 50
	tickRate  int = 10
)

var controlCommands = []int{
	keyboard.ArrowLeft,
	keyboard.ArrowRight,
	keyboard.ArrowUp,
	keyboard.ArrowDown,
	keyboard.W,
	keyboard.S,
}
