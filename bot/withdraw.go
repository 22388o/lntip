package bot

import (
	"context"

	"github.com/aureleoules/lntip/cfg"
	"github.com/aureleoules/lntip/lnclient"
	"github.com/aureleoules/lntip/models"
	"github.com/bwmarrin/discordgo"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/zpay32"
	"go.uber.org/zap"
)

func withdrawHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) != 1 {
		s.ChannelMessageSend(m.ChannelID, "Usage: `!withdraw <invoice>`")
		return
	}

	invoice, err := zpay32.Decode(args[0], cfg.ChainParams())
	if err != nil {
		zap.S().Errorw("Error decoding invoice", "error", err)
		s.ChannelMessageSend(m.ChannelID, "Invalid invoice")
		return
	}

	if invoice.MilliSat == nil {
		zap.S().Errorw("Error decoding invoice", "error", err)
		s.ChannelMessageSend(m.ChannelID, "Invalid invoice")
		return
	}

	user, err := models.GetUser(m.Author.ID)
	if err != nil {
		return
	}

	if int64(invoice.MilliSat.ToSatoshis()) > user.Balance {
		s.ChannelMessageSend(m.ChannelID, "You don't have enough funds.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Withdrawing...")

	err = models.UpdateUserBalance(m.Author.ID, user.Balance-int64(invoice.MilliSat.ToSatoshis()))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		return
	}

	resp, err := lnclient.Client.SendPaymentSync(context.Background(), &lnrpc.SendRequest{
		PaymentRequest: args[0],
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		return
	}

	if resp.PaymentError != "" {
		s.ChannelMessageSend(m.ChannelID, "Error: "+resp.PaymentError)
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Withdrawal successful!")
}
