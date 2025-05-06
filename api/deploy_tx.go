package api

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tr1sm0s1n/trxrt/helpers"
)

func DeployTx(client *ethclient.Client, key, bytecode string) (common.Address, error) {
	log.Println("\033[32m>>> Deployment Transaction: BEGIN <<<\033[0m")

	pkey, _, err := helpers.InitializeAccount(key)
	if err != nil {
		log.Println("Failed to generate key:", err)
		return common.Address{}, err
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Println("Failed to retrieve the chain ID:", err)
		return common.Address{}, err
	}

	auth := bind.NewKeyedTransactor(pkey, chainID)

	addr, deployTx, err := bind.DeployContract(auth, common.FromHex(bytecode), client, nil)
	if err != nil {
		log.Println("Failed to deploy contract:", err)
		return common.Address{}, err
	}

	if err = helpers.WaitForReceipt(client, deployTx); err != nil {
		log.Println("Failed to generate receipt:", err)
		return common.Address{}, err
	}

	log.Println("\033[32m>>> Deployment Transaction: END <<<\033[0m")
	return addr, nil
}
