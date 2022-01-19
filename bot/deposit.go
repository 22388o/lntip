package bot

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/aureleoules/lntip/lnclient"
	"github.com/bwmarrin/discordgo"
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

	s.ChannelFileSendWithMessage(m.ChannelID, "Deposit: "+invoice.PaymentRequest, "qrcode.png", b)
}
