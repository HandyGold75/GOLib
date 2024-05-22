package cfg

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

var (
	// Filename for the config file to Load from and Dump to.
	FileName string = "config.json"
)

// Loads config.json to Data.
//
// config.json gets created from Data if it does not exists.
//
// Missing keys in config.json will be filled in from the Data.
// Additional keys in config.json will be ignored.
// Values in config.json take priority over values in Data
func Load(Data any) (err error) {
	filePath, err := os.Executable()
	if err != nil {
		return err
	}

	filePathSplit := strings.Split(filePath, "/")
	filePath = strings.Join(filePathSplit[:len(filePathSplit)-1], "/") + "/" + FileName

	bytes, err := readFile(filePath)
	if err != nil {
		return err
	}

	if len(bytes) <= 0 {
		Dump(Data)

		bytes, err = readFile(filePath)
		if err != nil {
			return err
		}
	}

	err = json.Unmarshal(bytes, Data)
	if err != nil {
		return err
	}

	return nil
}

// Dumps Data to config.json.
//
// config.json gets created from Data if it does not exists.
//
// Missing keys in config.json will be filled in from the Data.
// Additional keys in config.json will be deleted.
// Values in Data take priority over values in config.json
func Dump(Data any) (err error) {
	filePath, err := os.Executable()
	if err != nil {
		return err
	}

	filePathSplit := strings.Split(filePath, "/")
	filePath = strings.Join(filePathSplit[:len(filePathSplit)-1], "/") + "/" + FileName

	configJson, err := json.MarshalIndent(Data, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, configJson, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		file.Close()

		file, err = os.Create(fileName)
		if err != nil {
			file.Close()
			return []byte{}, err
		}

		_, err = file.WriteString("{}")
		if err != nil {
			file.Close()
			return []byte{}, err
		}
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
