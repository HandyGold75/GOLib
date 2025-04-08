package tui

import (
	"bytes"
	"encoding/hex"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
)

type (
	item interface {
		// Get the name of the item.
		String() string

		// Get the value of the item.
		//
		// If the value is not implemented return a empty string.
		Value() string

		// Get the type of the item.
		//
		// This can be one of: `menu`, `text`, `action`, `list`, `digit`, `ipv4`, `ipv6`, `ipadd6`
		Type() string

		// Get the editing state of the item.
		Editing() bool

		enter() error
	}

	menu struct {
		mm            *MainMenu
		name          string
		BackText      string
		Color         color
		AccentColor   color
		SelectColor   color
		SelectBGColor color
		ValueColor    color
		Align         align
		Items         []item
		selected      int
		back          *menu
	}
	text struct {
		mm      *MainMenu
		name    string
		editing bool
		chars   charSet
		value   string
	}
	action struct {
		mm       *MainMenu
		name     string
		callback func()
	}
	list struct {
		mm       *MainMenu
		name     string
		editing  bool
		values   []string
		selected int
	}
	digit struct {
		mm      *MainMenu
		name    string
		editing bool
		value   int
		minimal int
		maximal int
	}
	ipv4 struct {
		mm       *MainMenu
		name     string
		editing  bool
		value    net.IP
		selected int
	}
	ipv6 struct {
		mm       *MainMenu
		name     string
		editing  bool
		value    net.IP
		selected int
	}
)

// Add a new menu to `m.Items`.
//
// Returns a pointer to the new menu.
//
// To set default colors set `tui.Defaults.Color`, `tui.Defaults.AccentColor`, `tui.Defaults.SelectColor`, `tui.Defaults.SelectBGColor`, `tui.Defaults.ValueColor` before creating menus.
//
// To set default alignment set `tui.Defaults.Align` before creating menus.
func (m *menu) NewMenu(name string) *menu {
	mn := &menu{
		mm:            m.mm,
		name:          name,
		BackText:      "Back",
		Color:         Defaults.Color,
		AccentColor:   Defaults.AccentColor,
		SelectColor:   Defaults.SelectColor,
		SelectBGColor: Defaults.SelectBGColor,
		ValueColor:    Defaults.ValueColor,
		Align:         Defaults.Align,
		Items:         []item{},
		selected:      0,
		back:          m,
	}
	m.Items = append(m.Items, mn)
	return mn
}

func (m *menu) String() string { return m.name }
func (m *menu) Value() string  { return "" }
func (m *menu) Type() string   { return "menu" }
func (m *menu) Editing() bool  { return false }

func (m *menu) enter() error {
	_ = m.mm.rdr.Render()
	for {
		in := make([]byte, 3)
		if _, err := os.Stdin.Read(in); err != nil {
			return err
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v keybind) bool { return slices.Equal(v, in) }) {
			return Errors.exit
		} else if slices.ContainsFunc(KeyBinds.Up, func(v keybind) bool { return slices.Equal(v, in) }) {
			m.selected = max(m.selected-1, 0)
			_ = m.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Down, func(v keybind) bool { return slices.Equal(v, in) }) {
			m.selected = min(m.selected+1, len(m.Items))
			_ = m.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Right, func(v keybind) bool { return slices.Equal(v, in) }) {
			if m.selected < len(m.Items) && m.selected >= 0 {
				if m.Items[m.selected].Type() == "menu" {
					m.mm.cur = m.Items[m.selected].(*menu)
					return nil
				}
				err := m.Items[m.selected].enter()
				return err
			}
			if m.back == nil {
				return Errors.exit
			}
			m.mm.cur = m.back
			return nil
		} else if slices.ContainsFunc(KeyBinds.Left, func(v keybind) bool { return slices.Equal(v, in) }) {
			if m.back == nil {
				break
			}
			m.mm.cur = m.back
			return nil
		} else if i := slices.IndexFunc(KeyBinds.Numbers, func(v keybind) bool { return slices.Equal(v, in) }); i != -1 {
			if i > len(m.Items) {
				continue
			}
			m.selected = i - 1

			if m.selected < len(m.Items) && m.selected >= 0 {
				if m.Items[m.selected].Type() == "menu" {
					m.mm.cur = m.Items[m.selected].(*menu)
					return nil
				}
				err := m.Items[m.selected].enter()
				return err
			}
			if m.back == nil {
				return Errors.exit
			}
			m.mm.cur = m.back
			return nil
		}
	}
	return nil
}

