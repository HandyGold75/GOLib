package main

import (
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"
)

type (
	menu struct {
		Title         string
		Color         color
		AccentColor   color
		SelectColor   color
		SelectBGColor color
		ValueColor    color
		Align         align
		selected      int
		back          *menu
		rdr           renderer
		Menus         []*menu
		Actions       []*action
		Lists         []*list
		Options       []*option
	}

	item interface {
		Get() string
		GetInt() (int, error)
		GetFloat() (float64, error)
		GetBool() bool
		edit()
	}

	action struct {
		Name     string
		callback func()
	}

	list struct {
		Name     string
		Allowed  []string
		selected int
		editing  bool
		rdr      renderer
	}

	option struct {
		Name    string
		Allowed string
		value   string
		editing bool
		rdr     renderer
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
		Title:         title,
		Color:         Defaults.Color,
		AccentColor:   Defaults.AccentColor,
		SelectColor:   Defaults.SelectColor,
		SelectBGColor: Defaults.SelectBGColor,
		ValueColor:    Defaults.ValueColor,
		Align:         Defaults.Align,
		back:          m,
		rdr:           m.rdr,
		Menus:         []*menu{},
		Actions:       []*action{},
		Lists:         []*list{},
		Options:       []*option{},
	}
	m.Menus = append(m.Menus, mn)
	return mn
}

func (m *menu) enter() (error, *menu) {
	if s := m.selected; s < len(m.Menus) && s >= 0 {
		return nil, m.Menus[s]
	} else if s := m.selected - len(m.Menus); s < len(m.Actions) && s >= 0 {
		m.Actions[s].callback()
		return Errors.Exit, nil
	} else if s := m.selected - len(m.Menus) - len(m.Actions); s < len(m.Lists) && s >= 0 {
		err := m.Lists[s].edit()
		return err, m
	} else if s := m.selected - len(m.Menus) - len(m.Actions) - len(m.Lists); s < len(m.Options) && s >= 0 {
		err := m.Options[s].edit()
		return err, m
	}
	if m.back == nil {
		return Errors.Exit, nil
	}
	return nil, m.back
}

func (m *menu) edit() (*menu, error) {
	var e error
	for {
		in := make([]byte, 3)
		if _, err := os.Stdin.Read(in); err != nil {
			e = err
			break
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v []byte) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Up, func(v []byte) bool { return slices.Equal(v, in) }) {
			m.selected = max(m.selected-1, 0)
			_ = m.rdr.Render()
			continue

		} else if slices.ContainsFunc(KeyBinds.Down, func(v []byte) bool { return slices.Equal(v, in) }) {
			m.selected = min(m.selected+1, len(m.Menus)+len(m.Actions)+len(m.Lists)+len(m.Options))
			_ = m.rdr.Render()
			continue

		} else if slices.ContainsFunc(KeyBinds.Right, func(v []byte) bool { return slices.Equal(v, in) }) {
			err, mn := m.enter()
			if err != nil {
				e = err
				break
			}
			return mn, nil

		} else if slices.ContainsFunc(KeyBinds.Left, func(v []byte) bool { return slices.Equal(v, in) }) {
			if m.back == nil {
				break
			}
			return m.back, nil

		} else if i := slices.IndexFunc(KeyBinds.Numbers, func(v []byte) bool { return slices.Equal(v, in) }); i != -1 {
			if i > len(m.Menus)+len(m.Actions)+len(m.Options) {
				continue
			}
			m.selected = i - 1
			err, mn := m.enter()
			if err != nil {
				e = err
				break
			}
			return mn, nil
		}
	}
	return nil, e
}

// Add a new action to `m.Actions`.
//
// `callback` is called when this actions is selected.
//
// Returns a pointer to the new action.
func (m *menu) NewAction(name string, callback func()) *action {
	act := &action{
		Name:     name,
		callback: callback,
	}
	m.Actions = append(m.Actions, act)
	return act
}

