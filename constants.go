package main

import (
	"time"

	"gobot.io/x/gobot/platforms/keyboard"
)

const (
	intensity       int           = 75
	tickRate        int           = 10                     // fps
	debounceTimeout time.Duration = 500 * time.Millisecond // ms
)

var controlCommands = []int{
	keyboard.ArrowLeft,
	keyboard.ArrowRight,
	keyboard.ArrowUp,
	keyboard.ArrowDown,
	keyboard.W,
	keyboard.A,
	keyboard.S,
	keyboard.D,
}
