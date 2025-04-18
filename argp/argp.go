package argp

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

func expandArgs[T any](s *T, args []string) []string {
	allSwitches := []string{}
	allPrefixes := []string{}
	forEachStructField(s, func(field reflect.StructField, value reflect.Value) {
		if field.Type.String() == "[]string" {
			return
		}
		prefix := field.Tag.Get("prefix")
		if prefix == "" {
			prefix = "-"
		}
		if !slices.Contains(allPrefixes, prefix) {
			allPrefixes = append(allPrefixes, prefix)
		}
		for _, swt := range strings.Split(field.Tag.Get("switch"), ",") {
			allSwitches = append(allSwitches, prefix+swt)
		}
	})

	for i := 0; i < len(args); i++ {
		if strings.Contains(args[i], "=") || slices.Contains(allSwitches, args[i]) {
			continue
		}

		prefix := ""
		for _, item := range allPrefixes {
			if strings.HasPrefix(args[i], item) {
				prefix = item
				break
			}
		}
		if prefix == "" {
			continue
		}

		newArgs := append([]string{}, args[:i]...) // Deepcopy
		for _, a := range strings.Split(strings.Replace(args[i], prefix, "", 1), "") {
			newArgs = append(newArgs, prefix+a)
		}
		args = append(newArgs, args[i+1:]...)
		i = i + (len(newArgs) - 1)
	}

	return args
}

func parseArgs[T any](s *T, args []string) []string {
	allPrefixes := []string{}
	forEachStructField(s, func(field reflect.StructField, value reflect.Value) {
		if field.Type.String() == "[]string" {
			return
		}
		prefix := field.Tag.Get("prefix")
		if prefix == "" {
			prefix = "-"
		}
		if !slices.Contains(allPrefixes, prefix) {
			allPrefixes = append(allPrefixes, prefix)
		}
	})

	forEachStructField(s, func(field reflect.StructField, value reflect.Value) {
		if field.Type.String() == "[]string" {
			return
		}

		prefix := field.Tag.Get("prefix")
		if prefix == "" {
			prefix = "-"
		}

		index := -1
		for _, swt := range strings.Split(field.Tag.Get("switch"), ",") {
			index = slices.IndexFunc(args, func(a string) bool { return a == prefix+swt || strings.HasPrefix(a, prefix+swt+"=") })
			if index > -1 {
				if slices.Contains(strings.Split(field.Tag.Get("opts"), ","), "help") {
					HelpMenu(*s, true)
					os.Exit(0)
				}
				break
			}
		}

		val := field.Tag.Get("default")
		if index > -1 {
			if strings.Contains(args[index], "=") {
				val = strings.Split(args[index], "=")[1]
			} else if field.Type.String() == "bool" {
				val = "true"
			} else if index+1 < len(args) && !slices.ContainsFunc(allPrefixes, func(item string) bool { return strings.HasPrefix(args[index+1], item) }) {
				val = args[index+1]
				args = slices.Delete(args, index, index+1)
			}
			args = slices.Delete(args, index, index+1)

		} else if slices.Contains(strings.Split(field.Tag.Get("opts"), ","), "posistional") {
			i := slices.IndexFunc(args, func(a string) bool {
				for _, p := range allPrefixes {
					if strings.HasPrefix(a, p) {
						return false
					}
				}
				return true
			})
			if i > -1 {
				val = args[i]
				args = slices.Delete(args, i, i+1)
			}
		}

		if val == "" && field.Type.String() != "bool" {
			if slices.Contains(strings.Split(field.Tag.Get("opts"), ","), "required") {
				fmt.Println("Missing required argument " + field.Name)
				HelpMenu(*s, false)
				os.Exit(0)
			} else if index > -1 {
				fmt.Println("Missing argument value for " + field.Name)
				HelpMenu(*s, false)
				os.Exit(0)
			}
			return
		}

		switch field.Type.String() {
		case "string":
			value.SetString(val)

		case "bool":
			value.SetBool(val == "true")

		case "int", "int8", "int16", "int32", "int64":
			v, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				fmt.Println("Invalid int format '" + val + "'")
				HelpMenu(*s, false)
				os.Exit(0)
			}
			value.SetInt(v)

		case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
			v, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				fmt.Println("Invalid uint format '" + val + "'")
				HelpMenu(*s, false)
				os.Exit(0)
			}
			value.SetUint(v)

		case "float32", "float64":
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				fmt.Println("Invalid float format '" + val + "'")
				HelpMenu(*s, false)
				os.Exit(0)
			}
			value.SetFloat(v)

		default:
			panic("unsuported type " + field.Type.String())
		}
	})

	forEachStructField(s, func(field reflect.StructField, value reflect.Value) {
		if field.Type.String() != "[]string" {
			return
		}
		if len(args) == 0 && slices.Contains(strings.Split(field.Tag.Get("opts"), ","), "required") {
			fmt.Println("Missing required argument " + field.Name)
			HelpMenu(*s, false)
			os.Exit(0)
		}
		value.Set(reflect.ValueOf(args))
		args = []string{}
	})

	return args
}

