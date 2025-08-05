package macro

import lua "github.com/yuin/gopher-lua"

type Macro struct {
	Name   string
	Hotkey string
	Action *lua.LFunction
}
