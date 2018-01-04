// Package main provides a basic example of using go-cleverbot.io.
package main

import (
	"fmt"
	"log"

	"github.com/CleverbotIO/go-cleverbot.io"
)

func main() {
	// The api key is given to you at https://cleverbot.io/keys.
	apiUser := "YOUR_API_USER"
	apiKey := "YOUR_API_KEY"

	// apiNick is optional.
	apiNick := ""

	// Initialize Cleverbot
	bot, err := cleverbot.New(apiUser, apiKey, apiNick)
	if err != nil {
		log.Fatal(err)
	}

	// Send Cleverbot a message.
	response, err := bot.Ask("hello world")
	if err != nil {
		log.Fatal(err)
	}

	// Print the response.
	fmt.Println(response)
	// "World? Who is world? My name is Timmy."
}
