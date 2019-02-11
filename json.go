package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func loadSettings(settingsFile string) (AppSettings, error) {
	jsonFile, err := os.Open(settingsFile)
	var settings = AppSettings{}
	if err != nil {
		fmt.Println(err)
		return settings, err
	}
	defer jsonFile.Close() //close file later
	data, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(data, &settings)
	return settings, err
}

func loadJSONFileMap(filepath string) (map[string]string, error) {
	result := make(map[string]string)
	jsonFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return result, err
	}
	defer jsonFile.Close()
	data, _ := ioutil.ReadAll(jsonFile)
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
