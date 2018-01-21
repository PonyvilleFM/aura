package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (cs *CommandSet) help(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	switch len(parv) {
	case 1:
		// print all help on all commands
		result := formHelp()

		s.ChannelMessageSend(m.ChannelID, result)

	default:
		return ErrParvCountMismatch
	}

	return nil
}

func (cs *CommandSet) dhelp(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	switch len(parv) {
	case 1:

		result := formHelp()

		authorChannel, err := s.UserChannelCreate(m.Author.ID)
		s.ChannelMessageSend(authorChannel.ID, result)

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("@%s check direct messages, help is there!", m.Author.Username))

	default:
		return ErrParvCountMismatch
	}

	return nil
}

func formHelp() string {
	result := "Bot commands: \n"

	for verb, cmd := range cs.cmds {
		result += fmt.Sprintf("%s%s: %s\n", cs.Prefix, verb, cmd.Helptext())
	}

	return (result + "If there's any problems please don't hesitate to ask a server admin for help.")
}
