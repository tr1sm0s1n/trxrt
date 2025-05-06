package api

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
	"github.com/tr1sm0s1n/trxrt/helpers"
)

func SetCodeTx(client *ethclient.Client, key, contract, authKey string, gas, maxFee, maxPriorityFee float64) error {
	log.Println("\033[32m>>> Type 0x4 Transaction: BEGIN <<<\033[0m")

	pkey, from, err := helpers.InitializeAccount(key)
	if err != nil {
		log.Println("Failed to generate key:", err)
		return err
	}

	signerKey, signerAddr, err := helpers.InitializeAccount(authKey)
	if err != nil {
		log.Println("Failed to generate auth signer key:", err)
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

	signerNonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		log.Println("Failed to fetch auth signer nonce:", err)
		return err
	}

	auth := types.SetCodeAuthorization{
		ChainID: *uint256.NewInt(chainID.Uint64()),
		Address: common.HexToAddress(contract),
		Nonce:   signerNonce,
	}

	auth, err = types.SignSetCode(signerKey, auth)
	if err != nil {
		log.Println("Failed to sign the authorization:", err)
		return err
	}

	signedTx, err := types.SignNewTx(pkey, types.LatestSignerForChainID(chainID), &types.SetCodeTx{
		Nonce:     nonce,
		To:        signerAddr,
		GasTipCap: uint256.NewInt((uint64(maxPriorityFee))),
		GasFeeCap: uint256.NewInt((uint64(maxFee * 1000000000))),
		Gas:       uint64(gas),
		Data:      nil,
		AuthList:  []types.SetCodeAuthorization{auth},
	})
	if err != nil {
		log.Println("Failed to create trx:", err)
		return err
	}

	if err = client.SendTransaction(context.Background(), signedTx); err != nil {
		log.Println("Failed to send trx:", err)
		return err
	}

	if err = helpers.WaitForReceipt(client, signedTx); err != nil {
		log.Println("Failed to generate receipt:", err)
		return err
	}

	log.Println("\033[32m>>> Type 0x4 Transaction: END <<<\033[0m")
	return nil
}
