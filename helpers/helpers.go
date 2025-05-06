package helpers

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func InitializeAccount(key string) (*ecdsa.PrivateKey, common.Address, error) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		return nil, common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, errors.New("error casting public key to ECDSA")
	}

	return privateKey, crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

func WaitForReceipt(client *ethclient.Client, trx *types.Transaction) error {
	for {
		r, err := client.TransactionReceipt(context.Background(), trx.Hash())
		if err != nil {
			if err == ethereum.NotFound {
				log.Println("Receipt isn't available")
				time.Sleep(5 * time.Second)
				continue
			}
			return err
		}

		if r.Status == types.ReceiptStatusSuccessful {
			log.Println("Transaction has been committed!!")
			log.Printf("Transaction Hash: \033[35m%s\033[0m\n", r.TxHash)
			break
		}

		log.Println("Transaction execution failed")
		return errors.New("failed to execute")
	}
	return nil
}
