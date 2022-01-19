package db

import (
	"errors"
	"fmt"
	"time"

	"github.com/aureleoules/lntip/cfg"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const maxTries = 10

var DB *sqlx.DB

// Open DB
func Open() error {
	zap.S().Info("Connecting to database...")
	connInfo := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4,utf8&parseTime=true&multiStatements=true", cfg.Config.Database.User, cfg.Config.Database.Password, cfg.Config.Database.Host, cfg.Config.Database.Port, cfg.Config.Database.Name)
	fmt.Println(connInfo)
	for i := 0; i < maxTries; i++ {
		db, err := sqlx.Connect("mysql", connInfo)
		if err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			continue
		}
		zap.S().Info("Connected to database.")
		DB = db
		return nil
	}

	return errors.New("could not connect to db")
}
