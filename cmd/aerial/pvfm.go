package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	derpiSearch "github.com/PonyvilleFM/aura/cmd/aerial/derpi"
	"github.com/PonyvilleFM/aura/internal/pvfm"
	"github.com/PonyvilleFM/aura/internal/pvfm/pvl"
	pvfmschedule "github.com/PonyvilleFM/aura/internal/pvfm/schedule"
	"github.com/PonyvilleFM/aura/internal/pvfm/station"
	"github.com/bwmarrin/discordgo"
	"github.com/tebeka/strftime"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// randomRange gives a random whole integer between the given integers [min, max)
func randomRange(min, max int) int {
	return rand.Intn(max-min) + min
}

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

func stats(s *discordgo.Session, m *discordgo.Message, parv []string) error {

	// Regular metadata info

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

	// Live DJ info

	// init variables
	cal, err := pvl.Get()
	if err != nil {
		return err
	}
	now := cal.Result[0]

	// times
	localTime := time.Now()
	thentime := time.Unix(now.StartTime, 0)

	// checks if the event is currently happening
	djInfo := "" // since we start with a conditional...
	if thentime.Unix() < localTime.Unix() {
		djInfo += fmt.Sprintf("**Currently live!**\n%s\n\n", now.Title)
		now = cal.Result[1]
	}

	// Prepare time string
	nowTime := time.Unix(now.StartTime, 0).UTC()
	zone, _ := nowTime.Zone()
	fmttime, _ := strftime.Format("%Y-%m-%d %H:%M:%S", nowTime)

	// Piece data together into the result
	djInfo += fmt.Sprintf("Next event:\n%s\n%s \x02%s\x02",
		now.Title,
		fmttime,
		zone,
	)

	outputEmbed := NewEmbed().
		SetTitle("Listener Statistics").
		SetDescription("Use `;streams` if you need a link to the radio!\nTotal listeners across all stations: " + strconv.Itoa(i.Listeners.Listeners) + " with a maximum  of " + strconv.Itoa(peak) + ".")

	outputEmbed.AddField("🎵 Main", strconv.Itoa(i.Main.Listeners)+" listeners.\n"+i.Main.Nowplaying)
	outputEmbed.AddField("🎵 Chill", strconv.Itoa(i.Secondary.Listeners)+" listeners.\n"+i.Secondary.Nowplaying)
	outputEmbed.AddField("🎵 Free! (no DJ sets)", strconv.Itoa(i.MusicOnly.Listeners)+" listeners.\n"+i.MusicOnly.Nowplaying)
	outputEmbed.AddField("🎛 Live DJs", djInfo)

	outputEmbed.InlineAllFields()

	s.ChannelMessageSendEmbed(m.ChannelID, outputEmbed.MessageEmbed)

	return nil
}

func schedule(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	schEntries, err := pvfmschedule.Get()
	if err != nil {
		return err
	}

	// Create embed object
	outputEmbed := NewEmbed().
		SetTitle("Upcoming Shows").
		SetDescription("These are the upcoming shows and events airing soon on PVFM 1.\n[Convert to your timezone](https://www.worldtimebuddy.com/?pl=1&lid=100&h=100)")

	for _, entry := range schEntries {

		// Format countdown timer
		startTimeUnix := time.Unix(int64(entry.StartUnix), 0)
		nowWithoutNanoseconds := time.Unix(time.Now().Unix(), 0)
		dur := startTimeUnix.Sub(nowWithoutNanoseconds)

		// Show "Live Now!" if the timer is less than 0h0m0s
		if dur > 0 {
			outputEmbed.AddField(":musical_note:  "+entry.Host+" - "+entry.Name, entry.StartTime+" "+entry.Timezone+"\nAirs in "+dur.String())
		} else {
			outputEmbed.AddField(":musical_note:  "+entry.Host+" - "+entry.Name, "Live now!")
		}
	}

	s.ChannelMessageSendEmbed(m.ChannelID, outputEmbed.MessageEmbed)
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
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The time currently is %s\nUse <https://www.worldtimebuddy.com/?pl=1&lid=100&h=100> to convert UTC to your local timezone.", time.Now().UTC().Format("2006-01-02 15:04:05 UTC")))

	return nil
}

