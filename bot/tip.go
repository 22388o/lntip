package bot

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aureleoules/lntip/models"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

const emojiPrefix = "lntip"

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
		return
	}

	if tipAmount < 1 {
		return
	}

	tip := models.Tip{
		UserID:    r.UserID,
		GuildID:   r.GuildID,
		ChannelID: r.ChannelID,
		Amount:    int64(tipAmount),
		MessageID: &r.MessageID,
		IsAward:   true,
	}

	err = tip.Create()
	if err != nil {
		zap.S().Errorw("Failed to create tip", "error", err)
	}
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

	_, err = models.CreateUserIfNoExists(tip.UserID)
	if err != nil {
		zap.S().Errorw("Failed to create user", "error", err)
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

	_, err = models.CreateUserIfNoExists(tip.ToUserID)
	if err != nil {
		zap.S().Errorw("Failed to create user", "error", err)
		return
	}

	err = tip.Create()
	if err != nil {
		zap.S().Errorw("Failed to create tip", "error", err)
	}

	channel, err := discord.UserChannelCreate(userID)
	if err != nil {
		zap.S().Error(err)
		return
	}

	discord.ChannelMessageSend(channel.ID, fmt.Sprintf("You've just been tipped %d sats by %s!", tipAmount, m.Author.Username))
	zap.S().Info("Tip success")
}
