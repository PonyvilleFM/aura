package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PonyvilleFM/aura/bot"
	"github.com/PonyvilleFM/aura/commands/source"
	"github.com/PonyvilleFM/aura/recording"
	"github.com/bwmarrin/discordgo"
	_ "github.com/joho/godotenv/autoload"
	hashids "github.com/speps/go-hashids"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type aura struct {
	cs *bot.CommandSet
	s  *discordgo.Session

	guildRecordings map[string]*rec
	state           *state
	hid             *hashids.HashID
}

type state struct {
	DownloadURLs map[string]string // Guild ID -> URL
	PermRoles    map[string]string // Guild ID -> needed role ID
	Shorturls    map[string]string // hashid -> partial route
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

type rec struct {
	*recording.Recording
	creator string
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
	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	gid := ch.GuildID
	creator := m.Author.Username

	member, err := s.GuildMember(gid, m.Author.ID)
	if err != nil {
		return err
	}

	if member.Nick != "" {
		creator = member.Nick
	}

	fname, err := genFname(creator)
	if err != nil {
		return err
	}

	_, ok := a.guildRecordings[gid]
	if ok {
		log.Println(a.guildRecordings)
		return errors.New("aura: another recording is already in progress")
	}

	os.Mkdir(path.Join(dataPrefix, gid), 0775)

	rr, err := recording.New(a.state.DownloadURLs[gid], path.Join(dataPrefix, gid, fname))
	if err != nil {
		return err
	}

	a.guildRecordings[gid] = &rec{
		Recording: rr,
		creator:   creator,
	}

	go func() {
		err := rr.Start()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("recording error: %v", err))
			return
		}
	}()

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Now recording: `%s`\n\n%s get in here fam", fname, os.Getenv("NOTIFICATION_SQUAD_ID")))

	inv, err := s.ChannelInviteCreate(gid, discordgo.Invite{
		MaxAge: 4800,
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	invurl := "http://discord.gg/" + inv.Code

	err = announce("Live DJ on-air: " + creator + "\nJoin our chat here: " + invurl)
	if err != nil {
		log.Println(err)
		return nil
	}

	go a.waitAndAnnounce(s, m, a.guildRecordings[gid], gid)

	return nil
}

func (a *aura) waitAndAnnounce(s *discordgo.Session, m *discordgo.Message, r *rec, gid string) {
	<-r.Done()

	defer delete(a.guildRecordings, gid)

	fname := r.OutputFilename()
	parts := strings.Split(fname, "/")

	recurl := fmt.Sprintf("https://%s/var/%s/%s", recordingDomain, parts[1], urlencode(parts[2]))
	id, err := a.hid.EncodeInt64([]int64{int64(rand.Int())})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("This state should be impossible. Recording saved but unknown short URL: %v", err))
		return
	}

	a.state.Shorturls[id] = recurl
	a.state.Save()

	slink := fmt.Sprintf("https://%s/id/%s", recordingDomain, id)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Recording complete (%s): %s", time.Now().Sub(r.StartTime()).String(), slink))

	sn := r.creator

	msg := fmt.Sprintf("New recording by %s: %s", sn, slink)
	err = announce(msg)
	if err != nil {
		log.Println(err)
		return
	}
}

func urlencode(inp string) string {
	return (&url.URL{Path: inp}).String()
}

func (a *aura) djoff(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		return err
	}

	gid := ch.GuildID

	r, ok := a.guildRecordings[gid]
	if r == nil || !ok {
		log.Println(a.guildRecordings)
		return errors.New("aura: no recording is currently in progress")
	}

	if r.Err == nil {
		s.ChannelMessageSend(m.ChannelID, "Finishing recording (waiting 30 seconds)")
		time.Sleep(30 * time.Second)

		r.Cancel()
	}

	return nil
}

func (a *aura) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	err := a.cs.Run(s, m.Message)
	if err != nil {
		log.Println(err)
	}
}

var (
	token           = os.Getenv("TOKEN")
	dataPrefix      = os.Getenv("DATA_PREFIX")
	recordingDomain = os.Getenv("RECORDING_DOMAIN")
	hashidsSalt     = os.Getenv("HASHIDS_SALT")
	port            = os.Getenv("PORT")
)

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	hid := hashids.NewData()
	hid.Salt = hashidsSalt

	a := &aura{
		cs:              bot.NewCommandSet(),
		s:               dg,
		guildRecordings: map[string]*rec{},

		hid: hashids.NewWithData(hid),

		state: &state{
			DownloadURLs: map[string]string{},
			PermRoles:    map[string]string{},
			Shorturls:    map[string]string{},
		},
	}

	err = a.state.Load()
	if err != nil {
		log.Println(err)
	}

	a.cs.AddCmd("roles", "", bot.NoPermissions, a.roles)
	a.cs.AddCmd("setup", setupHelp, bot.NoPermissions, a.setup)
	a.cs.AddCmd("djon", djonHelp, a.Permissons, a.djon)
	a.cs.AddCmd("djoff", djoffHelp, a.Permissons, a.djoff)
	a.cs.AddCmd("source", "Source code information", bot.NoPermissions, source.Source)

	dg.AddHandler(a.Handle)
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("ready")

	http.Handle("/id/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.RequestURI[4:]

		redir, ok := a.state.Shorturls[id]
		if !ok {
			http.Error(w, "not found, sorry", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, redir, http.StatusFound)
	}))

	http.HandleFunc("/links.json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(a.state.Shorturls)
	})

	http.Handle("/var/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./index.html")
	})

	http.ListenAndServe(":"+port, nil)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Print message to stdout.
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
}
