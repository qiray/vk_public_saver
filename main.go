package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getRequest() {
	response, err := http.Get("http://golang.org/")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
	}
}

func login(settings AppSettings) {
	path := "https://oauth.vk.com/authorize?client_id=" + settings.AppID + "&scope" +
		settings.Settings + "&v=" + settings.APIVersion + "&redirect_uri=" + settings.RedirectURL +
		"&display=" + settings.Display + "&response_type=token"
	print(path)
}

func main() {
	settings, err := loadSettings("settings.json")
	if err != nil {
		fmt.Println("Settings load failed. Closing...")
		os.Exit(1)
	}
	login(settings)
	// getRequest()
	// parseJSON()
	// dbExample()
}
