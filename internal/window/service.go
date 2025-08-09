package window

import (
	"log"
	"os"
)

type WindowInfo struct {
	Class string
	Title string
}

type WindowService interface {
	Close() error
	ActiveWindow() *WindowInfo
}

func DetermineAndCreateWindowService() (WindowService, error) {
	wm := os.Getenv("XDG_CURRENT_DESKTOP")

	if wm == "Hyprland" {
		return NewHyprWindowService()
	} else {
		log.Printf("[WindowService] window manager %v not supported, using NoopWindowService", wm)
		return NewNoopWindowService()
	}
}
