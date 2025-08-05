package input

import (
	"log"

	"github.com/holoplot/go-evdev"
)

type InputReceiver struct {
	// Channel input events are dispatched to (only type EV_KEY)
	events chan *evdev.InputEvent

	// List of channels used to close reciever goroutines
	closes []chan bool
}

func NewInputReceiver() (*InputReceiver, error) {
	device, err := evdev.ListDevicePaths()

	if err != nil {
		return nil, err
	}

	closes := make([]chan bool, 0)
	events := make(chan *evdev.InputEvent)

	for _, path := range device {
		device, err := evdev.Open(path.Path)
		if err != nil {
			log.Printf("[InputReceiver] Failed to open device %s: %v", path.Name, err)
			continue
		}

		// Create a channel to signal when to close the receiver
		closeChan := make(chan bool)
		closes = append(closes, closeChan)

		go func() {
			for {
				select {
				case <-closeChan:
					log.Printf("[InputReceiver] Closing receiver for device %s", path.Name)
					device.Close()
					return
				default:
					// Read events from the device
					event, err := device.ReadOne()

					if err != nil {
						log.Printf("[InputReceiver] Error reading from device %s: %v", path.Name, err)
						continue
					}

					if event.Type == evdev.EV_KEY {
						events <- event
					}
				}
			}
		}()

	}

	return &InputReceiver{
		events: events,
		closes: closes,
	}, err
}

func (r *InputReceiver) Events() <-chan *evdev.InputEvent {
	return r.events
}
