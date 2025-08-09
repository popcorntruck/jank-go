package main

import (
	"github.com/popcorntruck/jank-go/internal/input"
	"github.com/popcorntruck/jank-go/internal/macro"
	"github.com/popcorntruck/jank-go/internal/window"
)

const (
	VALUE_KEY_PRESSED  = 1
	VALUE_KEY_HELD     = 2
	VALUE_KEY_RELEASED = 0
)

func main() {
	_, err := window.DetermineAndCreateWindowService()
	if err != nil {
		panic(err)
	}

	engine := macro.NewMacroEngine()
	defer engine.Close()

	engine.RunScriptFile("macro.lua")

	// Initialize the input receiver
	recv, err := input.NewInputReceiver()

	if err != nil {
		panic(err)
	}

	for {
		event := <-recv.Events()

		if event.Value == VALUE_KEY_PRESSED {
			engine.TryCallByHotkey(event.CodeName())
		}
	} // Block forever until interrupted (e.g., Ctrl+C)
}
