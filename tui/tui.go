package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"

	"golang.org/x/term"
)

type (
	tuiErrors struct{ Exit, NotATerm, TuiStarted, TuiNotStarted error }

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

	defaults struct {
		Color, AccentColor, ValueColor, SelectColor, SelectBGColor color
		Align                                                      align
	}

	keybinds struct{ Up, Down, Right, Left, Exit, Numbers, Confirm, Delete [][]byte }

	MainMenu struct {
		Menu   *menu
		cur    *menu
		exit   chan error
		active bool
	}
)

var (
	Errors = tuiErrors{
		Exit:          errors.New("Tui is exiting"),
		NotATerm:      errors.New("stdin/ stdout should be a terminal"),
		TuiStarted:    errors.New("Tui is already Started"),
		TuiNotStarted: errors.New("Tui is not yet Started"),
	}

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

	Defaults = defaults{
		Color:         Colors.Red,
		AccentColor:   Colors.Black,
		ValueColor:    Colors.BrightWhite,
		SelectColor:   Colors.Black,
		SelectBGColor: Colors.BGWhite,
		Align:         Aligns.Middle,
	}

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

// Get a new main menu.
//
// Only 1 main menu should be active (started) at a time.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor`, `tui.Defaults.SelectBGColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func NewMenu(title string) (*MainMenu, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, Errors.NotATerm
	}

	mn := &menu{
		Title:         title,
		Color:         Defaults.Color,
		AccentColor:   Defaults.AccentColor,
		SelectColor:   Defaults.SelectColor,
		SelectBGColor: Defaults.SelectBGColor,
		Align:         Defaults.Align,
		selected:      0,
		back:          nil,
		trm: term.NewTerminal(struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}, ""),
		Menus:   []*menu{},
		Actions: []*action{},
		Options: []*option{},
		Lists:   []*list{},
	}
	return &MainMenu{
		Menu:   mn,
		cur:    mn,
		exit:   make(chan error),
		active: false,
	}, nil
}

// Start tui, this will render and handle user input in a goroutine.
//
// Term should be in raw mode, the previous state returned by `term.MakeRaw` should be passed to `state`.
// `state` can also be `nil` and term will automaticlly be put term in raw mode.
//
// Restores term to `state` after the tui stops.
// When panicking outside of the goroutine the restore will not happen and will need to be done in the goroutine that is panicking.
//
// `mm.Join` should always be called to ensure the goroutine joins back.
func (mm *MainMenu) Start(state *term.State) error {
	if mm.active {
		return Errors.TuiStarted
	}
	mm.active = true

	if state == nil {
		s, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		state = s
	}

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
				if mm.cur.back == nil {
					break
				}
				mm.cur = mm.cur.back
				_ = mm.cur.Render()
				continue

			} else if i := slices.IndexFunc(KeyBinds.Numbers, func(v []byte) bool { return slices.Equal(v, in) }); i != -1 {
				if i > len(mm.cur.Menus)+len(mm.cur.Actions)+len(mm.cur.Options) {
					continue
				}
				mm.cur.selected = i - 1
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

		if e != nil {
			if e == Errors.Exit {
				e = nil
			}
			mm.exit <- e
		} else if _, err := mm.cur.trm.Write([]byte("\033[2J\033[0;0H")); err != nil {
			mm.exit <- err
		}
		close(mm.exit)
	}()

	return nil
}

// Waits until `mm.Start` has finished and returns errors generated by `mm.Start`.
func (mm *MainMenu) Join() error {
	if !mm.active {
		return Errors.TuiNotStarted
	}

	var e error
	for err := range mm.exit {
		e = err
	}

	mm.exit = make(chan error)
	mm.active = false

	return e
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()

	mn, err := NewMenu("Some Title")
	if err != nil {
		fmt.Println(err)
		return
	}

	mn1 := mn.Menu.NewMenu("A Sub menu")
	mn1.NewOption("some Option 1", "val")
	mn2 := mn.Menu.NewMenu("Antoher Sub menu")
	mn2.NewOption("some Option 2", "val")
	mn.Menu.NewOption("some Option", "val")
	mn.Menu.NewOption("somemore Option", "val")
	mn.Menu.NewAction("a Action", func() {})
	mn.Menu.NewAction("evenmore Action", func() {})

	if err := mn.Start(oldState); err != nil {
		fmt.Println(err)
		return
	}

	if err := mn.Join(); err != nil {
		fmt.Println(err)
		return
	}

	if _, err := mn.cur.trm.Write([]byte("\033[2J\033[0;0H")); err != nil {
		panic(err)
	}
}
