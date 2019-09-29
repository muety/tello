package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/bep/debounce"
	"gobot.io/x/gobot/platforms/dji/tello"
	"gobot.io/x/gobot/platforms/keyboard"
)

func handleKeyboardInput(drone *tello.Driver) func(interface{}) {
	return func(data interface{}) {
		if !connected || landing {
			return
		}

		key := data.(keyboard.KeyEvent).Key

		// Case 1: Pressed key is a steering command
		for _, k := range controlCommands {
			if k == key {
				cc := atomic.LoadInt32(&currentControl)

				if cc == -1 {
					steeringDebounce = debounce.New(500 * time.Millisecond)
				}

				atomic.StoreInt32(&currentControl, int32(key))

				steeringDebounce(func() {
					atomic.StoreInt32(&currentControl, -1)
				})

				return
			}
		}

		// Case 2: Pressed key is any other commmand or none
		if key == keyboard.Spacebar {
			// Space
			if !flying {
				fmt.Println("Taking off.")
				if !dry {
					if err := drone.TakeOff(); err != nil {
						fmt.Println(err)
					}
				}
			} else {
				fmt.Println("Landing.")
				if !dry {
					if err := drone.Land(); err != nil {
						fmt.Println(err)
					}
				}
			}
		} else {
			fmt.Println("Unknown command...")
		}
	}
}

func handleFlightData(drone *tello.Driver) func(interface{}) {
	return func(data interface{}) {
		flightData := data.(*tello.FlightData)

		atomic.AddUint64(&dataCounter, 1)

		if atomic.LoadUint64(&dataCounter)%10 == 0 {
			fmt.Printf("Ground Speed: %.2f, Battery: %d %%, Height: %.2f m\n", flightData.GroundSpeed(), flightData.BatteryPercentage, float32(flightData.Height)/10.0)
		}

		flying = flightData.Flying
		if landing && flightData.Height <= 0 {
			fmt.Println("Drone has landed.")
			landing = false
		}
	}
}

func handleConnected(drone *tello.Driver) func(interface{}) {
	return func(data interface{}) {
		fmt.Println("Drone connected.")

		connected = true

		go func() {
			for connected {
				tick(drone)
				time.Sleep(time.Duration((1000.0 / tickRate)) * time.Millisecond)
			}
		}()
	}
}

func handleLanding(drone *tello.Driver) func(interface{}) {
	return func(data interface{}) {
		fmt.Println("Drone is landing ...")
	}
}
