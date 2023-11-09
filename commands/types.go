package commands

import "github.com/bwmarrin/discordgo"

type CommandExt struct {
	Definition discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
