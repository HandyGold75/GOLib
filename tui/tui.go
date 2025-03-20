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

	keybinds struct{ Up, Down, Right, Left, Exit, Numbers, Confirm, Delete [][]byte }

	MainMenu struct {
		Menu *menu
		cur  *menu
		exit chan error
	}

	menu struct {
		Title         string
		Color         color
		AccentColor   color
		SelectColor   color
		SelectBGColor color
		Align         align
		Selected      int
		Editing       bool
		Back          *menu
		Menus         []*menu
		Actions       []*action
		Options       []*option
	}

	action struct {
		Name     string
		Color    color
		Callback func()
	}

	option struct {
		Name        string
		Color       color
		AccentColor color
		ValueColor  color
		Allowed     string
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

	KeyBinds = keybinds{
		Up:      [][]byte{{119, 0, 0}, {107, 0, 0}, {27, 91, 65}},                                                                                 // W, K, UP
		Down:    [][]byte{{115, 0, 0}, {106, 0, 0}, {27, 91, 66}},                                                                                 // S, J, DOWN
		Right:   [][]byte{{100, 0, 0}, {108, 0, 0}, {27, 91, 67}, {13, 0, 0}},                                                                     // D, L, RIGHT, RETURN
		Left:    [][]byte{{97, 0, 0}, {104, 0, 0}, {27, 91, 68}, {113, 0, 0}, {127, 0, 0}},                                                        // A, H, LEFT, Q, BACKSPACE
		Exit:    [][]byte{{27, 0, 0}, {3, 0, 0}, {4, 0, 0}},                                                                                       // ESC, CTRL_C, CTRL_D,
		Numbers: [][]byte{{48, 0, 0}, {49, 0, 0}, {50, 0, 0}, {51, 0, 0}, {52, 0, 0}, {53, 0, 0}, {54, 0, 0}, {55, 0, 0}, {56, 0, 0}, {57, 0, 0}}, // 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
		Confirm: [][]byte{{13, 0, 0}},                                                                                                             // RETURN
		Delete:  [][]byte{{127, 0, 0}, {27, 91, 51}},                                                                                              // BACKSPACE, DEL
	}
)

// To set default colors set `tui.DefaultColor`, `tui.DefaultAccentColor` before creating menus.
//
// To set default alignment set `tui.DefaultAlign` before creating menus.
//
// Allowed colors are `tui.Colors.*`.
func NewMenu(title string) *MainMenu {
	mn := &menu{
		Title:         title,
		Color:         DefaultColor,
		AccentColor:   DefaultAccentColor,
		SelectColor:   DefaultSelectColor,
		SelectBGColor: DefaultSelectBGColor,
		Align:         DefaultAlign,
		Selected:      0,
		Back:          nil,
		Menus:         []*menu{},
		Actions:       []*action{},
		Options:       []*option{},
	}
	return &MainMenu{
		Menu: mn,
		cur:  mn,
		exit: make(chan error),
	}
}

// Add a new menu to m.Menus
//
// Returns a pointer to the new menu.
func (m *menu) NewMenu(title string) *menu {
	mn := &menu{
		Title:         title,
		Color:         DefaultColor,
		AccentColor:   DefaultAccentColor,
		SelectColor:   DefaultSelectColor,
		SelectBGColor: DefaultSelectBGColor,
		Align:         DefaultAlign,
		Back:          m,
		Menus:         []*menu{},
		Actions:       []*action{},
		Options:       []*option{},
	}
	m.Menus = append(m.Menus, mn)
	return mn
}

// Add a new action to m.Actions
//
// Returns a pointer to the new action.
//
// To set default colors set `tui.DefaultColor` before creating options.
func (m *menu) NewAction(name string, callback func()) *action {
	act := &action{
		Name:     name,
		Color:    DefaultColor,
		Callback: callback,
	}
	m.Actions = append(m.Actions, act)
	return act
}

// Add a new option to m.Options
//
// Returns a pointer to the new option.
//
// To set default colors set `tui.DefaultColor`, `tui.DefaultAccentColor`, `tui.DefaultValueColor` before creating options.
//
// To set default types set `tui.DefaultType` before creating options.
func (m *menu) NewOption(name string, value string) *option {
	opt := &option{
		Name:        name,
		Color:       DefaultColor,
		AccentColor: DefaultAccentColor,
		ValueColor:  DefaultValueColor,
		Allowed:     "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		Type:        DefaultType,
		Value:       value,
	}
	m.Options = append(m.Options, opt)
	return opt
}

func (m *menu) up() {
	m.Selected = max(m.Selected-1, 0)
	_ = m.Render()
}

func (m *menu) down() {
	m.Selected = min(m.Selected+1, len(m.Menus)+len(m.Actions)+len(m.Options))
	_ = m.Render()
}

func (m *menu) right() (error, *menu) {
	if s := m.Selected; s < len(m.Menus) && s >= 0 {
		return nil, m.Menus[s]
	} else if s := m.Selected - len(m.Menus); s < len(m.Actions) && s >= 0 {
		m.Actions[s].Callback()
		return errors.New("exit"), nil
	} else if s := m.Selected - len(m.Menus) - len(m.Actions); s < len(m.Options) && s >= 0 {
		if err := m.editOption(m.Options[s]); err != nil {
			return err, nil
		}
		return nil, m
	}
	if m.Back == nil {
		return errors.New("exit"), nil
	}
	return nil, m.Back
}

