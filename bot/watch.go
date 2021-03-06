package bot

import (
	"context"
	"fmt"

	"github.com/aureleoules/lntip/lnclient"
	"github.com/aureleoules/lntip/models"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/lightningnetwork/lnd/lnrpc"
	"go.uber.org/zap"
)

func watchInvoices() {
	s, err := lnclient.Client.SubscribeInvoices(context.Background(), &lnrpc.InvoiceSubscription{})
	if err != nil {
		zap.S().Fatal(err)
	}

	for {
		invoice, err := s.Recv()
		if err != nil {
			zap.S().Fatal(err)
		}

		zap.S().Info(invoice)

		if invoice.State == lnrpc.Invoice_SETTLED {
			var userID string
			fmt.Sscanf(invoice.Memo, "LNTIP-%s", &userID)

			user, err := models.GetUser(userID)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			err = models.UpdateUserBalance(user.ID, user.Balance+invoice.AmtPaidSat)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			channel, err := discord.UserChannelCreate(userID)
			if err != nil {
				zap.S().Error(err)
				continue
			}

			discord.ChannelMessageSendEmbed(channel.ID, &discordgo.MessageEmbed{
				Title: "Deposit received",
				Color: 0xFFFF00,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "Amount",
						Value: humanize.Comma(invoice.AmtPaidSat) + " sats",
					},
				},
			})

			zap.S().Info("Invoice settled")
		}
	}
}
