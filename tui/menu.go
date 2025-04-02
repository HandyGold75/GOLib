package main

import (
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"
)

type (
	item interface {
		// Get the value as a string.
		//
		// If the value is not supported then the title or name is used.
		String() string
		// Get the value as a int.
		//
		// If the value is not supported then the title or name is used.
		Int() (int, error)
		// Get the value as a float.
		//
		// If the value is not supported then the title or name is used.
		Float() (float64, error)
		// Get the value as a bool.
		//
		// If the value is not supported then the title or name is used.
		//
		// It return true opon one of these values (case insensitive): `1`, `t`, `true`, `y`, `yes`
		Bool() bool

		// Get the type of the item.
		//
		// This can be one of: `menu`, `action`, `list`, `option`
		Type() string

		enter() error
	}

	menu struct {
		mm            *MainMenu
		Title         string
		Color         color
		AccentColor   color
		SelectColor   color
		SelectBGColor color
		ValueColor    color
		Align         align
		Items         []item
		selected      int
		back          *menu
	}

	action struct {
		mm       *MainMenu
		Name     string
		callback func()
	}

	list struct {
		mm       *MainMenu
		Name     string
		Allowed  []string
		selected int
		editing  bool
	}

	option struct {
		mm      *MainMenu
		Name    string
		Allowed string
		value   string
		editing bool
	}
)

// Add a new menu to `m.Menus`.
//
// Returns a pointer to the new menu.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor`, `tui.Defaults.SelectBGColor`, `tui.Defaults.ValueColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func (m *menu) NewMenu(title string) *menu {
	mn := &menu{
		mm:            m.mm,
		Title:         title,
		Color:         Defaults.Color,
		AccentColor:   Defaults.AccentColor,
		SelectColor:   Defaults.SelectColor,
		SelectBGColor: Defaults.SelectBGColor,
		ValueColor:    Defaults.ValueColor,
		Align:         Defaults.Align,
		Items:         []item{},
		selected:      0,
		back:          m,
	}
	m.Items = append(m.Items, mn)
	return mn
}

// Get the title as a string.
func (m *menu) String() string { return m.Title }

// Get the title as a int.
func (m *menu) Int() (int, error) { return strconv.Atoi(m.Title) }

// Get the title as a float.
func (m *menu) Float() (float64, error) { return strconv.ParseFloat(m.Title, 64) }

// Get the title as a bool.
// It return true opon one of these values (case insensitive): 1, t, true, y, yes
func (m *menu) Bool() bool {
	return slices.Contains([]string{"1", "t", "true", "y", "yes"}, strings.ToLower(m.Title))
}

// Get the type of the item.
//
// Returns `menu`.
func (m *menu) Type() string { return "menu" }

func (m *menu) enter() error {
	_ = m.mm.rdr.Render()
	for {
		in := make([]byte, 3)
		if _, err := os.Stdin.Read(in); err != nil {
			return err
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v []byte) bool { return slices.Equal(v, in) }) {
			return Errors.exit
		} else if slices.ContainsFunc(KeyBinds.Up, func(v []byte) bool { return slices.Equal(v, in) }) {
			m.selected = max(m.selected-1, 0)
			_ = m.mm.rdr.Render()
			continue

		} else if slices.ContainsFunc(KeyBinds.Down, func(v []byte) bool { return slices.Equal(v, in) }) {
			m.selected = min(m.selected+1, len(m.Items))
			_ = m.mm.rdr.Render()
			continue

		} else if slices.ContainsFunc(KeyBinds.Right, func(v []byte) bool { return slices.Equal(v, in) }) {
			if m.selected < len(m.Items) && m.selected >= 0 {
				if m.Items[m.selected].Type() == "menu" {
					m.mm.cur = m.Items[m.selected].(*menu)
					return nil
				}
				err := m.Items[m.selected].enter()
				return err
			}
			if m.back == nil {
				return Errors.exit
			}
			m.mm.cur = m.back
			return nil

		} else if slices.ContainsFunc(KeyBinds.Left, func(v []byte) bool { return slices.Equal(v, in) }) {
			if m.back == nil {
				break
			}
			m.mm.cur = m.back
			return nil

		} else if i := slices.IndexFunc(KeyBinds.Numbers, func(v []byte) bool { return slices.Equal(v, in) }); i != -1 {
			if i > len(m.Items) {
				continue
			}
			m.selected = i - 1

			if m.selected < len(m.Items) && m.selected >= 0 {
				if m.Items[m.selected].Type() == "menu" {
					m.mm.cur = m.Items[m.selected].(*menu)
					return nil
				}
				err := m.Items[m.selected].enter()
				return err
			}
			if m.back == nil {
				return Errors.exit
			}
			m.mm.cur = m.back
			return nil
		}
	}
	return nil
}

// Add a new action to `m.Actions`.
//
// `callback` is called when this actions is selected.
//
// Returns a pointer to the new action.
func (m *menu) NewAction(name string, callback func()) *action {
	act := &action{
		mm:       m.mm,
		Name:     name,
		callback: callback,
	}
	m.Items = append(m.Items, act)
	return act
}

// Get the name as a string.
func (a *action) String() string { return a.Name }

// Get the name as a int.
func (a *action) Int() (int, error) { return strconv.Atoi(a.Name) }

// Get the name as a float.
func (a *action) Float() (float64, error) { return strconv.ParseFloat(a.Name, 64) }

