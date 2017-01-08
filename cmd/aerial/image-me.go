package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"github.com/paddycarey/gophy"
)

func helperExtractTerms(commandName string, messageText string) string {
	sentence := strings.SplitAfterN(messageText, commandName, 2)

	return sentence[1]
}

func giphySearch(search string) (string, error) {
	giphyCO := &gophy.ClientOptions{}
	giphyClient := gophy.NewClient(giphyCO)

	escapedStr := url.QueryEscape(search)

	gifs, _, gifErr := giphyClient.SearchGifs(escapedStr, "pg-13", 10, 0)
	if gifErr != nil {
		return "", gifErr
	}

	gifCount := len(gifs)
	if gifCount == 0 {
		return fmt.Sprintf("Sorry, Giphy didn't return any results for _%s_", search), nil
	}

	selectedGif := randomdata.Number(0, gifCount-1)
	return gifs[selectedGif].Images.Original.URL, nil
}

func imageMe(search string) (string, error) {
	APIpath := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?q=%s&key=%s&cx=%s&searchtype=image&num=10",
		url.QueryEscape(search), gClientID, gClientSecret)

	// Create the HTTP Request and Headers for Kong to work.
	req, err := http.NewRequest("GET", APIpath, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 403 {
		logrus.Warn("403 Forbidden Status received.")
		return fmt.Sprintf("We appear to have run out of Google Search credits today. These will reset within 24 hours.\n Try not to be so colletively greedy tomorrow, maybe? You can always use _animate me_ for your memes."), nil
	}

	if resp.StatusCode != 200 {
		logrus.Error("Other Non-200 status received.")
		return "Sorry, something went wrong!", nil
	}

	var images GoogleSearch
	err = json.Unmarshal(bodyContent, &images)
	if err != nil {
		return "Google messed up", err
	}

	imgCount := len(images.Items)
	if imgCount == 0 {
		return fmt.Sprintf("Sorry, Google didn't return any results for _%s_", search), nil
	}

	selectedImg := randomdata.Number(0, imgCount-1)

	return images.Items[selectedImg].Pagemap.Imageobject[0].URL, nil
}

func imageMeEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	body := strings.ToLower(m.Content)

	var result string
	var err error
	if strings.Contains(body, "animate me") {
		result, err = giphySearch(helperExtractTerms("animate me", body))
	} else if strings.Contains(body, "image me") {
		result, err = imageMe(helperExtractTerms("image me", body))
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
	}

	s.ChannelMessageSend(m.ChannelID, result)
}

type GoogleSearch struct {
	Items []struct {
		DisplayLink      string `json:"displayLink"`
		FormattedURL     string `json:"formattedUrl"`
		HTMLFormattedURL string `json:"htmlFormattedUrl"`
		HTMLSnippet      string `json:"htmlSnippet"`
		HTMLTitle        string `json:"htmlTitle"`
		Kind             string `json:"kind"`
		Link             string `json:"link"`
		Pagemap          struct {
			Imageobject []struct {
				Caption         string `json:"caption"`
				Copyrightholder string `json:"copyrightholder"`
				Description     string `json:"description"`
				Height          string `json:"height"`
				URL             string `json:"url"`
				Width           string `json:"width"`
			} `json:"imageobject"`
		} `json:"pagemap"`
		Snippet string `json:"snippet"`
		Title   string `json:"title"`
	} `json:"items"`
}
