package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PonyvilleFM/aura/pvfm"
	"github.com/PonyvilleFM/aura/pvfm/pvl"
	pvfmschedule "github.com/PonyvilleFM/aura/pvfm/schedule"
	"github.com/PonyvilleFM/aura/pvfm/station"
	"github.com/bwmarrin/discordgo"
	"github.com/tebeka/strftime"
)

func pesterLink(s *discordgo.Session, m *discordgo.MessageCreate) {
	if musicLinkRegex.Match([]byte(m.Content)) {
		i, err := pvfm.GetStats()
		if err != nil {
			log.Println(err)
			return
		}

		if i.IsDJLive() && m.ChannelID == youtubeSpamRoomID {
			s.ChannelMessageSend(m.ChannelID, "Please be mindful sharing links to music when a DJ is performing. Thanks!")
		}
	}
}

func np(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	i, err := pvfm.GetStats()
	if err != nil {
		log.Printf("Can't get info: %v, failing over to plan b", err)
		return doStationRequest(s, m, parv)
	}

	result := []string{}

	if i.Main.Nowplaying == "Fetching info..." {
		log.Println("Main information was bad, fetching from station directly...")

		err := doStationRequest(s, m, parv)
		if err != nil {
			return err
		}

		return nil
	} else {
		result = append(result, "📻 **Now Playing on PVFM**\n")

		result = append(result, fmt.Sprintf(
			"Main 🎵 %s\n",
			i.Main.Nowplaying,
		))
		result = append(result, fmt.Sprintf(
			"Chill 🎵 %s\n",
			i.Secondary.Nowplaying,
		))
		result = append(result, fmt.Sprintf(
			"Free! 🎵 %s",
			i.MusicOnly.Nowplaying,
		))
	}

	s.ChannelMessageSend(m.ChannelID, strings.Join(result, "\n"))
	return nil
}

func dj(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	cal, err := pvl.Get()
	if err != nil {
		return err
	}

	now := cal.Result[0]
	result := []string{}

	localTime := time.Now()
	thentime := time.Unix(now.StartTime, 0)
	if thentime.Unix() < localTime.Unix() {
		result = append(result, fmt.Sprintf("Currently live: %s\n", now.Title))
		now = cal.Result[1]
	}

	nowTime := time.Unix(now.StartTime, 0).UTC()
	zone, _ := nowTime.Zone()
	fmttime, _ := strftime.Format("%Y-%m-%d %H:%M:%S", nowTime)

	result = append(result, fmt.Sprintf("Next event: %s at %s \x02%s\x02",
		now.Title,
		fmttime,
		zone,
	))

	s.ChannelMessageSend(m.ChannelID, strings.Join(result, "\n"))
	return nil
}

func stats(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	i, err := pvfm.GetStats()
	if err != nil {
		log.Printf("Error getting the station info: %v, falling back to plan b", err)
		return doStatsFromStation(s, m, parv)
	}

	st, err := station.GetStats()
	if err != nil {
		return err
	}

	var l int
	var peak int

	for _, source := range st.Icestats.Source {
		l = l + source.Listeners
		peak = peak + source.ListenerPeak
	}

	result := []string{
		fmt.Sprintf(
			"Current listeners across all streams: %d with a maximum of %d!",
			i.Listeners.Listeners, peak,
		),
		fmt.Sprintf(
			"Detailed: Main: %d listeners, Two: %d listeners, Free: %d listeners",
			i.Main.Listeners, i.Secondary.Listeners, i.MusicOnly.Listeners,
		),
	}

	s.ChannelMessageSend(m.ChannelID, strings.Join(result, "\n"))

	return nil
}

func schedule(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	result := []string{}
	schEntries, err := pvfmschedule.Get()
	if err != nil {
		return err
	}

	for _, entry := range schEntries {
		result = append(result, entry.String())
	}

	s.ChannelMessageSend(m.ChannelID, strings.Join(result, "\n"))
	return nil
}

func doStationRequest(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	stats, err := station.GetStats()
	if err != nil {
		return err
	}

	result := fmt.Sprintf(
		"Now playing: %s - %s on Ponyville FM!",
		stats.Icestats.Source[0].Title,
		stats.Icestats.Source[0].Artist,
	)

	s.ChannelMessageSend(m.ChannelID, result)
	return nil
}

func doStatsFromStation(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	st, err := station.GetStats()
	if err != nil {
		return err
	}

	var l int
	var peak int

	for _, source := range st.Icestats.Source {
		l = l + source.Listeners
		peak = peak + source.ListenerPeak
	}

	result := []string{
		fmt.Sprintf("Current listeners: %d with a maximum of %d!", l, peak),
	}

	s.ChannelMessageSend(m.ChannelID, strings.Join(result, "\n"))
	return nil
}

func curTime(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The time currently is %s", time.Now().UTC().Format("2006-01-02 15:04:05 UTC")))

	return nil
}

func streams(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	currentMeta, metaErr := station.GetStats()
	if metaErr != nil {
		s.ChannelMessageSend(m.ChannelID, "Error receiving pvfm metadata")
		return metaErr
	}

	outputString := "**PVFM Servers:**\n"

	for _, element := range currentMeta.Icestats.Source {
		outputString += ":musical_note: " + element.ServerDescription + ":\n`" + strings.Replace(element.Listenurl, "aerial", "dj.bronyradio.com", -1) + "`\n"
	}

	outputString += "\n:cd: DJ Recordings:\n`http://darkling.darkwizards.com/wang/BronyRadio/?M=D`"

	s.ChannelMessageSend(m.ChannelID, outputString)

	return nil
}
