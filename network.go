package main

import (
	"encoding/json"
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
	//TODO: return bool and use cookies
	path := "https://oauth.vk.com/authorize?client_id=" + settings.AppID + "&scope=" +
		settings.Settings + "&v=" + settings.APIVersion + "&redirect_uri=" + settings.RedirectURL +
		"&display=wap&response_type=token"
	cookies := make(map[string][]*http.Cookie)

	jar, _ := cookiejar.New(nil)
	settings.client = &http.Client{
		Jar: jar,
	}
	resp, err := settings.client.Get(path)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	args, u := parseForm(resp.Body)

	args.Add("email", settings.userdata["email"])
	args.Add("pass", settings.userdata["pass"])

	resp, err = settings.client.PostForm(u, args)
	if err != nil {
		return
	}

	if resp.Request.URL.Path != "/blank.html" {
		args, u := parseForm(resp.Body)
		resp, err := settings.client.PostForm(u, args)
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
	updateCookies(cookies, jar, resp)
	saveCookies("cookies.txt", cookies)

}

func updateCookies(cookies map[string][]*http.Cookie, jar *cookiejar.Jar, resp *http.Response) {
	cookies[resp.Request.URL.String()] = jar.Cookies(resp.Request.URL)
}

func saveCookies(filepath string, cookies map[string][]*http.Cookie) {
	data, _ := json.Marshal(cookies)
	ioutil.WriteFile(filepath, data, 0644)
}

func loadCookies(filepath string, jar *cookiejar.Jar, resp *http.Response) {
	// result, err := loadJSONFileMap(filepath)
	// if err != nil {

	// }
}

func wallGet(settings AppSettings) {
	path := "https://api.vk.com/method/wall.get?owner_id=-89009548&v=" + settings.APIVersion + "&access_token=" + settings.token
	resp, err := settings.client.Get(path)
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
}
