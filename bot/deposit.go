package bot

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/aureleoules/lntip/lnclient"
	"github.com/aureleoules/lntip/models"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"go.uber.org/zap"
)

type bufferAdaptor struct {
	*bytes.Buffer
}

func (b bufferAdaptor) Close() error {
	return nil
}

func (b bufferAdaptor) Write(p []byte) (int, error) {
	return b.Buffer.Write(p)
}

func depositHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 1 {
		return
	}

	amount, err := strconv.Atoi(args[0])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Amount must be a number.")
		return
	}

	if amount < 1 {
		s.ChannelMessageSend(m.ChannelID, "Amount must be greater than 0.")
		return
	}

	_, err = models.CreateUserIfNoExists(m.Author.ID)
	if err != nil {
		zap.S().Errorw("create user failed", "err", err)
		return
	}

	invoice, err := lnclient.Client.AddInvoice(context.Background(), &lnrpc.Invoice{
		Memo:  "LNTIP-" + m.Author.ID,
		Value: int64(amount),
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: "+err.Error())
		return
	}

	qrc, err := qrcode.New(invoice.PaymentRequest)
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return
	}

	b := bufferAdaptor{Buffer: bytes.NewBuffer(nil)}
	w := standard.NewWithWriter(b)

	err = qrc.Save(w)
	if err != nil {
		zap.S().Errorf("could not generate QRCode: %v", err)
		return
	}

	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Deposit",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Amount",
					Value: fmt.Sprintf("%s sats", humanize.Comma(int64(amount))),
				},
				{
					Name:  "Invoice",
					Value: invoice.PaymentRequest,
				},
			},
			Description: "Scan this QR code with your mobile wallet to deposit funds.",
			Color:       0xFFFF00,
			Type:        "rich",
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://qrcode.png",
			},
		},
		Files: []*discordgo.File{
			{
				Reader:      b,
				Name:        "qrcode.png",
				ContentType: "image/png",
			},
		},
	})

	if err != nil {
		zap.S().Errorf("could not send message: %v", err)
		return
	}
}
