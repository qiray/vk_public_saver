package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"time"
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

func getPosts(db *sql.DB, settings AppSettings) {
	count := 50
	offset := 0
	pertime := 20
	totalNumber := count * pertime
	numberOfPosts := 2147483647
	finished := false
	publicID := settings.userdata["source"]

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
		savePostsResponse(db, p)
		responseCount := 0
		for _, val := range p.Response {
			responseCount++
			if len(val.Items) == 0 {
				finished = true
				break
			}
		}
		if responseCount == 0 {
			return //Return if no response
		}
		offset += totalNumber
		time.Sleep(250 * time.Millisecond)
	}
	print("Done\n")
}
