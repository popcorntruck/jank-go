package macro

import (
	"log"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type LogEntry struct {
	TimeStamp time.Time
	Message   string
}

type MacroLogger struct {
	entries []LogEntry
}

func NewMacroLogger() *MacroLogger {
	return &MacroLogger{
		entries: make([]LogEntry, 0),
	}
}

func (l *MacroLogger) Log(message string) {
	entry := LogEntry{
		TimeStamp: time.Now(),
		Message:   message,
	}
	l.entries = append(l.entries, entry)
}

// Lua function to handle script logging, bound to the global `print` function
func (l *MacroLogger) LHandleScriptLog(ls *lua.LState) int {
	value := ls.Get(1)

	logMessage := value.String()
	l.Log(logMessage)
	log.Printf("[lua::print]: %s", value.String())

	ls.Push(lua.LNil)
	return 1
}
