package main

import (
	"flag"
	"fmt"
	"sync/atomic"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gobot.io/x/gobot/platforms/keyboard"
)

var (
	dry bool

	connected bool
	flying    bool
	landing   bool

	previousControl int32 = -1
	currentControl  int32 = -1
	tickCounter     uint64
	dataCounter     uint64

	steeringDebounce func(func())
)

func tick(drone *tello.Driver) {
	cc, pc := atomic.LoadInt32(&currentControl), atomic.LoadInt32(&previousControl)

	if !connected || landing || cc == pc {
		return
	}

	if cc == keyboard.A {
		fmt.Println("Going left.")
		if !dry {
			drone.Left(intensity)
		}
	} else if cc == keyboard.D {
		fmt.Println("Going right.")
		if !dry {
			drone.Right(intensity)
		}
	} else if cc == keyboard.W {
		fmt.Println("Going up.")
		if !dry {
			drone.Up(intensity)
		}
	} else if cc == keyboard.S {
		fmt.Println("Going down.")
		if !dry {
			drone.Down(intensity)
		}
	} else if cc == keyboard.ArrowUp {
		fmt.Println("Going forward.")
		if !dry {
			drone.Forward(intensity)
		}
	} else if cc == keyboard.ArrowDown {
		fmt.Println("Going backward.")
		if !dry {
			drone.Backward(intensity)
		}
	} else if cc == keyboard.ArrowLeft {
		fmt.Println("Rotating counter-clockwise.")
		if !dry {
			drone.CounterClockwise(intensity)
		}
	} else if cc == keyboard.ArrowRight {
		fmt.Println("Rotating clockwise.")
		if !dry {
			drone.Clockwise(intensity)
		}
	} else {
		fmt.Println("Resetting steering.")
		resetSteering(drone)
	}

	atomic.StoreInt32(&previousControl, cc)
	atomic.AddUint64(&tickCounter, 1)
}

func init() {
	var dryFlag = flag.Bool("dry", false, "Perform a dry run (don't send any actual control commands to Tello")

	flag.Parse()

	dry = *dryFlag

	if dry {
		fmt.Println("Running in dry mode.")
	}
}

func main() {
	// Init Gobot drivers
	keys := keyboard.NewDriver()
	drone := tello.NewDriver("8890")

	work := func() {
		// Handle keyboard inputs
		keys.On(keyboard.Key, handleKeyboardInput(drone))

		// Handle drone events
		drone.On(tello.FlightDataEvent, handleFlightData(drone))
		drone.On(tello.ConnectedEvent, handleConnected(drone))
		drone.On(tello.LandingEvent, handleLanding(drone))
		drone.On(tello.VideoFrameEvent, handleVideo(drone))
	}

	robot := gobot.NewRobot(
		"tello",
		[]gobot.Connection{},
		[]gobot.Device{keys, drone},
		work,
	)

	robot.Start()
}
