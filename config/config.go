package config

import (
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

func DialClient() (*ethclient.Client, error) {
	eth, err := ethclient.Dial(os.Getenv("RPC_URL"))
	if err != nil {
		return nil, err
	}

	return eth, nil
}
