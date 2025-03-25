package main

import (
	"os"
	"slices"
)

type menu struct {
	Title         string
	Color         color
	AccentColor   color
	SelectColor   color
	SelectBGColor color
	Align         align
	selected      int
	back          *menu
	renderer      *func() error
	Menus         []*menu
	Actions       []*action
	Lists         []*list
	Options       []*option
}

// Add a new menu to `m.Menus`.
//
// Returns a pointer to the new menu.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor`, `tui.Defaults.SelectBGColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func (m *menu) NewMenu(title string) *menu {
	mn := &menu{
		Title:         title,
		Color:         Defaults.Color,
		AccentColor:   Defaults.AccentColor,
		SelectColor:   Defaults.SelectColor,
		SelectBGColor: Defaults.SelectBGColor,
		Align:         Defaults.Align,
		back:          m,
		renderer:      m.renderer,
		Menus:         []*menu{},
		Actions:       []*action{},
		Lists:         []*list{},
		Options:       []*option{},
	}
	m.Menus = append(m.Menus, mn)
	return mn
}

func (m *menu) right() (error, *menu) {
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
			_ = (*m.renderer)()
			continue

		} else if slices.ContainsFunc(KeyBinds.Down, func(v []byte) bool { return slices.Equal(v, in) }) {
			m.selected = min(m.selected+1, len(m.Menus)+len(m.Actions)+len(m.Lists)+len(m.Options))
			_ = (*m.renderer)()
			continue

		} else if slices.ContainsFunc(KeyBinds.Right, func(v []byte) bool { return slices.Equal(v, in) }) {
			err, mn := m.right()
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
			err, mn := m.right()
			if err != nil {
				e = err
				break
			}
			return mn, nil
		}
	}
	return nil, e
}
