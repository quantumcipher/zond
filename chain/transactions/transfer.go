package transactions

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/theQRL/go-qrllib/xmss"
	"github.com/theQRL/zond/misc"
	"github.com/theQRL/zond/protos"
	"github.com/theQRL/zond/state"

	log "github.com/sirupsen/logrus"
)

type Transfer struct {
	Transaction
}

func (tx *Transfer) AddrsTo() [][]byte {
	return tx.pbData.GetTransfer().AddrsTo
}

func (tx *Transfer) Amounts() []uint64 {
	return tx.pbData.GetTransfer().Amounts
}

func (tx *Transfer) SlavePKs() [][]byte {
	return tx.pbData.GetTransfer().SlavePks
}

func (tx *Transfer) Message() []byte {
	return tx.pbData.GetTransfer().Message
}

func (tx *Transfer) TotalAmounts() uint64 {
	totalAmount := uint64(0)
	for _, amount := range tx.Amounts() {
		totalAmount += amount
	}
	return totalAmount
}

func (tx *Transfer) GetSigningHash() []byte {
	tmp := new(bytes.Buffer)
	binary.Write(tmp, binary.BigEndian, tx.NetworkID())
	tmp.Write(tx.MasterAddr())
	binary.Write(tmp, binary.BigEndian, tx.Fee())
	binary.Write(tmp, binary.BigEndian, tx.Nonce())

	for i := range tx.AddrsTo() {
		tmp.Write(tx.AddrsTo()[i])
		binary.Write(tmp, binary.BigEndian, tx.Amounts()[i])
	}

	for _, slavePK := range tx.SlavePKs() {
		tmp.Write(slavePK)
	}

	tmp.Write(tx.Message())

	h := sha256.New()
	h.Write(tmp.Bytes())

	return h.Sum(nil)
}

func (tx *Transfer) validateData(stateContext *state.StateContext) bool {
	txHash := tx.TxHash(tx.GetSigningHash())

	// TODO: Common check that the length of all bytes field such as addressfrom, slavepks, tx hash are of even length
	for _, amount := range tx.Amounts() {
		if amount < 0 {
			log.Warnf("Transfer [%s] Invalid Amount = %d", hex.EncodeToString(txHash), amount)
			return false
		}
	}

	if tx.Fee() < 0 {
		log.Warnf("Transfer [%s] Invalid Fee = %d", hex.EncodeToString(txHash), tx.Fee())
		return false
	}

	addressState, err := stateContext.GetAddressState(hex.EncodeToString(tx.AddrFrom()))
	if err != nil {
		log.Warnf("Transfer [%s] Address missing into state context", tx.AddrFrom())
		return false
	}

	if tx.Nonce() != addressState.Nonce() {
		log.Warn(fmt.Sprintf("Transfer [%s] Invalid Nonce %d, Expected Nonce %d",
			hex.EncodeToString(txHash), tx.Nonce(), addressState.Nonce()))
		return false
	}

	balance := addressState.Balance()
	if balance < tx.TotalAmounts()+tx.Fee() {
		log.Warn("Insufficient balance",
			"txhash", hex.EncodeToString(txHash),
			"balance", balance,
			"fee", tx.Fee())
		return false
	}

	// TODO: Max number of destination addresses allowed per transaction
	//if len(tx.AddrsTo()) > int(tx.config.Dev.Transaction.MultiOutputLimit) {
	//	log.Warn("[Transfer] Number of addresses exceeds max limit'")
	//	log.Warn(">> Length of addrsTo %s", len(tx.AddrsTo()))
	//	log.Warn(">> Length of amounts %s", len(tx.Amounts()))
	//	return false
	//}

	if len(tx.AddrsTo()) != len(tx.Amounts()) {
		log.Warn("[Transfer] Mismatch number of addresses to & amounts")
		log.Warnf(">> Length of addrsTo %d", len(tx.AddrsTo()))
		log.Warnf(">> Length of amounts %d", len(tx.Amounts()))
		return false
	}

	// TODO: Move to some common validation
	if !xmss.IsValidXMSSAddress(misc.UnSizedXMSSAddressToSizedXMSSAddress(tx.AddrFrom())) {
		log.Warnf("[Transfer] Invalid address addr_from: %s", tx.AddrFrom())
		return false
	}

	for _, addrTo := range tx.AddrsTo() {
		if !xmss.IsValidXMSSAddress(misc.UnSizedXMSSAddressToSizedXMSSAddress(addrTo)) {
			log.Warnf("[Transfer] Invalid address addr_to: %s", tx.AddrsTo())
			return false
		}
	}

	slavePksLen := len(tx.SlavePKs())
	// TODO: Move 100 to the config
	if slavePksLen > 100 {
		log.Warnf("Slave Public Keys length beyond length limit: %d", slavePksLen)
		return false
	}

	for _, slavePK := range tx.SlavePKs() {
		binAddress := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK(slavePK))
		if !xmss.IsValidXMSSAddress(binAddress) {
			log.Warnf("Slave public key %s is invalid", hex.EncodeToString(slavePK))
			return false
		}
		slaveState := stateContext.GetSlaveState(hex.EncodeToString(tx.AddrFrom()), hex.EncodeToString(slavePK))
		if slaveState != nil {
			log.Warnf("Slave public key %s already exist for address %s", slavePK, tx.AddrFrom())
			return false
		}
	}

	msgLen := len(tx.Message())
	// TODO: Move 100 to the config
	if msgLen > 100 {
		log.Warnf("Message length beyond message length limit: %d", msgLen)
		return false
	}

	return true
}

