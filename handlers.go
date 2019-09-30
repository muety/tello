package main

import (
	"fmt"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/bep/debounce"
	"gobot.io/x/gobot"
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
					steeringDebounce = debounce.New(debounceTimeout)
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

		drone.SetVideoEncoderRate(2)
		gobot.Every(100*time.Millisecond, func() {
			drone.StartVideo()
		})

		gobot.Every(time.Duration((1000.0/tickRate))*time.Millisecond, func() {
			if connected {
				tick(drone)
			}
		})
	}
}

func handleLanding(drone *tello.Driver) func(interface{}) {
	return func(data interface{}) {
		fmt.Println("Drone is landing ...")
	}
}

func handleVideo(drone *tello.Driver) func(interface{}) {
	mplayer := exec.Command("mplayer", "-fps", "20", "-")

	videoIn, err := mplayer.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if err := mplayer.Start(); err != nil {
		fmt.Println(err)
		return nil
	}

	return func(data interface{}) {
		pkt := data.([]byte)
		if _, err := videoIn.Write(pkt); err != nil {
			fmt.Println(err)
		}
	}
}
