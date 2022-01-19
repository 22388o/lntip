package models

import (
	"database/sql"
	"time"

	"github.com/aureleoules/lntip/db"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID      string `json:"id" db:"id"`
	Balance int64  `json:"balance" db:"balance"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) Create() error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)
	return u.create(tx)
}

func (u *User) create(tx *sqlx.Tx) error {
	q, args, err := builder.Insert("users").
		Columns("id", "balance").
		Values(u.ID, u.Balance).
		ToSql()

	if err != nil {
		return err
	}

	_, err = tx.Exec(q, args...)
	return err
}

func GetUser(id string) (*User, error) {
	var u User

	q, args, err := builder.Select("*").
		From("users").
		Where("id = ?", id).
		Limit(1).
		ToSql()

	if err != nil {
		return nil, err
	}

	err = db.DB.Get(&u, q, args...)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func CreateUserIfNoExists(id string) (*User, error) {
	u, err := GetUser(id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if u != nil {
		return u, nil
	}

	u = &User{
		ID:      id,
		Balance: 0,
	}

	return u, u.Create()
}

func UpdateUserBalance(id string, balance int64) error {
	tx, err := db.DB.Beginx()
	if err != nil {
		return err
	}

	defer checkErr(tx, err)
	return updateUserBalance(tx, id, balance)
}

func updateUserBalance(tx *sqlx.Tx, id string, balance int64) error {
	q, args, err := builder.Update("users").
		Set("balance", balance).
		Where("id = ?", id).
		ToSql()

	if err != nil {
		return err
	}

	_, err = tx.Exec(q, args...)
	return err
}
