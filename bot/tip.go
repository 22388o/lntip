package bot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aureleoules/lntip/models"
	"github.com/aureleoules/lntip/rates"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"go.uber.org/zap"
)

const emojiPrefix = "lntip"

const rewardLimit = 2000

func reactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(r.Emoji.Name, "lntip") {
		return
	}

	n := strings.Index(r.Emoji.Name, "lntip")
	tipAmount, err := strconv.Atoi(r.Emoji.Name[n+len("lntip"):])
	if err != nil {
		zap.S().Warnw("Failed to parse tip amount", "error", err)
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
		return
	}

	if tipAmount < 1 {
		return
	}

	ok, err := models.HasTipped(r.UserID, r.MessageID, tipAmount)
	if err != nil {
		zap.S().Errorw("Failed to check if user has tipped", "error", err)
		return
	}

	if ok {
		zap.S().Infow("User has already tipped", "user", r.UserID, "message", r.MessageID, "amount", tipAmount)
		return
	}

	user, err := models.GetUser(r.UserID)
	if err != nil {
		zap.S().Errorw("Failed to get user", "error", err)
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
		return
	}

	channel, err := discord.UserChannelCreate(r.UserID)
	if err != nil {
		zap.S().Error(err)
		return
	}
	msg, err := discord.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		zap.S().Errorw("Failed to get message", "error", err)
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
		return
	}

	if msg.Author.ID == r.UserID {
		zap.S().Infow("User is author of message", "user", r.UserID, "message", r.MessageID)
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
		return
	}

	if user.Balance < int64(tipAmount) {
		s.ChannelMessageSendEmbed(channel.ID, &discordgo.MessageEmbed{
			Title: "Not enough sats",
			Description: fmt.Sprintf("You don't have enough sats to tip %s sats. You have %s sats.",
				humanize.Comma(int64(tipAmount)), humanize.Comma(user.Balance)),
			Color: 0xFF0000,
		})

		err = s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
		if err != nil {
			zap.S().Errorw("Failed to remove reaction", "error", err)
		}

		return
	}

	if tipAmount > rewardLimit {
		s.ChannelMessageSendEmbed(channel.ID, &discordgo.MessageEmbed{
			Title: "Tip too large",
			Description: fmt.Sprintf("You can't tip more than %s sats for now.",
				humanize.Comma(int64(rewardLimit))),
			Color: 0xFF0000,
		})
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
		return
	}

	tip := models.Tip{
		UserID:    r.UserID,
		ToUserID:  msg.Author.ID,
		GuildID:   r.GuildID,
		ChannelID: r.ChannelID,
		Amount:    int64(tipAmount),
		MessageID: &r.MessageID,
		IsAward:   true,
	}

	dUser, err := discord.User(r.UserID)
	if err != nil {
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
		zap.S().Error(err)
		return
	}
	zap.S().Info("Sending tip...")

	sendTip(s, &tip, r.ChannelID, dUser.Username)
}

func lntipHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		return
	}

	var userID string
	reg := regexp.MustCompile(`<@!?(\d+)>`)
	if reg.MatchString(args[0]) {
		userID = reg.FindStringSubmatch(args[0])[1]
	}

	if userID == "" {
		return
	}

	tipAmount, err := strconv.Atoi(args[1])
	if err != nil {
		return
	}

	if tipAmount < 1 {
		return
	}

	user, err := models.GetUser(m.Author.ID)
	if err != nil {
		zap.S().Errorw("Failed to get user", "error", err)
		return
	}

	if user.Balance < int64(tipAmount) {
		s.ChannelMessageSend(m.ChannelID, "You don't have enough to tip that much.")
		return
	}

	tip := models.Tip{
		UserID:    m.Author.ID,
		ToUserID:  userID,
		GuildID:   m.GuildID,
		ChannelID: m.ChannelID,
		Amount:    int64(tipAmount),
	}

	sendTip(s, &tip, m.ChannelID, m.Author.Username)
}

func sendTip(s *discordgo.Session, tip *models.Tip, channelID string, fromUsername string) error {
	_, err := models.CreateUserIfNoExists(tip.UserID)
	if err != nil {
		zap.S().Errorw("Failed to create user", "error", err)
		return err
	}

	_, err = models.CreateUserIfNoExists(tip.ToUserID)
	if err != nil {
		zap.S().Errorw("Failed to create user", "error", err)
		return err
	}

	err = tip.Create()
	if err != nil {
		zap.S().Errorw("Failed to create tip", "error", err)
		return err
	}

	channel, err := discord.UserChannelCreate(tip.ToUserID)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	title := "Tip"
	if tip.MessageID != nil {
		title = "Reward"
	}

	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Amount",
			Value:  humanize.Comma(tip.Amount) + " sats",
			Inline: true,
		},
		{
			Name:   "Value",
			Value:  fmt.Sprintf("%f €", (float64(tip.Amount)/1e8)*float64(rates.Price)),
			Inline: true,
		},
		{
			Name:  "From",
			Value: fromUsername,
		},
	}

	if tip.MessageID != nil {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Original message",
			Value: fmt.Sprintf("[Click](https://discord.com/channels/%s/%s/%s)", tip.GuildID, tip.ChannelID, *tip.MessageID),
		})
	}

	discord.ChannelMessageSendEmbed(channel.ID, &discordgo.MessageEmbed{
		Title:  title,
		Fields: fields,

		Description: fmt.Sprintf("You've just received a tip!"),
		Color:       0xFFFF00,
	})

	zap.S().Info("Tip success")

	return nil
}
