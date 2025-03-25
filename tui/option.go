package main

import "strconv"

type option struct {
	Name        string
	Color       color
	AccentColor color
	ValueColor  color
	Allowed     string
	value       string
}

// Add a new option to `m.Options`.
//
// Only characters in `o.Allowed` can be entered.
// `value` is the default value and is not checked against `o.Allowed`.
//
// Returns a pointer to the new option.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.ValueColor` before creating options.
func (m *menu) NewOption(name string, value string) *option {
	opt := &option{
		Name:        name,
		Color:       Defaults.Color,
		AccentColor: Defaults.AccentColor,
		ValueColor:  Defaults.ValueColor,
		Allowed:     "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		value:       value,
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
