package tui

import (
	"errors"
	"os"

	"golang.org/x/term"
)

type (
	color   string
	align   int
	charSet string

	tuiErrors struct{ NotATerm, TuiStarted, TuiNotStarted, MainMenuNotHooked, exit error }

	defaults struct {
		Color, AccentColor, ValueColor, SelectColor color
		Align                                       align
	}

	keybinds struct{ Up, Down, Right, Left, Exit, Numbers, Confirm, Delete []keybind }
	keybind  []byte

	// TODO: Document
	MainMenu struct {
		Menu       *Menu
		statusline string
		cur        *Menu
		rdr        Renderer
		exit       chan error
		active     bool
	}
)

const (
	Reset color = "\033[0m"

	Bold            color = "\033[1m"
	Faint           color = "\033[2m"
	Italic          color = "\033[3m"
	Underline       color = "\033[4m"
	StrikeTrough    color = "\033[9m"
	DubbleUnderline color = "\033[21m"

	Black   color = "\033[30m"
	Red     color = "\033[31m"
	Green   color = "\033[32m"
	Yellow  color = "\033[33m"
	Blue    color = "\033[34m"
	Magenta color = "\033[35m"
	Cyan    color = "\033[36m"
	White   color = "\033[37m"

	BGBlack   color = "\033[40m"
	BGRed     color = "\033[41m"
	BGGreen   color = "\033[42m"
	BGYellow  color = "\033[43m"
	BGBlue    color = "\033[44m"
	BGMagenta color = "\033[45m"
	BGCyan    color = "\033[46m"
	BGWhite   color = "\033[47m"

	BrightBlack   color = "\033[90m"
	BrightRed     color = "\033[91m"
	BrightGreen   color = "\033[92m"
	BrightYellow  color = "\033[93m"
	BrightBlue    color = "\033[94m"
	BrightMagenta color = "\033[95m"
	BrightCyan    color = "\033[96m"
	BrightWhite   color = "\033[97m"

	BGBrightBlack   color = "\033[100m"
	BGBrightRed     color = "\033[101m"
	BGBrightGreen   color = "\033[102m"
	BGBrightYellow  color = "\033[103m"
	BGBrightBlue    color = "\033[104m"
	BGBrightMagenta color = "\033[105m"
	BGBrightCyan    color = "\033[106m"
	BGBrightWhite   color = "\033[107m"
)

const (
	AlignLeft align = iota
	AlignMiddle
	AlignRight
)

const (
	Letters        charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits         charSet = "0123456789"
	Hex            charSet = "0123456789abcdefABCDEF"
	WhiteSpace     charSet = " "
	Punctuation    charSet = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	GeneralCharSet charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

var (
	// TODO: Document
	Errors = tuiErrors{
		NotATerm:          errors.New("stdin/ stdout should be a terminal"),
		TuiStarted:        errors.New("tui is already Started"),
		TuiNotStarted:     errors.New("tui is not yet Started"),
		MainMenuNotHooked: errors.New("main menu not hooked to renderer"),

		exit: errors.New("tui is exiting"),
	}

	// TODO: Document
	Defaults = defaults{
		Color:       White,
		AccentColor: Yellow,
		ValueColor:  Bold + Blue,
		SelectColor: BGYellow + Black,
		Align:       AlignMiddle,
	}

	// TODO: Document
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
)

// Get a new main menu with a custom renderer.
//
// The renderer will automaticlly be hooked to the main menu.
//
// Only 1 main menu should be active (started) at a time.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor`, `tui.Defaults.ValueColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func NewMenu(name string, rdr Renderer) (*MainMenu, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, Errors.NotATerm
	}
	mn := &Menu{
		mm:          nil,
		name:        name,
		BackText:    "Exit",
		Color:       Defaults.Color,
		AccentColor: Defaults.AccentColor,
		SelectColor: Defaults.SelectColor,
		ValueColor:  Defaults.ValueColor,
		Align:       Defaults.Align,
		Items:       []MenuItem{},
		selected:    0,
		back:        nil,
	}
	main := &MainMenu{
		Menu:       mn,
		statusline: "",
		cur:        mn,
		rdr:        rdr,
		exit:       make(chan error),
		active:     false,
	}
	mn.mm = main
	rdr.HookMainMenu(main)

	return main, nil
}

// Get a new main menu with a basic renderer.
//
// Only 1 main menu should be active (started) at a time.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func NewMenuBasic(title string) (*MainMenu, error) {
	return NewMenu(title, newRendererBasic())
}

// Get a new main menu with a bulky renderer.
//
// Only 1 main menu should be active (started) at a time.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func NewMenuBulky(title string) (*MainMenu, error) {
	return NewMenu(title, newRendererBulky())
}

// Set the statusline.
//
// Will cause a rerender of the current menu.
func (mm *MainMenu) StatusLine(status string) {
	mm.statusline = status
	_ = mm.rdr.Render()
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
			err := mm.cur.Enter()
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
