package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	ulock          sync.Mutex
	raffleUsers    map[string]map[string]struct{}
	raffleSessions map[string]string
)

func raffleStart(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	ulock.Lock()
	defer ulock.Unlock()

	raffleSessions[m.ChannelID] = m.Author.ID
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s> created a new raffle! Type %senter to enter!", m.Author.ID, ";"))

	return nil
}

func rafflePop(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	ulock.Lock()
	defer ulock.Unlock()

	creator, ok := raffleSessions[m.ChannelID]
	if !ok {
		return errors.New("no raffle currently running")
	}

	if m.Author.ID != creator {
		return fmt.Errorf("you are not <@%s>, you cannot end this raffle", creator)
	}

	cmap, ok := raffleUsers[m.ChannelID]
	if !ok {
		cmap = map[string]struct{}{}
	}

	users := []string{}
	for u := range cmap {
		users = append(users, u)
	}
	i := rand.Intn(len(users))
	u := users[i]

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Congrats <@%s>, you won!", u))

	delete(cmap, u)
	raffleUsers[m.ChannelID] = cmap
}

func raffleEnd(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	ulock.Lock()
	defer ulock.Unlock()

	creator, ok := raffleSessions[m.ChannelID]
	if !ok {
		return errors.New("no raffle currently running")
	}

	if m.Author.ID != creator {
		return fmt.Errorf("you are not <@%s>, you cannot end this raffle", creator)
	}

	delete(raffleUsers, m.ChannelID)
	delete(raffleSessions, m.ChannelID)

	s.ChannelMessageSend(m.ChannelID, "Raffle destroyed.")

	return nil
}

func raffle(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	ulock.Lock()
	defer ulock.Unlock()

	_, ok := raffleSessions[m.ChannelID]
	if !ok {
		return errors.New("no raffle currently running")
	}

	cmap, ok := raffleUsers[m.ChannelID]
	if !ok {
		cmap = map[string]struct{}{}
	}

	_, ok = cmap[m.Author.ID]
	if ok {
		return errors.New("you already joined this raffle, you can't join it again")
	}

	cmap[m.Author.ID] = struct{}{}
	raffleUsers[m.ChannelID] = cmap

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%s>: your entry has been recorded!", m.Author.ID))
	return nil
}
