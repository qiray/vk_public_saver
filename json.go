package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func parseJSON() {
	// TODO: improve parsing different types?
	type mmn struct {
		Rows []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			URL  string `json:"url"`
			Src  string `json:"src"`
		} `json:"rows"`
	}

	body1 := []byte(`{
		"rows": [
			{
				"id": "01ae6145-90a3-11e7-7a69-8f55000cda4b",
				"type": "url",
				"url": "google.ru",
				"modificationsCount": 0
			},
			{
				"id": "11",
				"type": "image",
				"src": "image.png",
				"quantity": 1
			}
		]
	}`)
	var app = mmn{}
	err1 := json.Unmarshal(body1, &app)
	if err1 != nil {
		log.Fatal("error")
	}

	for _, row := range app.Rows {
		info := row.ID + " " + row.Type
		if row.Type == "url" {
			info += " " + row.URL
		} else if row.Type == "image" {
			info += " " + row.Src
		}
		fmt.Println(info)
	}
}

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
