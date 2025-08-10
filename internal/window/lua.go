package window

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func CreateLWinClassActive(ws WindowService) func(ls *lua.LState) int {
	return func(ls *lua.LState) int {
		targetClass := ls.CheckString(1)

		aw := ws.GetActiveWindow()
		if aw == nil {
			log.Printf("[CreateLWinClassActive] No active window found")
			ls.Push(lua.LBool(false))
			return 1
		}

		ls.Push(lua.LBool(aw.Class == targetClass))
		return 1
	}
}