// Get the name as a bool.
// It return true opon one of these values (case insensitive): 1, t, true, y, yes
func (a *action) Bool() bool {
	return slices.Contains([]string{"1", "t", "true", "y", "yes"}, strings.ToLower(a.Name))
}

// Get the type of the item.
//
// Returns `action`.
func (a *action) Type() string { return "action" }

func (a *action) enter() error {
	a.callback()
	return Errors.exit
}

// Add a new list to `m.Lists`.
//
// Only options in `l.Allowed` can be selected.
//
// Returns a pointer to the new list.
func (m *menu) NewList(name string) *list {
	lst := &list{
		mm:       m.mm,
		Name:     name,
		Allowed:  []string{"Yes", "No"},
		selected: 0,
	}
	m.Items = append(m.Items, lst)
	return lst
}

// Get the value as a string.
func (l *list) String() string { return l.Allowed[l.selected] }

// Get the value as a int.
func (l *list) Int() (int, error) { return strconv.Atoi(l.Allowed[l.selected]) }

// Get the value as a float.
func (l *list) Float() (float64, error) { return strconv.ParseFloat(l.Allowed[l.selected], 64) }

// Get the value as a bool.
// It return true opon one of these values (case insensitive): 1, t, true, y, yes
func (l *list) Bool() bool {
	return slices.Contains([]string{"1", "t", "true", "y", "yes"}, strings.ToLower(l.Allowed[l.selected]))
}

func (l *list) enter() error {
	l.editing = true
	_ = l.mm.rdr.Render()
	for {
		in := make([]byte, 3)
		if _, err := os.Stdin.Read(in); err != nil {
			l.editing = false
			_ = l.mm.rdr.Render()
			return err
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v []byte) bool { return slices.Equal(v, in) }) || slices.ContainsFunc(KeyBinds.Confirm, func(v []byte) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Up, func(v []byte) bool { return slices.Equal(v, in) }) {
			l.selected = max(l.selected-1, 0)
			_ = l.mm.rdr.Render()
			continue

		} else if slices.ContainsFunc(KeyBinds.Down, func(v []byte) bool { return slices.Equal(v, in) }) {
			l.selected = min(l.selected+1, len(l.Allowed)-1)
			_ = l.mm.rdr.Render()
			continue

		} else if slices.ContainsFunc(KeyBinds.Right, func(v []byte) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Left, func(v []byte) bool { return slices.Equal(v, in) }) {
			break
		} else if i := slices.IndexFunc(KeyBinds.Numbers, func(v []byte) bool { return slices.Equal(v, in) }); i != -1 {
			if i > len(l.Allowed)-1 {
				continue
			}
			l.selected = i
			_ = l.mm.rdr.Render()
			continue
		}

		for i, str := range l.Allowed {
			if !strings.HasPrefix(strings.ToLower(str), string(bytes.Trim(in, "\x00")[:])) {
				continue
			}
			l.selected = i
			_ = l.mm.rdr.Render()
			continue
		}

		for i, str := range l.Allowed {
			if !strings.HasPrefix(strings.ToLower(str), strings.ToLower(string(bytes.Trim(in, "\x00")[:]))) {
				continue
			}
			l.selected = i
			_ = l.mm.rdr.Render()
			continue
		}
	}
	l.editing = false
	_ = l.mm.rdr.Render()
	return nil
}

// Get the type of the item.
//
// Returns `list`.
func (l *list) Type() string { return "list" }

// Add a new option to `m.Options`.
//
// Only characters in `o.Allowed` can be entered.
// `value` is the default value and is not checked against `o.Allowed`.
//
// Returns a pointer to the new option.
func (m *menu) NewOption(name string, value string) *option {
	opt := &option{
		mm:      m.mm,
		Name:    name,
		Allowed: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		value:   value,
	}
	m.Items = append(m.Items, opt)
	return opt
}

// Get the value as a string.
func (o *option) String() string { return o.value }

// Get the value as a int.
func (o *option) Int() (int, error) { return strconv.Atoi(o.value) }

// Get the value as a float.
func (o *option) Float() (float64, error) { return strconv.ParseFloat(o.value, 64) }

// Get the value as a bool.
// It return true opon one of these values (case insensitive): 1, t, true, y, yes
func (o *option) Bool() bool {
	return slices.Contains([]string{"1", "t", "true", "y", "yes"}, strings.ToLower(o.value))
}

// Get the type of the item.
//
// Returns `option`.
func (o *option) Type() string { return "option" }

func (o *option) enter() error {
	o.editing = true
	_ = o.mm.rdr.Render()
	for {
		in := make([]byte, 3)
		if _, err := os.Stdin.Read(in); err != nil {
			o.editing = false
			_ = o.mm.rdr.Render()
			return err
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v []byte) bool { return slices.Equal(v, in) }) || slices.ContainsFunc(KeyBinds.Confirm, func(v []byte) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Delete, func(v []byte) bool { return slices.Equal(v, in) }) {
			if len(o.value) > 0 {
				o.value = o.value[:len(o.value)-1]
				_ = o.mm.rdr.Render()
				continue
			}
		}

		if strings.ContainsAny(o.Allowed, string(in[:])) {
			o.value += string(bytes.Trim(in, "\x00")[:])
			_ = o.mm.rdr.Render()
		}
	}

	o.editing = false
	_ = o.mm.rdr.Render()
	return nil
}