func (m *menu) editOption(o *option) error {
	var e error
	m.Editing = true
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
			if len(o.Value) > 0 {
				o.Value = o.Value[:len(o.Value)-1]
				_ = m.Render()
				continue
			}
		}

		if strings.ContainsAny(o.Allowed, string(in[:])) {
			o.Value += string(bytes.Trim(in, "\x00")[:])
			_ = m.Render()
		}
	}

	m.Editing = false
	_ = m.Render()
	return e
}

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
			if itemLen == m.Selected {
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
			if itemLen == m.Selected {
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
			if itemLen == m.Selected {
				if m.Editing {
					lines = append(lines, slices.Concat(getCursorPos(len(opt.Name)+3+len(opt.Value), Aligns.Middle), opt.Color, []byte(opt.Name), opt.AccentColor, []byte(" ▷ "), m.SelectBGColor, m.SelectColor, []byte(opt.Value), Colors.Reset))
				} else {
					lines = append(lines, slices.Concat(getCursorPos(len(opt.Name)+3+len(opt.Value), Aligns.Middle), m.SelectBGColor, m.SelectColor, []byte(opt.Name), Colors.Reset, opt.AccentColor, []byte(" ▷ "), opt.ValueColor, []byte(opt.Value), Colors.Reset))
				}
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(opt.Name)+3+len(opt.Value), Aligns.Middle), opt.Color, []byte(opt.Name), opt.AccentColor, []byte(" ▷ "), opt.ValueColor, []byte(opt.Value), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	backText := "Exit"
	if m.Back != nil {
		backText = "Back"
	}
	itemLen += 1
	if itemLen == m.Selected {
		lines = append(lines, slices.Concat(getCursorPos(len(backText)+2, Aligns.Middle), m.AccentColor, []byte("◀ "), m.SelectBGColor, m.SelectColor, []byte(backText), Colors.Reset))
	} else {
		lines = append(lines, slices.Concat(getCursorPos(len(backText)+2, Aligns.Middle), m.AccentColor, []byte("◀ "), m.Color, []byte(backText), Colors.Reset))
	}

	if _, err := trm.Write(slices.Concat([]byte("\033[2J\033[0;0H"), bytes.Join(lines, []byte("\r\n")))); err != nil {
		return err
	}
	return nil
}

// Restores term to `state` after the tui stops.
func (mm *MainMenu) Start(state *term.State) {
	go func() {
		defer func() { _ = term.Restore(int(os.Stdin.Fd()), state) }()

		var e error
		_ = mm.cur.Render()

		for {
			in := make([]byte, 3)
			if _, err := os.Stdin.Read(in); err != nil {
				e = err
				break
			}

			if slices.ContainsFunc(KeyBinds.Exit, func(v []byte) bool { return slices.Equal(v, in) }) {
				break
			} else if slices.ContainsFunc(KeyBinds.Up, func(v []byte) bool { return slices.Equal(v, in) }) {
				mm.cur.up()
				continue

			} else if slices.ContainsFunc(KeyBinds.Down, func(v []byte) bool { return slices.Equal(v, in) }) {
				mm.cur.down()
				continue

			} else if slices.ContainsFunc(KeyBinds.Right, func(v []byte) bool { return slices.Equal(v, in) }) {
				err, mn := mm.cur.right()
				if err != nil {
					e = err
					break
				}
				mm.cur = mn
				_ = mm.cur.Render()
				continue

			} else if slices.ContainsFunc(KeyBinds.Left, func(v []byte) bool { return slices.Equal(v, in) }) {
				if mm.cur.Back == nil {
					break
				}
				mm.cur = mm.cur.Back
				_ = mm.cur.Render()
				continue

			} else if i := slices.IndexFunc(KeyBinds.Numbers, func(v []byte) bool { return slices.Equal(v, in) }); i != -1 {
				if i > len(mm.cur.Menus)+len(mm.cur.Actions)+len(mm.cur.Options) {
					continue
				}
				mm.cur.Selected = i - 1
				err, mn := mm.cur.right()
				if err != nil {
					e = err
					break
				}
				mm.cur = mn
				_ = mm.cur.Render()
				continue
			}
		}

		if _, err := trm.Write([]byte("\033[2J\033[0;0H")); err != nil {
			mm.exit <- err
			return
		}
		mm.exit <- e
	}()
}

func (mm *MainMenu) Join() error {
	err := <-mm.exit
	close(mm.exit)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		panic(errors.New("stdin/ stdout should be a terminal"))
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	mn := NewMenu("Some Title")

	mn1 := mn.Menu.NewMenu("A Sub menu")
	mn1.NewOption("some Option 1", "value")
	mn2 := mn.Menu.NewMenu("Antoher Sub menu")
	mn2.NewOption("some Option 2", "value")
	mn.Menu.NewOption("some Option", "value")
	mn.Menu.NewOption("somemore Option", "value")
	mn.Menu.NewAction("a Action", func() {})
	mn.Menu.NewAction("evenmore Action", func() {})

	mn.Start(oldState)
	if err := mn.Join(); err != nil {
		fmt.Println(err)
	}

	if _, err := trm.Write([]byte("\033[2J\033[0;0H")); err != nil {
		panic(err)
	}
}
