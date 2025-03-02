package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/term"
)

type (
	colors struct{ Black, Red, Green, Yellow, Blue, Magenta, Cyan, White, Reset color }
	color  []byte

	aligns struct{ Left, Middle, Right align }
	align  string

	optionTypes struct{ String, Int, Float, Bool optionType }
	optionType  string

	Menu struct {
		Title       string
		Color       color
		AccentColor color
		Align       align
		Menus       []*Menu
		Options     []*Option
	}

	Option struct {
		Name        string
		Color       color
		AccentColor color
		ValueColor  color
		Type        optionType
		Value       string
	}
)

var (
	trm = term.NewTerminal(struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}, "")

	Colors = colors{
		Black:   trm.Escape.Black,
		Red:     trm.Escape.Red,
		Green:   trm.Escape.Green,
		Yellow:  trm.Escape.Yellow,
		Blue:    trm.Escape.Blue,
		Magenta: trm.Escape.Magenta,
		Cyan:    trm.Escape.Cyan,
		White:   trm.Escape.White,
		Reset:   trm.Escape.Reset,
	}

	Aligns = aligns{
		Left:   "Left",
		Middle: "Middle",
		Right:  "Right",
	}

	Types = optionTypes{
		String: "String",
		Int:    "Int",
		Float:  "Float",
		Bool:   "Bool",
	}

	DefaultColor       = Colors.Red
	DefaultAccentColor = Colors.Black
	DefaultValueColor  = Colors.White

	DefaultAlign = Aligns.Middle
	DefaultType  = Types.String
)

// To set default colors set `tui.DefaultColor`, `tui.DefaultAccentColor` before creating menus.
//
// To set default alignment set `tui.DefaultAlign` before creating menus.
//
// Allowed colors are `tui.Colors.*`.
func NewMenu(title string) *Menu {
	return &Menu{
		Title:       title,
		Color:       DefaultColor,
		AccentColor: DefaultAccentColor,
		Align:       DefaultAlign,
		Menus:       []*Menu{},
		Options:     []*Option{},
	}
}

// Add a new menu to m.Menus
//
// Returns a pointer to the new menu.
func (m *Menu) NewMenu(title string) *Menu {
	menu := &Menu{
		Title:       title,
		Color:       DefaultColor,
		AccentColor: DefaultAccentColor,
		Align:       DefaultAlign,
		Menus:       []*Menu{},
		Options:     []*Option{},
	}
	m.Menus = append(m.Menus, menu)
	return menu
}

// Add a new option to m.Options
//
// Returns a pointer to the new option.
//
// To set default colors set `tui.DefaultColor`, `tui.DefaultAccentColor`, `tui.DefaultValueColor` before creating options.
//
// To set default alignment set `tui.DefaultAlign` before creating options.
//
// To set default types set `tui.DefaultType` before creating options.
func (m *Menu) NewOption(name string, value string) *Option {
	option := &Option{
		Name:        name,
		Color:       DefaultColor,
		AccentColor: DefaultAccentColor,
		ValueColor:  DefaultValueColor,
		Type:        DefaultType,
		Value:       value,
	}
	m.Options = append(m.Options, option)
	return option
}

func (m *Menu) Render() error {
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

	lines := append([][]byte{}, slices.Concat(getCursorPos(len(m.Title), Aligns.Middle), m.Color, []byte(m.Title), Colors.Reset))
	lines = append(lines, slices.Concat(m.AccentColor, []byte(strings.Repeat("─", x)), Colors.Reset))
	for _, menu := range m.Menus {
		lines = append(lines, slices.Concat(getCursorPos(len(menu.Title)+2, Aligns.Middle), menu.Color, []byte(menu.Title), m.AccentColor, []byte(" 🞂"), Colors.Reset))
	}
	lines = append(lines, []byte{})
	for _, option := range m.Options {
		lines = append(lines, slices.Concat(getCursorPos(len(option.Name)+3+len(option.Value), Aligns.Middle), option.Color, []byte(option.Name), option.AccentColor, []byte(" ▷ "), option.ValueColor, []byte(option.Value), Colors.Reset))
	}

	if _, err := trm.Write(slices.Concat([]byte("\033[2J\033[0;0H"), bytes.Join(lines, []byte("\r\n")))); err != nil {
		return err
	}
	return nil
}

func main() {
	// if !term.IsTerminal(int(os.Stdin.Fd())) {
	// 	return errors.New("stdin/ stdout should be a terminal")
	// }

	// oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() { term.Restore(int(os.Stdin.Fd()), oldState) }()

	menu := NewMenu("Some Title")
	menu.NewMenu("A Sub menu")
	menu.NewMenu("Antoher Sub menu")
	menu.NewOption("some Option", "value")
	menu.NewOption("someother Option", "value")
	menu.NewOption("somemore Option", "value")
	menu.Render()

	fmt.Println("")
}
