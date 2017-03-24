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

func setup() (*slack.Client, []slack.User) {

	config, err := ioutil.ReadFile("config.yml")
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			fmt.Println("Please create a config file with your Slack api key")
		}
	}
	token := strings.TrimPrefix(string(config), "token: ")
	token = strings.TrimSpace(token)
	apiClient := slack.New(token)
	//api.SetDebug(true)
	users, err := apiClient.GetUsers()
	if err != nil {
		fmt.Println("Slack api error: ", err)
	}
	return apiClient, users
}

func boltSetup(users []slack.User) *bolt.DB {

	db, err := bolt.Open("slacktracker.db", 0600, nil)
	if err != nil {
		fmt.Println("Bolt DB error when opening: ", err)
	}
	spew.Dump(db)
	defer db.Close()

	fmt.Printf("date, ")
	for _, user := range users {
		if user.Presence != "" {
			// Create a table for each user we find
			// Do we need to react to users we haven't seen before differently?
			// Probably only want to run this sporadically
			db.Update(func(tx *bolt.Tx) error {
				userBucket, err := tx.CreateBucketIfNotExists([]byte(user.Name))
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
	return db
}

func main() {
	var apiClient, users = setup()
	var boltDB = boltSetup(users)
	spew.Dump(boltDB)

	// Idea being to open 1 channel per user (is this sustainable for large numbers?)
	// DOes it work with bolt more importantly?

	starttime := time.Now()
	runlength := 5000 * time.Second
	// Note, any value less than 15 seconds for frequency may get time out from slack API
	frequency := 15 * time.Second
	for time.Since(starttime) < runlength {
		time.Sleep(frequency)
		fmt.Println("")
		timestamp := time.Now().Format(time.RFC3339)
		fmt.Printf(timestamp)
		fmt.Printf(",")
		// Refresh users with new data
		users, err := apiClient.GetUsers()
		if err != nil {
			spew.Dump(users)
			spew.Dump(err)
		}
		for _, user := range users {
			// Ignore bots
			if !user.IsBot {
				// Check if presence empty (slackbot is a special Bot case apparently)
				if user.Presence != "" {
					boltDB.Update(func(tx *bolt.Tx) error {
						bucket := tx.Bucket([]byte(user.Name))
						err := bucket.Put([]byte(timestamp), []byte(user.Presence))
						return err
					})
					fmt.Printf(user.Presence)
					// Todo don't print if this is last user
					fmt.Printf(", ")
				}
			}
		}
	}
}
