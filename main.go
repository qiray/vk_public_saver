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
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var version = Version{"vk_public_saver", 1, 0, 0}

func getVersion() string {
	return fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Build)
}

func getAboutInfo() string {
	return "\n" + version.Name + " " + getVersion() + " Copyright (C) 2019 Yaroslav Zotov.\n" +
		"This program comes with ABSOLUTELY NO WARRANTY.\n" +
		"This is free software under GNU GPL3; see the source for copying conditions\n"
}

func setCredentials(settings *AppSettings, tokenFlag bool) {
	reader := bufio.NewReader(os.Stdin)
	username, password := "", ""

	if tokenFlag {
		fmt.Print("Enter token: ")
		settings.token, _ = reader.ReadString('\n')
	} else {
		fmt.Print("Enter Username: ")
		username, _ = reader.ReadString('\n')

		fmt.Print("Enter Password: ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		print("\n") //Add newline
		if err != nil {
			fmt.Println("Failed to get password")
			return
		}
		password = string(bytePassword)
	}

	fmt.Print("Enter public or user id to save it's wall data: ")
	source, _ := reader.ReadString('\n')
	settings.userdata = make(map[string]string)
	settings.userdata["email"] = strings.TrimSpace(username)
	settings.userdata["pass"] = strings.TrimSpace(password)
	settings.userdata["source"] = strings.TrimSpace(source)
}

func main() {
	userdataFlag := flag.Bool("userdata", false, "Use userdata.json for email and password")
	tokenFlag := flag.Bool("token", false, "Use token instead of email and password")
	aboutFlag := flag.Bool("about", false, "Show about info")
	flag.Parse()
	if *aboutFlag {
		fmt.Println(getAboutInfo())
		return
	}
	settings, err := loadSettings("settings.json")
	if err != nil {
		fmt.Println("Settings load failed. Closing...")
		os.Exit(1)
	}
	if *userdataFlag {
		userdata, err := loadJSONFileMap("userdata.json")
		if err != nil {
			fmt.Println("Userdata load failed. Closing...")
			os.Exit(1)
		}
		settings.userdata = userdata
	} else {
		setCredentials(&settings, *tokenFlag)
	}

	_ = login(&settings)
	print(settings.token, "\n")
	dbPath := fmt.Sprintf("./data_%s.db", settings.userdata["source"])
	db := initDataBase(dbPath)
	getPosts(db, settings)
	closeDatabase(db)
}
