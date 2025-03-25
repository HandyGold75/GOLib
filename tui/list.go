package main

import (
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
	value       string
}

// Add a new list to `m.Lists`.
//
// Only options in `l.Allowed` can be selected.
// `value` is the default value and is not checked against `l.Allowed`.
//
// Returns a pointer to the new list.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.ValueColor` before creating options.
func (m *menu) NewList(name string, value string) *list {
	lst := &list{
		Name:        name,
		Color:       Defaults.Color,
		AccentColor: Defaults.AccentColor,
		ValueColor:  Defaults.ValueColor,
		Allowed:     []string{"Yes", "No"},
		value:       value,
	}
	m.Lists = append(m.Lists, lst)
	return lst
}

// Get the value as a string.
func (l *list) Get() string { return l.value }

// Get the value as a int.
func (l *list) GetInt() (int, error) { return strconv.Atoi(l.value) }

// Get the value as a float.
func (l *list) GetFloat() (float64, error) { return strconv.ParseFloat(l.value, 64) }

// Get the value as a bool.
// It return true opon one of these values (case insensitive): 1, t, true
func (l *list) GetBool() bool {
	return slices.Contains([]string{"1", "t", "true"}, strings.ToLower(l.value))
}
