package argp

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

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
		if tOfField.Tag.Get("switch") == "" {
			panic("no switch specified for field " + tOfField.Name)
		}
		handler(tOfField, vOfField)
	}
}

// Shows help menu related to s. Panics if s is not of type struct. Private struct fields are ignored.
func HelpMenu[T any](s T, details bool) {
	maxLenSwts := 0
	forEachStructField(&s, func(field reflect.StructField, _ reflect.Value) {
		lenSwt := len(field.Tag.Get("switch")) + max(1, len(field.Tag.Get("prefix"))) + (len(strings.Split(field.Tag.Get("switch"), ",")))
		if lenSwt > maxLenSwts {
			maxLenSwts = lenSwt
		}
	})

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

		switch field.Type.Name() {
		case "string":
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + " <" + field.Type.Name() + ">]"

		case "bool":
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + "]"

		case "int", "int8", "int16", "int32", "int64":
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + " <" + field.Type.Name() + ">]"

		case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + " <" + field.Type.Name() + ">]"

		case "float32", "float64":
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + " <" + field.Type.Name() + ">]"

		default:
			argsStr += " [" + prefix + strings.Split(field.Tag.Get("switch"), ",")[0] + "]"
		}

		helpMenuStr += fmt.Sprintf("%v\r\n %-"+strconv.Itoa(maxLenSwts)+"v  %-9v", field.Name, swts, "<"+field.Type.Name()+">")
		if field.Tag.Get("opts") != "" {
			helpMenuStr += fmt.Sprintf(" (%v)", field.Tag.Get("opts"))
		}
		helpMenuStr += "\r\n"
		if field.Tag.Get("help") != "" {
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
		fmt.Print("\r\n" + helpMenuStr)
	}
}

// Parses os.Args into s. Panics if s is not of type struct or a public field doesn't contain a switch tag. Private struct fields are ignored.
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
//	opts:    Optional parameters [posistional|required|help].
//	default: Optional default value.
//	help:    Help message.
//
// Opts:
//
//	posistional: Field will be assigned using the first argument not tied to a switch.
//	required:    Field will be required, shows help menu if missing, if a default is given required is ignored.
//	help:        Quick opt for help menu implementation, when a switch from this field is present args.HelpMenu gets called followed by os.Exit(0)
//
// Example:
//
//	arg struct {
//		Help   bool    `switch:"h,-help"   opts:"help"`
//		FieldA float64 `switch:"a"         opts:"posistional"`
//		FieldB bool    `switch:"b"`
//		FieldC int     `switch:"c,-cc"     default:"5"`
//		FieldD uint    `switch:"d,-dd-ddd" default:"7" opts:"required"`
//		FieldE string  `switch:"e,-ee"     prefix:"/" help:"Some help for e."`
//	}
//
// Usage: exec -a=-1.1 -bc 10 --dd-ddd=2 /e "some message"
func Parse[T any](s T) T {
	args := os.Args[1:]

	allSwitches := [][2]string{}
	allPrefixes := []string{}
	forEachStructField(&s, func(field reflect.StructField, value reflect.Value) {
		prefix := field.Tag.Get("prefix")
		if prefix == "" {
			prefix = "-"
		}

		for _, swt := range strings.Split(field.Tag.Get("switch"), ",") {
			allSwitches = append(allSwitches, [2]string{prefix, swt})
			if !slices.Contains(allPrefixes, prefix) {
				allPrefixes = append(allPrefixes, prefix)
			}
		}
	})

	for i := 0; i < len(args); i++ {
		if strings.Contains(args[i], "=") {
			continue
		}
		if slices.ContainsFunc(allSwitches, func(item [2]string) bool { return item[0]+item[1] == args[i] }) {
			continue
		}

		prefix := ""
		for _, item := range allSwitches {
			if strings.HasPrefix(args[i], item[0]) {
				prefix = item[0]
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

	forEachStructField(&s, func(field reflect.StructField, value reflect.Value) {
		prefix := field.Tag.Get("prefix")
		if prefix == "" {
			prefix = "-"
		}

		index := -1
		for _, swt := range strings.Split(field.Tag.Get("switch"), ",") {
			index = slices.IndexFunc(args, func(a string) bool { return a == prefix+swt || strings.HasPrefix(a, prefix+swt+"=") })
			if index > -1 {
				if slices.Contains(strings.Split(field.Tag.Get("opts"), ","), "help") {
					HelpMenu(s, true)
					os.Exit(0)
				}
				break
			}
		}

		val := field.Tag.Get("default")
		if index > -1 {
			if strings.Contains(args[index], "=") {
				val = strings.Split(args[index], "=")[1]
			} else if index+1 < len(args) && !slices.ContainsFunc(allPrefixes, func(item string) bool { return strings.HasPrefix(args[index+1], item) }) {
				val = args[index+1]
				args = slices.Delete(args, index, index+1)
			} else if field.Type.Name() == "bool" {
				val = strconv.FormatBool(index != -1)
			}
			args = slices.Delete(args, index, index+1)
		}

		if val == "" && field.Type.Name() != "bool" {
			if slices.Contains(strings.Split(field.Tag.Get("opts"), ","), "required") {
				fmt.Println("Missing required argument " + field.Name)
				HelpMenu(s, false)
				os.Exit(0)
			} else if index > -1 {
				fmt.Println("Missing argument value for " + field.Name)
				HelpMenu(s, false)
				os.Exit(0)
			}
			return
		}

		switch field.Type.Name() {
		case "string":
			value.SetString(val)

		case "bool":
			value.SetBool(val == "true")

		case "int", "int8", "int16", "int32", "int64":
			v, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				fmt.Println("Invalid int format '" + val + "'")
				HelpMenu(s, false)
				os.Exit(0)
			}
			value.SetInt(v)

		case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
			v, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				fmt.Println("Invalid uint format '" + val + "'")
				HelpMenu(s, false)
				os.Exit(0)
			}
			value.SetUint(v)

		case "float32", "float64":
			v, err := strconv.ParseFloat(val, 64)
			if err != nil {
				fmt.Println("Invalid float format '" + val + "'")
				HelpMenu(s, false)
				os.Exit(0)
			}
			value.SetFloat(v)

		default:
			panic("unsuported type " + field.Type.Name())
		}

	})

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
