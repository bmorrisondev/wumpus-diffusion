package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bmorrisondev/aibot/commands"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const prefix string = "!mf"

func main() {
	godotenv.Load()

	token := os.Getenv("BOT_TOKEN")
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	// Handles chat stream
	sess.AddHandler(ChatStreamHandler)

	sess.Identify.Intents = discordgo.IntentsAll

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	RegisterSlashCommands(sess)

	fmt.Println("the bot is online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

var commandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

func RegisterSlashCommands(sess *discordgo.Session) {
	sess.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	registeredCommands := []commands.CommandExt{
		commands.ImgCommand,
	}

	for _, el := range registeredCommands {
		// Register definition
		commandHandlers[el.Definition.Name] = el.Handler

		// Register handler
		_, err := sess.ApplicationCommandCreate(sess.State.User.ID, os.Getenv("GUILD_ID"), &el.Definition)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ChatStreamHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Chat stream handlers
	args := strings.Split(m.Content, " ")

	if args[0] != prefix {
		return
	}

	if args[1] == "hello" {
		s.ChannelMessageSend(m.ChannelID, "world!")
	}
}
