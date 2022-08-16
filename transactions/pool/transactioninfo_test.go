package pool

import (
	"encoding/hex"
	"testing"

	"github.com/theQRL/zond/api/view"
	"github.com/theQRL/zond/ntp"
)

func TestCreateTransactionInfo(t *testing.T) {
	tx := view.PlainTransferTransaction{}
	txn, _ := tx.ToTransferTransactionObject()
	txHash := txn.TxHash(txn.GetSigningHash())
	blockNumber := uint64(10)
	timestamp := ntp.GetNTP().Time()

	transactionInfo := CreateTransactionInfo(txn, txHash, blockNumber, timestamp)

	if string(transactionInfo.txHash) != string(txHash) {
		t.Errorf("expected transaction in create transactionInfo (%v), got (%v)", hex.EncodeToString(txHash), hex.EncodeToString(transactionInfo.TxHash()))
	}
}

// func TestCheckOTSExist(t *testing.T) {
// 	tx := view.PlainTransferTransaction{}
// 	txn, _ := tx.ToTransferTransactionObject()
// 	txHash := txn.TxHash(txn.GetSigningHash())
// 	blockNumber := uint64(10)
// 	timestamp := ntp.GetNTP().Time()

// 	tx2 := view.PlainTransferTransaction{}
// 	txn2, _ := tx2.ToTransferTransactionObject()

// 	transactionInfo := CreateTransactionInfo(txn, txHash, blockNumber, timestamp)
// 	output := transactionInfo.CheckOTSExist(txn2)
// 	if output {
// 		t.Errorf("expected output of checkots exist to be (%v), got (%v)", false, output)
// 	}
// }

func TestIsStale(t *testing.T) {
	tx := view.PlainTransferTransaction{}
	txn, _ := tx.ToTransferTransactionObject()
	txHash := txn.TxHash(txn.GetSigningHash())
	blockNumber := uint64(30)
	timestamp := ntp.GetNTP().Time()

	transactionInfo := CreateTransactionInfo(txn, txHash, blockNumber, timestamp)

	testCases := []struct {
		name               string
		currentBlockHeight uint64
		TransactionInfo    *TransactionInfo
		expectedOutput     bool
	}{
		{
			name:               "ok",
			currentBlockHeight: 46,
			TransactionInfo:    transactionInfo,
			expectedOutput:     true,
		},
		{
			name:               "not stale",
			currentBlockHeight: 11,
			TransactionInfo:    transactionInfo,
			expectedOutput:     false,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			output := tc.TransactionInfo.IsStale(tc.currentBlockHeight)
			if output != tc.expectedOutput {
				t.Errorf("expected output of isstale function to be (%v) but returned (%v)", tc.expectedOutput, output)
			}

		})
	}
}

func TestUpdateBlockNumber(t *testing.T) {
	tx := view.PlainTransferTransaction{}
	txn, _ := tx.ToTransferTransactionObject()
	txHash := txn.TxHash(txn.GetSigningHash())
	blockNumber := uint64(30)
	timestamp := ntp.GetNTP().Time()

	transactionInfo := CreateTransactionInfo(txn, txHash, blockNumber, timestamp)

	transactionInfo.UpdateBlockNumber(40)
	output := transactionInfo.BlockNumber()
	if output != uint64(40) {
		t.Errorf("unable to update block number, expected (%v), got (%v)", 40, output)
	}
}
