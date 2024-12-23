package pbar

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/term"
)

var (
	// Done actions to caculate the % completion from.
	// Can be set manually or progressivaly with Next and Back.
	ActionsDone = 0

	// Total actions to caculate the % completion from.
	// Should be set manually.
	ActionsTotal = 0

	// Verbosity 1: extra info is added and msg is appended to the progress bar in Log.
	// The progress bar will be dynamicly size, progress bar and msg will be cut off respecting the terminal width.
	//
	// Verbosity 2 or larger progress bar will act like fmt.Printf.
	Verbose = 0
)

// Prints a progress bar respecting ActionsDone and ActionsTotal to define the progress.
//
// By default msg, format and v are ignored.
// The progress bar will be dynamicly size and will be cut off if nessisairy respecting the terminal width.
func Log(msg string, format string, v ...any) {
	if Verbose >= 2 {
		fmt.Printf(format, v...)
		return
	}

	progress := 0.0
	if ActionsTotal != 0 {
		progress = float64(ActionsDone) / float64(ActionsTotal)
	}

	width, _, _ := term.GetSize(0)
	visualProgress := float64(progress) * float64(width) / 4
	visualProgressStr := strings.Repeat("█", int(visualProgress))
	visualProgressDesimal := visualProgress - float64(int(visualProgress))

	if visualProgressDesimal < 0.25 && progress != 1 {
		visualProgressStr += " "
	} else if visualProgressDesimal < 0.75 && progress != 1 {
		if int(visualProgress)%2 == 0 {
			visualProgressStr += "▄"
		} else {
			visualProgressStr += "▀"
		}
	} else {
		visualProgressStr += "█"
	}

	if Verbose >= 1 {
		msg = fmt.Sprintf("\r|%-"+strconv.Itoa(int(float64(width)/4)+1)+"v| %.1f%% (%v/%v) -> %v", visualProgressStr, progress*100, ActionsDone, ActionsTotal, msg)
	} else {
		msg = fmt.Sprintf("\r|%-"+strconv.Itoa(int(float64(width)/4)+1)+"v| %.1f%%", visualProgressStr, progress*100)
	}

	fmt.Printf("\r%"+strconv.Itoa(width)+"v", "")
	if len([]rune(msg)) > width {
		fmt.Printf("%."+strconv.Itoa(width-3)+"s...", msg)
	} else {
		fmt.Printf("%."+strconv.Itoa(width)+"s", msg)
	}
}

// Same as Log but increments ActionsDone before loging if the verbosity is lower then 2
func Next(msg string, format string, v ...any) {
	if Verbose < 2 {
		ActionsDone += 1
	}
	Log(msg, format, v...)
}

// Same as Log but decrements ActionsDone before loging if the verbosity is lower then 2
func Back(msg string, format string, v ...any) {
	if Verbose < 2 {
		ActionsDone -= 1
	}
	Log(msg, format, v...)
}
