package main

import (
	"bytes"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/term"
)

type renderer interface {
	// Render the current menu of the hooked main menu.
	Render() error
	// Set the statusline, will be displayed on the next call to `rdr.Render`.
	StatusLine(string)
	// Instantly clear the screen.
	Clear() error
	// Hook a main menu to the renderer, this is required before calling `rdr.Render`
	HookMainMenu(*MainMenu)
}

type Basic struct {
	mm         *MainMenu
	trm        *term.Terminal
	statusline string
}

// A basic renderer
func NewBasic() *Basic {
	return &Basic{
		mm: nil,
		trm: term.NewTerminal(struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}, ""),
		statusline: "",
	}
}

func (rdr *Basic) Render() error {
	if rdr.mm == nil {
		return Errors.MainMenuNotHooked
	}
	x, y, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	getLine := func(text string, accentText string, valueText string, isSelected bool, isEditing bool, alignment align) []byte {
		tR, atR, vtR := []rune(text), []rune(accentText), []rune(valueText)
		if len(tR)+len(atR)+len(vtR) > x {
			if len(vtR) > 3 {
				diffCap := min(max((len(tR)+len(atR)+len(vtR))-x, 0), 3)
				vtR = append([]rune(strings.Repeat(".", diffCap)), vtR[max(0, min(len(vtR)-3, (len(tR)+len(atR)+len(vtR))-(x-diffCap))):]...)
			}
			if len(tR)+len(atR)+len(vtR) > x && len(tR) > 3 {
				diffCap := min(max((len(tR)+len(atR)+len(vtR))-x, 0), 3)
				tR = append(tR[:max(0, min(len(tR), x-(len(atR)+len(vtR)+diffCap)))], []rune(strings.Repeat(".", diffCap))...)
			}
		}

		line := []byte{}
		if alignment == Aligns.Left {
			line = []byte{}
		} else if alignment == Aligns.Middle {
			i := int((float64(x) / 2) - (float64(len(tR)+len(atR)+len(vtR)) / 2))
			if i > 0 {
				line = []byte("\033[" + strconv.Itoa(i) + "C")
			}
		} else if alignment == Aligns.Right {
			i := x - (len(tR) + len(atR) + len(vtR))
			if i > 0 {
				line = []byte("\033[" + strconv.Itoa(i) + "C")
			}
		}

		if isSelected {
			if isEditing {
				return slices.Concat(line, rdr.mm.cur.Color, []byte(string(tR)), Colors.Reset, rdr.mm.cur.AccentColor, []byte(string(atR)), Colors.Reset, rdr.mm.cur.SelectBGColor, rdr.mm.cur.SelectColor, []byte(string(vtR)), Colors.Reset)
			} else {
				return slices.Concat(line, rdr.mm.cur.SelectBGColor, rdr.mm.cur.SelectColor, []byte(string(tR)), Colors.Reset, rdr.mm.cur.AccentColor, []byte(string(atR)), Colors.Reset, rdr.mm.cur.ValueColor, []byte(string(vtR)), Colors.Reset)
			}
		}
		return slices.Concat(line, rdr.mm.cur.Color, []byte(string(tR)), Colors.Reset, rdr.mm.cur.AccentColor, []byte(string(atR)), Colors.Reset, rdr.mm.cur.ValueColor, []byte(string(vtR)), Colors.Reset)
	}

	lines := append([][]byte{getLine(rdr.mm.cur.Title, "", "", false, false, rdr.mm.cur.Align)}, slices.Concat(rdr.mm.cur.AccentColor, []byte(strings.Repeat("─", x)), Colors.Reset))
	itemLen := 0

	if len(rdr.mm.cur.Menus) > 0 {
		for i, men := range rdr.mm.cur.Menus {
			lines = append(lines, getLine(men.Title, " 🞂", "", itemLen+i == rdr.mm.cur.selected, false, rdr.mm.cur.Align))
		}
		lines = append(lines, []byte{})
		itemLen += len(rdr.mm.cur.Menus)
	}

	if len(rdr.mm.cur.Actions) > 0 {
		for i, act := range rdr.mm.cur.Actions {
			lines = append(lines, getLine(act.Name, "", "", itemLen+i == rdr.mm.cur.selected, false, rdr.mm.cur.Align))
		}
		lines = append(lines, []byte{})
		itemLen += len(rdr.mm.cur.Actions)
	}

	if len(rdr.mm.cur.Lists) > 0 {
		for i, lst := range rdr.mm.cur.Lists {
			lines = append(lines, getLine(lst.Name, " ▷ ", lst.Get(), itemLen+i == rdr.mm.cur.selected, lst.editing, rdr.mm.cur.Align))
		}
		lines = append(lines, []byte{})
		itemLen += len(rdr.mm.cur.Lists)
	}

	if len(rdr.mm.cur.Options) > 0 {
		for i, opt := range rdr.mm.cur.Options {
			lines = append(lines, getLine(opt.Name, " ▷ ", opt.value, itemLen+i == rdr.mm.cur.selected, opt.editing, rdr.mm.cur.Align))
		}
		lines = append(lines, []byte{})
		itemLen += len(rdr.mm.cur.Options)
	}

	backText := "Exit"
	if rdr.mm.cur.back != nil {
		backText = "Back"
	}
	lines = append(lines, getLine("", "◀ ", backText, itemLen == rdr.mm.cur.selected, itemLen == rdr.mm.cur.selected, rdr.mm.cur.Align))

	if _, err := rdr.trm.Write(slices.Concat([]byte("\033[2J\033[0;0H"), bytes.Join(lines[:min(len(lines), y-1)], []byte("\r\n")), []byte("\r\n"+rdr.statusline[:min(len(rdr.statusline), y-1)]))); err != nil {
		return err
	}
	return nil
}

