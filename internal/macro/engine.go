package macro

import (
	"log"
	"os/exec"
	"time"

	"github.com/popcorntruck/jank-go/internal/window"
	lua "github.com/yuin/gopher-lua"
)

type MacroEngine struct {
	collection    *macroCollection
	windowService window.WindowService

	luaState *lua.LState
	logger   *MacroLogger
}

func NewMacroEngine() *MacroEngine {
	ws, _ := window.DetermineAndCreateWindowService()

	e := &MacroEngine{
		collection:    newMacroCollection(),
		luaState:      lua.NewState(),
		logger:        NewMacroLogger(),
		windowService: ws,
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

	log.Printf("[MacroEngine] macro '%s' bound to '%s'", macro.Name, hotkey)

	if macro.Action != nil {
		go func() {
			rt, _ := e.luaState.NewThread()

			e.luaState.Resume(rt, macro.Action)
		}()
		// // WE NEED TO RUN THESE IN A COROUTINE
		// if err := rt.PCall(0, 0, nil); err != nil {
		// 	log.Printf("[MacroEngine] Error executing macro '%s': %v", macro.Name, err)
		// 	return err
		// }
	}

	return nil
}

func (e *MacroEngine) initLuaState() {
	e.luaState.SetGlobal("macro", e.luaState.NewFunction(e.lHandleRegisterMacro))

	e.luaState.SetGlobal("print", e.luaState.NewFunction(e.logger.lHandleScriptLog))
	e.luaState.SetGlobal("sleep", e.luaState.NewFunction(e.lSleep))
	e.luaState.SetGlobal("send_notification", e.luaState.NewFunction(lHandleSendNotification))

	e.luaState.SetGlobal("win_class_active", e.luaState.NewFunction(window.CreateLWinClassActive(e.windowService)))
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

func (m *MacroEngine) lSleep(ls *lua.LState) int {
	ms := ls.CheckInt(1)
	time.Sleep(time.Duration(ms) * time.Millisecond)
	return 0
}

func lHandleSendNotification(ls *lua.LState) int {
	value := ls.CheckString(1)

	exec.Command("notify-send", value).Run()

	ls.Push(lua.LNil)
	return 1
}
