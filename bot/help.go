package bot

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func helpHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var fields []*discordgo.MessageEmbedField
	for _, v := range commands {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  v.usage,
			Value: v.description,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  prefix + "lntip",
		Value: "Shows this message",
	})

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Help",
		Description: "lntip is a bot that allows you to tip Discord users.\n\n[Here](https://www.walletofsatoshi.com/) is a good Lightning wallet for beginners if you need one.",
		Fields:      fields,
		Type:        discordgo.EmbedTypeRich,
		Color:       0xFFFF00,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Built by @Nuf",
		},
	})

	if err != nil {
		zap.S().Error(err)
		return
	}

}
