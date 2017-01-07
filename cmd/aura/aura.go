package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/PonyvilleFM/aura/bot"
	"github.com/PonyvilleFM/aura/recording"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
)

type aura struct {
	cs *bot.CommandSet
	s  *discordgo.Session

	guildRecordings map[string]*recording.Recording
	state           *state
}

type state struct {
	DownloadURLs map[string]string // Guild ID -> URL
	PermRoles    map[string]string // Guild ID -> needed role ID
}

func (s *state) Save() error {
	fout, err := os.Create(path.Join(dataPrefix, "state.json"))
	if err != nil {
		return err
	}
	defer fout.Close()

	return json.NewEncoder(fout).Encode(s)
}

func (s *state) Load() error {
	fin, err := os.Open(path.Join(dataPrefix, "state.json"))
	if err != nil {
		return err
	}
	defer fin.Close()

	return json.NewDecoder(fin).Decode(s)
}

const (
	djonHelp  = ``
	djoffHelp = ``
	setupHelp = ``
)

func (a *aura) Permissons(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	gid := ch.GuildID
	role := a.state.PermRoles[gid]

	gu, err := s.GuildMember(gid, m.Author.ID)
	if err != nil {
		return err
	}

	found := false
	for _, r := range gu.Roles {
		if r == role {
			found = true
			break
		}
	}

	if !found {
		return errors.New("aura: no permissions")
	}

	return nil
}

func (a *aura) roles(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	log.Println("got here")
	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	gid := ch.GuildID

	result := "Roles in this group:\n"

	roles, err := s.GuildRoles(gid)
	if err != nil {
		return err
	}

	for _, r := range roles {
		result += fmt.Sprintf("- %s: %s\n", r.ID, r.Name)
	}

	s.ChannelMessageSend(m.ChannelID, result)
	return nil
}

func (a *aura) setup(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	if len(parv) != 3 {
		return errors.New("aura: wrong number of params for setup")
	}

	role := parv[1]
	url := parv[2]

	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	gid := ch.GuildID

	roles, err := s.GuildRoles(gid)
	if err != nil {
		return err
	}

	found := false
	for _, r := range roles {
		if r.ID == role {
			found = true
			break
		}
	}

	if !found {
		return errors.New("aura: Role not found")
	}

	a.state.PermRoles[gid] = role
	a.state.DownloadURLs[gid] = url

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Guild %s set up for recording url %s controlled by role %s", gid, url, role))

	a.state.Save()
	return nil
}

func (a *aura) djon(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	fname := fmt.Sprintf("%s - %s.mp3", m.Author.Username, time.Now().Format(time.ANSIC))

	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	gid := ch.GuildID

	_, ok := a.guildRecordings[gid]
	if ok {
		return errors.New("aura: another recording is already in progress")
	}

	r, err := recording.New(a.state.DownloadURLs[gid], path.Join(dataPrefix, fname))
	if err != nil {
		return err
	}

	a.guildRecordings[gid] = r
	go r.Start()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Now recording: %s", fname))

	return nil
}

func (a *aura) djoff(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	gid := ch.GuildID

	r, ok := a.guildRecordings[gid]
	if !ok {
		return errors.New("aura: no recording is currently in progress")
	}

	r.Cancel()
	<-r.Done()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Recording complete: %s", r.OutputFilename()))

	return nil
}

func (a *aura) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := a.cs.Run(s, m.Message)
	if err != nil {
		log.Println(err)
	}
}

var (
	token      = os.Getenv("TOKEN")
	dataPrefix = os.Getenv("DATA_PREFIX")
)

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	a := &aura{
		cs:              bot.NewCommandSet(),
		s:               dg,
		guildRecordings: map[string]*recording.Recording{},

		state: &state{
			DownloadURLs: map[string]string{},
			PermRoles:    map[string]string{},
		},
	}

	err = a.state.Load()
	if err != nil {
		log.Println(err)
	}

	a.cs.Add(bot.NewBasicCommand("roles", "", bot.NoPermissions, a.roles))
	a.cs.Add(bot.NewBasicCommand("setup", setupHelp, bot.NoPermissions, a.setup))
	a.cs.Add(bot.NewBasicCommand("djon", djonHelp, a.Permissons, a.djon))
	a.cs.Add(bot.NewBasicCommand("djoff", djoffHelp, a.Permissons, a.djoff))

	dg.AddHandler(a.Handle)
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ready")
	<-make(chan struct{})
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Print message to stdout.
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
}
