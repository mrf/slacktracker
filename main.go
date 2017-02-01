package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
	"io/ioutil"
)

func main() {
	config, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(string(config))
	api := slack.New(token)
	users, err := api.GetUsers()
	if err != nil {
		fmt.Println("errr: ", err)
		return
	}
	spew.Dump(users)
}
