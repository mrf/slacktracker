package main

import (
	"fmt"
	//"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	config, err := ioutil.ReadFile("config.yml")
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			fmt.Println("Please create a config file with your Slack api key")
		}
		return
	}
	token := strings.TrimPrefix(string(config), "token: ")
	token = strings.TrimSpace(token)
	api := slack.New(token)
	users, err := api.GetUsers()
	if err != nil {
		fmt.Println("errr: ", err)
		return
	}
	for _, user := range users {
		// Ignore bots
		if !user.IsBot {
			// Check if presence empty (slackbot is special case)
			if user.Presence != "" {
				fmt.Println(user.Name)
				fmt.Println(user.Presence)
			}
		}
	}
}
