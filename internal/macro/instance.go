package macro

import (
	lua "github.com/yuin/gopher-lua"
)

type MacroConfig struct {
	// Max times the macro can be ran concurrently
	MaxThreads uint8
}

func DefaultMacroConfig() *MacroConfig {
	return &MacroConfig{
		MaxThreads: 1,
	}
}

func MacroConfigFromTable(table *lua.LTable) *MacroConfig {
	config := DefaultMacroConfig()
	if table == nil {
		return config
	}

	if maxThreads := table.RawGetString("max_threads"); maxThreads != lua.LNil {
		config.MaxThreads = uint8(maxThreads.(lua.LNumber))
	}

	return config
}

type Macro struct {
	Name   string
	Hotkey string
	Action *lua.LFunction
	Config *MacroConfig
}
