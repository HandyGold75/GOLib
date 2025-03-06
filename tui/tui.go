package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/term"
)

type (
	colors struct {
		Black, Red, Green, Yellow, Blue, Magenta, Cyan, White                                                                 color
		BrightBlack, BrightRed, BrightGreen, BrightYellow, BrightBlue, BrightMagenta, BrightCyan, BrightWhite                 color
		BGBlack, BGRed, BGGreen, BGYellow, BGBlue, BGMagenta, BGCyan, BGWhite                                                 color
		BGBrightBlack, BGBrightRed, BGBrightGreen, BGBrightYellow, BGBrightBlue, BGBrightMagenta, BGBrightCyan, BGBrightWhite color
		Reset                                                                                                                 color
	}
	color []byte

	aligns struct{ Left, Middle, Right align }
	align  string

	optionTypes struct{ String, Int, Float, Bool optionType }
	optionType  string

	keybinds struct{ Up, Down, Select, Back, Exit, Numbers [][]byte }

	MainMenu struct {
		Menu     *Menu
		current  *Menu
		selected int
		exit     chan error
	}

	Menu struct {
		Title         string
		Color         color
		AccentColor   color
		SelectColor   color
		SelectBGColor color
		Align         align
		Back          *Menu
		Menus         []*Menu
		Actions       []*Action
		Options       []*Option
	}

	Action struct {
		Name     string
		Color    color
		Callback func()
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
		Black:   []byte{27, '[', '3', '0', 'm'},
		Red:     []byte{27, '[', '3', '1', 'm'},
		Green:   []byte{27, '[', '3', '2', 'm'},
		Yellow:  []byte{27, '[', '3', '3', 'm'},
		Blue:    []byte{27, '[', '3', '4', 'm'},
		Magenta: []byte{27, '[', '3', '5', 'm'},
		Cyan:    []byte{27, '[', '3', '6', 'm'},
		White:   []byte{27, '[', '3', '7', 'm'},

		BrightBlack:   []byte{27, '[', '9', '0', 'm'},
		BrightRed:     []byte{27, '[', '9', '1', 'm'},
		BrightGreen:   []byte{27, '[', '9', '2', 'm'},
		BrightYellow:  []byte{27, '[', '9', '3', 'm'},
		BrightBlue:    []byte{27, '[', '9', '4', 'm'},
		BrightMagenta: []byte{27, '[', '9', '5', 'm'},
		BrightCyan:    []byte{27, '[', '9', '6', 'm'},
		BrightWhite:   []byte{27, '[', '9', '7', 'm'},

		BGBlack:   []byte{27, '[', '4', '0', 'm'},
		BGRed:     []byte{27, '[', '4', '1', 'm'},
		BGGreen:   []byte{27, '[', '4', '2', 'm'},
		BGYellow:  []byte{27, '[', '4', '3', 'm'},
		BGBlue:    []byte{27, '[', '4', '4', 'm'},
		BGMagenta: []byte{27, '[', '4', '5', 'm'},
		BGCyan:    []byte{27, '[', '4', '6', 'm'},
		BGWhite:   []byte{27, '[', '4', '7', 'm'},

		BGBrightBlack:   []byte{27, '[', '1', '0', '0', 'm'},
		BGBrightRed:     []byte{27, '[', '1', '0', '1', 'm'},
		BGBrightGreen:   []byte{27, '[', '1', '0', '2', 'm'},
		BGBrightYellow:  []byte{27, '[', '1', '0', '3', 'm'},
		BGBrightBlue:    []byte{27, '[', '1', '0', '4', 'm'},
		BGBrightMagenta: []byte{27, '[', '1', '0', '5', 'm'},
		BGBrightCyan:    []byte{27, '[', '1', '0', '6', 'm'},
		BGBrightWhite:   []byte{27, '[', '1', '0', '7', 'm'},

		Reset: []byte{27, '[', '0', 'm'},
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
	DefaultValueColor  = Colors.BrightWhite

	DefaultSelectColor   = Colors.Black
	DefaultSelectBGColor = Colors.BGWhite

	DefaultAlign = Aligns.Middle
	DefaultType  = Types.String

	// ESC: {27, 0, 0}, RETURN: {13, 0, 0}
	// CTRL_C: {3, 0, 0}, CTRL_D: {4, 0, 0}, Q: {113, 0, 0}
	// W: {119, 0, 0}, D: {100, 0, 0}, S: {115, 0, 0}, A: {97, 0, 0}
	// K: {108, 0, 0}, L: {107, 0, 0}, J: {106, 0, 0}, H: {104, 0, 0}
	// UP: {27, 91, 65}, RIGHT: {27, 91, 67}, DOWN: {27, 91, 66}, LEFT: {27, 91, 68}
	// Zero: {48, 0, 0}, One: {49, 0, 0}, Two: {50, 0, 0}, Three: {51, 0, 0}, For: {52, 0, 0}, Five: {53, 0, 0}, Six: {54, 0, 0}, Seven: {55, 0, 0}, Eight: {56, 0, 0}, Nine: {57, 0, 0}
	KeyBinds = keybinds{
		Up:      [][]byte{{119, 0, 0}, {108, 0, 0}, {27, 91, 65}},                                                                                 // W, K, UP
		Down:    [][]byte{{115, 0, 0}, {106, 0, 0}, {27, 91, 66}},                                                                                 // S, J, DOWN
		Select:  [][]byte{{100, 0, 0}, {107, 0, 0}, {27, 91, 67}, {13, 0, 0}},                                                                     // D, L, RIGHT, RETURN
		Back:    [][]byte{{97, 0, 0}, {104, 0, 0}, {27, 91, 68}},                                                                                  // A, H, LEFT
		Exit:    [][]byte{{27, 0, 0}, {3, 0, 0}, {4, 0, 0}, {113, 0, 0}},                                                                          // ESC, CTRL_C, CTRL_D, Q
		Numbers: [][]byte{{48, 0, 0}, {49, 0, 0}, {50, 0, 0}, {51, 0, 0}, {52, 0, 0}, {53, 0, 0}, {54, 0, 0}, {55, 0, 0}, {56, 0, 0}, {57, 0, 0}}, // 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
	}
)

// To set default colors set `tui.DefaultColor`, `tui.DefaultAccentColor` before creating menus.
//
// To set default alignment set `tui.DefaultAlign` before creating menus.
//
// Allowed colors are `tui.Colors.*`.
func NewMenu(title string) *MainMenu {
	menu := &Menu{
		Title:         title,
		Color:         DefaultColor,
		AccentColor:   DefaultAccentColor,
		SelectColor:   DefaultSelectColor,
		SelectBGColor: DefaultSelectBGColor,
		Align:         DefaultAlign,
		Back:          nil,
		Menus:         []*Menu{},
		Actions:       []*Action{},
		Options:       []*Option{},
	}
	return &MainMenu{
		Menu:     menu,
		current:  menu,
		selected: 0,
		exit:     make(chan error),
	}
}

// Add a new menu to m.Menus
//
// Returns a pointer to the new menu.
func (m *Menu) NewMenu(title string) *Menu {
	menu := &Menu{
		Title:         title,
		Color:         DefaultColor,
		AccentColor:   DefaultAccentColor,
		SelectColor:   DefaultSelectColor,
		SelectBGColor: DefaultSelectBGColor,
		Align:         DefaultAlign,
		Back:          m,
		Menus:         []*Menu{},
		Actions:       []*Action{},
		Options:       []*Option{},
	}
	m.Menus = append(m.Menus, menu)
	return menu
}

// Add a new action to m.Actions
//
// Returns a pointer to the new action.
//
// To set default colors set `tui.DefaultColor` before creating options.
func (m *Menu) NewAction(name string, callback func()) *Action {
	action := &Action{
		Name:     name,
		Color:    DefaultColor,
		Callback: callback,
	}
	m.Actions = append(m.Actions, action)
	return action
}

// Add a new option to m.Options
//
// Returns a pointer to the new option.
//
// To set default colors set `tui.DefaultColor`, `tui.DefaultAccentColor`, `tui.DefaultValueColor` before creating options.
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

func (m *Menu) Render(selected int) error {
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
		for _, menu := range m.Menus {
			itemLen += 1
			if itemLen == selected {
				lines = append(lines, slices.Concat(getCursorPos(len(menu.Title)+2, Aligns.Middle), m.SelectBGColor, m.SelectColor, []byte(menu.Title), Colors.Reset, m.AccentColor, []byte(" 🞂"), Colors.Reset))
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(menu.Title)+2, Aligns.Middle), menu.Color, []byte(menu.Title), m.AccentColor, []byte(" 🞂"), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	if len(m.Actions) > 0 {
		for _, action := range m.Actions {
			itemLen += 1
			if itemLen == selected {
				lines = append(lines, slices.Concat(getCursorPos(len(action.Name), Aligns.Middle), m.SelectBGColor, m.SelectColor, []byte(action.Name), Colors.Reset))
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(action.Name), Aligns.Middle), action.Color, []byte(action.Name), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	if len(m.Options) > 0 {
		for _, option := range m.Options {
			itemLen += 1
			if itemLen == selected {
				lines = append(lines, slices.Concat(getCursorPos(len(option.Name)+3+len(option.Value), Aligns.Middle), m.SelectBGColor, m.SelectColor, []byte(option.Name), Colors.Reset, option.AccentColor, []byte(" ▷ "), option.ValueColor, []byte(option.Value), Colors.Reset))
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(option.Name)+3+len(option.Value), Aligns.Middle), option.Color, []byte(option.Name), option.AccentColor, []byte(" ▷ "), option.ValueColor, []byte(option.Value), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	backText := "Exit"
	if m.Back != nil {
		backText = "Back"
	}
	itemLen += 1
	if itemLen == selected {
		lines = append(lines, slices.Concat(getCursorPos(len(backText)+2, Aligns.Middle), m.AccentColor, []byte("◀ "), m.SelectBGColor, m.SelectColor, []byte(backText), Colors.Reset))
	} else {
		lines = append(lines, slices.Concat(getCursorPos(len(backText)+2, Aligns.Middle), m.AccentColor, []byte("◀ "), m.Color, []byte(backText), Colors.Reset))
	}

	if _, err := trm.Write(slices.Concat([]byte("\033[2J\033[0;0H"), bytes.Join(lines, []byte("\r\n")))); err != nil {
		return err
	}
	return nil
}

func (m *Menu) Select(index int) (*Menu, *Action, *Option) {
	itemLen := -1
	for _, menu := range m.Menus {
		itemLen += 1
		if itemLen == index {
			return menu, nil, nil
		}
	}

	for _, action := range m.Actions {
		itemLen += 1
		if itemLen == index {
			return nil, action, nil
		}
	}

	for _, option := range m.Options {
		itemLen += 1
		if itemLen == index {
			return nil, nil, option
		}
	}

	return nil, nil, nil
}

// Restores term to `state` after the tui stops.
func (mm *MainMenu) Start(state *term.State) {
	go func() {
		defer term.Restore(int(os.Stdin.Fd()), state)

		var e error
		mm.current.Render(mm.selected)

		for {
			in := make([]byte, 3)
			if _, err := os.Stdin.Read(in); err != nil {
				e = err
				break
			}

			if slices.ContainsFunc(KeyBinds.Exit, func(v []byte) bool { return slices.Equal(v, in) }) {
				break
			} else if slices.ContainsFunc(KeyBinds.Up, func(v []byte) bool { return slices.Equal(v, in) }) {
				mm.selected = max(mm.selected-1, 0)
				mm.current.Render(mm.selected)
				continue
			} else if slices.ContainsFunc(KeyBinds.Down, func(v []byte) bool { return slices.Equal(v, in) }) {
				mm.selected = min(mm.selected+1, len(mm.current.Menus)+len(mm.current.Actions)+len(mm.current.Options))
				mm.current.Render(mm.selected)
				continue

			} else if slices.ContainsFunc(KeyBinds.Select, func(v []byte) bool { return slices.Equal(v, in) }) {
				menuTmp, actionTmp, optionTmp := mm.current.Select(mm.selected)
				if menuTmp != nil {
					mm.current = menuTmp
				} else if actionTmp != nil {
					actionTmp.Callback()
					break
				} else if optionTmp != nil {
				} else {
					if mm.current.Back == nil {
						break
					}
					mm.current = mm.current.Back
				}

				mm.current.Render(mm.selected)
				continue

			} else if slices.ContainsFunc(KeyBinds.Back, func(v []byte) bool { return slices.Equal(v, in) }) {
				if mm.current.Back == nil {
					break
				}
				mm.current = mm.current.Back
				mm.current.Render(mm.selected)
				continue

			} else if i := slices.IndexFunc(KeyBinds.Numbers, func(v []byte) bool { return slices.Equal(v, in) }); i != -1 {
				if i == 0 {
					if mm.current.Back == nil {
						break
					}
					mm.current = mm.current.Back
					mm.current.Render(mm.selected)
					continue
				}

				menuTmp, actionTmp, optionTmp := mm.current.Select(i - 1)
				if menuTmp != nil {
					mm.current = menuTmp
				} else if actionTmp != nil {
					actionTmp.Callback()
					break
				} else if optionTmp != nil {
				}

				mm.current.Render(mm.selected)
				continue
			}
		}

		if _, err := trm.Write([]byte("\033[2J\033[0;0H")); err != nil {
			mm.exit <- err
			return
		}
		mm.exit <- e
		return
	}()
}

func (mm *MainMenu) Join() {
	err := <-mm.exit
	if err != nil {
		fmt.Println(err)
	}
	close(mm.exit)
}

func main() {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		panic(errors.New("stdin/ stdout should be a terminal"))
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { term.Restore(int(os.Stdin.Fd()), oldState) }()

	menu := NewMenu("Some Title")

	menu1 := menu.Menu.NewMenu("A Sub menu")
	menu1.NewOption("some Option 1", "value")
	menu2 := menu.Menu.NewMenu("Antoher Sub menu")
	menu2.NewOption("some Option 2", "value")
	menu.Menu.NewOption("some Option", "value")
	menu.Menu.NewOption("somemore Option", "value")
	menu.Menu.NewAction("a Action", func() {})
	menu.Menu.NewAction("evenmore Action", func() {})

	menu.Start(oldState)
	menu.Join()

	if _, err := trm.Write([]byte("\033[2J\033[0;0H")); err != nil {
		panic(err)
	}
}