// Add a new list to `m.Lists`.
//
// Only options in `l.Allowed` can be selected.
//
// Returns a pointer to the new list.
func (m *menu) NewList(name string) *list {
	lst := &list{
		Name:     name,
		Allowed:  []string{"Yes", "No"},
		selected: 0,
		rdr:      m.rdr,
	}
	m.Lists = append(m.Lists, lst)
	return lst
}

// Get the value as a string.
func (l *list) Get() string { return l.Allowed[l.selected] }

// Get the value as a int.
func (l *list) GetInt() (int, error) { return strconv.Atoi(l.Allowed[l.selected]) }

// Get the value as a float.
func (l *list) GetFloat() (float64, error) { return strconv.ParseFloat(l.Allowed[l.selected], 64) }

// Get the value as a bool.
// It return true opon one of these values (case insensitive): 1, t, true, y, yes
func (l *list) GetBool() bool {
	return slices.Contains([]string{"1", "t", "true", "y", "yes"}, strings.ToLower(l.Allowed[l.selected]))
}

func (l *list) edit() error {
	var e error
	l.editing = true
	_ = l.rdr.Render()

	for {
		in := make([]byte, 3)
		if _, err := os.Stdin.Read(in); err != nil {
			e = err
			break
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v []byte) bool { return slices.Equal(v, in) }) || slices.ContainsFunc(KeyBinds.Confirm, func(v []byte) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Up, func(v []byte) bool { return slices.Equal(v, in) }) {
			l.selected = max(l.selected-1, 0)
			_ = l.rdr.Render()
			continue

		} else if slices.ContainsFunc(KeyBinds.Down, func(v []byte) bool { return slices.Equal(v, in) }) {
			l.selected = min(l.selected+1, len(l.Allowed)-1)
			_ = l.rdr.Render()
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
			_ = l.rdr.Render()
			continue
		}

		for i, str := range l.Allowed {
			if !strings.HasPrefix(strings.ToLower(str), string(bytes.Trim(in, "\x00")[:])) {
				continue
			}
			l.selected = i
			_ = l.rdr.Render()
			continue
		}

		for i, str := range l.Allowed {
			if !strings.HasPrefix(strings.ToLower(str), strings.ToLower(string(bytes.Trim(in, "\x00")[:]))) {
				continue
			}
			l.selected = i
			_ = l.rdr.Render()
			continue
		}
	}

	l.editing = false
	_ = l.rdr.Render()
	return e
}

// Add a new option to `m.Options`.
//
// Only characters in `o.Allowed` can be entered.
// `value` is the default value and is not checked against `o.Allowed`.
//
// Returns a pointer to the new option.
func (m *menu) NewOption(name string, value string) *option {
	opt := &option{
		Name:    name,
		Allowed: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		value:   value,
		rdr:     m.rdr,
	}
	m.Options = append(m.Options, opt)
	return opt
}

// Get the value as a string.
func (o *option) Get() string { return o.value }

// Get the value as a int.
func (o *option) GetInt() (int, error) { return strconv.Atoi(o.value) }

// Get the value as a float.
func (o *option) GetFloat() (float64, error) { return strconv.ParseFloat(o.value, 64) }

// Get the value as a bool.
// It return true opon one of these values (case insensitive): 1, t, true, y, yes
func (o *option) GetBool() bool {
	return slices.Contains([]string{"1", "t", "true", "y", "yes"}, strings.ToLower(o.value))
}

func (o *option) edit() error {
	var e error
	o.editing = true
	_ = o.rdr.Render()

	for {
		in := make([]byte, 3)
		if _, err := os.Stdin.Read(in); err != nil {
			e = err
			break
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v []byte) bool { return slices.Equal(v, in) }) || slices.ContainsFunc(KeyBinds.Confirm, func(v []byte) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Delete, func(v []byte) bool { return slices.Equal(v, in) }) {
			if len(o.value) > 0 {
				o.value = o.value[:len(o.value)-1]
				_ = o.rdr.Render()
				continue
			}
		}

		if strings.ContainsAny(o.Allowed, string(in[:])) {
			o.value += string(bytes.Trim(in, "\x00")[:])
			_ = o.rdr.Render()
		}
	}

	o.editing = false
	_ = o.rdr.Render()
	return e
}
