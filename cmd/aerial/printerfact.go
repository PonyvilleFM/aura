package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
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
	resp, err := http.Get("https://printerfacts.cetacean.club/fact")
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