func (rdr *Basic) StatusLine(str string) { rdr.statusline = str }

func (rdr *Basic) Clear() error {
	_, err := rdr.trm.Write([]byte("\033[2J\033[0;0H"))
	return err
}

func (rdr *Basic) HookMainMenu(mm *MainMenu) { rdr.mm = mm }

type Big struct {
	mm               *MainMenu
	trm              *term.Terminal
	statusline       string
	renderBigCharMap map[rune][][2]int
}

// A bulky renderer
func NewBulky() *Big {
	return &Big{
		mm: nil,
		trm: term.NewTerminal(struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}, ""),
		statusline: "",

		//	r := [][2]int{
		//		{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0},
		//		{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1},
		//		{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2},
		//		{0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3},
		//		{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4},
		//	}
		renderBigCharMap: map[rune][][2]int{
			'A':  {{1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {0, 3}, {4, 3}, {0, 4}, {4, 4}},
			'B':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {0, 2}, {1, 2}, {2, 2}, {3, 2}, {0, 3}, {4, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}},
			'C':  {{1, 0}, {2, 0}, {3, 0}, {4, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
			'D':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {0, 2}, {4, 2}, {0, 3}, {4, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}},
			'E':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}, {3, 2}, {0, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
			'F':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}, {3, 2}, {0, 3}, {0, 4}},
			'G':  {{1, 0}, {2, 0}, {3, 0}, {4, 0}, {0, 1}, {0, 2}, {3, 2}, {4, 2}, {0, 3}, {4, 3}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
			'H':  {{0, 0}, {4, 0}, {0, 1}, {4, 1}, {0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {0, 3}, {4, 3}, {0, 4}, {4, 4}},
			'I':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {2, 1}, {2, 2}, {2, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
			'J':  {{1, 0}, {2, 0}, {3, 0}, {4, 0}, {3, 1}, {3, 2}, {0, 3}, {3, 3}, {1, 4}, {2, 4}},
			'K':  {{0, 0}, {4, 0}, {0, 1}, {3, 1}, {0, 2}, {1, 2}, {2, 2}, {0, 3}, {3, 3}, {0, 4}, {4, 4}},
			'L':  {{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
			'M':  {{0, 0}, {1, 0}, {3, 0}, {4, 0}, {0, 1}, {2, 1}, {4, 1}, {0, 2}, {2, 2}, {4, 2}, {0, 3}, {4, 3}, {0, 4}, {4, 4}},
			'N':  {{0, 0}, {1, 0}, {4, 0}, {0, 1}, {2, 1}, {4, 1}, {0, 2}, {2, 2}, {4, 2}, {0, 3}, {2, 3}, {4, 3}, {0, 4}, {3, 4}, {4, 4}},
			'O':  {{1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {0, 2}, {4, 2}, {0, 3}, {4, 3}, {1, 4}, {2, 4}, {3, 4}},
			'P':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {0, 2}, {4, 2}, {0, 3}, {1, 3}, {2, 3}, {3, 3}, {0, 3}, {0, 4}},
			'Q':  {{1, 0}, {2, 0}, {3, 0}, {0, 2}, {4, 2}, {0, 3}, {3, 3}, {0, 1}, {4, 1}, {1, 4}, {2, 4}, {4, 4}},
			'R':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {0, 2}, {4, 2}, {0, 3}, {1, 3}, {2, 3}, {3, 3}, {0, 4}, {4, 4}},
			'S':  {{1, 0}, {2, 0}, {3, 0}, {4, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}},
			'T':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}},
			'U':  {{0, 0}, {4, 0}, {0, 1}, {4, 1}, {0, 2}, {4, 2}, {0, 3}, {4, 3}, {1, 4}, {2, 4}, {3, 4}},
			'V':  {{0, 0}, {4, 0}, {0, 1}, {4, 1}, {0, 2}, {4, 2}, {1, 3}, {3, 3}, {2, 4}},
			'W':  {{0, 0}, {4, 0}, {0, 1}, {4, 1}, {0, 2}, {2, 2}, {4, 2}, {0, 3}, {2, 3}, {4, 3}, {1, 4}, {3, 4}},
			'X':  {{0, 0}, {4, 0}, {1, 1}, {3, 1}, {2, 2}, {1, 3}, {3, 3}, {0, 4}, {4, 4}},
			'Y':  {{0, 0}, {4, 0}, {1, 1}, {3, 1}, {2, 2}, {2, 3}, {2, 4}},
			'Z':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {4, 1}, {1, 2}, {2, 2}, {3, 2}, {0, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
			'0':  {{1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {0, 2}, {2, 2}, {4, 2}, {0, 3}, {4, 3}, {1, 4}, {2, 4}, {3, 4}},
			'1':  {{2, 0}, {1, 1}, {2, 1}, {0, 2}, {2, 2}, {2, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
			'2':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 1}, {1, 2}, {2, 2}, {3, 2}, {0, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
			'3':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 1}, {1, 2}, {2, 2}, {3, 2}, {4, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}},
			'4':  {{0, 0}, {3, 0}, {0, 1}, {3, 1}, {0, 2}, {3, 2}, {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3}, {3, 4}},
			'5':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}},
			'6':  {{1, 0}, {2, 0}, {3, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}, {3, 2}, {0, 3}, {4, 3}, {1, 4}, {2, 4}, {3, 4}},
			'7':  {{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {4, 1}, {3, 2}, {2, 3}, {2, 4}},
			'8':  {{1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {1, 2}, {2, 2}, {3, 2}, {0, 3}, {4, 3}, {1, 4}, {2, 4}, {3, 4}},
			'9':  {{1, 0}, {2, 0}, {3, 0}, {0, 1}, {4, 1}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {4, 3}, {1, 4}, {2, 4}, {3, 4}},
			' ':  {},
			'!':  {{2, 0}, {2, 1}, {2, 2}, {2, 4}},
			'?':  {{1, 0}, {2, 0}, {0, 1}, {3, 1}, {2, 2}, {2, 4}},
			'@':  {{3, 0}, {1, 1}, {4, 1}, {0, 2}, {2, 2}, {4, 2}, {0, 3}, {1, 3}, {2, 3}, {4, 3}, {0, 4}, {2, 4}, {3, 4}},
			'#':  {{1, 0}, {3, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {1, 2}, {3, 2}, {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3}, {1, 4}, {3, 4}},
			'$':  {{1, 0}, {2, 0}, {3, 0}, {4, 0}, {0, 1}, {2, 1}, {0, 2}, {1, 2}, {2, 2}, {3, 2}, {2, 3}, {4, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}},
			'%':  {{0, 0}, {1, 0}, {4, 0}, {0, 1}, {1, 1}, {3, 1}, {2, 2}, {1, 3}, {3, 3}, {4, 3}, {0, 4}, {3, 4}, {4, 4}},
			'^':  {{2, 0}, {1, 1}, {3, 1}, {0, 2}, {4, 2}},
			'&':  {{1, 0}, {0, 1}, {2, 1}, {1, 2}, {2, 2}, {4, 2}, {0, 3}, {3, 3}, {1, 4}, {2, 4}, {4, 4}},
			'=':  {{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3}},
			'*':  {{1, 1}, {3, 1}, {2, 2}, {1, 3}, {3, 3}},
			'+':  {{2, 1}, {1, 2}, {2, 2}, {3, 2}, {2, 3}},
			'-':  {{1, 2}, {2, 2}, {3, 2}},
			'_':  {{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
			'(':  {{2, 0}, {1, 1}, {1, 2}, {1, 3}, {2, 4}},
			')':  {{2, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 4}},
			'{':  {{2, 0}, {1, 1}, {2, 2}, {1, 3}, {2, 4}},
			'}':  {{2, 0}, {3, 1}, {2, 2}, {3, 3}, {2, 4}},
			'[':  {{1, 0}, {2, 0}, {1, 1}, {1, 2}, {1, 3}, {1, 4}, {2, 4}},
			']':  {{2, 0}, {3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 4}, {3, 4}},
			'\'': {{1, 0}, {2, 0}, {1, 1}, {2, 1}, {2, 2}, {1, 3}},
			'"':  {{0, 0}, {1, 0}, {3, 0}, {4, 0}, {0, 1}, {1, 1}, {3, 1}, {4, 1}, {1, 2}, {4, 2}, {0, 3}, {3, 3}},
			':':  {{1, 0}, {2, 0}, {1, 1}, {2, 1}, {1, 3}, {2, 3}, {1, 4}, {2, 4}},
			';':  {{1, 0}, {2, 0}, {1, 1}, {2, 1}, {2, 3}, {1, 4}, {2, 4}},
			',':  {{2, 3}, {1, 4}, {2, 4}},
			'.':  {{1, 3}, {2, 3}, {1, 4}, {2, 4}},
			'`':  {{1, 0}, {2, 0}, {2, 1}},
			'~':  {{1, 1}, {0, 2}, {2, 2}, {4, 2}, {3, 3}},
			'\\': {{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}},
			'/':  {{4, 0}, {3, 1}, {2, 2}, {1, 3}, {0, 4}},
			'|':  {{2, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}},
			'<':  {{3, 0}, {4, 0}, {1, 1}, {2, 1}, {0, 2}, {1, 3}, {2, 3}, {3, 4}, {4, 4}},
			'>':  {{0, 0}, {1, 0}, {2, 1}, {3, 1}, {4, 2}, {2, 3}, {3, 3}, {0, 4}, {1, 4}},
		},
	}
}

func (rdr *Big) Render() error {
	if rdr.mm == nil {
		return Errors.MainMenuNotHooked
	}
	x, y, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	x = int(float64(x) / 12)

	addBlocks := func(lines [][]byte, str string) {
		for end, r := range str {
			cords, ok := rdr.renderBigCharMap[unicode.ToUpper(r)]
			if !ok {
				return
			}
			block := [][][]byte{}
			for i := 0; i < 5; i++ {
				block = append(block, [][]byte{[]byte("  "), []byte("  "), []byte("  "), []byte("  "), []byte("  ")})
			}
			for _, cord := range cords {
				block[cord[1]] = slices.Concat(block[cord[1]][:cord[0]], [][]byte{[]byte("██")}, block[cord[1]][(cord[0])+1:])
			}

			for i, line := range block {
				lines[i] = slices.Concat(lines[i], bytes.Join(line, []byte{}))
				if end != len(str)-1 {
					lines[i] = append(lines[i], []byte("  ")...)
				}
			}
		}
	}

	getLines := func(text string, accentText string, valueText string, isSelected bool, isEditing bool, alignment align) [][]byte {
		tR, atR, vtR := []rune(text), []rune(accentText), []rune(valueText)
		if len(tR)+len(atR)+len(vtR) > x {
			if len(vtR) > 3 {
				diffCap := min(max((len(tR)+len(atR)+len(vtR))-x, 0), 3)
				vtR = append([]rune(strings.Repeat(".", diffCap)), vtR[max(0, min(len(vtR)-3, (len(tR)+len(atR)+len(vtR))-(x-diffCap))):]...)
			}
			if len(tR)+len(atR)+len(vtR) > x && len(tR) > 3 {
				diffCap := min(max((len(tR)+len(atR)+len(vtR))-x, 0), 3)
				tR = append(tR[:max(0, min(len(tR), x-(len(atR)+len(vtR)+diffCap)))], []rune(strings.Repeat(".", diffCap))...)
			}
		}

		lines := [][]byte{{}, {}, {}, {}, {}}

		for i := range lines {
			if isSelected {
				if isEditing {
					lines[i] = slices.Concat(lines[i], rdr.mm.cur.Color)
				} else {
					lines[i] = slices.Concat(lines[i], rdr.mm.cur.SelectBGColor, rdr.mm.cur.SelectColor)
				}
			} else {
				lines[i] = slices.Concat(lines[i], rdr.mm.cur.Color)
			}
		}

		addBlocks(lines, string(tR))

		for i := range lines {
			lines[i] = slices.Concat(lines[i], Colors.Reset, rdr.mm.cur.AccentColor)
		}

		addBlocks(lines, string(atR))

		for i := range lines {
			lines[i] = slices.Concat(lines[i], Colors.Reset, rdr.mm.cur.ValueColor)
			if isSelected {
				if isEditing {
					lines[i] = slices.Concat(lines[i], rdr.mm.cur.SelectBGColor, rdr.mm.cur.SelectColor)
				} else {
					lines[i] = slices.Concat(lines[i], rdr.mm.cur.ValueColor)
				}
			} else {
				lines[i] = slices.Concat(lines[i], rdr.mm.cur.ValueColor)
			}
		}

		addBlocks(lines, string(vtR))

		for i := range lines {
			lines[i] = append(lines[i], Colors.Reset...)
		}

		return lines
	}

	lines := append(getLines(rdr.mm.cur.Title, "", "", false, false, rdr.mm.cur.Align), slices.Concat(rdr.mm.cur.AccentColor, []byte(strings.Repeat("─", x*12)), Colors.Reset))
	itemLen := 0

	if len(rdr.mm.cur.Menus) > 0 {
		for i, men := range rdr.mm.cur.Menus {
			if itemLen+i < rdr.mm.cur.selected {
				continue
			}
			lines = slices.Concat(lines, getLines(men.Title, " >", "", itemLen+i == rdr.mm.cur.selected, false, rdr.mm.cur.Align), [][]byte{{}})
		}
		itemLen += len(rdr.mm.cur.Menus)
		if itemLen > rdr.mm.cur.selected {
			lines = append(lines, [][]byte{{}, {}, {}, {}}...)
		}
	}

	if len(rdr.mm.cur.Actions) > 0 {
		for i, act := range rdr.mm.cur.Actions {
			if itemLen+i < rdr.mm.cur.selected {
				continue
			}
			lines = slices.Concat(lines, getLines(act.Name, "", "", itemLen+i == rdr.mm.cur.selected, false, rdr.mm.cur.Align), [][]byte{{}})
		}
		itemLen += len(rdr.mm.cur.Actions)
		if itemLen > rdr.mm.cur.selected {
			lines = append(lines, [][]byte{{}, {}, {}, {}}...)
		}
	}

	if len(rdr.mm.cur.Lists) > 0 {
		for i, lst := range rdr.mm.cur.Lists {
			if itemLen+i < rdr.mm.cur.selected {
				continue
			}
			lines = slices.Concat(lines, getLines(lst.Name, ": ", lst.Get(), itemLen+i == rdr.mm.cur.selected, lst.editing, rdr.mm.cur.Align), [][]byte{{}})
		}
		itemLen += len(rdr.mm.cur.Lists)
		if itemLen > rdr.mm.cur.selected {
			lines = append(lines, [][]byte{{}, {}, {}, {}}...)
		}
	}

	if len(rdr.mm.cur.Options) > 0 {
		for i, opt := range rdr.mm.cur.Options {
			if itemLen+i < rdr.mm.cur.selected {
				continue
			}
			lines = slices.Concat(lines, getLines(opt.Name, ": ", opt.value, itemLen+i == rdr.mm.cur.selected, opt.editing, rdr.mm.cur.Align), [][]byte{{}})
		}
		itemLen += len(rdr.mm.cur.Options)
		if itemLen > rdr.mm.cur.selected {
			lines = append(lines, [][]byte{{}, {}, {}, {}}...)
		}
	}

	backText := "Exit"
	if rdr.mm.cur.back != nil {
		backText = "Back"
	}
	lines = append(lines, getLines("", "< ", backText, itemLen == rdr.mm.cur.selected, itemLen == rdr.mm.cur.selected, rdr.mm.cur.Align)...)

	if _, err := rdr.trm.Write(slices.Concat([]byte("\033[2J\033[0;0H"), bytes.Join(lines[:min(len(lines), y-1)], []byte("\r\n")), []byte("\r\n"+rdr.statusline[:min(len(rdr.statusline), y-1)]))); err != nil {
		return err
	}
	return nil
}

func (rdr *Big) StatusLine(str string) { rdr.statusline = str }

func (rdr *Big) Clear() error {
	_, err := rdr.trm.Write([]byte("\033[2J\033[0;0H"))
	return err
}

func (rdr *Big) HookMainMenu(mm *MainMenu) { rdr.mm = mm }
