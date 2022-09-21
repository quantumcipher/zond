package pool

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/zond/api/view"
	"github.com/theQRL/zond/chain/block"
	"github.com/theQRL/zond/ntp"
	"github.com/theQRL/zond/protos"
)

func TestAdd(t *testing.T) {
	pool := CreateTransactionPool()
	tx := view.PlainTransferTransaction{}
	txn, _ := tx.ToTransferTransactionObject()
	txHash := txn.TxHash(txn.GetSigningHash())
	blockNumber := uint64(30)
	timestamp := ntp.GetNTP().Time()

	err := pool.Add(txn, txHash, blockNumber, timestamp)
	if err != nil {
		t.Error("got unexpected error while adding transaction to pool ", err)
	}
}

func TestAddTxFromBlock(t *testing.T) {
	networkId := uint64(1)
	timestamp := ntp.GetNTP().Time()
	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	slotNumber := uint64(120)
	parentHeaderHash := sha256.New().Sum([]byte("parentHeaderHash"))

	var txs []*protos.Transaction
	tx := view.PlainTransferTransaction{}
	txn1, _ := tx.ToTransferTransactionObject()
	txs = append(txs, txn1.PBData())
	fmt.Print(txs[0].Fee)
	protocolTxs := make([]*protos.ProtocolTransaction, 100)
	lastCoinBaseNonce := uint64(10)

	newBlock := block.NewBlock(networkId, timestamp, blockProposerPK[:], slotNumber, parentHeaderHash, txs, protocolTxs, lastCoinBaseNonce)

	pool := CreateTransactionPool()
	pool.AddTxFromBlock(newBlock, 30)

	if pool.txPool.Len() != 1 {
		t.Errorf("expected pool length after adding transaction from block to be 1, got (%v) ", pool.txPool.Len())
	}
}