// Add a new text to `m.Items`.
//
// Only characters in `chars` can be present in `value`.
//
// Returns a pointer to the new value.
func (m *menu) NewText(name string, chars charSet, value string) *text {
	opt := &text{
		mm:      m.mm,
		name:    name,
		editing: false,
		chars:   chars,
		value:   "",
	}
	for _, char := range value {
		if strings.ContainsAny(string(opt.chars), string(char)) {
			opt.value += string(char)
		}
	}
	m.Items = append(m.Items, opt)
	return opt
}

func (v *text) String() string { return v.name }
func (v *text) Value() string  { return v.value }
func (v *text) Type() string   { return "text" }
func (v *text) Editing() bool  { return v.editing }

func (v *text) enter() error {
	v.editing = true
	_ = v.mm.rdr.Render()
	for {
		in := make([]byte, 3)
		if _, err := os.Stdin.Read(in); err != nil {
			v.editing = false
			_ = v.mm.rdr.Render()
			return err
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v keybind) bool { return slices.Equal(v, in) }) || slices.ContainsFunc(KeyBinds.Confirm, func(v keybind) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Delete, func(v keybind) bool { return slices.Equal(v, in) }) {
			if len(v.value) > 0 {
				v.value = v.value[:len(v.value)-1]
				_ = v.mm.rdr.Render()
				continue
			}
		}

		if strings.ContainsAny(string(v.chars), string(in[:])) {
			v.value += string(bytes.Trim(in, "\x00")[:])
			_ = v.mm.rdr.Render()
		}
	}

	v.editing = false
	_ = v.mm.rdr.Render()
	return nil
}

// Add a new action to `m.Items`.
//
// `callback` is called when this actions is selected.
//
// Returns a pointer to the new action.
func (m *menu) NewAction(name string, callback func()) *action {
	act := &action{
		mm:       m.mm,
		name:     name,
		callback: callback,
	}
	m.Items = append(m.Items, act)
	return act
}

func (a *action) String() string { return a.name }
func (a *action) Value() string  { return "" }
func (a *action) Type() string   { return "action" }
func (a *action) Editing() bool  { return false }

func (a *action) enter() error {
	a.callback()
	return Errors.exit
}

// Add a new list to `m.Items`.
//
// Only options in `values` can be selected.
//
// Returns a pointer to the new list.
func (m *menu) NewList(name string, values []string) *list {
	lst := &list{
		mm:       m.mm,
		name:     name,
		editing:  false,
		values:   values,
		selected: 0,
	}
	m.Items = append(m.Items, lst)
	return lst
}

func (l *list) String() string { return l.name }
func (l *list) Value() string {
	if l.selected < 0 || l.selected > len(l.values)-1 {
		return ""
	}
	return l.values[l.selected]
}
func (l *list) Type() string  { return "list" }
func (l *list) Editing() bool { return l.editing }

func (l *list) enter() error {
	l.editing = true
	_ = l.mm.rdr.Render()
	for {
		in := make([]byte, 3)
		if _, err := os.Stdin.Read(in); err != nil {
			l.editing = false
			_ = l.mm.rdr.Render()
			return err
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v keybind) bool { return slices.Equal(v, in) }) || slices.ContainsFunc(KeyBinds.Confirm, func(v keybind) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Up, func(v keybind) bool { return slices.Equal(v, in) }) {
			l.selected = max(l.selected-1, 0)
			_ = l.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Down, func(v keybind) bool { return slices.Equal(v, in) }) {
			l.selected = min(l.selected+1, len(l.values)-1)
			_ = l.mm.rdr.Render()
			continue

		} else if slices.ContainsFunc(KeyBinds.Right, func(v keybind) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Left, func(v keybind) bool { return slices.Equal(v, in) }) {
			break
		} else if i := slices.IndexFunc(KeyBinds.Numbers, func(v keybind) bool { return slices.Equal(v, in) }); i != -1 {
			if i > len(l.values)-1 {
				continue
			}
			l.selected = i
			_ = l.mm.rdr.Render()
			continue
		}

		for i, str := range l.values {
			if !strings.HasPrefix(strings.ToLower(str), string(bytes.Trim(in, "\x00")[:])) {
				continue
			}
			l.selected = i
			_ = l.mm.rdr.Render()
			continue
		}

		for i, str := range l.values {
			if !strings.HasPrefix(strings.ToLower(str), strings.ToLower(string(bytes.Trim(in, "\x00")[:]))) {
				continue
			}
			l.selected = i
			_ = l.mm.rdr.Render()
			continue
		}
	}
	l.editing = false
	_ = l.mm.rdr.Render()
	return nil
}

// Add a new digit to `m.Items`.
//
// `value` can only be between or equal to `d.minimal` and `d.maximal`.
//
// Returns a pointer to the new list.
func (m *menu) NewDigit(name string, value int, minimal int, maximal int) *digit {
	dgt := &digit{
		mm:      m.mm,
		name:    name,
		editing: false,
		value:   min(max(value, minimal), maximal),
		minimal: minimal,
		maximal: maximal,
	}
	m.Items = append(m.Items, dgt)
	return dgt
}

