package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func getRequest(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(2)
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}

func login(settings AppSettings, userdata map[string]string) {
	path := "https://oauth.vk.com/authorize?client_id=" + settings.AppID + "&scope" +
		settings.Settings + "&v=" + settings.APIVersion + "&redirect_uri=" + settings.RedirectURL +
		"&display=" + settings.Display + "&response_type=token"
	// print(path, "\n")
	contents, err := getRequest(path)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(3)
	}
	contentsString := string(contents)
	iphRegex := regexp.MustCompile("ip_h.*value=\"(.*?)\"")
	lghRegex := regexp.MustCompile("lg_h.*value=\"(.*?)\"")
	res := iphRegex.FindStringSubmatch(contentsString)
	var iph, lgh string
	if res != nil {
		iph = res[1]
	}
	res = lghRegex.FindStringSubmatch(contentsString)
	if res != nil {
		lgh = res[1]
	}
	url := "https://login.vk.com/?act=login&_origin=http://m.vk.com&ip_h=" + iph + "&lg_h=" + lgh + "&role=pda&utf8=1"
	fmt.Printf("%s\n", url)
}

func main() {
	settings, err := loadSettings("settings.json")
	if err != nil {
		fmt.Println("Settings load failed. Closing...")
		os.Exit(1)
	}
	userdata, _ := loadJSONFileMap("userdata.json")
	fmt.Println(userdata)
	login(settings, userdata)

	// parseJSON()
	// dbExample()
}
