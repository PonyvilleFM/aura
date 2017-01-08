package station

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	latestInfo Wrapper

	bugTime = flag.Int("station-poke-delay", 15, "how stale the info can get")
)

type Wrapper struct {
	Age  time.Time
	Info Info
}

type Info struct {
	Icestats struct {
		Admin              string `json:"admin"`
		Host               string `json:"host"`
		Location           string `json:"location"`
		ServerID           string `json:"server_id"`
		ServerStart        string `json:"server_start"`
		ServerStartIso8601 string `json:"server_start_iso8601"`
		Source             []struct {
			Artist             string      `json:"artist"`
			AudioBitrate       int         `json:"audio_bitrate"`
			AudioChannels      int         `json:"audio_channels"`
			AudioInfo          string      `json:"audio_info"`
			AudioSamplerate    int         `json:"audio_samplerate"`
			Channels           int         `json:"channels"`
			Genre              string      `json:"genre"`
			IceBitrate         int         `json:"ice-bitrate"`
			ListenerPeak       int         `json:"listener_peak"`
			Listeners          int         `json:"listeners"`
			Listenurl          string      `json:"listenurl"`
			Quality            string      `json:"quality"`
			Samplerate         int         `json:"samplerate"`
			ServerDescription  string      `json:"server_description"`
			ServerName         string      `json:"server_name"`
			ServerType         string      `json:"server_type"`
			ServerURL          string      `json:"server_url"`
			StreamStart        string      `json:"stream_start"`
			StreamStartIso8601 string      `json:"stream_start_iso8601"`
			Subtype            string      `json:"subtype"`
			Title              string      `json:"title"`
			Dummy              interface{} `json:"dummy"`
		} `json:"source"`
	} `json:"icestats"`
}

func GetStats() (Info, error) {
	now := time.Now()
	if now.Before(latestInfo.Age.Add(time.Second * time.Duration(*bugTime))) {
		return latestInfo.Info, nil
	}

	i := Info{}

	resp, err := http.Get("http://dj.bronyradio.com:7090/status-json.xsl")
	if err != nil {
		return Info{}, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Info{}, err
	}

	err = json.Unmarshal(content, &i)
	if err != nil {
		return Info{}, err
	}

	latestInfo.Info = i
	latestInfo.Age = now

	return latestInfo.Info, nil
}
