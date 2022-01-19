package lnclient

import (
	"github.com/aureleoules/lntip/cfg"
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc"
	"go.uber.org/zap"
)

var Client lnrpc.LightningClient

func Init() {
	var err error
	Client, err = lndclient.NewBasicClient(cfg.Config.LND.Host, cfg.Config.LND.TLSPath, cfg.Config.LND.MacaroonPath, "mainnet")
	if err != nil {
		zap.S().Fatal(err)
	}

	zap.S().Info("Connected to LND")

}
