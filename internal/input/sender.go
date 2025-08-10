package input

import (
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type InputSender interface {
	Close() error
	SendClick() error
	SendKeyPress(keyCode string) error
}

type NoopInputSender struct {
}

func NewNoopInputSender() (*NoopInputSender, error) {
	return &NoopInputSender{}, nil
}

func (s *NoopInputSender) Close() error {
	return nil
}

func (s *NoopInputSender) SendClick() error {
	return nil
}

func (s *NoopInputSender) SendKeyPress(keyCode string) error {
	return nil
}

func GetPlatformInputSender() (InputSender, error) {
	sp := path.Join(os.Getenv("HOME"), ".ydotool_socket")

	if path, err := exec.LookPath("ydotool"); err == nil {
		if err := ensureYdotooldRunning(sp, 1*time.Second); err == nil {
			return NewYdotoolInputSender(path, sp)
		}

		log.Println("[GetPlatformInputSender] ydotoold is not available or not running. Ensure ydotoold started. falling back to NoopInputSender.")
		return NewNoopInputSender()
	}

	log.Println("[GetPlatformInputSender] ydotool is not available on the system, falling back to NoopInputSender.")
	return NewNoopInputSender()
}

func CreateLClick(is InputSender) func(ls *lua.LState) int {
	return func(ls *lua.LState) int {
		is.SendClick()

		ls.Push(lua.LNil)
		return 1
	}
}
