package bot

import (
	"fmt"

	"github.com/aureleoules/lntip/models"
	"github.com/aureleoules/lntip/rates"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"go.uber.org/zap"
)

func balanceHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	user, err := models.GetUser(m.Author.ID)
	if err != nil {
		zap.S().Error(err)
	}

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Your balance",
		Color: 0xFFFF00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Sats",
				Value: humanize.Comma(user.Balance),
			},
			{
				Name:  "Euro",
				Value: fmt.Sprintf("%f €", (float32(user.Balance)/100000000.)*rates.Price),
			},
		},
	})

}
