// +build !pvfm

package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func getFname(u *discordgo.User) (string, error) {
	return fmt.Sprintf("%s - %s.mp3", u.Username, time.Now().Format(time.RFC822)), nil
}
