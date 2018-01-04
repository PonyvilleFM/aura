// Package main provides a basic example of using go-cleverbot.io.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/CleverbotIO/go-cleverbot.io"
)

func main() {
	// The api key is given to you at https://cleverbot.io/keys.
	apiUser := "YOUR_API_USER"
	apiKey := "YOUR_API_KEY"

	// apiNick is optional.
	apiNick := ""

	// Initialize Cleverbot.
	bot, err := cleverbot.New(apiUser, apiKey, apiNick)
	if err != nil {
		log.Fatal(err)
	}

	// Start the chat loop.
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Press CTRL-C to exit.")
	for scanner.Scan() {
		response, err := Talk(bot, scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Cleverbot: " + response)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

// Talk sends Cleverbot a message, returns Cleverbot's response.
func Talk(bot *cleverbot.Session, input string) (response string, err error) {
	response, err = bot.Ask(input)
	if err != nil {
		return
	}
	return
}
