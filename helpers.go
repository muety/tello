package main

import "gobot.io/x/gobot/platforms/dji/tello"

func resetSteering(drone *tello.Driver) {
	drone.Left(0)
	drone.Right(0)
	drone.Up(0)
	drone.Down(0)
	drone.Forward(0)
	drone.Backward(0)
	drone.Clockwise(0)
	drone.CounterClockwise(0)
}
