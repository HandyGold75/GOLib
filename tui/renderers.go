package tui

import (
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/term"
)

type (
	// The Renderer interface may be parsed to a MainMenu to be used as Renderer.
	//
	// The MainMenu is reponsible for all logic while the Renderer is responseble for visializing the MainMenu.
	//
	// The MainMenu should always call HookMainMenu first to allow the Renderer to hook into the MainMenu.
	Renderer interface {
		// Gets called when a rerender is requested.
		Render() error
		// Gets called when screen clear is requested.
		Clear() error
		// Gets called before any calls to `rdr.Render` to hook into the current `MainMenu`.
		HookMainMenu(*MainMenu)
	}

	rendererBasic struct {
		mm  *MainMenu
		trm *term.Terminal
	}
	rendererBulky struct {
		mm               *MainMenu
		trm              *term.Terminal
		renderBigCharMap map[rune][][2]int
	}
)

func newRendererBasic() *rendererBasic {
	return &rendererBasic{
		mm: nil,
		trm: term.NewTerminal(struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}, ""),
	}
}

func (rdr *rendererBasic) Render() error {
	if rdr.mm == nil {
		return Errors.MainMenuNotHooked
	}
	x, y, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	getLine := func(text string, accentText string, valueText string, isSelected bool, isEditing bool) string {
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

		line := ""
		switch rdr.mm.cur.Align {
		case AlignLeft:
			line = ""
		case AlignMiddle:
			i := int((float64(x) / 2) - (float64(len(tR)+len(atR)+len(vtR)) / 2))
			if i > 0 {
				line = "\033[" + strconv.Itoa(i) + "C"
			}
		case AlignRight:
			i := x - (len(tR) + len(atR) + len(vtR))
			if i > 0 {
				line = "\033[" + strconv.Itoa(i) + "C"
			}
		}

		if isSelected {
			if isEditing {
				return line + string(rdr.mm.cur.Color) + string(tR) + string(Reset+rdr.mm.cur.AccentColor) + string(atR) + string(Reset+rdr.mm.cur.SelectColor) + string(vtR) + string(Reset)
			} else {
				return line + string(rdr.mm.cur.SelectColor) + string(tR) + string(Reset+rdr.mm.cur.AccentColor) + string(atR) + string(Reset+rdr.mm.cur.ValueColor) + string(vtR) + string(Reset)
			}
		}
		return line + string(rdr.mm.cur.Color) + string(tR) + string(Reset+rdr.mm.cur.AccentColor) + string(atR) + string(Reset+rdr.mm.cur.ValueColor) + string(vtR) + string(Reset)
	}

	lines := []string{getLine(rdr.mm.cur.String(), "", "", false, false), string(rdr.mm.cur.AccentColor) + strings.Repeat("─", x) + string(Reset)}

	for i, itm := range rdr.mm.cur.Items {
		switch itm.Type() {
		case "menu":
			lines = append(lines, getLine(itm.String(), " ▶", "", rdr.mm.cur.selected == i, false))
		case "text", "list", "digit", "ipv4", "ipv6":
			lines = append(lines, getLine(itm.String(), " ▷ ", itm.Value(), rdr.mm.cur.selected == i, itm.Editing()))
		default:
			lines = append(lines, getLine(itm.String(), "", "", rdr.mm.cur.selected == i, itm.Editing()))
		}
	}

	lines = append(lines, getLine("", "◀ ", rdr.mm.cur.BackText, rdr.mm.cur.selected >= len(rdr.mm.cur.Items), rdr.mm.cur.selected >= len(rdr.mm.cur.Items)))

	if _, err := rdr.trm.Write([]byte("\033[2J\033[0;0H" + strings.Join(lines[:min(len(lines), y-1)], "\r\n") + "\033[" + strconv.Itoa(y) + ";0H" + rdr.mm.statusline[:min(len(rdr.mm.statusline), x-1)])); err != nil {
		return err
	}
	return nil
}

func (rdr *rendererBasic) Clear() error {
	_, err := rdr.trm.Write([]byte("\033[2J\033[0;0H"))
	return err
}

func (rdr *rendererBasic) HookMainMenu(mm *MainMenu) { rdr.mm = mm }

