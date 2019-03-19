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

var version = Version{1, 0, 0}

func getVersion() string {
	return fmt.Sprintf("%d,%d,%d", version.Major, version.Minor, version.Build)
}

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
	fmt.Print("Enter public or user id: ")
	source, _ := reader.ReadString('\n')
	settings.userdata = make(map[string]string)
	settings.userdata["email"] = strings.TrimSpace(username)
	settings.userdata["pass"] = strings.TrimSpace(password)
	settings.userdata["source"] = strings.TrimSpace(source)
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
		userdata, err := loadJSONFileMap("userdata.json")
		if err != nil {
			fmt.Println("Userdata load failed. Closing...")
			os.Exit(1)
		}
		settings.userdata = userdata
	} else {
		setCredentials(&settings)
	}

	_ = login(&settings)
	print(settings.token, "\n")
	dbPath := fmt.Sprintf("./data_%s.db", settings.userdata["source"])
	db := initDataBase(dbPath)
	getPosts(db, settings)
	closeDatabase(db)
}
