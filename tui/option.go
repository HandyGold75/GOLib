package main

import (
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"
)

type option struct {
	Name          string
	Color         color
	AccentColor   color
	ValueColor    color
	SelectColor   color
	SelectBGColor color
	Allowed       string
	value         string
	editing       bool
	renderer      *func() error
}

// Add a new option to `m.Options`.
//
// Only characters in `o.Allowed` can be entered.
// `value` is the default value and is not checked against `o.Allowed`.
//
// Returns a pointer to the new option.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor`, `tui.Defaults.SelectBGColor`, `tui.Defaults.ValueColor` before creating menus.
func (m *menu) NewOption(name string, value string) *option {
	opt := &option{
		Name:          name,
		Color:         Defaults.Color,
		AccentColor:   Defaults.AccentColor,
		SelectColor:   Defaults.SelectColor,
		SelectBGColor: Defaults.SelectBGColor,
		ValueColor:    Defaults.ValueColor,
		Allowed:       "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		value:         value,
		renderer:      m.renderer,
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

func (o *option) edit() error {
	var e error
	o.editing = true
	_ = (*o.renderer)()

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
				_ = (*o.renderer)()
				continue
			}
		}

		if strings.ContainsAny(o.Allowed, string(in[:])) {
			o.value += string(bytes.Trim(in, "\x00")[:])
			_ = (*o.renderer)()
		}
	}

	o.editing = false
	_ = (*o.renderer)()
	return e
}
