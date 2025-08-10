package window

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	EVENT_ACTIVE_WINDOW = "activewindow"
)

type HyprWindowService struct {
	socket       net.Conn
	activeWindow *WindowInfo
}

func NewHyprWindowService() (*HyprWindowService, error) {
	xdg := os.Getenv("XDG_RUNTIME_DIR")
	his := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	location := fmt.Sprintf("%s/hypr/%s/.socket2.sock", xdg, his)

	socket, err := net.Dial("unix", location)
	if err != nil {
		log.Println("[WindowService:hypr] Error connecting to socket:", err)
		return nil, err
	}

	ws := &HyprWindowService{
		socket:       socket,
		activeWindow: nil,
	}

	go func() {
		reader := bufio.NewReader(socket)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					// @todo HANDLE THIS PROPERTLY
					log.Println("[WindowService:hypr] Socket closed:", err)
					return
				}

				log.Println("[WindowService:hypr] Error reading from socket:", err)
				return
			}
			event, data := parseEvent(line)

			if event == EVENT_ACTIVE_WINDOW {
				windowData := strings.SplitN(data, ",", 2)
				if len(windowData) < 2 {
					log.Println("Invalid activewindow data:", data)
					continue
				}

				window := WindowInfo{
					Class: windowData[0],
					Title: windowData[1],
				}

				ws.activeWindow = &window

				// log.Printf("[WindowService:hypr] activewindow >> %+v\n", window)
			}
		}
	}()

	return ws, nil
}

func (s *HyprWindowService) Close() error {
	s.socket.Close()
	return nil
}

// Gets the active window information.
// Will be nil if the socket can't connect, or a activewindow event hasn't been sent
// since the WindowService was created.
func (s *HyprWindowService) GetActiveWindow() *WindowInfo {
	return s.activeWindow
}

// returns (event_name, data)
// defined at https://wiki.hypr.land/IPC/
func parseEvent(line string) (string, string) {
	// remove newline
	line = line[:len(line)-1]
	info := strings.SplitN(line, ">>", 2)
	return info[0], info[1]
}
