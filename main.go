package main

import (
	"bytes"
	"flag"
	"fmt"
	api "github.com/GwentAPI/GwentBot/api"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

/*
Set this variable with go build with the -ldflags="-X main.version=<value>" parameter.
*/
var version = "undefined"

// Variables used for command line parameters
var (
	Token  string
	Client *http.Client
)

func init() {

	versionFlag := flag.Bool("v", false, "Prints current version")
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}
}

func main() {
	if Token == "" {
		fmt.Println("Token not set.")
		return
	}
	Client = &http.Client{Timeout: 10 * time.Second}
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by bots, including himself
	// This isn't required in this specific example but it's a good practice.
	// Also ignore lengthy messages
	if m.Author.Bot || len(m.Content) > 100 {
		return
	}

	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}
	if strings.HasPrefix(m.Content, "!card") {
		cardQuery := strings.TrimPrefix(m.Content, "!card")
		cardQuery = strings.TrimSpace((cardQuery))
		search(s, g, c, m.Author, cardQuery)
		return
	}

}

func search(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, u *discordgo.User, q string) {
	page, err := api.RequestPage(Client, q)

	if err != nil {
		s.ChannelMessageSend(c.ID, err.Error())
		return
	}
	var buffer bytes.Buffer
	card, err := api.RequestCard(Client, page.Results[0].Href)

	if err != nil {
		s.ChannelMessageSend(c.ID, err.Error())
		return
	}
	buffer.WriteString("``")
	buffer.WriteString(card.String())
	buffer.WriteString("``")
	s.ChannelMessageSend(c.ID, buffer.String())
	return
}
