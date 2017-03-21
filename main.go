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
	//api.SetDebug(true)
	users, err := api.GetUsers()
	if err != nil {
		fmt.Println("Slack api error: ", err)
		return
	}
	db, err := bolt.Open("slacktracker.db", 0600, nil)
	if err != nil {
		fmt.Println("Bolt DB not working, maybe forgot to run go get? Err: ", err)
	}
	defer db.Close()

	// TODO actual meat move?
	fmt.Printf("date, ")
	for _, user := range users {
		if user.Presence != "" {
			// Create a table for each user we find
			// Do we need to react to users we haven't seen before differently?
			// Probably only want to run this sporadically
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
	// DOes it work with bolt more importantly?

	starttime := time.Now()
	var runlength time.Duration = 5000 * time.Second
	// Note, any value less than 15 seconds for frequency may get time out from slack API
	var frequency time.Duration = 15 * time.Second
	for time.Since(starttime) < runlength {
		time.Sleep(frequency)
		fmt.Println("")
		fmt.Printf(time.Now().Format(time.RFC3339))
		fmt.Printf(",")
		users, err := api.GetUsers()
		if err != nil {
			spew.Dump(users)
			spew.Dump(err)
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
