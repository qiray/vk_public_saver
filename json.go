// vk_public_saver - tool for saving walls data from vk.com
// Copyright (c) 2019, Yaroslav Zotov, https://github.com/qiray/
// All rights reserved.

// This file is part of vk_public_saver.

// vk_public_saver is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// vk_public_saver is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with vk_public_saver.  If not, see <https://www.gnu.org/licenses/>.

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
