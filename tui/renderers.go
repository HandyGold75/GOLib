package main

import (
	"bytes"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/term"
)

// Basic renderer for redering the current menu.
func RenderBasic(mm *MainMenu) error {
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
	lines := append([][]byte{}, slices.Concat(getCursorPos(len(mm.cur.Title), Aligns.Middle), mm.cur.Color, []byte(mm.cur.Title), Colors.Reset))
	lines = append(lines, slices.Concat(mm.cur.AccentColor, []byte(strings.Repeat("─", x)), Colors.Reset))

	if len(mm.cur.Menus) > 0 {
		for _, mn := range mm.cur.Menus {
			itemLen += 1
			if itemLen == mm.cur.selected {
				lines = append(lines, slices.Concat(getCursorPos(len(mn.Title)+2, Aligns.Middle), mn.SelectBGColor, mn.SelectColor, []byte(mn.Title), Colors.Reset, mn.AccentColor, []byte(" 🞂"), Colors.Reset))
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(mn.Title)+2, Aligns.Middle), mn.Color, []byte(mn.Title), mn.AccentColor, []byte(" 🞂"), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	if len(mm.cur.Actions) > 0 {
		for _, act := range mm.cur.Actions {
			itemLen += 1
			if itemLen == mm.cur.selected {
				lines = append(lines, slices.Concat(getCursorPos(len(act.Name), Aligns.Middle), act.SelectBGColor, act.SelectColor, []byte(act.Name), Colors.Reset))
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(act.Name), Aligns.Middle), act.Color, []byte(act.Name), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	if len(mm.cur.Lists) > 0 {
		for _, lst := range mm.cur.Lists {
			itemLen += 1
			v := lst.Get()
			if itemLen == mm.cur.selected {
				if lst.editing {
					lines = append(lines, slices.Concat(getCursorPos(len(lst.Name)+3+len(v), Aligns.Middle), lst.Color, []byte(lst.Name), lst.AccentColor, []byte(" ▷ "), lst.SelectBGColor, lst.SelectColor, []byte(v), Colors.Reset))
				} else {
					lines = append(lines, slices.Concat(getCursorPos(len(lst.Name)+3+len(v), Aligns.Middle), lst.SelectBGColor, lst.SelectColor, []byte(lst.Name), Colors.Reset, lst.AccentColor, []byte(" ▷ "), lst.ValueColor, []byte(v), Colors.Reset))
				}
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(lst.Name)+3+len(v), Aligns.Middle), lst.Color, []byte(lst.Name), lst.AccentColor, []byte(" ▷ "), lst.ValueColor, []byte(v), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	if len(mm.cur.Options) > 0 {
		for _, opt := range mm.cur.Options {
			itemLen += 1
			if itemLen == mm.cur.selected {
				if opt.editing {
					lines = append(lines, slices.Concat(getCursorPos(len(opt.Name)+3+len(opt.value), Aligns.Middle), opt.Color, []byte(opt.Name), opt.AccentColor, []byte(" ▷ "), opt.SelectBGColor, opt.SelectColor, []byte(opt.value), Colors.Reset))
				} else {
					lines = append(lines, slices.Concat(getCursorPos(len(opt.Name)+3+len(opt.value), Aligns.Middle), opt.SelectBGColor, opt.SelectColor, []byte(opt.Name), Colors.Reset, opt.AccentColor, []byte(" ▷ "), opt.ValueColor, []byte(opt.value), Colors.Reset))
				}
			} else {
				lines = append(lines, slices.Concat(getCursorPos(len(opt.Name)+3+len(opt.value), Aligns.Middle), opt.Color, []byte(opt.Name), opt.AccentColor, []byte(" ▷ "), opt.ValueColor, []byte(opt.value), Colors.Reset))
			}
		}
		lines = append(lines, []byte{})
	}

	backText := "Exit"
	if mm.cur.back != nil {
		backText = "Back"
	}
	itemLen += 1
	if itemLen == mm.cur.selected {
		lines = append(lines, slices.Concat(getCursorPos(len(backText)+2, Aligns.Middle), mm.cur.AccentColor, []byte("◀ "), mm.cur.SelectBGColor, mm.cur.SelectColor, []byte(backText), Colors.Reset))
	} else {
		lines = append(lines, slices.Concat(getCursorPos(len(backText)+2, Aligns.Middle), mm.cur.AccentColor, []byte("◀ "), mm.cur.Color, []byte(backText), Colors.Reset))
	}

	if _, err := mm.trm.Write(slices.Concat([]byte("\033[2J\033[0;0H"), bytes.Join(lines, []byte("\r\n")), []byte("\r\n"))); err != nil {
		return err
	}
	return nil
}
