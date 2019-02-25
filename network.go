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
	"strconv"
	"time"

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

func login(settings *AppSettings) bool {

	path := "https://oauth.vk.com/authorize?client_id=" + settings.AppID + "&scope=" +
		settings.Settings + "&v=" + settings.APIVersion + "&redirect_uri=" + settings.RedirectURL +
		"&display=wap&response_type=token"
	cookies := make(map[string][]*http.Cookie)

	jar, _ := cookiejar.New(nil)
	// loadCookies("cookies.txt", jar) //TODO: use cookies (maybe error appears because of expires = 0)
	settings.client = &http.Client{
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

func updateCookies(cookies map[string][]*http.Cookie, jar *cookiejar.Jar, resp *http.Response) {
	url := resp.Request.URL
	cookies[url.Scheme+"://"+url.Hostname()] = jar.Cookies(resp.Request.URL)
}

func saveCookies(filepath string, cookies map[string][]*http.Cookie) {
	data, _ := json.Marshal(cookies)
	ioutil.WriteFile(filepath, data, 0644)
}

func loadCookies(filepath string, jar *cookiejar.Jar) {
	data, err := readFile(filepath)
	if err != nil {
		return
	}
	cookies := make(map[string][]*http.Cookie)
	err = json.Unmarshal([]byte(data), &cookies)
	if err != nil {
		return
	}
	for key, val := range cookies {
		url, err := url.Parse(key)
		if err != nil {
			print(111)
			continue
		}
		jar.SetCookies(url, val)
	}
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

func getPosts(settings AppSettings, publicID string) {
	count := 50
	offset := 0
	pertime := 20
	totalNumber := count * pertime
	numberOfPosts := 2147483647
	finished := false

	for !finished {
		print("Saving posts, offset: ", offset, "\n")
		code := `
		var result = [];
		var i = 0;
		var max_posts = ` + strconv.Itoa(numberOfPosts) + `;
		while (i < ` + strconv.Itoa(pertime) + ` && max_posts > 0) {
			result.push(API.wall.get({
				owner_id: ` + publicID + `,
				count: ` + strconv.Itoa(count) + `,
				offset: ` + strconv.Itoa(count) + `*i+` + strconv.Itoa(offset) + `,
				filter: "all",
				v : ` + settings.APIVersion + `,
				}));
			max_posts = max_posts - ` + strconv.Itoa(count) + `;
			i = i+1;
		};
		return result;`
		numberOfPosts -= totalNumber
		code = url.QueryEscape(code)
		path := "https://api.vk.com/method/execute?code=" + code + "&v=" + settings.APIVersion + "&access_token=" + settings.token
		resp, err := settings.client.Get(path)
		if err != nil {
			fmt.Println(err, 1)
			return
		}
		defer resp.Body.Close()
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err, 2)
			return
		}
		responseString := string(responseData)
		var p PostsResponse
		err = json.Unmarshal([]byte(responseString), &p)
		if err != nil {
			fmt.Println(err, 3)
		}
		for _, val := range p.Response { //TODO: save data in database
			if len(val.Items) == 0 {
				finished = true
				break
			}
		}
		offset += totalNumber
		time.Sleep(250 * time.Millisecond)
	}
	print("Done\n")
}
