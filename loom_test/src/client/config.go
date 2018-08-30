package client

import (
	"ethcontract"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/loomnetwork/go-loom/client/plasma_cash"
	loom_ethcontract "github.com/loomnetwork/go-loom/client/plasma_cash/eth/ethcontract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

// Loads plasma-config.yml or equivalent from the cwd
func parseConfig() (*viper.Viper, error) {
	// When running "go test" the cwd is set to the package dir, not the root dir
	// where the config is, so gotta do a bit more work to figure out the config dir...
	_, filename, _, _ := runtime.Caller(0)
	cfgDir := filepath.Join(filepath.Dir(filename), "../..")

	v := viper.New()
	v.SetConfigName("plasma-config")
	v.AddConfigPath(cfgDir)
	return v, v.ReadInConfig()
}

func GetTestAccountHexKey(name string) string {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}
	return cfg.GetString(name)
}

func GetContractHexAddress(name string) string {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}
	return cfg.GetString(name)
}

func GetTokenContract(name string) plasma_cash.TokenContract {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}
	tokenAddr := common.HexToAddress(cfg.GetString("token_contract"))
	privKeyHexStr := cfg.GetString(name)
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(privKeyHexStr, "0x"))
	if err != nil {
		log.Fatalf("failed to load private key for %s: %v", name, err)
	}
	tokenContract, err := ethcontract.NewCards(tokenAddr, conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	return NewTokenContract(name, privKey, tokenContract)
}

func GetRootChain(name string) plasma_cash.RootChainClient {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}
	contractAddr := common.HexToAddress(cfg.GetString("root_chain"))
	privKeyHexStr := cfg.GetString(name)
	privKey, err := crypto.HexToECDSA(strings.TrimPrefix(privKeyHexStr, "0x"))
	if err != nil {
		log.Fatalf("failed to load private key for %s: %v", name, err)
	}
	plasmaContract, err := loom_ethcontract.NewRootChain(contractAddr, conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	return NewRootChainService(name, privKey, plasmaContract)
}