func (d *digit) String() string { return d.name }
func (d *digit) Value() string  { return strconv.Itoa(d.value) }
func (d *digit) Type() string   { return "digit" }
func (d *digit) Editing() bool  { return d.editing }

func (d *digit) enter() error {
	d.editing = true
	_ = d.mm.rdr.Render()
	for {
		in := make([]byte, 3)
		if _, err := os.Stdin.Read(in); err != nil {
			d.editing = false
			_ = d.mm.rdr.Render()
			return err
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v keybind) bool { return slices.Equal(v, in) }) || slices.ContainsFunc(KeyBinds.Confirm, func(v keybind) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Delete, func(v keybind) bool { return slices.Equal(v, in) }) {
			vStr := strconv.Itoa(d.value)
			if len(vStr) <= 0 {
				continue
			}
			if len(strings.Replace(vStr, "-", "", 1)) == 1 {
				d.value = min(max(0, d.minimal), d.maximal)
				_ = d.mm.rdr.Render()
				continue
			}
			v, err := strconv.Atoi(vStr[:len(vStr)-1])
			if err != nil {
				d.editing = false
				_ = d.mm.rdr.Render()
				return err
			}
			d.value = min(max(v, d.minimal), d.maximal)
			_ = d.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Up, func(v keybind) bool { return slices.Equal(v, in) }) {
			d.value = min(d.value+1, d.maximal)
			_ = d.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Down, func(v keybind) bool { return slices.Equal(v, in) }) {
			d.value = max(d.value-1, d.minimal)
			_ = d.mm.rdr.Render()
			continue
		}

		if strings.ContainsAny(string(CharSets.Digits), string(in[:])) {
			v, err := strconv.Atoi(strconv.Itoa(d.value) + string(bytes.Trim(in, "\x00")[:]))
			if err != nil {
				d.editing = false
				_ = d.mm.rdr.Render()
				return err
			}
			d.value = min(max(v, d.minimal), d.maximal)
			_ = d.mm.rdr.Render()
		}
	}
	d.editing = false
	_ = d.mm.rdr.Render()
	return nil
}

// Add a new ipv4 to `m.Items`.
//
// `value` must be a valid IPv4 address.
//
// Returns a pointer to the new list.
func (m *menu) NewIPv4(name string, value string) *ipv4 {
	ip4 := &ipv4{
		mm:       m.mm,
		name:     name,
		value:    net.ParseIP(value).To4(),
		selected: 0,
		editing:  false,
	}
	if ip4.value == nil {
		ip4.value = net.ParseIP("0.0.0.0").To4()
	}
	m.Items = append(m.Items, ip4)
	return ip4
}

func (p *ipv4) String() string { return p.name }
func (p *ipv4) Value() string  { return p.value.String() }
func (p *ipv4) Type() string   { return "ipv4" }
func (p *ipv4) Editing() bool  { return p.editing }

func (p *ipv4) enter() error {
	p.editing = true
	carry := []byte{}
	_ = p.mm.rdr.Render()
	for {
		in := make([]byte, 3)
		if len(carry) > 0 {
			in = carry
			carry = []byte{}
		} else if _, err := os.Stdin.Read(in); err != nil {
			p.editing = false
			_ = p.mm.rdr.Render()
			return err
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v keybind) bool { return slices.Equal(v, in) }) || slices.ContainsFunc(KeyBinds.Confirm, func(v keybind) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Delete, func(v keybind) bool { return slices.Equal(v, in) }) {
			vStr := strconv.Itoa(int(p.value[p.selected]))
			if len(vStr) <= 0 {
				continue
			}
			if len(vStr) == 1 {
				p.value[p.selected] = byte(0)
				_ = p.mm.rdr.Render()
				continue
			}
			v, err := strconv.Atoi(vStr[:len(vStr)-1])
			if err != nil {
				p.editing = false
				_ = p.mm.rdr.Render()
				return err
			}
			p.value[p.selected] = byte(v)
			_ = p.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Right, func(v keybind) bool { return slices.Equal(v, in) }) {
			p.selected = min(p.selected+1, 3)
			_ = p.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Left, func(v keybind) bool { return slices.Equal(v, in) }) {
			p.selected = max(p.selected-1, 0)
			_ = p.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Up, func(v keybind) bool { return slices.Equal(v, in) }) {
			p.value[p.selected] = p.value[p.selected] + 1
			_ = p.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Down, func(v keybind) bool { return slices.Equal(v, in) }) {
			p.value[p.selected] = p.value[p.selected] - 1
			_ = p.mm.rdr.Render()
			continue
		}

		if strings.ContainsAny(string(CharSets.Digits), string(in[:])) {
			v, err := strconv.Atoi(strconv.Itoa(int(p.value[p.selected])) + string(bytes.Trim(in, "\x00")[:]))
			if err != nil {
				p.editing = false
				_ = p.mm.rdr.Render()
				return err
			}
			if v > 255 {
				if p.selected >= 3 {
					continue
				}
				p.selected += 1
				carry = in
				continue
			}
			p.value[p.selected] = byte(v)
			_ = p.mm.rdr.Render()
		}
	}
	p.editing = false
	_ = p.mm.rdr.Render()
	return nil
}

