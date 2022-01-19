package bot

import (
	"errors"
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

func reactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	fmt.Println(r.Emoji.Name)
	zap.S().Info("Reaction added")
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
		return
	}

	if tipAmount < 1 {
		return
	}

	msg, err := discord.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		zap.S().Errorw("Failed to get message", "error", err)
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

	user, err := models.GetUser(tip.UserID)
	if err != nil {
		zap.S().Errorw("Failed to get user", "error", err)
		return err
	}

	if user.Balance < tip.Amount {
		s.ChannelMessageSend(channelID, "You don't have enough to tip that much.")
		return errors.New("Not enough balance")
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

	discord.ChannelMessageSend(channel.ID, fmt.Sprintf("You've just been tipped %s sats by %s! (%f €)", humanize.Comma(tip.Amount), fromUsername, (float64(tip.Amount)/1e8)*float64(rates.Price)))
	zap.S().Info("Tip success")

	return nil
}
