package api

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tr1sm0s1n/project-wallet-x/helpers"
)

func DynamicFeeTx(client *ethclient.Client, key, to string, amount int64, gas, maxFee, maxPriorityFee float64) error {
	log.Println("\033[32m>>> Type 0x2 Transaction: BEGIN <<<\033[0m")

	pkey, from, err := helpers.InitializeAccount(key)
	if err != nil {
		log.Println("Failed to generate key:", err)
		return err
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("Failed to retrieve the chain ID:", err)
		return err
	}

	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		log.Println("Failed to fetch nonce:", err)
		return err
	}

	receiver := common.HexToAddress(to)
	signedTx, _ := types.SignNewTx(pkey, types.LatestSignerForChainID(chainID), &types.DynamicFeeTx{
		Nonce:     nonce,
		To:        &receiver,
		GasTipCap: big.NewInt(int64(maxPriorityFee)),
		GasFeeCap: big.NewInt(int64(maxFee * 1000000000)),
		Gas:       uint64(gas),
		Value:     new(big.Int).Mul(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), big.NewInt(amount)),
		Data:      nil,
	})

	if err = client.SendTransaction(context.Background(), signedTx); err != nil {
		log.Println("Failed to send trx:", err)
		return err
	}

	if err = helpers.WaitForReceipt(client, signedTx); err != nil {
		log.Println("Failed to generate receipt:", err)
		return err
	}

	br, _ := client.BalanceAt(context.Background(), receiver, nil)
	bs, _ := client.BalanceAt(context.Background(), from, nil)
	log.Println("Balance of receiver:", br)
	log.Println("Balance of sender:", bs)

	log.Println("\033[32m>>> Type 0x2 Transaction: END <<<\033[0m")
	return nil
}
