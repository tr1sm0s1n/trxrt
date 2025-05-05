package api

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tr1sm0s1n/trxrt/helpers"
)

func DeployTx(client *ethclient.Client, key, bytecode string, amount int64, gas, maxFee, maxPriorityFee float64) (common.Address, error) {
	log.Println("\033[32m>>> Deployment Transaction: BEGIN <<<\033[0m")

	pkey, from, err := helpers.InitializeAccount(key)
	if err != nil {
		log.Println("Failed to generate key:", err)
		return common.Address{}, err
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal("Failed to retrieve the chain ID:", err)
		return common.Address{}, err
	}

	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		log.Println("Failed to fetch nonce:", err)
		return common.Address{}, err
	}

	code, err := hex.DecodeString(bytecode[2:])
	if err != nil {
		log.Println("Failed to decode bytecode:", err)
		return common.Address{}, err
	}

	signedTx, err := types.SignNewTx(pkey, types.LatestSignerForChainID(chainID), &types.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: big.NewInt(int64(maxPriorityFee)),
		GasFeeCap: big.NewInt(int64(maxFee * 1000000000)),
		Gas:       uint64(gas),
		Value:     new(big.Int).Mul(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil), big.NewInt(amount)),
		Data:      code,
	})

	if err != nil {
		log.Println("Failed to create/sign a transaction:", err)
		return common.Address{}, err
	}

	if err = client.SendTransaction(context.Background(), signedTx); err != nil {
		log.Println("Failed to send trx:", err)
		return common.Address{}, err
	}

	if err = helpers.WaitForReceipt(client, signedTx); err != nil {
		log.Println("Failed to generate receipt:", err)
		return common.Address{}, err
	}

	log.Println("\033[32m>>> Deployment Transaction: END <<<\033[0m")
	return crypto.CreateAddress(from, signedTx.Nonce()), nil
}
