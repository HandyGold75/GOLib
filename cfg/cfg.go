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
	return LoadAbs(file+"/golib/"+name+".json", data)
}

// Loads config file into data.
//
// If config file is not present then it tries creating it with content of data.
//
// Config file is stored in `./<name>.json` relative to `os.Executable`.
func LoadRel(name string, data any) error {
	file, err := os.Executable()
	if err != nil {
		return err
	}
	fileSplit := strings.Split(strings.ReplaceAll(file, "\\", "/"), "/")
	return LoadAbs(strings.Join(fileSplit[:len(fileSplit)-1], "/")+"/"+name+".json", data)
}

// Loads config file into data.
//
// If config file is not present then it tries creating it with content of data.
func LoadAbs(file string, data any) error {
	bytes, err := os.ReadFile(file)
	if os.IsNotExist(err) || len(bytes) == 0 {
		fileSplit := strings.Split(strings.ReplaceAll(file, "\\", "/"), "/")
		err := os.MkdirAll(strings.Join(fileSplit[:len(fileSplit)-1], "/"), os.ModePerm)
		if err != nil {
			return err
		}
		bytes, err = json.MarshalIndent(data, "", "\t")
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
	return DumpAbs(file+"/golib/"+name+".json", data)
}

// Dumps data to config file.
//
// If config file is not present then it tries creating it with content of data.
//
// Config file is stored in `./<name>.json` relative to `os.Executable`.
func DumpRel(name string, data any) error {
	file, err := os.Executable()
	if err != nil {
		return err
	}
	fileSplit := strings.Split(strings.ReplaceAll(file, "\\", "/"), "/")
	return DumpAbs(strings.Join(fileSplit[:len(fileSplit)-1], "/")+"/"+name+".json", data)
}

// Dumps data to config file.
//
// If config file is not present then it tries creating it with content of data.
func DumpAbs(file string, data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(file, bytes, os.ModePerm); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
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
	return nil
}

// Returns path if file exists, else returns empty string.
//
// File is checked in `./golib/<name>.json` relative to `os.UserConfigDir`.
func Check(name string) string {
	file, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return CheckAbs(file + "/golib/" + name + ".json")
}

// Returns path if file exists, else returns empty string.
//
// File is checked in `./<name>.json` relative to `os.Executable`.
func CheckRel(name string) string {
	file, err := os.Executable()
	if err != nil {
		return ""
	}
	fileSplit := strings.Split(strings.ReplaceAll(file, "\\", "/"), "/")
	return CheckAbs(strings.Join(fileSplit[:len(fileSplit)-1], "/") + "/" + name + ".json")
}

// Returns path if file exists, else returns empty string.
func CheckAbs(file string) string {
	f, err := os.Stat(file)
	if err != nil || f.IsDir() {
		return ""
	}
	return file
}

// Returns path if dir exists, else returns empty string.
//
// Dir is checked in `./golib/<name>` relative to `os.UserConfigDir`.
func CheckDir(name string) string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return CheckDirAbs(dir + "/golib/" + name)
}

// Returns path if dir exists, else returns empty string.
//
// Dir is checked in `./<name>` relative to `os.Executable`.
func CheckDirRel(name string) string {
	dir, err := os.Executable()
	if err != nil {
		return ""
	}
	fileSplit := strings.Split(strings.ReplaceAll(dir, "\\", "/"), "/")
	return CheckDirAbs(strings.Join(fileSplit[:len(fileSplit)-1], "/") + "/" + name)
}

// Returns path if dir exists, else returns empty string.
func CheckDirAbs(path string) string {
	f, err := os.Stat(path)
	if err != nil || !f.IsDir() {
		return ""
	}
	return path
}
