package main

type action struct {
	Name     string
	Color    color
	callback func()
}

// Add a new action to `m.Actions`.
//
// `callback` is called when this actions is selected.
//
// Returns a pointer to the new action.
//
// To set default colors set `tui.Defaults.Color` before creating options.
func (m *menu) NewAction(name string, callback func()) *action {
	act := &action{
		Name:     name,
		Color:    Defaults.Color,
		callback: callback,
	}
	m.Actions = append(m.Actions, act)
	return act
}
