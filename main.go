package main

import (
	"github.com/aureleoules/lntip/bot"
	"github.com/aureleoules/lntip/cfg"
	"github.com/aureleoules/lntip/db"
	"github.com/aureleoules/lntip/lnclient"
	"github.com/aureleoules/lntip/rates"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	c := zap.NewDevelopmentConfig()
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ = c.Build()
	zap.ReplaceGlobals(logger)
}

func main() {
	cfg.Load("config.yml")
	db.Open()
	lnclient.Init()

	go rates.Run()

	bot.Run()
}
