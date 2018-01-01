package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"strings"
)

func printerFact(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	fact, err := getPrinterFact()
	if err != nil {
		return err
	}

	newFact := strings.Replace(fact, "kitten", "scanner", -1) // replace kitten with scanner because the API doesn't already do this

	s.ChannelMessageSend(m.ChannelID, newFact)
	return nil
}

func getPrinterFact() (string, error) {
	resp, err := http.Get("https://xena.stdlib.com/printerfacts")
	if err != nil {
		return "", err
	}

	factStruct := &struct {
		Facts []string `json:"facts"`
	}{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	json.Unmarshal(body, factStruct)

	text := fmt.Sprintf("%s", factStruct.Facts[0])

	return text, nil
}
