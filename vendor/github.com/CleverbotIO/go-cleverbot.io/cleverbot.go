// Package cleverbot implements wrapper for the cleverbot.io API.
package cleverbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// API Endpoints.
const (
	baseURL   = "https://cleverbot.io/1.0/"
	createURL = baseURL + "create"
	askURL    = baseURL + "ask"
)

// New bot instance.
// "nick" is optional if you did not specify it, a random one is generated for you.
// A successful call returns err == nil.
func New(user, key string, nick ...string) (s *Session, err error) {
	var sessionName string
	if len(nick) > 0 {
		sessionName = nick[0]
	}

	s = &Session{
		User: user,
		Key:  key,
		Nick: sessionName,
	}

	params, err := json.Marshal(s)
	if err != nil {
		return
	}

	response, err := http.Post(createURL, "application/json", bytes.NewBuffer(params))
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	m := map[string]string{}
	err = json.Unmarshal([]byte(body), &m)
	if err != nil {
		return
	}

	if m["status"] == "success" {
		s.Nick = m["nick"]
	} else {
		err = fmt.Errorf(m["status"])
		return
	}

	return
}

// Ask Cleverbot a question, returns Cleverbots response.
// A successful call returns err == nil.
func (s *Session) Ask(text string) (output string, err error) {

	s.Text = text

	params, err := json.Marshal(s)
	if err != nil {
		return
	}

	response, err := http.Post(askURL, "application/json", bytes.NewBuffer(params))
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	m := map[string]string{}
	err = json.Unmarshal([]byte(body), &m)
	if err != nil {
		return
	}

	if m["status"] != "success" {
		err = fmt.Errorf(m["status"])
		return
	}

	// return the bots asnwer.
	return m["response"], nil
}
