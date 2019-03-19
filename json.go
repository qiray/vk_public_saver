package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func readFile(filepath string) ([]byte, error) {
	dataFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer dataFile.Close() //close file later
	data, err := ioutil.ReadAll(dataFile)
	return data, err
}

func loadSettings(settingsFile string) (AppSettings, error) {
	data, err := readFile(settingsFile)
	var settings = AppSettings{}
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(data, &settings)
	return settings, err
}

func loadJSONFileMap(filepath string) (map[string]string, error) {
	result := make(map[string]string)
	data, err := readFile(filepath)
	if err != nil {
		return result, err
	}
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(data, &objmap)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	for key, value := range objmap {
		var stringValue string
		err = json.Unmarshal(*value, &stringValue)
		if err != nil {
			fmt.Println(err)
			return result, err
		}
		result[key] = stringValue
	}
	return result, err
}
