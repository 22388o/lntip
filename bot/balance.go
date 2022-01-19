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

	str := fmt.Sprintf("%s, you have %s sats (%f €).", m.Author.Mention(), humanize.Comma(user.Balance), (float32(user.Balance)/100000000.)*rates.Price)
	s.ChannelMessageSend(m.ChannelID, str)
}
