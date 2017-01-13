package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	sayings = []string{
		"throws a boat at $PERSON",
		"gives $PERSON a trip to flavortown",
		// etc
	}
)

func snack(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	saying := sayings[rand.Intn(len(sayings))]
	person := fmt.Sprintf("<@%s>", m.Author.ID)
	pf := func(x string) string {
		if x == "PERSON" {
			return person
		}
		return ""
	}
	result := os.Expand(saying, pf)
	_, err := s.ChannelMessageSend(m.ChannelID, result)
	return err
}
