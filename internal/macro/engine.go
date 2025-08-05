package macro

import (
	"log"
	"os/exec"

	lua "github.com/yuin/gopher-lua"
)

type MacroEngine struct {
	collection *macroCollection

	luaState *lua.LState
	logger   *MacroLogger
}

func NewMacroEngine() *MacroEngine {
	e := &MacroEngine{
		collection: newMacroCollection(),
		luaState:   lua.NewState(),
		logger:     NewMacroLogger(),
	}

	e.initLuaState()

	return e
}

func (e *MacroEngine) Close() {
	if e.luaState != nil {
		e.luaState.Close()
	}
}

func (e *MacroEngine) RunScriptFile(path string) error {
	if err := e.luaState.DoFile(path); err != nil {
		panic(err)
	}

	return nil
}

func (e *MacroEngine) TryCallByHotkey(hotkey string) error {
	macro := e.collection.getByHotkey(hotkey)
	if macro == nil {
		return nil // No macro registered for this hotkey
	}

	log.Printf("[MacroEngine] macro '%s' bpund to '%s'", macro.Name, hotkey)

	if macro.Action != nil {
		e.luaState.Push(macro.Action)
		if err := e.luaState.PCall(0, 0, nil); err != nil {
			log.Printf("[MacroEngine] Error executing macro '%s': %v", macro.Name, err)
			return err
		}
	}

	return nil
}

func (e *MacroEngine) initLuaState() {
	e.luaState.SetGlobal("macro", e.luaState.NewFunction(e.lHandleRegisterMacro))

	e.luaState.SetGlobal("print", e.luaState.NewFunction(e.logger.LHandleScriptLog))
	e.luaState.SetGlobal("send_notification", e.luaState.NewFunction(lHandleSendNotification))
	e.luaState.SetGlobal("exec", e.luaState.NewFunction(lHandleExec))
}

func (e *MacroEngine) lHandleRegisterMacro(ls *lua.LState) int {
	name := ls.CheckString(1)
	config := ls.CheckTable(2)

	if name == "" {
		ls.RaiseError("Macro name cannot be empty")
		return 0
	}
	if e.collection.getByName(name) != nil {
		ls.RaiseError("Macro with name '%s' already exists", name)
		return 0
	}

	macro := &Macro{Name: name}

	if hotkey := config.RawGetString("hotkey"); hotkey != lua.LNil {
		macro.Hotkey = hotkey.String()
	}

	if action := config.RawGetString("action"); action != lua.LNil {
		if fn, ok := action.(*lua.LFunction); ok {
			macro.Action = fn
		}
	}

	e.collection.add(macro)
	log.Printf("[MacroEngine] Registered macro: %s (hotkey: %s)\n", name, macro.Hotkey)
	return 0
}

func lHandleSendNotification(ls *lua.LState) int {
	value := ls.CheckString(1)

	exec.Command("notify-send", value).Run()

	ls.Push(lua.LNil)
	return 1
}

func lHandleExec(ls *lua.LState) int {
	command := ls.CheckString(1)

	// not everyone uses bash btw
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()

	if err != nil {
		ls.Push(lua.LNil)
		ls.Push(lua.LString(err.Error()))
		return 2
	}

	ls.Push(lua.LString(string(output)))
	ls.Push(lua.LNil)
	return 2
}
