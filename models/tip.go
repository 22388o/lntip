package models

import (
	"github.com/aureleoules/lntip/db"
	"github.com/jmoiron/sqlx"
)

type Tip struct {
	UserID    string  `json:"user_id" db:"user_id"`
	ToUserID  string  `json:"to_user_id" db:"to_user_id"`
	MessageID *string `json:"message_id" db:"message_id"`
	GuildID   string  `json:"guild_id" db:"guild_id"`
	ChannelID string  `json:"channel_id" db:"channel_id"`
	Amount    int64   `json:"amount" db:"amount"`
	CreatedAt int64   `json:"created_at" db:"created_at"`
	IsAward   bool    `json:"is_award" db:"is_award"`
}

func (u *Tip) Create() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)
	return u.create(tx)
}

func (t *Tip) create(tx *sqlx.Tx) error {
	q, args, err := builder.Insert("tips").
		Columns("user_id", "to_user_id", "message_id", "guild_id", "channel_id", "amount", "is_award").
		Values(t.UserID, t.ToUserID, t.MessageID, t.GuildID, t.ChannelID, t.Amount, t.IsAward).
		ToSql()

	if err != nil {
		return err
	}

	_, err = tx.Exec(q, args...)
	return err
}

func HasTipped(userID, messageID string, amount int) (bool, error) {
	q, args, err := builder.Select("count(*)").
		From("tips").
		Where("user_id = ? AND message_id = ? AND amount = ?", userID, messageID, amount).
		Limit(1).
		ToSql()

	if err != nil {
		return false, err
	}

	var count int
	err = db.DB.Get(&count, q, args...)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
