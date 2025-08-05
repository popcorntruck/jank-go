package macro

type macroCollection struct {
	byName   map[string]*Macro
	byHotkey map[string]*Macro
}

func newMacroCollection() *macroCollection {
	return &macroCollection{
		byName:   make(map[string]*Macro),
		byHotkey: make(map[string]*Macro),
	}
}

func (c *macroCollection) add(macro *Macro) {
	if macro == nil {
		return
	}

	c.byName[macro.Name] = macro

	if macro.Hotkey != "" {
		c.byHotkey[macro.Hotkey] = macro
	}
}

func (c *macroCollection) getByName(name string) *Macro {
	if name == "" {
		return nil
	}

	return c.byName[name]
}

func (c *macroCollection) getByHotkey(hotkey string) *Macro {
	if hotkey == "" {
		return nil
	}

	return c.byHotkey[hotkey]
}
