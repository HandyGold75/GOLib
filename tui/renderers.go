package main

import (
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/term"
)

// Basic renderer for redering the current menu.
func RenderBasic(mm *MainMenu) error {
	x, y, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	getLine := func(text string, accentText string, valueText string, isSelected bool, isEditing bool, alignment align) []byte {
		tR, atR, vtR := []rune(text), []rune(accentText), []rune(valueText)
		tRLen, atRLen, vtRLen := len(tR), len(atR), len(vtR)
		if tRLen+atRLen+vtRLen > x {
			diffCap := min(max((tRLen+atRLen+vtRLen)-x, 0), 3)
			if vtRLen > 3 {
				vtR = append([]rune(strings.Repeat(".", diffCap)), vtR[max(0, min(vtRLen-3, (tRLen+atRLen+vtRLen)-(x-diffCap))):]...)
			}
			if tRLen+atRLen+vtRLen > x && tRLen > 3 {
				tR = append(tR[:max(0, min(tRLen, x-(atRLen+vtRLen+diffCap)))], []rune(strings.Repeat(".", diffCap))...)
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
				return slices.Concat(line, mm.cur.Color, []byte(string(tR)), Colors.Reset, mm.cur.AccentColor, []byte(string(atR)), Colors.Reset, mm.cur.SelectBGColor, mm.cur.SelectColor, []byte(string(vtR)), Colors.Reset)
			} else {
				return slices.Concat(line, mm.cur.SelectBGColor, mm.cur.SelectColor, []byte(string(tR)), Colors.Reset, mm.cur.AccentColor, []byte(string(atR)), Colors.Reset, mm.cur.ValueColor, []byte(string(vtR)), Colors.Reset)
			}
		}
		return slices.Concat(line, mm.cur.Color, []byte(string(tR)), Colors.Reset, mm.cur.AccentColor, []byte(string(atR)), Colors.Reset, mm.cur.ValueColor, []byte(string(vtR)), Colors.Reset)
	}

	lines := append([][]byte{getLine(mm.cur.Title, "", "", false, false, mm.cur.Align)}, slices.Concat(mm.cur.AccentColor, []byte(strings.Repeat("─", x)), Colors.Reset))
	itemLen := 0

	if len(mm.cur.Menus) > 0 {
		for i, mn := range mm.cur.Menus {
			lines = append(lines, getLine(mn.Title, " 🞂", "", itemLen+i == mm.cur.selected, false, mm.cur.Align))
		}
		lines = append(lines, []byte{})
		itemLen += len(mm.cur.Menus)
	}

	if len(mm.cur.Actions) > 0 {
		for i, act := range mm.cur.Actions {
			lines = append(lines, getLine(act.Name, "", "", itemLen+i == mm.cur.selected, false, mm.cur.Align))
		}
		lines = append(lines, []byte{})
		itemLen += len(mm.cur.Actions)
	}

	if len(mm.cur.Lists) > 0 {
		for i, lst := range mm.cur.Lists {
			lines = append(lines, getLine(lst.Name, " ▷ ", lst.Get(), itemLen+i == mm.cur.selected, lst.editing, mm.cur.Align))
		}
		lines = append(lines, []byte{})
		itemLen += len(mm.cur.Lists)
	}

	if len(mm.cur.Options) > 0 {
		for i, opt := range mm.cur.Options {
			lines = append(lines, getLine(opt.Name, " ▷ ", opt.value, itemLen+i == mm.cur.selected, opt.editing, mm.cur.Align))
		}
		lines = append(lines, []byte{})
		itemLen += len(mm.cur.Options)
	}

	backText := "Exit"
	if mm.cur.back != nil {
		backText = "Back"
	}
	lines = append(lines, getLine("", "◀ ", backText, itemLen == mm.cur.selected, itemLen == mm.cur.selected, mm.cur.Align))

	if _, err := mm.trm.Write(slices.Concat([]byte("\033[2J\033[0;0H"), bytes.Join(lines[:min(len(lines), y-1)], []byte("\r\n")), []byte("\r\n"))); err != nil {
		return err
	}
	return nil
}

// Format
//
//	// [2]int{x, y}
//
//	r := [][2]int{
//		{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0},
//		{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1},
//		{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2},
//		{0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3},
//		{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4},
//	}
var renderBigCharMap = map[rune][][2]int{
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
}

