package cfg

import (
	"encoding/json"
	"os"
	"strings"
)

// Loads config file into data.
//
// If config file is not present then it tries creating it with content of data.
//
// Config file is stored in `./golib/<name>.json` relative to `os.UserConfigDir`.
func Load(name string, data any) error {
	file, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	file += "/golib/" + name + ".json"

	bytes, err := os.ReadFile(file)
	if os.IsNotExist(err) || len(bytes) == 0 {
		fileSplit := strings.Split(strings.ReplaceAll(file, "\\", "/"), "/")
		err := os.MkdirAll(strings.Join(fileSplit[:len(fileSplit)-1], "/"), os.ModePerm)
		if err != nil {
			return err
		}
		bytes, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			return err
		}
		err = os.WriteFile(file, bytes, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		return err
	}
	return nil
}

// Dumps data to config file.
//
// If config file is not present then it tries creating it with content of data.
//
// Config file is stored in `./golib/<name>.json` relative to `os.UserConfigDir`.
func Dump(name string, data any) error {
	file, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	file += "/golib/" + name + ".json"

	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(file, bytes, os.ModePerm)
	if os.IsNotExist(err) {
		fileSplit := strings.Split(strings.ReplaceAll(file, "\\", "/"), "/")
		err := os.MkdirAll(strings.Join(fileSplit[:len(fileSplit)-1], "/"), os.ModePerm)
		if err != nil {
			return err
		}
		err = os.WriteFile(file, bytes, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return err
}
