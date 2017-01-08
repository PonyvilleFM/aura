// +build !pvfm

package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func genFname(u *discordgo.User) (string, error) {
	return fmt.Sprintf("%s - %s.mp3", u.Username, time.Now().Format(time.RFC822)), nil
}

func announce(msg string) error { return nil }

func showName() (string, error) { return "", nil }
