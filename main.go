package main

import (
	"fmt"
	"os"
)

//TODO: use one http.clent (create APP object with settings, userdata etc) and cookies

func main() {
	settings, err := loadSettings("settings.json")
	if err != nil {
		fmt.Println("Settings load failed. Closing...")
		os.Exit(1)
	}
	userdata, _ := loadJSONFileMap("userdata.json")
	settings.userdata = userdata
	login(&settings)
	print(settings.token, "\n")
	wallGet(settings)

	// parseJSON()
	// dbExample()
}
