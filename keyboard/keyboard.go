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

const (
	evSyn uint16 = 0x00
	evKey uint16 = 0x01
)

type inputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

var inputEventSize = int(unsafe.Sizeof(inputEvent{}))

func (i *inputEvent) String() string  { return keyCodeMap[i.Code] }           // Returns key as a string.
func (i *inputEvent) IsPress() bool   { return i.Value == int32(KeyPress) }   // Returns true if event is a press event.
func (i *inputEvent) IsRelease() bool { return i.Value == int32(KeyRelease) } // Returns true if event is a release event.

type keyEvent int32

const (
	KeyPress   keyEvent = 1
	KeyRelease keyEvent = 0
)

type KeyBoard struct {
	name string
	fd   *os.File
}

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
			return &KeyBoard{name: device, fd: fd}, nil
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
			ret = append(ret, &KeyBoard{name: device, fd: fd})
		}
	}
	if len(ret) <= 0 {
		return ret, Errors.NoKeyBoardFound
	}
	return ret, nil
}

// Returns the keyboard name.
func (k *KeyBoard) Name() string { return k.name }

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
func (k *KeyBoard) Send(direction keyEvent, key string) error {
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
	if err := binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: evKey, Code: code, Value: int32(direction)}); err != nil {
		return err
	}
	return binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: evSyn, Code: 0, Value: 0})
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
	for _, i := range []keyEvent{KeyPress, KeyRelease} {
		if err := binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: evKey, Code: code, Value: int32(i)}); err != nil {
			return err
		}
	}
	return binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: evSyn, Code: 0, Value: 0})
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

	if err := binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: evKey, Code: codeMod, Value: int32(KeyPress)}); err != nil {
		return err
	}
	for _, i := range []keyEvent{KeyPress, KeyRelease} {
		if err := binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: evKey, Code: code, Value: int32(i)}); err != nil {
			return err
		}
	}
	if err := binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: evKey, Code: codeMod, Value: int32(KeyRelease)}); err != nil {
		return err
	}

	return binary.Write(k.fd, binary.LittleEndian, inputEvent{Type: evSyn, Code: 0, Value: 0})
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
	0:  "RESERVED",
	1:  "ESC",
	2:  "1",
	3:  "2",
	4:  "3",
	5:  "4",
	6:  "5",
	7:  "6",
	8:  "7",
	9:  "8",
	10: "9",
	11: "0",
	12: "-",
	13: "=",
	14: "BACKSPACE",
	15: "TAB",
	16: "Q",
	17: "W",
	18: "E",
	19: "R",
	20: "T",
	21: "Y",
	22: "U",
	23: "I",
	24: "O",
	25: "P",
	26: "[",
	27: "]",
	28: "ENTER",
	29: "LEFTCTRL",
	30: "A",
	31: "S",
	32: "D",
	33: "F",
	34: "G",
	35: "H",
	36: "J",
	37: "K",
	38: "L",
	39: ";",
	40: "'",
	41: "`",
	42: "LEFTSHIFT",
	43: "\\",
	44: "Z",
	45: "X",
	46: "C",
	47: "V",
	48: "B",
	49: "N",
	50: "M",
	51: ",",
	52: ".",
	53: "/",
	54: "RIGHTSHIFT",
	55: "KPASTERISK",
	56: "LEFTALT",
	57: " ",
	58: "CAPSLOCK",
	59: "F1",
	60: "F2",
	61: "F3",
	62: "F4",
	63: "F5",
	64: "F6",
	65: "F7",
	66: "F8",
	67: "F9",
	68: "F10",
	69: "NUMLOCK",
	70: "SCROLLLOCK",
	71: "KP7",
	72: "KP8",
	73: "KP9",
	74: "KPMINUS",
	75: "KP4",
	76: "KP5",
	77: "KP6",
	78: "KPPLUS",
	79: "KP1",
	80: "KP2",
	81: "KP3",
	82: "KP0",
	83: "KPDOT",

	85:  "ZENKAKUHANKAKU",
	86:  "102ND",
	87:  "F11",
	88:  "F12",
	89:  "RO",
	90:  "KATAKANA",
	91:  "HIRAGANA",
	92:  "HENKAN",
	93:  "KATAKANAHIRAGANA",
	94:  "MUHENKAN",
	95:  "KPJPCOMMA",
	96:  "KPENTER",
	97:  "RIGHTCTRL",
	98:  "KPSLASH",
	99:  "SYSRQ",
	100: "RIGHTALT",
	101: "LINEFEED",
	102: "HOME",
	103: "UP",
	104: "PAGEUP",
	105: "LEFT",
	106: "RIGHT",
	107: "END",
	108: "DOWN",
	109: "PAGEDOWN",
	110: "INSERT",
	111: "DELETE",
	112: "MACRO",
	113: "MUTE",
	114: "VOLUMEDOWN",
	115: "VOLUMEUP",
	116: "POWER",
	117: "KPEQUAL",
	118: "KPPLUSMINUS",
	119: "PAUSE",
	120: "SCALE",

	121: "KPCOMMA",
	122: "HANGEUL",
	123: "HANJA",
	124: "YEN",
	125: "LEFTMETA",
	126: "RIGHTMETA",
	127: "COMPOSE",

	128: "STOP",
	129: "AGAIN",
	130: "PROPS",
	131: "UNDO",
	132: "FRONT",
	133: "COPY",
	134: "OPEN",
	135: "PASTE",
	136: "FIND",
	137: "CUT",
	138: "HELP",
	139: "MENU",
	140: "CALC",
	141: "SETUP",
	142: "SLEEP",
	143: "WAKEUP",
	144: "FILE",
	145: "SENDFILE",
	146: "DELETEFILE",
	147: "XFER",
	148: "PROG1",
	149: "PROG2",
	150: "WWW",
	151: "MSDOS",
	152: "COFFEE",
	153: "ROTATE_DISPLAY",
	154: "CYCLEWINDOWS",
	155: "MAIL",
	156: "BOOKMARKS",
	157: "COMPUTER",
	158: "BACK",
	159: "FORWARD",
	160: "CLOSECD",
	161: "EJECTCD",
	162: "EJECTCLOSECD",
	163: "NEXTSONG",
	164: "PLAYPAUSE",
	165: "PREVIOUSSONG",
	166: "STOPCD",
	167: "RECORD",
	168: "REWIND",
	169: "PHONE",
	170: "ISO",
	171: "CONFIG",
	172: "HOMEPAGE",
	173: "REFRESH",
	174: "EXIT",
	175: "MOVE",
	176: "EDIT",
	177: "SCROLLUP",
	178: "SCROLLDOWN",
	179: "KPLEFTPAREN",
	180: "KPRIGHTPAREN",
	181: "NEW",
	182: "REDO",

	183: "F13",
	184: "F14",
	185: "F15",
	186: "F16",
	187: "F17",
	188: "F18",
	189: "F19",
	190: "F20",
	191: "F21",
	192: "F22",
	193: "F23",
	194: "F24",

	200: "PLAYCD",
	201: "PAUSECD",
	202: "PROG3",
	203: "PROG4",
	204: "ALL_APPLICATIONS",
	205: "SUSPEND",
	206: "CLOSE",
	207: "PLAY",
	208: "FASTFORWARD",
	209: "BASSBOOST",
	210: "PRINT",
	211: "HP",
	212: "CAMERA",
	213: "SOUND",
	214: "QUESTION",
	215: "EMAIL",
	216: "CHAT",
	217: "SEARCH",
	218: "CONNECT",
	219: "FINANCE",
	220: "SPORT",
	221: "SHOP",
	222: "ALTERASE",
	223: "CANCEL",
	224: "BRIGHTNESSDOWN",
	225: "BRIGHTNESSUP",
	226: "MEDIA",

	227: "SWITCHVIDEOMODE",

	228: "KBDILLUMTOGGLE",
	229: "KBDILLUMDOWN",
	230: "KBDILLUMUP",

	231: "SEND",
	232: "REPLY",
	233: "FORWARDMAIL",
	234: "SAVE",
	235: "DOCUMENTS",

	236: "BATTERY",

	237: "BLUETOOTH",
	238: "WLAN",
	239: "UWB",

	240: "UNKNOWN",

	241: "VIDEO_NEXT",
	242: "VIDEO_PREV",
	243: "BRIGHTNESS_CYCLE",
	244: "BRIGHTNESS_AUTO",
	245: "DISPLAY_OFF",

	246: "WWAN",
	247: "RFKILL",

	248: "MICMUTE",
}
