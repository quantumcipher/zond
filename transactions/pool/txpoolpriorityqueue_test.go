package pool

import (
	"testing"

	"github.com/theQRL/zond/api/view"
	"github.com/theQRL/zond/ntp"
)

func TestPush(t *testing.T) {
	tx := view.PlainTransferTransaction{}
	txn, _ := tx.ToTransferTransactionObject()
	txHash := txn.Hash()
	blockNumber := uint64(30)
	timestamp := ntp.GetNTP().Time()

	transactionInfo := CreateTransactionInfo(txn, txHash, blockNumber, timestamp)

	var queue PriorityQueue
	queue.Push(transactionInfo)

	if len(queue) != 1 {
		t.Error("unable to push data to priority queue")
	}
}

func TestPop(t *testing.T) {
	tx := view.PlainTransferTransaction{}
	txn, _ := tx.ToTransferTransactionObject()
	txHash := txn.Hash()
	blockNumber := uint64(30)
	timestamp := ntp.GetNTP().Time()

	transactionInfo := CreateTransactionInfo(txn, txHash, blockNumber, timestamp)

	var queue PriorityQueue
	queue.Push(transactionInfo)
	queue.Pop()

	if len(queue) != 0 {
		t.Error("unable to pop data from priority queue")
	}
}

func TestRemove(t *testing.T) {
	tx := view.PlainTransferTransaction{}
	txn, _ := tx.ToTransferTransactionObject()
	txHash := txn.Hash()
	blockNumber := uint64(30)
	timestamp := ntp.GetNTP().Time()

	transactionInfo := CreateTransactionInfo(txn, txHash, blockNumber, timestamp)

	var queue PriorityQueue
	queue.Push(transactionInfo)
	queue.Remove(transactionInfo.tx, transactionInfo.txHash)

	if len(queue) != 0 {
		t.Error("unable to remove data from priority queue")
	}
}

func TestRemoveByIndex(t *testing.T) {
	tx := view.PlainTransferTransaction{}
	txn, _ := tx.ToTransferTransactionObject()
	txHash := txn.Hash()
	blockNumber := uint64(30)
	timestamp := ntp.GetNTP().Time()

	transactionInfo := CreateTransactionInfo(txn, txHash, blockNumber, timestamp)

	var queue PriorityQueue
	queue.Push(transactionInfo)
	queue.removeByIndex(0)

	if len(queue) != 0 {
		t.Error("unable to remove data from priority queue")
	}
}

func TestContains(t *testing.T) {
	tx := view.PlainTransferTransaction{}
	txn, _ := tx.ToTransferTransactionObject()
	txHash := txn.Hash()
	blockNumber := uint64(30)
	timestamp := ntp.GetNTP().Time()

	transactionInfo := CreateTransactionInfo(txn, txHash, blockNumber, timestamp)

	var queue PriorityQueue
	queue.Push(transactionInfo)

	if !queue.Contains(transactionInfo) {
		t.Error("transaction exists in priority queue but not detected")
	}
}
