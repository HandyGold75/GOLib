package keyboard

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var Errors = struct{ NoKeyBoardFound error }{NoKeyBoardFound: errors.New("no keyboard found")}

type inputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

var inputEventSize = int(unsafe.Sizeof(inputEvent{}))

func (i *inputEvent) String() string  { return keyCodeMap[i.Code] } // Returns key as a string.
func (i *inputEvent) IsPress() bool   { return i.Value == 1 }       // Returns true if event is a press event.
func (i *inputEvent) IsRelease() bool { return i.Value == 0 }       // Returns true if event is a release event.

type KeyBoard struct{ fd *os.File }

// Returns the first keyboard containing `name`, if name is empty then uses `keyboard` as `name`.
//
// Only returns an error if no keyboards are found matching `name`.
func NewKeyboard(name string) (*KeyBoard, error) {
	if name == "" {
		name = "keyboard"
	}
	for i := range 255 {
		buff, err := os.ReadFile(fmt.Sprintf("/sys/class/input/event%d/device/name", i))
		if err != nil {
			continue
		}
		device := strings.ToLower(string(buff))

		if strings.Contains(device, strings.ToLower(name)) {
			fd, err := os.OpenFile(fmt.Sprintf("/dev/input/event%d", i), os.O_RDWR, os.ModeCharDevice)
			if err != nil {
				continue
			}
			return &KeyBoard{fd: fd}, nil
		}
	}
	return &KeyBoard{}, Errors.NoKeyBoardFound
}

// Returns the all keyboards containing `name`, if name is empty then uses `keyboard` as `name`.
//
// Only returns an error if no keyboards are found matching `name`.
func NewKeyboards(name string) ([]*KeyBoard, error) {
	if name == "" {
		name = "keyboard"
	}
	ret := []*KeyBoard{}
	for i := range 255 {
		buff, err := os.ReadFile(fmt.Sprintf("/sys/class/input/event%d/device/name", i))
		if err != nil {
			continue
		}
		device := strings.ToLower(string(buff))

		if strings.Contains(device, strings.ToLower(name)) {
			fd, err := os.OpenFile(fmt.Sprintf("/dev/input/event%d", i), os.O_RDWR, os.ModeCharDevice)
			if err != nil {
				continue
			}
			ret = append(ret, &KeyBoard{fd: fd})
		}
	}
	if len(ret) <= 0 {
		return ret, Errors.NoKeyBoardFound
	}
	return ret, nil
}

// Returns channel where events can be read from.
//
// Caller is responsible for closing the channel when finish.
func (k *KeyBoard) Read() chan inputEvent {
	event := make(chan inputEvent)
	go func(event chan inputEvent) {
		for {
			buffer := make([]byte, inputEventSize)
			n, err := k.fd.Read(buffer)
			if err != nil {
				break
			}
			if n <= 0 {
				continue
			}
			e := &inputEvent{}
			err = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, e)
			if err != nil {
				break
			}
			event <- *e
		}
	}(event)
	return event
}

// Press or release key on the keyboard.
func (k *KeyBoard) Send(press bool, key string) error {
	dir := int32(0)
	if press {
		dir = 1
	}
	key = strings.ToUpper(key)
	code := uint16(0)
	for c, k := range keyCodeMap {
		if k == key {
			code = c
			break
		}
	}
	if code == 0 {
		return fmt.Errorf("%s key not found in key code map", key)
	}
	if err := binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: 0x01, Code: code, Value: dir}); err != nil {
		return err
	}
	return binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: 0x00, Code: 0, Value: 0})
}

// Press and release a key on the keyboard.
func (k *KeyBoard) Press(key string) error {
	key = strings.ToUpper(key)
	code := uint16(0)
	for c, k := range keyCodeMap {
		if k == key {
			code = c
			break
		}
	}
	if code == 0 {
		return fmt.Errorf("%s key not found in key code map", key)
	}
	for _, i := range []int32{1, 0} {
		if err := binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: 0x01, Code: code, Value: i}); err != nil {
			return err
		}
	}
	return binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: 0x00, Code: 0, Value: 0})
}

