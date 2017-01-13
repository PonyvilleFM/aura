package bot

import (
	"errors"
	"log"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	verb     string
	helptext string
}

func (c *Command) Verb() string {
	return c.verb
}

func (c *Command) Helptext() string {
	return c.helptext
}

type Handler func(*discordgo.Session, *discordgo.Message, []string) error

type CommandHandler interface {
	Verb() string
	Helptext() string

	Handler(*discordgo.Session, *discordgo.Message, []string) error
	Permissions(*discordgo.Session, *discordgo.Message, []string) error
}

type basicCommand struct {
	*Command
	handler     Handler
	permissions Handler
}

func (bc *basicCommand) Handler(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	return bc.handler(s, m, parv)
}

func (bc *basicCommand) Permissions(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	return bc.permissions(s, m, parv)
}

var (
	DefaultCommandSet = NewCommandSet()

	ErrAlreadyExists     = errors.New("bot: command already exists")
	ErrNoSuchCommand     = errors.New("bot: no such command exists")
	ErrNoPermissions     = errors.New("bot: you do not have permissions for this command")
	ErrParvCountMismatch = errors.New("bot: parameter count mismatch")
)

const (
	DefaultPrefix = "."
)

// NewCommand creates an anonymous command and adds it to the default CommandSet.
func NewCommand(verb, helptext string, handler, permissions Handler) error {
	return DefaultCommandSet.Add(NewBasicCommand(verb, helptext, handler, permissions))
}

func NewBasicCommand(verb, helptext string, permissions, handler Handler) CommandHandler {
	return &basicCommand{
		Command: &Command{
			verb:     verb,
			helptext: helptext,
		},
		handler:     handler,
		permissions: permissions,
	}
}

type CommandSet struct {
	sync.Mutex
	cmds map[string]CommandHandler

	Prefix string
}

func NewCommandSet() *CommandSet {
	cs := &CommandSet{
		cmds:   map[string]CommandHandler{},
		Prefix: DefaultPrefix,
	}

	cs.AddCmd("help", "Shows help for the bot", NoPermissions, cs.help)

	return cs
}

func NoPermissions(s *discordgo.Session, m *discordgo.Message, parv []string) error {
	return nil
}

func (cs *CommandSet) AddCmd(verb, helptext string, permissions, handler Handler) error {
	log.Printf("Added command %s: %s", verb, helptext)
	return cs.Add(NewBasicCommand(verb, helptext, permissions, handler))
}

func (cs *CommandSet) Add(h CommandHandler) error {
	cs.Lock()
	defer cs.Unlock()

	v := strings.ToLower(h.Verb())

	if _, ok := cs.cmds[v]; ok {
		return ErrAlreadyExists
	}

	cs.cmds[v] = h

	return nil
}

func (cs *CommandSet) Run(s *discordgo.Session, msg *discordgo.Message) error {
	cs.Lock()
	defer cs.Unlock()

	if strings.HasPrefix(msg.Content, cs.Prefix) {
		params := strings.Fields(msg.Content)
		verb := strings.ToLower(params[0][1:])

		cmd, ok := cs.cmds[verb]
		if !ok {
			return ErrNoSuchCommand
		}

		err := cmd.Permissions(s, msg, params)
		if err != nil {
			log.Printf("Permissions error: %s: %v", msg.Author.Username, err)
			s.ChannelMessageSend(msg.ChannelID, "You don't have permissions for that, sorry.")
			return ErrNoPermissions
		}

		log.Println("calling command handler")
		err = cmd.Handler(s, msg, params)
		if err != nil {
			log.Printf("command handler error: %v", err)
			s.ChannelMessageSend(msg.ChannelID, "error: "+err.Error())
			return err
		}
	}

	return nil
}