// Add a new ipv6 to `m.Items`.
//
// `value` must be a valid IPv4 address.
//
// Returns a pointer to the new list.
func (m *menu) NewIPv6(name string, value string) *ipv6 {
	ip4 := &ipv6{
		mm:       m.mm,
		name:     name,
		value:    net.ParseIP(value).To16(),
		selected: 0,
		editing:  false,
	}
	if ip4.value == nil {
		ip4.value = net.ParseIP("::").To16()
	}
	m.Items = append(m.Items, ip4)
	return ip4
}

func (p *ipv6) String() string { return p.name }
func (p *ipv6) Value() string  { return p.value.String() }
func (p *ipv6) Type() string   { return "ipv6" }
func (p *ipv6) Editing() bool  { return p.editing }

func (p *ipv6) enter() error {
	p.editing = true
	carry := []byte{}
	_ = p.mm.rdr.Render()
	for {
		in := make([]byte, 3)
		if len(carry) > 0 {
			in = carry
			carry = []byte{}
		} else if _, err := os.Stdin.Read(in); err != nil {
			p.editing = false
			_ = p.mm.rdr.Render()
			return err
		}

		if slices.ContainsFunc(KeyBinds.Exit, func(v keybind) bool { return slices.Equal(v, in) }) || slices.ContainsFunc(KeyBinds.Confirm, func(v keybind) bool { return slices.Equal(v, in) }) {
			break
		} else if slices.ContainsFunc(KeyBinds.Delete, func(v keybind) bool { return slices.Equal(v, in) }) {
			vStr := hex.EncodeToString(p.value[p.selected : p.selected+2])
			if len(vStr) <= 0 {
				continue
			}
			v, err := hex.DecodeString(strings.Repeat("0", 5-len(vStr)) + vStr[:len(vStr)-1])
			if err != nil {
				p.editing = false
				_ = p.mm.rdr.Render()
				return err
			}
			p.value[p.selected], p.value[p.selected+1] = v[0], v[1]
			_ = p.mm.rdr.Render()
			continue
		} else if strings.ContainsAny(string(CharSets.Hex), string(in[:])) {
			vStr := strings.Replace(hex.EncodeToString(p.value[p.selected:p.selected+2])+string(bytes.Trim(in, "\x00")[:]), "0", "", 1)
			if len(vStr) > 4 {
				// vStr = "ffff"
				if p.selected >= 14 {
					continue
				}
				p.selected += 2
				carry = in
				continue
			}
			v, err := hex.DecodeString(vStr)
			if err != nil {
				p.editing = false
				_ = p.mm.rdr.Render()
				return err
			}
			p.value[p.selected], p.value[p.selected+1] = v[0], v[1]
			_ = p.mm.rdr.Render()
		} else if slices.ContainsFunc(KeyBinds.Right, func(v keybind) bool { return slices.Equal(v, in) }) {
			p.selected = min(p.selected+2, 14)
			_ = p.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Left, func(v keybind) bool { return slices.Equal(v, in) }) {
			p.selected = max(p.selected-2, 0)
			_ = p.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Up, func(v keybind) bool { return slices.Equal(v, in) }) {
			if int(p.value[p.selected+1]) < 255 {
				p.value[p.selected+1] = p.value[p.selected+1] + 1
			} else if int(p.value[p.selected]) < 255 {
				p.value[p.selected] = p.value[p.selected] + 1
			} else {
				p.value[p.selected], p.value[p.selected+1] = byte(0), byte(0)
			}
			_ = p.mm.rdr.Render()
			continue
		} else if slices.ContainsFunc(KeyBinds.Down, func(v keybind) bool { return slices.Equal(v, in) }) {
			if int(p.value[p.selected]) > 0 {
				p.value[p.selected] = p.value[p.selected] - 1
			} else if int(p.value[p.selected+1]) > 0 {
				p.value[p.selected+1] = p.value[p.selected+1] - 1
			} else {
				p.value[p.selected], p.value[p.selected+1] = byte(255), byte(255)
			}
			_ = p.mm.rdr.Render()
			continue
		}

	}
	p.editing = false
	_ = p.mm.rdr.Render()
	return nil
}
