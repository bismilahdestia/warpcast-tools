package main

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/tokinaa/warpcast-tools/degen"
	"github.com/tokinaa/warpcast-tools/warpcast"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type ConfigStruct struct {
	Accounts  []string `json:"accounts"`
	DelayLike int      `json:"delayLike"`
}

var (
	myConfig = LoadConfig()
)

func LoadConfig() ConfigStruct {
	// Load from config.json
	openFile, err := os.Open("config.json")
	if err != nil {
		return ConfigStruct{}
	}

	defer openFile.Close()

	var config ConfigStruct
	jsonParser := json.NewDecoder(openFile)
	jsonParser.Decode(&config)

	return config
}

func init() {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		file, _ := json.MarshalIndent(ConfigStruct{
			Accounts:  []string{},
			DelayLike: 1000,
		}, "", " ")
		_ = os.WriteFile("config.json", file, 0644)
	}
	openLoadConfig := LoadConfig()
	myConfig = openLoadConfig
}

func checkingError(err error) {
	if err != nil {
		if err.Error() == "interrupt" {
			os.Exit(0)
		} else {
			fmt.Printf("ERROR: %s\n", err)
		}
	}
}

func showPressEnter() {
	fmt.Print("Press Enter to Back...")
	var input string
	fmt.Scanln(&input)
}

func multiAccountsManagement() {
	fmt.Print("\033[H\033[2J")
	fmt.Println()
	fmt.Println("1. Add Account")
	fmt.Println("2. Back")
	fmt.Println()

	inputMenu := ""
	inputMenuError := survey.AskOne(&survey.Input{
		Message: "Select Menu:",
	}, &inputMenu, survey.WithValidator(survey.Required))

	checkingError(inputMenuError)

	switch inputMenu {
	case "1":
		fmt.Print("\033[H\033[2J")
		fmt.Println("Add Account")
		fmt.Println()

		inputAccount := ""
		inputAccountError := survey.AskOne(&survey.Input{
			Message: "Authorization Token:",
		}, &inputAccount, survey.WithValidator(survey.Required))

		checkingError(inputAccountError)

		myConfig.Accounts = append(myConfig.Accounts, inputAccount)

		file, _ := json.MarshalIndent(myConfig, "", " ")
		_ = os.WriteFile("config.json", file, 0644)

		fmt.Println("Account Added")
		fmt.Println()

		showPressEnter()

		fmt.Print("\033[H\033[2J")
		multiAccountsManagement()
	case "2":
		fmt.Print("\033[H\033[2J")
		main()
	}
}

func targetLike() {
	fmt.Print("\033[H\033[2J")

	fmt.Println("Target Like")
	fmt.Println()

	inputSelectAccount := 0
	inputSelectAccountError := survey.AskOne(&survey.Select{
		Message: "Select Account:",
		Options: myConfig.Accounts,
	}, &inputSelectAccount, survey.WithValidator(survey.Required))

	checkingError(inputSelectAccountError)

	targetProfiles := ""
	targetProfilesError := survey.AskOne(&survey.Input{
		Message: "Enter target profiles (comma separated):",
	}, &targetProfiles, survey.WithValidator(survey.Required))

	checkingError(targetProfilesError)

	profiles := strings.Split(targetProfiles, ",")

	for _, profile := range profiles {
		profile = strings.TrimSpace(profile)

		fmt.Printf("[TARGET] Processing profile: %s\n", profile)

		casts, err := warpcast.GetUserCasts(myConfig.Accounts[inputSelectAccount], profile)
		if err != nil {
			fmt.Printf("[GET CASTS] ERROR : %s\n", err)
			continue
		}

		for _, cast := range casts.Result.Items {
			fmt.Printf("[CAST] [https://warpcast.com/%s/%s] ", cast.Author.Username, cast.Hash)

			if cast.ViewerContext.Reacted {
				fmt.Printf("[LIKE] ALREADY\n")
			} else {
				_, err := warpcast.Like(myConfig.Accounts[inputSelectAccount], cast.Hash)
				if err != nil {
					fmt.Printf("[LIKE] ERROR : %s\n", err)
				} else {
					fmt.Printf("[LIKE] SUCCESS\n")
					delayLike := time.Duration(myConfig.DelayLike) * time.Millisecond
					time.Sleep(delayLike)
				}
			}
		}
	}

	showPressEnter()
	main()
}

func main() {
	fmt.Println("Warpcast Tools")
	fmt.Println("Author : @x0xdead / Wildaann")
	fmt.Println()
	fmt.Println("1. Multi Accounts Management")
	fmt.Println("2. Target Like")
	fmt.Println()

	inputMenu := ""
	inputMenuError := survey.AskOne(&survey.Input{
		Message: "Select Menu:",
	}, &inputMenu, survey.WithValidator(survey.Required))

	checkingError(inputMenuError)

	switch inputMenu {
	case "1":
		multiAccountsManagement()
	case "2":
		targetLike()
	}
}
