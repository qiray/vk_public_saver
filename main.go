package main

import (
	"fmt"
	"os"
)

func main() {
	settings, err := loadSettings("settings.json")
	if err != nil {
		fmt.Println("Settings load failed. Closing...")
		os.Exit(1)
	}
	userdata, _ := loadJSONFileMap("userdata.json")
	settings.userdata = userdata //TODO: maybe read from input
	login(&settings)
	print(settings.token, "\n")
	wallGet(settings)

	// dbExample()
}