func forEachStructField[T any](s *T, handler func(reflect.StructField, reflect.Value)) {
	tOf := reflect.TypeOf(*s)
	if tOf.Kind() != reflect.Struct {
		panic("s is not a stuct")
	}
	vOf := reflect.ValueOf(s).Elem()
	if vOf.Kind() != reflect.Struct {
		panic("s is not a stuct")
	}

	for i := 0; i < tOf.NumField(); i++ {
		tOfField := tOf.Field(i)
		vOfField := vOf.Field(i)
		if !vOfField.CanSet() {
			continue
		}
		if tOfField.Tag.Get("switch") == "" && tOfField.Type.String() != "[]string" {
			panic("no switch specified for field " + tOfField.Name)
		}
		handler(tOfField, vOfField)
	}
}

// Shows help menu related to `s`. Panics if `s` is not of type struct or a public field doesn't contain a switch tag. Private struct fields are ignored.
func HelpMenu[T any](s T, details bool) {
	maxLenSwts := 0
	forEachStructField(&s, func(field reflect.StructField, _ reflect.Value) {
		if field.Type.String() == "[]string" {
			return
		}
		lenSwt := len(field.Tag.Get("switch")) + max(1, len(field.Tag.Get("prefix"))) + (len(strings.Split(field.Tag.Get("switch"), ",")))
		if lenSwt > maxLenSwts {
			maxLenSwts = lenSwt
		}
	})

	mainHelpMsg := ""
	helpMenuStr := ""
	argsStr := ""
	forEachStructField(&s, func(field reflect.StructField, _ reflect.Value) {
		prefix := field.Tag.Get("prefix")
		if prefix == "" {
			prefix = "-"
		}

		swts := ""
		for _, swt := range strings.Split(field.Tag.Get("switch"), ",") {
			swts += " " + prefix + swt
		}

		switch field.Type.String() {
		case "string":
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + " <" + field.Type.String() + ">]"

		case "bool":
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + "]"

		case "int", "int8", "int16", "int32", "int64":
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + " <" + field.Type.String() + ">]"

		case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + " <" + field.Type.String() + ">]"

		case "float32", "float64":
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + " <" + field.Type.String() + ">]"

		case "[]string":
			argsStr = " [" + field.Name + "...]" + argsStr
			swts = ""

		default:
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + "]"
		}

		helpMenuStr += fmt.Sprintf("%v\r\n %-"+strconv.Itoa(maxLenSwts)+"v  %-9v", field.Name, swts, "<"+field.Type.String()+">")
		if field.Tag.Get("opts") != "" {
			helpMenuStr += fmt.Sprintf(" (%v)", field.Tag.Get("opts"))
		}
		helpMenuStr += "\r\n"
		if field.Tag.Get("help") != "" {
			if slices.Contains(strings.Split(field.Tag.Get("opts"), ","), "help") {
				mainHelpMsg = fmt.Sprintf("	%v\r\n", field.Tag.Get("help"))
				return
			}
			helpMenuStr += fmt.Sprintf("	%v\r\n", field.Tag.Get("help"))
		}
	})

	execPath, err := os.Executable()
	if err != nil {
		execPath = "exec"
	}
	execPathSplit := strings.Split(strings.ReplaceAll(execPath, "\\", "/"), "/")
	fmt.Print("Usage: " + execPathSplit[len(execPathSplit)-1] + argsStr + "\r\n")
	if details {
		fmt.Print(mainHelpMsg + "\r\n" + helpMenuStr)
	}
}

