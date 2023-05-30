package main

import (
	"fmt"
	"time"

	"github.com/PonyvilleFM/aura/internal/pvfm/pvl"
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
