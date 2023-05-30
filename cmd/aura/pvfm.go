package main

import (
	"fmt"
	"os"
	"time"

	"github.com/PonyvilleFM/aura/internal/pvfm/pvl"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func genFname(username string) (string, error) {
	return fmt.Sprintf("%s - %s.mp3", username, time.Now().Format(time.RFC3339)), nil
}

func showName() (string, error) {
	cal, err := pvl.Get()
	if err != nil {
		return "", nil
	}

	now := cal.Result[0]
	return now.Title, nil
}
