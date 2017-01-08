// +build pvfm

package main

import (
	"github.com/PonyvilleFM/aura/pvfm"
	"github.com/bwmarrin/discordgo"
)

func genFname(u *discordgo.User) (string, error) {
	return pvfm.GenFilename()
}
