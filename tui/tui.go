package main

import (
	"errors"
	"fmt"
	"io"
	"os"

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
		Menu     *menu
		cur      *menu
		trm      *term.Terminal
		renderer *func() error
		exit     chan error
		active   bool
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
		Color:         Colors.White,
		AccentColor:   Colors.Yellow,
		ValueColor:    Colors.Blue,
		SelectColor:   Colors.Black,
		SelectBGColor: Colors.BGYellow,
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
		renderer:      nil,
		Menus:         []*menu{},
		Actions:       []*action{},
		Options:       []*option{},
		Lists:         []*list{},
	}
	main := &MainMenu{
		Menu: mn,
		cur:  mn,
		trm: term.NewTerminal(struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}, ""),
		renderer: nil,
		exit:     make(chan error),
		active:   false,
	}

	r := func() error { return RenderBasic(main) }
	main.renderer = &r
	mn.renderer = &r

	return main, nil
}

// Get a new main menu with a custom renderer.
//
// Only 1 main menu should be active (started) at a time.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor`, `tui.Defaults.SelectBGColor`, `tui.Defaults.ValueColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func NewMenuCustom(title string, renderer func(*MainMenu) error) (*MainMenu, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, Errors.NotATerm
	}

	mn := &menu{
		Title:         title,
		Color:         Defaults.Color,
		AccentColor:   Defaults.AccentColor,
		SelectColor:   Defaults.SelectColor,
		SelectBGColor: Defaults.SelectBGColor,
		ValueColor:    Defaults.ValueColor,
		Align:         Defaults.Align,
		selected:      0,
		back:          nil,
		renderer:      nil,
		Menus:         []*menu{},
		Actions:       []*action{},
		Options:       []*option{},
		Lists:         []*list{},
	}
	main := &MainMenu{
		Menu: mn,
		cur:  mn,
		trm: term.NewTerminal(struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}, ""),
		renderer: nil,
		exit:     make(chan error),
		active:   false,
	}

	r := func() error { return renderer(main) }
	main.renderer = &r
	mn.renderer = &r

	return main, nil
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

		_ = (*mm.renderer)()
		for {
			mn, err := mm.cur.edit()
			if err != nil {
				if err == Errors.Exit {
					err = nil
				}
				mm.exit <- err
				close(mm.exit)
				return
			}
			if mn == nil {
				break
			}
			mm.cur = mn
			_ = (*mm.renderer)()
		}

		if _, err := mm.trm.Write([]byte("\033[2J\033[0;0H")); err != nil {
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

	mn, err := NewMenuCustom("Some Title", RenderBig)
	if err != nil {
		fmt.Println(err)
		return
	}

	mn1 := mn.Menu.NewMenu("A Sxb menu")
	mn1.NewOption("some Option 1", "val")
	mn2 := mn.Menu.NewMenu("Antoher Sub menu")
	mn2.NewOption("some Option 2", "val")
	lst := mn.Menu.NewList("some List")
	lst.Allowed = append(lst.Allowed, "Maybe")
	mn.Menu.NewOption("some Option a very nemeeeee", "someverutylongoon")
	mn.Menu.NewOption("somemore Option", "!!")
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

	if _, err := mn.trm.Write([]byte("\033[2J\033[0;0H")); err != nil {
		panic(err)
	}
}
