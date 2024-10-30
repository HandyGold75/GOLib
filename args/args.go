package main

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

type (
	// Field:  `switch:"[<switch>,...]" prefix:"-" opts:"[posistional,required]" default:"<value>" help:"<message>"`
	// Example: `switch:"s,-switch,-switch-me" opts:"required" default:"hello" help:"hello world"`
	// Usage:   -<switch> <value>
	//          --<switch>=<value>
	arg struct {
		FieldA string  `switch:"a"         help:"Some help for a."`
		FieldB bool    `switch:"b"         prefix:"/" opts:"required"`
		FieldC int     `switch:"c,-cc"     default:"5"`
		FieldD uint    `switch:"d,-dd-ddd" default:"7" opts:"required"`
		FieldE float64 `switch:"e,-ee"     opts:"posistional"`
		Help   bool    `switch:"h,-help"   opts:"help"`
	}
)

// Shows help menu related to s. Panics if s is not of type struct. Private struct fields are ignored.
func HelpMenu[T any](s T) {
	tOf := reflect.TypeOf(s)
	if tOf.Kind() != reflect.Struct {
		panic("s is not a stuct")
	}
	vOf := reflect.ValueOf(&s).Elem()
	if tOf.Kind() != reflect.Struct {
		panic("s is not a stuct")
	}

	maxLenSwts := 0
	for i := 0; i < tOf.NumField(); i++ {
		tOfField := tOf.Field(i)
		lenSwt := len(tOfField.Tag.Get("switch")) + max(1, len(tOfField.Tag.Get("prefix"))) + (len(strings.Split(tOfField.Tag.Get("switch"), ",")))
		if lenSwt > maxLenSwts {
			maxLenSwts = lenSwt
		}
	}

	for i := 0; i < tOf.NumField(); i++ {
		tOfField := tOf.Field(i)
		vOfField := vOf.Field(i)
		if !vOfField.CanSet() {
			continue
		}

		prefix := tOfField.Tag.Get("prefix")
		if prefix == "" {
			prefix = "-"
		}

		swts := ""
		for _, swt := range strings.Split(tOfField.Tag.Get("switch"), ",") {
			swts += " " + prefix + swt
		}

		fmt.Println(tOfField.Name)
		fmt.Printf(" %-"+strconv.Itoa(maxLenSwts)+"v  %-9v", swts, "<"+tOfField.Type.Name()+">")
		if tOfField.Tag.Get("opts") != "" {
			fmt.Printf(" (%v)", tOfField.Tag.Get("opts"))
		}
		fmt.Println()
		if tOfField.Tag.Get("help") != "" {
			fmt.Printf("	%v\r\n", tOfField.Tag.Get("help"))
		}
	}
}

// Parses os.Args into s. Panics if s is not of type struct. Private struct fields are ignored.
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
//		FieldA string  `switch:"a"         help:"Some help for a."`
//		FieldB bool    `switch:"b"         prefix:"/" opts:"required"`
//		FieldC int     `switch:"c,-cc"     default:"5"`
//		FieldD uint    `switch:"d,-dd-ddd" default:"7" opts:"required"`
//		FieldE float64 `switch:"e,-ee"     opts:"posistional"`
//		Help   bool    `switch:"h,-help"   opts:"help"`
//	}
//
// Usage: exec -a=somemessage /b --cc 10 --dd-ddd=2
func ParseArgs[T any](s T) T {
	args := os.Args[1:]
	tOf := reflect.TypeOf(s)
	if tOf.Kind() != reflect.Struct {
		panic("s is not a stuct")
	}
	vOf := reflect.ValueOf(&s).Elem()
	if tOf.Kind() != reflect.Struct {
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

		prefix := tOfField.Tag.Get("prefix")
		if prefix == "" {
			prefix = "-"
		}

		index := -1
		for _, swt := range strings.Split(tOfField.Tag.Get("switch"), ",") {
			index = slices.IndexFunc(args, func(a string) bool { return a == prefix+swt || strings.HasPrefix(a, prefix+swt+"=") })
			if index > -1 {
				if slices.Contains(strings.Split(tOfField.Tag.Get("opts"), ","), "help") {
					HelpMenu(s)
					os.Exit(0)
				}
				break
			}
		}

		value := tOfField.Tag.Get("default")
		if index > -1 {
			if strings.Contains(args[index], "=") {
				value = strings.Split(args[index], "=")[1]
			} else if index+1 < len(args) && !strings.HasPrefix(args[index+1], prefix) {
				value = args[index+1]
				args = slices.Delete(args, index, index+1)
			} else if tOfField.Type.Name() == "bool" {
				value = strconv.FormatBool(index != -1)
			}
			args = slices.Delete(args, index, index+1)
		}

		if value == "" && tOfField.Type.Name() != "bool" {
			if slices.Contains(strings.Split(tOfField.Tag.Get("opts"), ","), "required") {
				fmt.Println("Missing required argument " + tOfField.Name + "\r\n")
				HelpMenu(s)
				os.Exit(0)
			}
			continue
		}

		switch tOfField.Type.Name() {
		case "string":
			vOfField.SetString(value)

		case "bool":
			vOfField.SetBool(value == "true")

		case "int", "int8", "int16", "int32", "int64":
			val, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				fmt.Println("Invalid int format " + value + "\r\n")
				HelpMenu(s)
				os.Exit(0)
			}
			vOfField.SetInt(val)

		case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
			val, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				fmt.Println("Invalid uint format " + value + "\r\n")
				HelpMenu(s)
				os.Exit(0)
			}
			vOfField.SetUint(val)

		case "float32", "float64":
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				fmt.Println("Invalid float format " + value + "\r\n")
				HelpMenu(s)
				os.Exit(0)
			}
			vOfField.SetFloat(val)

		default:
			panic("unsuported type " + tOfField.Type.Name())
		}
	}

	return s
}

func main() {
	arg := ParseArgs(arg{})
	fmt.Println(arg)
}
