package api

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/holiman/uint256"
	"github.com/tr1sm0s1n/trxrt/helpers"
)

func BlobTx(client *ethclient.Client, key, to string, data []string, gas, maxFee, maxPriorityFee, blobFee float64) error {
	log.Println("\033[32m>>> Type 0x3 Transaction: BEGIN <<<\033[0m")

	var sidecar types.BlobTxSidecar
	for i := range min(len(data), 6) {
		blob := new(kzg4844.Blob)
		copy(blob[:], []byte(data[i]))

		blobCommit, _ := kzg4844.BlobToCommitment(blob)
		blobProof, _ := kzg4844.ComputeBlobProof(blob, blobCommit)

		sidecar.Blobs = append(sidecar.Blobs, *blob)
		sidecar.Commitments = append(sidecar.Commitments, blobCommit)
		sidecar.Proofs = append(sidecar.Proofs, blobProof)
	}

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
	signedTx, _ := types.SignNewTx(pkey, types.LatestSignerForChainID(chainID), &types.BlobTx{
		Nonce:      nonce,
		To:         receiver,
		GasTipCap:  uint256.NewInt(uint64(maxPriorityFee)),
		GasFeeCap:  uint256.NewInt(uint64(maxFee * 1000000000)),
		Gas:        uint64(gas),
		BlobFeeCap: uint256.NewInt(uint64(blobFee)),
		BlobHashes: sidecar.BlobHashes(),
		Sidecar:    &sidecar,
	})

	if err = client.SendTransaction(context.Background(), signedTx); err != nil {
		log.Println("Failed to send trx:", err)
		return err
	}

	if err = helpers.WaitForReceipt(client, signedTx); err != nil {
		log.Println("Failed to generate receipt:", err)
		return err
	}

	btx, _, err := client.TransactionByHash(context.Background(), signedTx.Hash())
	if err != nil {
		log.Fatal("Failed to fetch trx:", err)
	}

	if btx.BlobHashes()[0] != sidecar.BlobHashes()[0] {
		log.Fatal("Failed to verify blob hashes")
	}

	log.Println("\033[32m>>> Type 0x3 Transaction: END <<<\033[0m")
	return nil
}