func newRendererBulky() *rendererBulky {
	return &rendererBulky{
		mm: nil,
		trm: term.NewTerminal(struct {
			io.Reader
			io.Writer
		}{os.Stdin, os.Stdout}, ""),

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
			'"':  {{0, 0}, {1, 0}, {3, 0}, {4, 0}, {0, 1}, {1, 1}, {3, 1}, {4, 1}, {1, 2}, {4, 2}, {0, 3}, {3, 3}},
			'#':  {{1, 0}, {3, 0}, {0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {1, 2}, {3, 2}, {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3}, {1, 4}, {3, 4}},
			'$':  {{1, 0}, {2, 0}, {3, 0}, {4, 0}, {0, 1}, {2, 1}, {0, 2}, {1, 2}, {2, 2}, {3, 2}, {2, 3}, {4, 3}, {0, 4}, {1, 4}, {2, 4}, {3, 4}},
			'%':  {{0, 0}, {1, 0}, {4, 0}, {0, 1}, {1, 1}, {3, 1}, {2, 2}, {1, 3}, {3, 3}, {4, 3}, {0, 4}, {3, 4}, {4, 4}},
			'&':  {{1, 0}, {0, 1}, {2, 1}, {1, 2}, {2, 2}, {4, 2}, {0, 3}, {3, 3}, {1, 4}, {2, 4}, {4, 4}},
			'\'': {{1, 0}, {2, 0}, {1, 1}, {2, 1}, {2, 2}, {1, 3}},
			'(':  {{2, 0}, {1, 1}, {1, 2}, {1, 3}, {2, 4}},
			')':  {{2, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 4}},
			'*':  {{1, 1}, {3, 1}, {2, 2}, {1, 3}, {3, 3}},
			'+':  {{2, 1}, {1, 2}, {2, 2}, {3, 2}, {2, 3}},
			',':  {{2, 3}, {1, 4}, {2, 4}},
			'-':  {{1, 2}, {2, 2}, {3, 2}},
			'.':  {{1, 3}, {2, 3}, {1, 4}, {2, 4}},
			'/':  {{4, 0}, {3, 1}, {2, 2}, {1, 3}, {0, 4}},
			':':  {{1, 0}, {2, 0}, {1, 1}, {2, 1}, {1, 3}, {2, 3}, {1, 4}, {2, 4}},
			';':  {{1, 0}, {2, 0}, {1, 1}, {2, 1}, {2, 3}, {1, 4}, {2, 4}},
			'<':  {{3, 0}, {4, 0}, {1, 1}, {2, 1}, {0, 2}, {1, 3}, {2, 3}, {3, 4}, {4, 4}},
			'=':  {{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3}},
			'>':  {{0, 0}, {1, 0}, {2, 1}, {3, 1}, {4, 2}, {2, 3}, {3, 3}, {0, 4}, {1, 4}},
			'?':  {{1, 0}, {2, 0}, {0, 1}, {3, 1}, {2, 2}, {2, 4}},
			'@':  {{3, 0}, {1, 1}, {4, 1}, {0, 2}, {2, 2}, {4, 2}, {0, 3}, {1, 3}, {2, 3}, {4, 3}, {0, 4}, {2, 4}, {3, 4}},
			'[':  {{1, 0}, {2, 0}, {1, 1}, {1, 2}, {1, 3}, {1, 4}, {2, 4}},
			'\\': {{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4}},
			']':  {{2, 0}, {3, 0}, {3, 1}, {3, 2}, {3, 3}, {2, 4}, {3, 4}},
			'^':  {{2, 0}, {1, 1}, {3, 1}, {0, 2}, {4, 2}},
			'_':  {{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}},
			'`':  {{1, 0}, {2, 0}, {2, 1}},
			'{':  {{2, 0}, {1, 1}, {2, 2}, {1, 3}, {2, 4}},
			'|':  {{2, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}},
			'}':  {{2, 0}, {3, 1}, {2, 2}, {3, 3}, {2, 4}},
			'~':  {{1, 1}, {0, 2}, {2, 2}, {4, 2}, {3, 3}},
		},
	}
}

