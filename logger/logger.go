package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

type (
	Logger struct {
		// When logging to file this file will be used.
		//
		// If `logger.DynamicFileName` is not nil this becomes the used path (.log Surfix is trimmed).
		FilePath string
		// Append the return to `logger.FilePath` (.log Surfix is trimmed from `logger.FilePath`).
		DynamicFileName func() string
		// Mapping of Vebosities to set allowed verbosities and their priority.
		Verbosities map[string]int
		// Minimal verbose priotity to log message to CLI.
		VerboseToCLI int
		// Minimal verbose priotity to log message to file.
		VerboseToFile int
		// Prepend logs with the date time.
		AppendDateTime bool
		// Prepend logs with the verbosity.
		AppendVerbosity bool
		// Prepend logs with this.
		PrepentCLI string
		// Called after every log to CLI.
		MessageCLIHook func(msg string)
		// Minimal char count a log part will take up.
		CharCountPerPart int
		// Minimal char space the verbosity part will take up (AppendVerbosity must be true to take effect).
		CharCountVerbosity int
		// When true RecordSepperator and EORSepperator are used when loggin to file, otherwise log the raw message.
		UseSeperators bool
		// Seperator string between message parts when logging to file (Logged message can not contain this string).
		RecordSepperator string
		// End of record string after a message when logging to file (Logged message can not contain this string).
		EORSepperator string
	}
)

var (
	// Default value for `logger.DynamicFileName`, does not effect exisiting loggers.
	DynamicFileName func() string = nil
	// Default value for `logger.Verbosities`, does not effect exisiting loggers.
	Verbosities = map[string]int{"high": 3, "medium": 2, "low": 1}
	// Default value for `logger.VerboseToCLI`, does not effect exisiting loggers.
	VerboseToCLI = 1
	// Default value for `logger.VerboseToFile`, does not effect exisiting loggers.
	VerboseToFile = 2
	// Default value for `logger.AppendDateTime`, does not effect exisiting loggers.
	AppendDateTime = true
	// Default value for `logger.AppendVerbosity`, does not effect exisiting loggers.
	AppendVerbosity = true
	// Default value for `logger.PrepentCLI`, does not effect exisiting loggers.
	PrepentCLI = ""
	// Default value for `logger.MessageCLIHook`, does not effect exisiting loggers.
	MessageCLIHook func(msg string) = nil
	// Default value for `logger.CharCountPerPart`, does not effect exisiting loggers.
	CharCountPerPart = 32
	// Default value for `logger.CharCountVerbosity`, does not effect exisiting loggers.
	CharCountVerbosity = 7
	// Default value for `logger.UseSeperators`, does not effect exisiting loggers.
	UseSeperators = true
	// Default value for `logger.RecordSepperator`, does not effect exisiting loggers.
	RecordSepperator = "<SEP>"
	// Default value for `logger.EORSepperator`, does not effect exisiting loggers.
	EORSepperator = "<EOR>\n"
)

func (logger Logger) logToCLI(verbosity string, msgs ...any) {
	width, _, _ := term.GetSize(0)
	msg := fmt.Sprintf(strings.Repeat("%-"+strconv.Itoa(min(logger.CharCountPerPart, int(float64(width)/float64(len(msgs)))))+"v", len(msgs)), msgs...)

	if logger.AppendVerbosity {
		msg = fmt.Sprintf("%-"+strconv.Itoa(logger.CharCountVerbosity)+"v ", verbosity) + msg
	}
	if logger.AppendDateTime {
		msg = "[" + time.Now().Format(time.DateTime) + "] " + msg
	}
	msg = logger.PrepentCLI + msg

	if len([]rune(msg)) > width {
		fmt.Printf("%."+strconv.Itoa(width-3)+"s...\n", msg)
	} else {
		fmt.Printf("%."+strconv.Itoa(width)+"s\n", msg)
	}

	if logger.MessageCLIHook != nil {
		logger.MessageCLIHook(msg)
	}
}

