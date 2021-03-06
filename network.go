// vk_public_saver - tool for saving walls data from vk.com
// Copyright (c) 2019, Yaroslav Zotov, https://github.com/qiray/
// All rights reserved.

// This file is part of vk_public_saver.

// vk_public_saver is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// vk_public_saver is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with vk_public_saver.  If not, see <https://www.gnu.org/licenses/>.

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

	"github.com/vorkytaka/easyvk-go/easyvk"
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

func login(settings *AppSettings) bool {
	vk, err := easyvk.WithAuth(settings.userdata["email"], settings.userdata["pass"], settings.AppID, "wall")
	if err != nil {
		fmt.Println(err)
		return false
	}
	settings.token = vk.AccessToken
	return true
}

func getPosts(db *sql.DB, settings AppSettings) {
	count := 50
	offset := 0
	pertime := 20
	totalNumber := count * pertime
	numberOfPosts := 2147483647
	finished := false
	publicID := settings.userdata["source"]
	client := &http.Client{}

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
		resp, err := client.Get(path)
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
