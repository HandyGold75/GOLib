package jsonp

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

type jsonParser struct {
	// Filename for the config file to Load from and Dump to.
	FilePath string
}

// Loads jsonParser.FilePath to Data.
//
// jsonParser.FilePath gets created from Data if it does not exists.
//
// Missing keys in jsonParser.FilePath will be filled in from the Data.
// Additional keys in jsonParser.FilePath will be ignored.
// Values in jsonParser.FilePath take priority over values in Data
func (jsonP jsonParser) Load(Data any) (err error) {
	bytes, err := readFile(jsonP.FilePath)
	if err != nil {
		return err
	}

	if len(bytes) <= 0 {
		jsonP.Dump(Data)

		bytes, err = readFile(jsonP.FilePath)
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

// Dumps Data to jsonParser.FilePath.
//
// jsonParser.FilePath gets created from Data if it does not exists.
//
// Missing keys in jsonParser.FilePath will be filled in from the Data.
// Additional keys in jsonParser.FilePath will be deleted.
// Values in Data take priority over values in jsonParser.FilePath
func (jsonP jsonParser) Dump(Data any) (err error) {
	configJson, err := json.MarshalIndent(Data, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(jsonP.FilePath, configJson, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func readFile(FilePath string) ([]byte, error) {
	file, err := os.Open(FilePath)
	if err != nil {
		file.Close()

		file, err = os.Create(FilePath)
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

func NewParser(filePath string) jsonParser {
	return jsonParser{FilePath: filePath}
}

func ExecPath() (string, error) {
	filePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	filePathSplit := strings.Split(filePath, "/")
	return strings.Join(filePathSplit[:len(filePathSplit)-1], "/"), nil
}
