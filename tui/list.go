package main

import (
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"
)

type list struct {
	Name        string
	Color       color
	AccentColor color
	ValueColor  color
	Allowed     []string
	selected    int
	editing     bool
	renderer    *func() error
}

// Add a new list to `m.Lists`.
//
// Only options in `l.Allowed` can be selected.
//
// Returns a pointer to the new list.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.ValueColor` before creating options.
func (m *menu) NewList(name string) *list {
	lst := &list{
		Name:        name,
		Color:       Defaults.Color,
		AccentColor: Defaults.AccentColor,
		ValueColor:  Defaults.ValueColor,
		Allowed:     []string{"Yes", "No"},
		selected:    0,
		renderer:    m.renderer,
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
	_ = (*l.renderer)()

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
			_ = (*l.renderer)()
			continue

		} else if slices.ContainsFunc(KeyBinds.Down, func(v []byte) bool { return slices.Equal(v, in) }) {
			l.selected = min(l.selected+1, len(l.Allowed)-1)
			_ = (*l.renderer)()
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
			_ = (*l.renderer)()
			continue
		}

		for i, str := range l.Allowed {
			if !strings.HasPrefix(strings.ToLower(str), string(bytes.Trim(in, "\x00")[:])) {
				continue
			}
			l.selected = i
			_ = (*l.renderer)()
			continue
		}

		for i, str := range l.Allowed {
			if !strings.HasPrefix(strings.ToLower(str), strings.ToLower(string(bytes.Trim(in, "\x00")[:]))) {
				continue
			}
			l.selected = i
			_ = (*l.renderer)()
			continue
		}
	}

	l.editing = false
	_ = (*l.renderer)()
	return e
}
