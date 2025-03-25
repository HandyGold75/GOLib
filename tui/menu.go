package main

import (
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/term"
)

type menu struct {
	Title         string
	Color         color
	AccentColor   color
	SelectColor   color
	SelectBGColor color
	Align         align
	selected      int
	editing       bool
	back          *menu
	trm           *term.Terminal
	Menus         []*menu
	Actions       []*action
	Options       []*option
	Lists         []*list
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
		trm:           m.trm,
		Menus:         []*menu{},
		Actions:       []*action{},
		Options:       []*option{},
		Lists:         []*list{},
	}
	m.Menus = append(m.Menus, mn)
	return mn
}

func (m *menu) up() {
	m.selected = max(m.selected-1, 0)
	_ = m.Render()
}

func (m *menu) down() {
	m.selected = min(m.selected+1, len(m.Menus)+len(m.Actions)+len(m.Options))
	_ = m.Render()
}

func (m *menu) right() (error, *menu) {
	if s := m.selected; s < len(m.Menus) && s >= 0 {
		return nil, m.Menus[s]
	} else if s := m.selected - len(m.Menus); s < len(m.Actions) && s >= 0 {
		m.Actions[s].callback()
		return Errors.Exit, nil
	} else if s := m.selected - len(m.Menus) - len(m.Actions); s < len(m.Options) && s >= 0 {
		if err := m.editOption(m.Options[s]); err != nil {
			return err, nil
		}
		return nil, m
	}
	if m.back == nil {
		return Errors.Exit, nil
	}
	return nil, m.back
}

func (m *menu) editOption(o *option) error {
	var e error
	m.editing = true
	_ = m.Render()

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
				_ = m.Render()
				continue
			}
		}

		if strings.ContainsAny(o.Allowed, string(in[:])) {
			o.value += string(bytes.Trim(in, "\x00")[:])
			_ = m.Render()
		}
	}

	m.editing = false
	_ = m.Render()
	return e
}

// Render the current menu.
func (m *menu) Render() error {
	x, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	getCursorPos := func(textWidth int, alignment align) []byte {
		if alignment == Aligns.Left {
			return []byte{}
		} else if alignment == Aligns.Middle {
			return []byte("\033[" + strconv.Itoa(int((float64(x)/2)-(float64(textWidth)/2))) + "C")
		} else if alignment == Aligns.Right {
			return []byte("\033[" + strconv.Itoa(x-textWidth) + "C")
		}
		return []byte{}
	}

	itemLen := -1
	lines := append([][]byte{}, slices.Concat(getCursorPos(len(m.Title), Aligns.Middle), m.Color, []byte(m.Title), Colors.Reset))
	lines = append(lines, slices.Concat(m.AccentColor, []byte(strings.Repeat("─", x)), Colors.Reset))

	if len(m.Menus) > 0 {
		for _, mn := range m.Menus {
			itemLen += 1
			if itemLen == m.selected {
				lines = append(lines, slices.Concat(getCursorPos(len(mn.Title)+2, Aligns.Middle), m.SelectBGColor, m.SelectColor, []byte(mn.Title), Colors.Reset, m.AccentColor, []byte(" 🞂"), Colors.Reset))
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(mn.Title)+2, Aligns.Middle), mn.Color, []byte(mn.Title), m.AccentColor, []byte(" 🞂"), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	if len(m.Actions) > 0 {
		for _, act := range m.Actions {
			itemLen += 1
			if itemLen == m.selected {
				lines = append(lines, slices.Concat(getCursorPos(len(act.Name), Aligns.Middle), m.SelectBGColor, m.SelectColor, []byte(act.Name), Colors.Reset))
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(act.Name), Aligns.Middle), act.Color, []byte(act.Name), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	if len(m.Options) > 0 {
		for _, opt := range m.Options {
			itemLen += 1
			if itemLen == m.selected {
				if m.editing {
					lines = append(lines, slices.Concat(getCursorPos(len(opt.Name)+3+len(opt.value), Aligns.Middle), opt.Color, []byte(opt.Name), opt.AccentColor, []byte(" ▷ "), m.SelectBGColor, m.SelectColor, []byte(opt.value), Colors.Reset))
				} else {
					lines = append(lines, slices.Concat(getCursorPos(len(opt.Name)+3+len(opt.value), Aligns.Middle), m.SelectBGColor, m.SelectColor, []byte(opt.Name), Colors.Reset, opt.AccentColor, []byte(" ▷ "), opt.ValueColor, []byte(opt.value), Colors.Reset))
				}
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(opt.Name)+3+len(opt.value), Aligns.Middle), opt.Color, []byte(opt.Name), opt.AccentColor, []byte(" ▷ "), opt.ValueColor, []byte(opt.value), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	if len(m.Lists) > 0 {
		for _, lst := range m.Lists {
			itemLen += 1
			if itemLen == m.selected {
				if m.editing {
					lines = append(lines, slices.Concat(getCursorPos(len(lst.Name)+3+len(lst.value), Aligns.Middle), lst.Color, []byte(lst.Name), lst.AccentColor, []byte(" ▷ "), m.SelectBGColor, m.SelectColor, []byte(lst.value), Colors.Reset))
				} else {
					lines = append(lines, slices.Concat(getCursorPos(len(lst.Name)+3+len(lst.value), Aligns.Middle), m.SelectBGColor, m.SelectColor, []byte(lst.Name), Colors.Reset, lst.AccentColor, []byte(" ▷ "), lst.ValueColor, []byte(lst.value), Colors.Reset))
				}
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(lst.Name)+3+len(lst.value), Aligns.Middle), lst.Color, []byte(lst.Name), lst.AccentColor, []byte(" ▷ "), lst.ValueColor, []byte(lst.value), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	backText := "Exit"
	if m.back != nil {
		backText = "Back"
	}
	itemLen += 1
	if itemLen == m.selected {
		lines = append(lines, slices.Concat(getCursorPos(len(backText)+2, Aligns.Middle), m.AccentColor, []byte("◀ "), m.SelectBGColor, m.SelectColor, []byte(backText), Colors.Reset))
	} else {
		lines = append(lines, slices.Concat(getCursorPos(len(backText)+2, Aligns.Middle), m.AccentColor, []byte("◀ "), m.Color, []byte(backText), Colors.Reset))
	}

	if _, err := m.trm.Write(slices.Concat([]byte("\033[2J\033[0;0H"), bytes.Join(lines, []byte("\r\n")))); err != nil {
		return err
	}
	return nil
}
