package main

import (
	"fmt"
	//	"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
	"io/ioutil"
	"os"
	"strings"
	"time"
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
	fmt.Printf("date, ")
	for _, user := range users {
		if !user.IsBot {
			fmt.Printf(user.Name)
			fmt.Printf(", ")
		}
	}

	// TODO fix this to not run forever
	t := time.Now()
	for t.Month() < 100 {
		t := time.Now()
		time.Sleep(100 * time.Millisecond)
		// run this every x seconds
		fmt.Println()
		fmt.Printf("%d\n", t.Second())
		// run this and keep it running
		for _, user := range users {
			// Ignore bots

			if !user.IsBot {
				// Check if presence empty (slackbot is special case)
				if user.Presence != "" {
					fmt.Printf(user.Presence)
					// Todo don't print if this is last user
					fmt.Printf(", ")
				}
			}
		}
	}
}
