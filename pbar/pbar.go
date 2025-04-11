package pbar

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/term"
)

var (
	// Completed actions.
	Done = 0

	// Total actions.
	Total = 0

	// Size multiplier (0.0 - 1.0).
	Size = 0.25

	// Verbosity.
	Verbose = 0
)

// Log progress bar.
func Log() {
	prog := 0.0
	if Total != 0 {
		prog = float64(Done) / float64(max(0, Total))
	}

	width, _, _ := term.GetSize(0)
	progLen := float64(prog) * float64(width) * min(1.0, max(0.0, Size))

	lastChar := "█"
	switch i := progLen - float64(int(progLen)); {
	case i < 0.25 && prog != 1:
		lastChar += " "
	case i < 0.5 && prog != 1:
		lastChar += "▄"
	case i < 0.75 && prog != 1:
		lastChar += "▀"
	}

	msg := fmt.Sprintf("\r|%-"+strconv.Itoa(int(float64(width)/4)+1)+"v| %.1f%%", strings.Repeat("█", int(progLen))+lastChar, prog*100)

	fmt.Printf("\r%"+strconv.Itoa(width)+"v", "")
	if len([]rune(msg)) > width {
		fmt.Printf("%."+strconv.Itoa(width-3)+"s...", msg)
	} else {
		fmt.Printf("%."+strconv.Itoa(width)+"s", msg)
	}
}

// Short for `pbar.Done += 1; pbar.Log()`.
func Next(msg string, format string, v ...any) { Done += 1; Log() }

// Short for `pbar.Done -= 1; pbar.Log()`.
func Back() { Done -= 1; Log() }

// Log progress bar with message capabilities.
//
// `msg`, `format` and `v` may be ignored based on `pbar.Verbose`.
//
// Verbosities
//
//	<= 0: Plain progress bar.
//	== 1: `msg` is appended to the progress bar.
//	>= 2: Progress bar is discarded, `msgLong` is forwarded to `fmt.Print`.
func LogMsg(msg string, msgLong string) {
	if Verbose >= 2 {
		fmt.Print(msgLong)
		return
	}

	prog := 0.0
	if Total != 0 {
		prog = float64(Done) / float64(max(0, Total))
	}

	width, _, _ := term.GetSize(0)
	progLen := float64(prog) * float64(width) * min(1.0, max(0.0, Size))

	lastChar := "█"
	switch i := progLen - float64(int(progLen)); {
	case i < 0.25 && prog != 1:
		lastChar += " "
	case i < 0.5 && prog != 1:
		lastChar += "▄"
	case i < 0.75 && prog != 1:
		lastChar += "▀"
	}

	if Verbose >= 1 {
		msg = fmt.Sprintf("\r|%-"+strconv.Itoa(int(float64(width)/4)+1)+"v| %.1f%% (%v/%v) -> %v", strings.Repeat("█", int(progLen))+lastChar, prog*100, Done, Total, msg)
	} else {
		msg = fmt.Sprintf("\r|%-"+strconv.Itoa(int(float64(width)/4)+1)+"v| %.1f%%", strings.Repeat("█", int(progLen))+lastChar, prog*100)
	}

	fmt.Printf("\r%"+strconv.Itoa(width)+"v", "")
	if len([]rune(msg)) > width {
		fmt.Printf("%."+strconv.Itoa(width-3)+"s...", msg)
	} else {
		fmt.Printf("%."+strconv.Itoa(width)+"s", msg)
	}
}

// Short for `pbar.Done += 1; pbar.LogMsg(msg, msgLong)`.
func NextMsg(msg string, msgLong string) { Done += 1; LogMsg(msg, msgLong) }

// Short for `pbar.Done -= 1; pbar.LogMsg(msg, msgLong)`.
func BackMsg(msg string, msgLong string) { Done -= 1; LogMsg(msg, msgLong) }