func (tx *Transfer) Validate(stateContext *state.StateContext) bool {
	/*
		TODO:
		1. Validate OTS Key Reuse
		2. Validate Slave [ done ]
		3. Validate Data in the fields [ done ]
		4. Validate Signature [done]
	*/
	txHash := tx.TxHash(tx.GetSigningHash())

	if !tx.ValidateSlave(stateContext) {
		return false
	}

	if !tx.validateData(stateContext) {
		log.Warn("Custom validation failed")
		return false
	}

	// TODO: Reject any transaction having destination to coinbase address
	// TODO: Reject any transaction having master address to coinbase address
	//if reflect.DeepEqual(tx.config.Dev.Genesis.CoinbaseAddress, tx.PK()) || reflect.DeepEqual(tx.config.Dev.Genesis.CoinbaseAddress, tx.MasterAddr()) {
	//	log.Warn("Coinbase Address only allowed to do Coinbase Transaction")
	//	return false
	//}

	expectedTransactionHash := tx.GenerateTxHash(tx.GetSigningHash())

	// TODO: Move to common function
	if !reflect.DeepEqual(expectedTransactionHash, txHash) {
		log.Warn("Invalid Transaction hash",
			"Expected Transaction hash", hex.EncodeToString(expectedTransactionHash),
			"Found Transaction hash", hex.EncodeToString(txHash))
		return false
	}

	// XMSS Signature Verification
	if !xmss.Verify(tx.GetSigningHash(), tx.Signature(), misc.UnSizedPKToSizedPK(tx.PK())) {
		log.Warn("XMSS Verification Failed")
		return false
	}
	return true
}

func (tx *Transfer) ApplyStateChanges(stateContext *state.StateContext) error {
	if err := tx.applyStateChangesForPK(stateContext); err != nil {
		return err
	}

	addrFrom := tx.AddrFrom()

	addrState, err := stateContext.GetAddressState(hex.EncodeToString(addrFrom))
	if err != nil {
		return err
	}
	total := tx.TotalAmounts()
	addrState.SubtractBalance(total)

	addrsTo := tx.AddrsTo()
	amounts := tx.Amounts()
	for index := range addrsTo {
		addrTo := addrsTo[index]
		amount := amounts[index]

		addrState, err := stateContext.GetAddressState(hex.EncodeToString(addrTo))
		if err != nil {
			return err
		}
		addrState.AddBalance(amount)
	}
	return nil
}

func (tx *Transfer) SetAffectedAddress(stateContext *state.StateContext) error {
	err := stateContext.PrepareAddressState(hex.EncodeToString(tx.AddrFrom()))
	if err != nil {
		return err
	}
	for _, address := range tx.AddrsTo() {
		err = stateContext.PrepareAddressState(hex.EncodeToString(address))
		if err != nil {
			return err
		}
	}

	for _, slavePKs := range tx.SlavePKs() {
		err = stateContext.PrepareSlaveMetaData(hex.EncodeToString(tx.AddrFrom()), hex.EncodeToString(slavePKs))
		if err != nil {
			return err
		}
	}
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK(tx.PK()))
	addrFromPK := hex.EncodeToString(address[:])

	err = stateContext.PrepareOTSIndexMetaData(addrFromPK, tx.OTSIndex())
	if err != nil {
		return err
	}

	err = stateContext.PrepareAddressState(addrFromPK)
	return err
}

func NewTransfer(networkID uint64, addrsTo [][]byte, amounts []uint64, fee uint64,
	slavesPKs [][]byte, message []byte, nonce uint64, xmssPK []byte,
	masterAddr []byte) *Transfer {
	tx := &Transfer{}

	tx.pbData = &protos.Transaction{}
	tx.pbData.NetworkId = networkID
	tx.pbData.Type = &protos.Transaction_Transfer{Transfer: &protos.Transfer{}}

	if masterAddr != nil {
		tx.pbData.MasterAddr = masterAddr
	}

	tx.pbData.Pk = xmssPK
	tx.pbData.Fee = fee
	tx.pbData.Nonce = nonce
	transferPBData := tx.pbData.GetTransfer()
	transferPBData.AddrsTo = addrsTo
	transferPBData.Amounts = amounts
	transferPBData.Message = message

	for _, slavePK := range slavesPKs {
		transferPBData.SlavePks = append(transferPBData.SlavePks, slavePK)
	}

	// TODO: Pass StateContext
	//if !tx.Validate(nil) {
	//	return nil
	//}

	return tx
}