func (rdr *rendererBulky) Render() error {
	if rdr.mm == nil {
		return Errors.MainMenuNotHooked
	}
	x, y, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	x = int(float64(x) / 12)

	addBlocks := func(lines []string, str string) {
		for end, r := range str {
			cords, ok := rdr.renderBigCharMap[unicode.ToUpper(r)]

			if !ok {
				return
			}
			block := [][]string{}
			for range 5 {
				block = append(block, []string{"  ", "  ", "  ", "  ", "  "})
			}
			for _, cord := range cords {
				block[cord[1]] = slices.Concat(block[cord[1]][:cord[0]], []string{"██"}, block[cord[1]][(cord[0])+1:])
			}

			for i, line := range block {
				lines[i] += strings.Join(line, "")
				if end != len(str)-1 {
					lines[i] += "  "
				}
			}
		}
	}

	getLines := func(text string, accentText string, valueText string, isSelected bool, isEditing bool) []string {
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

		lines := []string{"", "", "", "", ""}
		switch rdr.mm.cur.Align {
		case AlignLeft:
			lines = []string{"", "", "", "", ""}
		case AlignMiddle:
			i := int((float64(x)/2)-(float64(len(tR)+len(atR)+len(vtR))/2)) * 12
			if i > 0 {
				cur := "\033[" + strconv.Itoa(i) + "C"
				lines = []string{cur, cur, cur, cur, cur}
			}
		case AlignRight:
			i := (x - (len(tR) + len(atR) + len(vtR))) * 12
			if i > 0 {
				cur := "\033[" + strconv.Itoa(i) + "C"
				lines = []string{cur, cur, cur, cur, cur}
			}
		}

		for i := range lines {
			if isSelected {
				if isEditing {
					lines[i] += string(rdr.mm.cur.Color)
				} else {
					lines[i] += string(rdr.mm.cur.SelectColor)
				}
			} else {
				lines[i] += string(rdr.mm.cur.Color)
			}
		}

		addBlocks(lines, string(tR))

		for i := range lines {
			lines[i] += string(Reset + rdr.mm.cur.AccentColor)
		}

		addBlocks(lines, string(atR))

		for i := range lines {
			lines[i] += string(Reset + rdr.mm.cur.ValueColor)
			if isSelected {
				if isEditing {
					lines[i] += string(rdr.mm.cur.SelectColor)
				} else {
					lines[i] += string(rdr.mm.cur.ValueColor)
				}
			} else {
				lines[i] += string(rdr.mm.cur.ValueColor)
			}
		}

		addBlocks(lines, string(vtR))

		for i := range lines {
			lines[i] += string(Reset)
		}

		lines = append(lines, "")
		return lines
	}

	lines := append(getLines(rdr.mm.cur.String(), "", "", false, false), string(rdr.mm.cur.AccentColor)+strings.Repeat("─", x*11)+string(Reset), "")

	for i, itm := range rdr.mm.cur.Items[max(0, rdr.mm.cur.selected):] {
		switch itm.Type() {
		case "menu":
			lines = append(lines, getLines(itm.String(), " >", "", i == 0, false)...)
		case "text", "list", "digit", "ipv4", "ipv6":
			lines = append(lines, getLines(itm.String(), ": ", itm.Value(), i == 0, itm.Editing())...)
		default:
			lines = append(lines, getLines(itm.String(), "", "", i == 0, false)...)
		}
	}

	lines = append(lines, getLines("", "< ", rdr.mm.cur.BackText, rdr.mm.cur.selected >= len(rdr.mm.cur.Items), rdr.mm.cur.selected >= len(rdr.mm.cur.Items))...)

	if _, err := rdr.trm.Write([]byte("\033[2J\033[0;0H" + strings.Join(lines[:min(len(lines), y-1)], "\r\n") + "\033[" + strconv.Itoa(y) + ";0H" + rdr.mm.statusline[:min(len(rdr.mm.statusline), (x*12)-1)])); err != nil {
		return err
	}
	return nil
}

func (rdr *rendererBulky) Clear() error {
	_, err := rdr.trm.Write([]byte("\033[2J\033[0;0H"))
	return err
}

// Hook a main menu to the renderer, this is required before calling `rdr.Render`
func (rdr *rendererBulky) HookMainMenu(mm *MainMenu) { rdr.mm = mm }