func (logger Logger) logToFile(verbosity string, msgs ...any) {
	var msg string

	if logger.UseSeperators {
		msg = fmt.Sprintf(strings.Repeat("%v"+logger.RecordSepperator, len(msgs)), msgs...)

		if logger.AppendVerbosity {
			msg = verbosity + logger.RecordSepperator + msg
		}
		if logger.AppendDateTime {
			msg = time.Now().Format(time.RFC3339Nano) + logger.RecordSepperator + msg
		}

		i := strings.LastIndex(msg, logger.RecordSepperator)
		msg = msg[:i] + strings.Replace(msg[i:], logger.RecordSepperator, "", 1) + logger.EORSepperator
	} else {
		msg = fmt.Sprintf(strings.Repeat("%-"+strconv.Itoa(logger.CharCountPerPart)+"v", len(msgs))+"\n", msgs...)

		if logger.AppendVerbosity {
			msg = fmt.Sprintf("%-"+strconv.Itoa(logger.CharCountVerbosity)+"v ", verbosity) + msg
		}
		if logger.AppendDateTime {
			msg = "[" + time.Now().Format(time.DateTime) + "] " + msg
		}
	}

	fp := logger.FilePath
	if logger.DynamicFileName != nil {
		fp = strings.TrimSuffix(fp, ".log") + "/" + logger.DynamicFileName()
	}

	if _, err := os.Stat(fp); os.IsNotExist(err) {
		fileSplit := strings.Split(strings.ReplaceAll(fp, "\\", "/"), "/")
		err := os.MkdirAll(strings.Join(fileSplit[:len(fileSplit)-1], "/"), os.ModePerm)
		if err != nil {
			logger.logToCLI("ERROR", "Failed creating logpath", err)
			return
		}
		if err := os.WriteFile(fp, []byte(msg), 0640); err != nil {
			logger.logToCLI("ERROR", "Failed creating logfile", err)
		}
		return
	}

	logFile, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.logToCLI("ERROR", "Failed opening logfile", err)
	}
	if _, err := logFile.Write([]byte(msg)); err != nil {
		logger.logToCLI("ERROR", "Failed writing to logfile", err)
	}
	if err := logFile.Close(); err != nil {
		logger.logToCLI("ERROR", "Failed closing logfile", err)
	}
}

// Create new logger instance.
//
// If log file is not present then it tries creating it.
//
// Log file is stored in `./golib/<name>.log` relative to `os.UserConfigDir`.
func New(name string) (*Logger, error) {
	file, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return NewAbs(file + "/golib/" + name + ".log"), nil
}

// Create new logger instance.
//
// If log file is not present then it tries creating it.
//
// Log file is stored in `./<name>.log` relative to `os.Executable`.
func NewRel(name string) (*Logger, error) {
	file, err := os.Executable()
	if err != nil {
		return nil, err
	}
	fileSplit := strings.Split(strings.ReplaceAll(file, "\\", "/"), "/")
	return NewAbs(strings.Join(fileSplit[:len(fileSplit)-1], "/") + "/" + name + ".log"), nil
}

// Create new logger instance.
//
// If log file is not present then it tries creating it.
func NewAbs(file string) *Logger {
	return &Logger{
		FilePath:           file,
		DynamicFileName:    DynamicFileName,
		Verbosities:        Verbosities,
		VerboseToCLI:       VerboseToCLI,
		VerboseToFile:      VerboseToFile,
		AppendDateTime:     AppendDateTime,
		AppendVerbosity:    AppendVerbosity,
		PrepentCLI:         PrepentCLI,
		MessageCLIHook:     MessageCLIHook,
		CharCountPerPart:   CharCountPerPart,
		CharCountVerbosity: CharCountVerbosity,
		UseSeperators:      UseSeperators,
		RecordSepperator:   RecordSepperator,
		EORSepperator:      EORSepperator,
	}
}

// Log an message.
//
// If `verbosity` is not present in `logger.Verbosities` then it is set to `{ "ERROR": 99 }`
//
// Message is logged to CLI if `verbosity >= logger.VerboseToCLI`.
//
// Message is logged to file if `verbosity >= logger.VerboseToFile`.
func (logger Logger) Log(verbosity string, msgs ...any) {
	verboseLevel, ok := logger.Verbosities[verbosity]
	if !ok {
		verbosity, verboseLevel = "ERROR", 99
	}

	if verboseLevel >= logger.VerboseToFile {
		for _, msg := range msgs {
			if strings.Contains(fmt.Sprintf("%v", msg), logger.RecordSepperator) {
				logger.logToCLI("ERROR", "Msg contains "+logger.RecordSepperator, msg)
				return
			}
			if strings.Contains(fmt.Sprintf("%v", msg), logger.EORSepperator) {
				logger.logToCLI("ERROR", "Msg contains "+strings.ReplaceAll(logger.EORSepperator, "\n", "\\n"), strings.ReplaceAll(msg.(string), "\n", "\\n"))
				return
			}
		}
		logger.logToFile(verbosity, msgs...)
	}
	if verboseLevel >= logger.VerboseToCLI {
		logger.logToCLI(verbosity, msgs...)
	}
}
