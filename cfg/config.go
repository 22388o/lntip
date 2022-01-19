package cfg

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config contains the globally accessable configuration of the app
var Config ServerConfig

// ServerConfig contains various configuration parameters
type ServerConfig struct {
	Bot struct {
		Token string `env:"BOT_TOKEN"`
	}
	Database struct {
		Name     string `yaml:"name" env:"DATABASE_NAME" env-default:"lnswap"`
		Host     string `yaml:"host" env:"DATABASE_HOST" env-default:"127.0.0.1"`
		Port     string `yaml:"port" env:"DATABASE_PORT" env-default:"3306"`
		User     string `yaml:"user" env:"DATABASE_USER" env-default:"root"`
		Password string `yaml:"password" env:"DATABASE_PASS" env-default:""`
	} `yaml:"database"`
	LND struct {
		Host         string `yaml:"host" env:"LND_HOST" env-default:"127.0.0.1"`
		TLSPath      string `yaml:"tls_path" env:"LND_TLS_PATH" env-default:"tls"`
		MacaroonPath string `yaml:"macaroon_path" env:"LND_MACAROON_PATH" env-default:"macaroon"`
		Network      string `yaml:"network" env:"LND_NETWORK" env-default:"regtest"`
	} `yaml:"lnd"`
}

// Load config from a config file
func Load(configFile string) error {
	err := cleanenv.ReadConfig(configFile, &Config)
	if err != nil {
		fmt.Println(err)
	}

	return cleanenv.ReadEnv(&Config)
}

func ChainParams() *chaincfg.Params {
	if Config.LND.Network == "mainnet" {
		return &chaincfg.MainNetParams
	}

	if Config.LND.Network == "testnet" {
		return &chaincfg.TestNet3Params
	}

	if Config.LND.Network == "regtest" {
		return &chaincfg.RegressionNetParams
	}

	return nil
}
