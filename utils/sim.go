package utils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/catalyst"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

type Backend struct {
	node   *node.Node
	beacon *catalyst.SimulatedBeacon
	client *ethclient.Client
	key    *ecdsa.PrivateKey
	addr   *common.Address
}

func NewBackend() (*Backend, error) {
	nodeConf := node.DefaultConfig
	nodeConf.DataDir = ""
	nodeConf.P2P = p2p.Config{NoDiscovery: true}

	ethConf := ethconfig.Defaults
	ethConf.Genesis = &core.Genesis{
		Config:   params.AllDevChainProtocolChanges,
		GasLimit: ethconfig.Defaults.Miner.GasCeil,
	}
	ethConf.SyncMode = ethconfig.FullSync
	ethConf.TxPool.NoLocals = true

	stack, err := node.New(&nodeConf)
	if err != nil {
		return nil, err
	}
	backend, err := newWithNode(stack, &ethConf, 0)
	if err != nil {
		return nil, err
	}
	return backend, nil
}

func newWithNode(stack *node.Node, conf *eth.Config, blockPeriod uint64) (*Backend, error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	address := crypto.PubkeyToAddress(key.PublicKey)
	conf.Genesis.Alloc = map[common.Address]types.Account{
		address: {Balance: new(big.Int).Exp(big.NewInt(10), big.NewInt(21), nil)}, // 1000 ETH
	}

	backend, err := eth.New(stack, conf)
	if err != nil {
		return nil, err
	}

	filterSystem := filters.NewFilterSystem(backend.APIBackend, filters.Config{})
	stack.RegisterAPIs([]rpc.API{{
		Namespace: "eth",
		Service:   filters.NewFilterAPI(filterSystem),
	}})

	if err := stack.Start(); err != nil {
		return nil, err
	}

	beacon, err := catalyst.NewSimulatedBeacon(blockPeriod, common.Address{}, backend)
	if err != nil {
		return nil, err
	}
	if err := beacon.Fork(backend.BlockChain().GetCanonicalHash(0)); err != nil {
		return nil, err
	}
	return &Backend{
		node:   stack,
		beacon: beacon,
		key:    key,
		addr:   &address,
		client: ethclient.NewClient(stack.Attach()),
	}, nil
}

func (n *Backend) Close() error {
	if n.client != nil {
		n.client.Close()
		n.client = nil
	}
	var err error
	if n.beacon != nil {
		err = n.beacon.Stop()
		n.beacon = nil
	}
	if n.node != nil {
		err = errors.Join(err, n.node.Close())
		n.node = nil
	}
	return err
}

func (n *Backend) Commit() common.Hash {
	return n.beacon.Commit()
}

func (n *Backend) Client() *ethclient.Client {
	return n.client
}

func (n *Backend) FaucetKey() string {
	k := crypto.FromECDSA(n.key)
	enc := make([]byte, len(k)*2)
	hex.Encode(enc, k)
	return string(enc)
}

func (n *Backend) FaucetAddr() string {
	return n.addr.Hex()
}