// Press and release a key on the keyboard while pressing and releasing another mod key around the key.
//
// The mod key does not need to be a modifier, this can be any key.
func (k *KeyBoard) PressWithMod(key string, mod string) error {
	key = strings.ToUpper(key)
	code := uint16(0)
	for c, k := range keyCodeMap {
		if k == key {
			code = c
			break
		}
	}
	if code == 0 {
		return fmt.Errorf("%s key not found in key code map", key)
	}
	mod = strings.ToUpper(mod)
	codeMod := uint16(0)
	for c, k := range keyCodeMap {
		if k == mod {
			codeMod = c
			break
		}
	}
	if codeMod == 0 {
		return fmt.Errorf("%s key not found in key code map", mod)
	}

	if err := binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: 0x01, Code: codeMod, Value: 1}); err != nil {
		return err
	}
	for _, i := range []int32{1, 0} {
		if err := binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: 0x01, Code: code, Value: i}); err != nil {
			return err
		}
	}
	if err := binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: 0x01, Code: codeMod, Value: 0}); err != nil {
		return err
	}

	return binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: 0x00, Code: 0, Value: 0})
}

// Close the keyboard.
func (k *KeyBoard) Close() error {
	if k.fd == nil {
		return nil
	}
	return k.fd.Close()
}

// https://raw.githubusercontent.com/torvalds/linux/master/include/uapi/linux/input-event-codes.h
var keyCodeMap = map[uint16]string{
	1:   "ESC",
	2:   "1",
	3:   "2",
	4:   "3",
	5:   "4",
	6:   "5",
	7:   "6",
	8:   "7",
	9:   "8",
	10:  "9",
	11:  "0",
	12:  "-",
	13:  "=",
	14:  "BS",
	15:  "TAB",
	16:  "Q",
	17:  "W",
	18:  "E",
	19:  "R",
	20:  "T",
	21:  "Y",
	22:  "U",
	23:  "I",
	24:  "O",
	25:  "P",
	26:  "[",
	27:  "]",
	28:  "ENTER",
	29:  "L_CTRL",
	30:  "A",
	31:  "S",
	32:  "D",
	33:  "F",
	34:  "G",
	35:  "H",
	36:  "J",
	37:  "K",
	38:  "L",
	39:  ";",
	40:  "'",
	41:  "`",
	42:  "L_SHIFT",
	43:  "\\",
	44:  "Z",
	45:  "X",
	46:  "C",
	47:  "V",
	48:  "B",
	49:  "N",
	50:  "M",
	51:  ",",
	52:  ".",
	53:  "/",
	54:  "R_SHIFT",
	55:  "*",
	56:  "L_ALT",
	57:  " ",
	58:  "CAPS_LOCK",
	59:  "F1",
	60:  "F2",
	61:  "F3",
	62:  "F4",
	63:  "F5",
	64:  "F6",
	65:  "F7",
	66:  "F8",
	67:  "F9",
	68:  "F10",
	69:  "NUM_LOCK",
	70:  "SCROLL_LOCK",
	71:  "HOME",
	72:  "UP_8",
	73:  "PGUP_9",
	74:  "-",
	75:  "LEFT_4",
	76:  "5",
	77:  "RT_ARROW_6",
	78:  "+",
	79:  "END_1",
	80:  "DOWN",
	81:  "PGDN_3",
	82:  "INS",
	83:  "DEL",
	84:  "",
	85:  "",
	86:  "",
	87:  "F11",
	88:  "F12",
	89:  "",
	90:  "",
	91:  "",
	92:  "",
	93:  "",
	94:  "",
	95:  "",
	96:  "R_ENTER",
	97:  "R_CTRL",
	98:  "/",
	99:  "PRT_SCR",
	100: "R_ALT",
	101: "",
	102: "HOME",
	103: "UP",
	104: "PGUP",
	105: "LEFT",
	106: "RIGHT",
	107: "END",
	108: "DOWN",
	109: "PGDN",
	110: "INSERT",
	111: "DEL",
	112: "",
	113: "",
	114: "",
	115: "",
	116: "",
	117: "",
	118: "",
	119: "PAUSE",
}
