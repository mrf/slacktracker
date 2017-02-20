package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	// Setup stuff TODO move to setup function
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
	db, err := bolt.Open("slacktracker.db", 0600, nil)
	if err != nil {
		fmt.Println("Db not werking")
	}
	defer db.Close()

	// TODO actual meat
	fmt.Printf("date, ")
	for _, user := range users {
		if user.Presence != "" {
			db.Update(func(tx *bolt.Tx) error {
				userBucket, err := tx.CreateBucket([]byte(user.Name))
				if err != nil {
					return fmt.Errorf("create bucket: %s", err)
				}
				spew.Dump(userBucket)
				return nil
			})
			fmt.Printf(user.Name)
			fmt.Printf(", ")
		}
	}

	// Idea being to open 1 channel per user (is this sustainable for large numbers?)

	// TODO fix this to not run forever
	t := time.Now()
	for t.Month() < 100 {
		t := time.Now()
		if time.After(25*time.Second) != nil {
			// run this every x seconds
			fmt.Println()
			// run this and keep it running
			fmt.Printf(t.Format(time.RFC3339))
			fmt.Printf(",")
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
						fmt.Printf(user.Presence)
						// Todo don't print if this is last user
						fmt.Printf(", ")
					}
				}
			}
		}
	}
}