// Parses args into s. Panics if s is not of type struct or a public field doesn't contain a switch tag. Private struct fields are ignored.
// Runs args.HelpMenu and exits gracefully if user input is invalid.
//
// Struct format:
//
//	struct { Field type `tag:"value"` }
//
// Supported types:
//
//	string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64
//
// Available tags:
//
//	switch:  List of switches that can map to this field (Required).
//	prefix:  Prefix for the switches (Default: "-").
//	opts:    Optional parameters [posistional,required,help].
//	default: Optional default value.
//	help:    Help message.
//
// Opts:
//
//	posistional: Field will be assigned using the first argument without a switch or prefix.
//	             Needs to be last in the struct, this is to ensure no value arguments are taken from other switches in cases " " is used to seperate key and value arguments.
//	             As a special case the type of a posistional may be []string to populate this filed with left over arguments, in this case the tags switch, prefix and default wil be ignored and order in stuct non-important.
//	required:    Field will be required, shows help menu if missing, if a default is given required is ignored.
//	help:        Quick opt for help menu implementation, when a switch from this field is present `args.HelpMenu` gets called followed by `os.Exit(0)`.
//	             Needs to be first in the struct, this is to ensure the help menu is show when a help switch is present in cases a panic is trow by other switches.
//	             The help tag message provided becomes the decription of the executable.
//
// Example:
//
//	args struct {
//		Help   bool     `switch:"h,-help"   opts:"help" help:"Some help for exec."`
//		FieldF []string `opts:"posistional,required"`
//		FieldA float64  `switch:"a"         prefix:"/" help:"Some help for e."`
//		FieldB bool     `switch:"b"`
//		FieldC int      `switch:"c,-cc"     default:"5"`
//		FieldD uint     `switch:"d,-dd-ddd" default:"7" opts:"required"`
//		FieldE string   `switch:"e,-ee"     opts:"posistional"`
//	}
//
// Usage: exec -a=-1.1 -bc 10 --dd-ddd=2 /e "some message"
func Parse[T any](s T, args []string) T {
	args = expandArgs(&s, args)
	args = parseArgs(&s, args)

	if len(args) > 0 {
		unknown := []string{}
		for _, arg := range args {
			if strings.Contains(arg, "=") {
				unknown = append(unknown, strings.Split(arg, "=")[0])
				continue
			}
			unknown = append(unknown, arg)
		}

		fmt.Println("Unknown arguments '" + strings.Join(unknown, "', '") + "'")
		HelpMenu(s, false)
		os.Exit(0)
	}

	return s
}

// Short for `argp.Parse(s, os.Args[1:])`
//
// Parses args into s. Panics if s is not of type struct or a public field doesn't contain a switch tag. Private struct fields are ignored.
// Runs args.HelpMenu and exits gracefully if user input is invalid.
//
// Struct format:
//
//	struct { Field type `tag:"value"` }
//
// Supported types:
//
//	string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64
//
// Available tags:
//
//	switch:  List of switches that can map to this field (Required).
//	prefix:  Prefix for the switches (Default: "-").
//	opts:    Optional parameters [posistional,required,help].
//	default: Optional default value.
//	help:    Help message.
//
// Opts:
//
//	posistional: Field will be assigned using the first argument without a switch or prefix.
//	             Needs to be last in the struct, this is to ensure no value arguments are taken from other switches in cases " " is used to seperate key and value arguments.
//	             As a special case the type of a posistional may be []string to populate this filed with left over arguments, in this case the tags switch, prefix and default wil be ignored and order in stuct non-important.
//	required:    Field will be required, shows help menu if missing, if a default is given required is ignored.
//	help:        Quick opt for help menu implementation, when a switch from this field is present `args.HelpMenu` gets called followed by `os.Exit(0)`.
//	             Needs to be first in the struct, this is to ensure the help menu is show when a help switch is present in cases a panic is trow by other switches.
//	             The help tag message provided becomes the decription of the executable.
//
// Example:
//
//	args struct {
//		Help   bool     `switch:"h,-help"   opts:"help" help:"Some help for exec."`
//		FieldF []string `opts:"posistional,required"`
//		FieldA float64  `switch:"a"         prefix:"/" help:"Some help for e."`
//		FieldB bool     `switch:"b"`
//		FieldC int      `switch:"c,-cc"     default:"5"`
//		FieldD uint     `switch:"d,-dd-ddd" default:"7" opts:"required"`
//		FieldE string   `switch:"e,-ee"     opts:"posistional"`
//	}
//
// Usage: exec -a=-1.1 -bc 10 --dd-ddd=2 /e "some message"
func ParseArgs[T any](s T) T { return Parse(s, os.Args[1:]) }
