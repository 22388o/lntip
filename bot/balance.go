package bot

import (
	"fmt"

	"github.com/aureleoules/lntip/models"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

func balanceHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	user, err := models.GetUser(m.Author.ID)
	if err != nil {
		zap.S().Error(err)
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, you have %d sats.", m.Author.Mention(), user.Balance))
}
