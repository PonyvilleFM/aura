package cleverbot

// Session ...
type Session struct {
	User string `json:"user"`
	Key  string `json:"key"`
	Nick string `json:"nick"`
	Text string `json:"text"`
}
