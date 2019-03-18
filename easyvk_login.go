package main

//TODO: replace it with call of easyvk library

//This file contains code with Apache license from https://github.com/Vorkytaka/easyvk-go/blob/master/easyvk/

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"golang.org/x/net/html"
)

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

func login(settings *AppSettings) bool {

	path := "https://oauth.vk.com/authorize?client_id=" + settings.AppID + "&scope=" +
		settings.Settings + "&v=" + settings.APIVersion + "&redirect_uri=" + settings.RedirectURL +
		"&display=wap&response_type=token"
	cookies := make(map[string][]*http.Cookie)

	jar, _ := cookiejar.New(nil)
	// loadCookies("cookies.txt", jar) //TODO: use cookies (maybe error appears because of expires = 0)
	settings.client = &http.Client{ //TODO: replace client with call of easyVK
		Jar: jar,
	}
	resp, err := settings.client.Get(path)
	if err != nil {
		fmt.Println(err, 1)
		return false
	}
	defer resp.Body.Close()
	args, u := parseForm(resp.Body)

	args.Add("email", settings.userdata["email"])
	args.Add("pass", settings.userdata["pass"])

	resp, err = settings.client.PostForm(u, args)
	if err != nil {
		fmt.Println("Failed to login.")
		return false
	}

	if resp.Request.URL.Path != "/blank.html" {
		args, u := parseForm(resp.Body)
		resp, err := settings.client.PostForm(u, args)
		if err != nil {
			fmt.Println(err, 3)
			return false
		}
		defer resp.Body.Close()

		if resp.Request.URL.Path != "/blank.html" {
			fmt.Println(resp.Request.URL, 4)
			return false
		}
	}

	urlArgs, err := url.ParseQuery(resp.Request.URL.Fragment)
	if err != nil {
		fmt.Println(err, 5)
		return false
	}
	token := urlArgs["access_token"][0]
	settings.token = token
	updateCookies(cookies, jar, resp)
	saveCookies("cookies.txt", cookies)
	return true
}
