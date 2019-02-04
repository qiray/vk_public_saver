package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"
)

//TODO: read https://github.com/Vorkytaka/easyvk-go

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
	print(path, "\n")
	contents, err := getRequest("http://m.vk.com")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(3)
	}
	contentsString := string(contents)
	loginRegex := regexp.MustCompile("action=\"(.*?)\"")
	res := loginRegex.FindStringSubmatch(contentsString)
	var loginURL string
	if res != nil {
		loginURL = res[1]
	}
	print(loginURL, "\n")
	// loginURL := "https://login.vk.com/?act=login&_origin=http://m.vk.com&ip_h=" + iph + "&lg_h=" + lgh + "&role=pda&utf8=1"
	// https://login.vk.com/?act=login&_origin=https://m.vk.com&ip_h=c6d1d61811616206c5&lg_h=c248a9dbf718faf449&role=pda&utf8=1"
	// response, err := http.PostForm(loginURL, url.Values{"email": {userdata["email"]}, "pass": {userdata["pass"]}})
	// if err != nil {
	// 	fmt.Printf("%s", err)
	// 	os.Exit(4)
	// }
	// defer response.Body.Close()
	// contents, _ = ioutil.ReadAll(response.Body)
	// fmt.Printf("%s\n", contents)

	params := url.Values{}
	params.Set("email", userdata["email"])
	params.Set("pass", userdata["pass"])
	postData := strings.NewReader(params.Encode())

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, CheckRedirect: nil}
	request, err := http.NewRequest("POST", loginURL, postData)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:26.0) Gecko/20100101 Firefox/26.0")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(4)
	}
	defer response.Body.Close()
	contents, _ = ioutil.ReadAll(response.Body)

	fmt.Printf("%d\n", response.StatusCode)
	fmt.Printf("%s\n", contents)

	// data := []byte(`{"foo":"bar"}`)
	// r := bytes.NewReader(data)
	// resp, err := http.Post("http://example.com/upload", "application/json", r)
	// if err != nil {
	// 	return
	// }
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
