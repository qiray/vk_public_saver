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

func setCredentials(settings *AppSettings) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, _ := reader.ReadString('\n')

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	print("\n") //Add newline
	if err != nil {
		fmt.Println("Failed to get password")
		return
	}
	password := string(bytePassword)
	settings.userdata = make(map[string]string)
	settings.userdata["email"] = strings.TrimSpace(username)
	settings.userdata["pass"] = strings.TrimSpace(password)
}

func main() {
	userdataFlag := flag.Bool("userdata", false, "Use userdata.json for email and pass")
	flag.Parse()
	settings, err := loadSettings("settings.json")
	if err != nil {
		fmt.Println("Settings load failed. Closing...")
		os.Exit(1)
	}
	if *userdataFlag {
		userdata, _ := loadJSONFileMap("userdata.json")
		settings.userdata = userdata
	} else {
		setCredentials(&settings)
	}

	_ = login(&settings)
	print(settings.token, "\n")
	dbPath := "./data.db"
	db := initDataBase(dbPath)
	getPosts(db, settings, "-89009548")
	// dbExample()
}
