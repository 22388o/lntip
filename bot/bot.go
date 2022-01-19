package bot

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"

	"github.com/aureleoules/lntip/cfg"
	"github.com/bwmarrin/discordgo"
)

const prefix = "!"

type command struct {
	name        string
	f           func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
	dmOnly      bool
	usage       string
	description string
}

var discord *discordgo.Session
var commands = []command{
	{
		name:        "deposit",
		usage:       "!deposit <amount>",
		f:           depositHandler,
		dmOnly:      true,
		description: "Deposit sats in your account.",
	},
	{
		name:        "lntip",
		usage:       "!lntip <@user> <amount>",
		f:           lntipHandler,
		description: "Send a tip to a user.\nYou can also reward users by reacting with a **lntip<amount>** emoji on their messages.",
	},
	{
		name:        "balance",
		usage:       "!balance",
		f:           balanceHandler,
		dmOnly:      true,
		description: "Check your balance.",
	},
	{
		name:        "withdraw",
		usage:       "!withdraw <amount>",
		f:           withdrawHandler,
		dmOnly:      true,
		description: "Withdraw sats to your wallet.",
	},
}

var usersMutex = make(map[string]*sync.Mutex)

func Run() {
	var err error
	discord, err = discordgo.New("Bot " + cfg.Config.Bot.Token)
	if err != nil {
		panic(err)
	}

	discord.AddHandler(messageCreate)
	discord.AddHandler(reactionAdd)

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsDirectMessageReactions | discordgo.IntentsGuildMessageReactions)

	err = discord.Open()
	if err != nil {
		panic(err)
	}

	err = discord.UpdateGameStatus(0, "!lntip")
	if err != nil {
		panic(err)
	}

	go watchInvoices()

	defer discord.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	fmt.Println("Shutdown...")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, prefix+"help") {
		helpHandler(s, m, nil)
		return
	}

	_, ok := usersMutex[m.Author.ID]
	if !ok {
		usersMutex[m.Author.ID] = &sync.Mutex{}
	}

	usersMutex[m.Author.ID].Lock()
	for _, c := range commands {
		space := regexp.MustCompile(`\s+`)
		if strings.HasPrefix(m.Content, prefix+c.name) {
			c.f(s, m, strings.Split(space.ReplaceAllString(m.Content, " "), " ")[1:])
			break
		}
	}
	usersMutex[m.Author.ID].Unlock()
}
