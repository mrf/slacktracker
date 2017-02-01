package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/nlopes/slack"
	"io/ioutil"
	"strings"
)

func main() {
	config, err := ioutil.ReadFile("config.yml")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(string(config))
	token := strings.TrimPrefix(string(config), "token: ")
	token = strings.TrimSpace(token)
	api := slack.New(token)
	spew.Dump(api)
	users, err := api.GetUsers()
	if err != nil {
		fmt.Println("errr: ", err)
		return
	}
	spew.Dump(users)
}