func RenderBig(mm *MainMenu) error {
	x, y, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	addBlocks := func(lines [][]byte, str string) {
		for end, r := range str {
			cords, ok := renderBigCharMap[unicode.ToUpper(r)]
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
		tRLen, atRLen, vtRLen := len(tR), len(atR), len(vtR)
		temp := float64(12)
		if tRLen+atRLen+vtRLen > int(float64(x)/temp) {
			diffCap := min(max((tRLen+atRLen+vtRLen)-int(float64(x)/temp), 0), 3)
			if vtRLen > 3 {
				vtR = append([]rune(strings.Repeat(".", diffCap)), vtR[max(0, min(vtRLen-3, (tRLen+atRLen+vtRLen)-(int(float64(x)/temp)-diffCap))):]...)
			}
			if tRLen+atRLen+vtRLen > int(float64(x)/temp) && tRLen > 3 {
				tR = append(tR[:max(0, min(tRLen, int(float64(x)/temp)-(atRLen+vtRLen+diffCap)))], []rune(strings.Repeat(".", diffCap))...)
			}
		}

		lines := [][]byte{{}, {}, {}, {}, {}}

		for i := range lines {
			if isSelected {
				if isEditing {
					lines[i] = slices.Concat(lines[i], mm.cur.Color)
				} else {
					lines[i] = slices.Concat(lines[i], mm.cur.SelectBGColor, mm.cur.SelectColor)
				}
			} else {
				lines[i] = slices.Concat(lines[i], mm.cur.Color)
			}
		}

		addBlocks(lines, string(tR))

		for i := range lines {
			lines[i] = slices.Concat(lines[i], Colors.Reset, mm.cur.AccentColor)
		}

		addBlocks(lines, string(atR))

		for i := range lines {
			lines[i] = slices.Concat(lines[i], Colors.Reset, mm.cur.ValueColor)
			if isSelected {
				if isEditing {
					lines[i] = slices.Concat(lines[i], mm.cur.SelectBGColor, mm.cur.SelectColor)
				} else {
					lines[i] = slices.Concat(lines[i], mm.cur.ValueColor)
				}
			} else {
				lines[i] = slices.Concat(lines[i], mm.cur.ValueColor)
			}
		}

		addBlocks(lines, string(vtR))

		for i := range lines {
			lines[i] = append(lines[i], Colors.Reset...)
		}

		return lines
	}

	lines := append(getLines(mm.cur.Title, "", "", false, false, mm.cur.Align), slices.Concat(mm.cur.AccentColor, []byte(strings.Repeat("─", x)), Colors.Reset))
	itemLen := 0

	if len(mm.cur.Menus) > 0 {
		for i, mn := range mm.cur.Menus {
			if itemLen+i < mm.cur.selected {
				continue
			}
			lines = slices.Concat(lines, getLines(mn.Title, " >", "", itemLen+i == mm.cur.selected, false, mm.cur.Align), [][]byte{{}})
		}
		itemLen += len(mm.cur.Menus)
		if itemLen > mm.cur.selected {
			lines = append(lines, [][]byte{{}, {}, {}, {}}...)
		}
	}

	if len(mm.cur.Actions) > 0 {
		for i, act := range mm.cur.Actions {
			if itemLen+i < mm.cur.selected {
				continue
			}
			lines = slices.Concat(lines, getLines(act.Name, "", "", itemLen+i == mm.cur.selected, false, mm.cur.Align), [][]byte{{}})
		}
		itemLen += len(mm.cur.Actions)
		if itemLen > mm.cur.selected {
			lines = append(lines, [][]byte{{}, {}, {}, {}}...)
		}
	}

	if len(mm.cur.Lists) > 0 {
		for i, lst := range mm.cur.Lists {
			if itemLen+i < mm.cur.selected {
				continue
			}
			lines = slices.Concat(lines, getLines(lst.Name, ": ", lst.Get(), itemLen+i == mm.cur.selected, lst.editing, mm.cur.Align), [][]byte{{}})
		}
		itemLen += len(mm.cur.Lists)
		if itemLen > mm.cur.selected {
			lines = append(lines, [][]byte{{}, {}, {}, {}}...)
		}
	}

	if len(mm.cur.Options) > 0 {
		for i, opt := range mm.cur.Options {
			if itemLen+i < mm.cur.selected {
				continue
			}
			lines = slices.Concat(lines, getLines(opt.Name, ": ", opt.value, itemLen+i == mm.cur.selected, opt.editing, mm.cur.Align), [][]byte{{}})
		}
		itemLen += len(mm.cur.Options)
		if itemLen > mm.cur.selected {
			lines = append(lines, [][]byte{{}, {}, {}, {}}...)
		}
	}

	backText := "Exit"
	if mm.cur.back != nil {
		backText = "Back"
	}
	lines = append(lines, getLines("", "< ", backText, itemLen == mm.cur.selected, itemLen == mm.cur.selected, mm.cur.Align)...)

	if _, err := mm.trm.Write(slices.Concat([]byte("\033[2J\033[0;0H"), bytes.Join(lines[:min(len(lines), y-1)], []byte("\r\n")), []byte("\r\n"))); err != nil {
		return err
	}
	return nil
}
