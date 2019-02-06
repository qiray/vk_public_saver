package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"golang.org/x/net/html"
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

//Based on https://github.com/Vorkytaka/easyvk-go/blob/master/easyvk/
func parseForm(body io.ReadCloser) (url.Values, string) {
	//Parse vk login form
	tokenizer := html.NewTokenizer(body)

	u := ""
	keys := []string{"_origin", "to", "ip_h", "lg_h"} //data for auth
	formData := map[string]string{}

	end := false
	for !end {
		tag := tokenizer.Next()

		switch tag {
		case html.ErrorToken:
			end = true
		case html.StartTagToken, html.SelfClosingTagToken:
			switch token := tokenizer.Token(); token.Data {
			case "form":
				for _, attr := range token.Attr {
					if attr.Key == "action" {
						u = attr.Val
					}
				}
			case "input":
				for _, key := range keys {
					if token.Attr[1].Val == key {
						formData[key] = token.Attr[2].Val
					}
				}
			}
		}
	}

	args := url.Values{}

	for key, val := range formData {
		args.Add(key, val)
	}

	return args, u
}

func login(settings *AppSettings) {
	path := "https://oauth.vk.com/authorize?client_id=" + settings.AppID + "&scope=" +
		settings.Settings + "&v=" + settings.APIVersion + "&redirect_uri=" + settings.RedirectURL +
		"&display=wap&response_type=token"

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	resp, err := client.Get(path)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	args, u := parseForm(resp.Body)

	args.Add("email", settings.userdata["email"])
	args.Add("pass", settings.userdata["pass"])

	resp, err = client.PostForm(u, args)
	if err != nil {
		return
	}

	if resp.Request.URL.Path != "/blank.html" {
		args, u := parseForm(resp.Body)
		resp, err := client.PostForm(u, args)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.Request.URL.Path != "/blank.html" {
			return
		}
	}

	urlArgs, err := url.ParseQuery(resp.Request.URL.Fragment)
	if err != nil {
		return
	}
	token := urlArgs["access_token"][0]
	settings.token = token
}

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

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	path := "https://api.vk.com/method/wall.get?owner_id=-89009548&v=" + settings.APIVersion + "&access_token=" + settings.token
	resp, err := client.Get(path)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseString := string(responseData)
	fmt.Println(responseString)

	// parseJSON()
	// dbExample()
}
