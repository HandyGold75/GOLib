package main

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

type (
	tuiErrors struct{ NotATerm, TuiStarted, TuiNotStarted, MainMenuNotHooked, exit error }

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

	keybinds struct{ Up, Down, Right, Left, Exit, Numbers, Confirm, Delete []keybind }
	keybind  []byte

	charSets struct{ Letters, Digits, Hex, WhiteSpace, Punctuation, General charSet }
	charSet  string

	MainMenu struct {
		Menu   *menu
		cur    *menu
		rdr    renderer
		exit   chan error
		active bool
	}
)

var (
	Errors = tuiErrors{
		NotATerm:          errors.New("stdin/ stdout should be a terminal"),
		TuiStarted:        errors.New("tui is already Started"),
		TuiNotStarted:     errors.New("tui is not yet Started"),
		MainMenuNotHooked: errors.New("main menu not hooked to renderer"),

		exit: errors.New("tui is exiting"),
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
		Up:      []keybind{{119, 0, 0}, {107, 0, 0}, {27, 91, 65}, {43, 0, 0}},                                                                     // W, K, UP
		Down:    []keybind{{115, 0, 0}, {106, 0, 0}, {27, 91, 66}, {45, 0, 0}},                                                                     // S, J, DOWN
		Right:   []keybind{{100, 0, 0}, {108, 0, 0}, {27, 91, 67}, {13, 0, 0}},                                                                     // D, L, RIGHT, RETURN
		Left:    []keybind{{97, 0, 0}, {104, 0, 0}, {27, 91, 68}, {113, 0, 0}, {127, 0, 0}},                                                        // A, H, LEFT, Q, BACKSPACE
		Exit:    []keybind{{27, 0, 0}, {3, 0, 0}, {4, 0, 0}},                                                                                       // ESC, CTRL_C, CTRL_D,
		Numbers: []keybind{{48, 0, 0}, {49, 0, 0}, {50, 0, 0}, {51, 0, 0}, {52, 0, 0}, {53, 0, 0}, {54, 0, 0}, {55, 0, 0}, {56, 0, 0}, {57, 0, 0}}, // 0, 1, 2, 3, 4, 5, 6, 7, 8, 9
		Confirm: []keybind{{13, 0, 0}},                                                                                                             // RETURN
		Delete:  []keybind{{127, 0, 0}, {27, 91, 51}},                                                                                              // BACKSPACE, DEL
	}

	CharSets = charSets{
		Letters:     "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Digits:      "0123456789",
		Hex:         "0123456789abcdefABCDEF",
		WhiteSpace:  " ",
		Punctuation: "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
		General:     "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
	}
)

// Get a new main menu with a custom renderer.
//
// The renderer will automaticlly be hooked to the main menu.
//
// Only 1 main menu should be active (started) at a time.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor`, `tui.Defaults.SelectBGColor`, `tui.Defaults.ValueColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func NewMenu(name string, rdr renderer) (*MainMenu, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, Errors.NotATerm
	}

	mn := &menu{
		mm:            nil,
		name:          name,
		BackText:      "Exit",
		Color:         Defaults.Color,
		AccentColor:   Defaults.AccentColor,
		SelectColor:   Defaults.SelectColor,
		SelectBGColor: Defaults.SelectBGColor,
		ValueColor:    Defaults.ValueColor,
		Align:         Defaults.Align,
		Items:         []item{},
		selected:      0,
		back:          nil,
	}
	main := &MainMenu{
		Menu:   mn,
		cur:    mn,
		rdr:    rdr,
		exit:   make(chan error),
		active: false,
	}
	mn.mm = main
	rdr.HookMainMenu(main)

	return main, nil
}

// Get a new main menu with the basic renderer.
//
// Only 1 main menu should be active (started) at a time.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor`, `tui.Defaults.SelectBGColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func NewMenuBasic(title string) (*MainMenu, error) {
	return NewMenu(title, NewBasic())
}

// Get a new main menu with the bulky renderer.
//
// Only 1 main menu should be active (started) at a time.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor`, `tui.Defaults.SelectBGColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func NewMenuBulky(title string) (*MainMenu, error) {
	return NewMenu(title, NewBulky())
}

// Set the statusline.
//
// Will cause a rerender of the current menu.
func (mm *MainMenu) StatusLine(status string) {
	mm.rdr.StatusLine(status)
	mm.rdr.Render()
}

// Start tui, this will render and handle user input in a goroutine.
//
// Term should be in raw mode, the previous state returned by `term.MakeRaw` should be passed to `state`.
// `state` can also be `nil` and term will automaticlly be put in raw mode.
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

		_ = mm.rdr.Render()
		for {
			err := mm.cur.enter()
			if err != nil {
				if err == Errors.exit {
					break
				}
				mm.exit <- err
				close(mm.exit)
				return
			}
			_ = mm.rdr.Render()
		}

		if err := mm.rdr.Clear(); err != nil {
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

	mm, err := NewMenuBasic("Some Title")
	if err != nil {
		fmt.Println(err)
		return
	}
	mm.StatusLine("A status in a line")

	_ = mm.Menu.NewMenu("A menu")
	mm.Menu.NewText("A text", CharSets.General, "val")
	mm.Menu.NewAction("A action", func() {})
	mm.Menu.NewList("A list", []string{"Yes", "No", "Maybe"})
	mm.Menu.NewDigit("A digit", 99, -128, 127)
	mm.Menu.NewIPv4("A ipv4", "127.0.0.1")
	mm.Menu.NewIPv6("A ipv4", "::1")

	if err := mm.Start(oldState); err != nil {
		fmt.Println(err)
		return
	}

	if err := mm.Join(); err != nil {
		fmt.Println(err)
		return
	}
}