func streams(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	currentMeta, metaErr := station.GetStats()
	if metaErr != nil {
		s.ChannelMessageSend(m.ChannelID, "Error receiving pvfm metadata")
		return metaErr
	}

	// start building custom embed
	outputEmbed := NewEmbed().
		SetTitle("Stream Links").
		SetDescription("These are direct feeds of the live streams; most browsers and media players can play them!")

	// this will dynamically build the list from station metadata
	pvfmList := ""
	for _, element := range currentMeta.Icestats.Source {
		pvfmList += element.ServerDescription + ":\n<" + strings.Replace(element.Listenurl, "aerial", "dj.bronyradio.com", -1) + ">\n"
	}

	// PVFM
	outputEmbed.AddField(":musical_note:  PVFM Servers", pvfmList)
	// Luna Radio
	outputEmbed.AddField(":musical_note:  Luna Radio Servers", "Luna Radio MP3 128Kbps Stream:\n<http://radio.ponyvillelive.com:8002/stream.mp3>\nLuna Radio Mobile MP3 64Kbps Stream:\n<http://radio.ponyvillelive.com:8002/mobile?;stream.mp3>\n")
	// Recordings
	outputEmbed.AddField(":cd:  DJ Recordings", "Archive\n<https://pvfmsets.cf/var/93252527679639552/>\nLegacy Archive\n<http://darkling.darkwizards.com/wang/BronyRadio/?M=D>")

	s.ChannelMessageSendEmbed(m.ChannelID, outputEmbed.MessageEmbed)

	// no errors yay!!!!
	return nil
}

func derpi(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	if m.ChannelID == "292755043684450304" {

		searchResults, err := derpiSearch.SearchDerpi(m.Content[7:len(m.Content)]) // Safe tag will be added in derpi/derpi.go
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "An error occured.")
			return err
		}
		if len(searchResults.Search) < 1 {
			s.ChannelMessageSend(m.ChannelID, "Error: No results")
			return nil
		}
		derpiImage := searchResults.Search[randomRange(0, len(searchResults.Search))]

		tags := strings.Split(derpiImage.Tags, ", ") // because this isn't an array for some reason

		// Check for artist tag
		artist := ""
		for _, tag := range tags {
			if strings.Contains(tag, "artist:") {
				artist = tag[7:]
			}
		}

		outputEmbed := NewEmbed().
			SetTitle("Derpibooru Image").
			SetURL("https://derpibooru.org/" + derpiImage.ID).
			SetDescription(derpiImage.Description).
			SetImage("http:" + derpiImage.Image).
			SetFooter("Image score: " + strconv.Itoa(derpiImage.Score) + " | Uploaded: " + derpiImage.CreatedAt.String())

		// Credit the artist!
		if artist == "" {
			outputEmbed.SetAuthor("No artist")
		} else {
			outputEmbed.SetAuthor("Artist: " + artist)
		}

		s.ChannelMessageSendEmbed(m.ChannelID, outputEmbed.MessageEmbed)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please use this command in <#292755043684450304> only.")
	}
	return nil
}

func weather(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	responses := []string{
		"Cloudy with a chance of meatballs.",
		"It's currently pouring down even more than Pinkie.",
		"It's the most overcast I've ever seen. In other words, same as always.",
		"Do you have a better conversation starter than that?",
		"There's at least 5 or 6 weather right now, my dude.",
		"It's soggy enough for Rainbow Dash to get fired, if she didn't have a literal deity keeping her in charge.",
		"Surprisingly, the weather is pretty alright.",
		"You'd be happy to know that it's hot enough to make a phoenix sweat.",
		"The weather right now is like you took London and stuck it in a dishwasher.",
		"The Crystal Empire is warmer than this weather.",
	}

	s.ChannelMessageSend(m.ChannelID, responses[randomRange(0, len(responses))])

	return nil
}
