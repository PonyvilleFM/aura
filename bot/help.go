package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (cs *CommandSet) help(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	switch len(parv) {
	case 1:
		// print all help on all commands
		result := "Bot commands: \n"

		for verb, cmd := range cs.cmds {
			result += fmt.Sprintf("%s%s: %s\n", cs.Prefix, verb, cmd.Helptext())
		}

		result += "If there's any problems please don't hesitate to ask a server admin for help."

		s.ChannelMessageSend(m.ChannelID, result)

	default:
		return ErrParvCountMismatch
	}

	return nil
}
