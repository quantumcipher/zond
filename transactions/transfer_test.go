package transactions

import (
	"encoding/hex"
	"testing"

	"github.com/theQRL/go-qrllib/xmss"
)

func TestNewTransfer(t *testing.T) {
	masterXmss := xmss.NewXMSSFromHeight(4, 0)
	masterXmssPK := masterXmss.GetPK()
	addrTo := xmss.NewXMSSFromHeight(10, 0).GetAddress()

	amount := uint64(30)
	fee := uint64(1)
	message := []byte("message")
	nonce := uint64(30)
	networkID := uint64(1)
	transfer := NewTransfer(networkID, addrTo[:], amount, fee, 0, message, nonce, masterXmssPK[:])

	if transfer.Value() != uint64(30) {
		t.Error("the total amount of transfer is incorrect")
	}
}

// func TestValidateTransferData(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	masterXmss := xmss.NewXMSSFromHeight(4, 0)
// 	masterXmssPK := masterXmss.GetPK()
// 	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((masterXmssPK[:])))

// 	masterXmss2 := xmss.NewXMSSFromHeight(6, 0)
// 	masterXmssPK2 := masterXmss2.GetPK()
// 	masterAddr2 := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((masterXmssPK2[:])))

// 	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
// 	slaveXmss1PK := slaveXmss1.GetPK()
// 	slaveXmss2 := xmss.NewXMSSFromHeight(8, 0)
// 	slaveXmss2PK := slaveXmss2.GetPK()
// 	var slavePKs [][]byte
// 	slavePKs = append(slavePKs, slaveXmss1PK[:])
// 	slavePKs = append(slavePKs, slaveXmss2PK[:])
// 	slavePKs_case6 := make([][]byte, 101)
// 	validatorDilithium := dilithium.New()
// 	validatorDilithiumPK := validatorDilithium.GetPK()
// 	validatorDilithium2 := dilithium.New()
// 	validatorDilithium2PK := validatorDilithium2.GetPK()

// 	addrTo1 := xmss.NewXMSSFromHeight(10, 0).GetAddress()
// 	addrTo2 := xmss.NewXMSSFromHeight(12, 0).GetAddress()
// 	var addrsTo [][]byte
// 	addrsTo = append(addrsTo, addrTo1[:])
// 	addrsTo = append(addrsTo, addrTo2[:])

// 	var amounts []uint64
// 	var amounts_case5 []uint64
// 	amounts = append(amounts, 10)
// 	amounts = append(amounts, 20)
// 	amounts_case5 = append(amounts_case5, 10)
// 	amounts_case5 = append(amounts_case5, 20)
// 	amounts_case5 = append(amounts_case5, 30)
// 	fee := uint64(1)
// 	message := []byte("message")
// 	message_case7 := []byte("jdncdfvfvkfdkfkfndkckvmfkvdkcmv c scasf cajs c sj cisncoianofcikaoek fcoaek coej flakw colakd fcalsaa")
// 	nonce := uint64(0)
// 	networkID := uint64(1)
// 	transfer := NewTransfer(networkID, addrsTo, amounts, fee, slavePKs, message, nonce, masterXmssPK[:], masterAddr[:])

// 	transfer_case2 := NewTransfer(networkID, addrsTo, amounts, fee, slavePKs, message, nonce, masterXmssPK2[:], masterAddr2[:])
// 	transfer_case3 := NewTransfer(networkID, addrsTo, amounts, fee, slavePKs, message, 10, masterXmssPK[:], masterAddr[:])
// 	transfer_case4 := NewTransfer(networkID, addrsTo, amounts, 200, slavePKs, message, nonce, masterXmssPK[:], masterAddr[:])
// 	transfer_case5 := NewTransfer(networkID, addrsTo, amounts_case5, fee, slavePKs, message, nonce, masterXmssPK[:], masterAddr[:])
// 	transfer_case6 := NewTransfer(networkID, addrsTo, amounts, fee, slavePKs_case6, message, nonce, masterXmssPK[:], masterAddr[:])
// 	transfer_case7 := NewTransfer(networkID, addrsTo, amounts, fee, slavePKs, message_case7, nonce, masterXmssPK[:], masterAddr[:])

// 	var validators [][]byte
// 	validators = append(validators, validatorDilithiumPK[:])
// 	validators = append(validators, validatorDilithium2PK[:])
// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
// 	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
// 	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
// 	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
// 		parentBlockHeaderHash, 0)
// 	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
// 	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
// 	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
// 	addressState := address.NewAddressState(masterAddr[:], nonce, 100)
// 	addressesState := make(map[string]*address.AddressState)
// 	addressesState[hex.EncodeToString(masterAddr[:])] = addressState
// 	addressesStateSerialized, _ := addressState.Serialize()
// 	dbAddressStateKey := address.GetAddressStateKey(masterAddr[:])

// 	store := mockdb.NewMockDB(ctrl)
// 	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
// 	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(dbAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(addressesStateSerialized, nil).AnyTimes()
// 	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
// 		blockSigningHash, epochMetadata)
// 	stateContext.PrepareAddressState(hex.EncodeToString(masterAddr[:]))
// 	if err != nil {
// 		t.Error("unexpected error while creating new statecontext ", err)
// 	}

// 	testCases := []struct {
// 		name           string
// 		transfer       *Transfer
// 		stateContext   state.StateContext
// 		expectedOutput bool
// 	}{
// 		{
// 			name:           "ok",
// 			transfer:       transfer,
// 			stateContext:   *stateContext,
// 			expectedOutput: true,
// 		},
// 		{
// 			name:           "transfer address missing",
// 			transfer:       transfer_case2,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "invalid nonce",
// 			transfer:       transfer_case3,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "insufficient balance",
// 			transfer:       transfer_case4,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "addresses and amounts mismatch",
// 			transfer:       transfer_case5,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "slavePKs length exceeded",
// 			transfer:       transfer_case6,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "message length exceeded",
// 			transfer:       transfer_case7,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			output := tc.transfer.validateData(&tc.stateContext)
// 			if output != tc.expectedOutput {
// 				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
// 			}

// 		})
// 	}

// }

// func TestValidate(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	masterXmss := xmss.NewXMSSFromHeight(4, 0)
// 	masterXmssPK := masterXmss.GetPK()
// 	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((masterXmssPK[:])))

// 	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
// 	slaveXmss1PK := slaveXmss1.GetPK()
// 	slaveXmss2 := xmss.NewXMSSFromHeight(8, 0)
// 	slaveXmss2PK := slaveXmss2.GetPK()
// 	var slavePKs [][]byte
// 	slavePKs = append(slavePKs, slaveXmss2PK[:])
// 	validatorDilithium := dilithium.New()
// 	validatorDilithiumPK := validatorDilithium.GetPK()
// 	validatorDilithium2 := dilithium.New()
// 	validatorDilithium2PK := validatorDilithium2.GetPK()

// 	addrTo1 := xmss.NewXMSSFromHeight(10, 0).GetAddress()
// 	addrTo2 := xmss.NewXMSSFromHeight(12, 0).GetAddress()
// 	var addrsTo [][]byte
// 	addrsTo = append(addrsTo, addrTo1[:])
// 	addrsTo = append(addrsTo, addrTo2[:])

// 	var amounts []uint64
// 	amounts = append(amounts, 10)
// 	amounts = append(amounts, 20)
// 	fee := uint64(1)
// 	message := []byte("message")
// 	nonce := uint64(10)
// 	networkID := uint64(1)
// 	transfer := NewTransfer(networkID, addrsTo, amounts, fee, slavePKs, message, nonce, slaveXmss1PK[:], masterAddr[:])
// 	transfer.Sign(slaveXmss1, transfer.GetSigningHash())

// 	var validators [][]byte
// 	validators = append(validators, validatorDilithiumPK[:])
// 	validators = append(validators, validatorDilithium2PK[:])
// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
// 	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
// 	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
// 	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
// 		parentBlockHeaderHash, 0)
// 	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
// 	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
// 	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
// 	addressState := address.NewAddressState(masterAddr[:], nonce, 100)
// 	addressesState := make(map[string]*address.AddressState)
// 	addressesState[hex.EncodeToString(masterAddr[:])] = addressState
// 	addressesStateSerialized, _ := addressState.Serialize()
// 	dbAddressStateKey := address.GetAddressStateKey(masterAddr[:])
// 	slaveMetadata1 := metadata.NewSlaveMetaData(sha256.New().Sum([]byte("transactionHash")), masterAddr[:], slaveXmss1PK[:])
// 	slaveMetadataSerialized, _ := slaveMetadata1.Serialize()
// 	slaveState := make(map[string]*metadata.SlaveMetaData)
// 	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:]))] = slaveMetadata1

// 	store := mockdb.NewMockDB(ctrl)
// 	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
// 	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(dbAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(addressesStateSerialized, nil).AnyTimes()
// 	store.EXPECT().
// 		GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).
// 		Return(slaveMetadataSerialized, nil).AnyTimes()

// 	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
// 		blockSigningHash, epochMetadata)
// 	stateContext.PrepareAddressState(hex.EncodeToString(masterAddr[:]))
// 	stateContext.PrepareSlaveMetaData(hex.EncodeToString(masterAddr[:]), hex.EncodeToString(slaveXmss1PK[:]))

// 	if err != nil {
// 		t.Error("unexpected error while creating new statecontext ", err)
// 	}

// 	testCases := []struct {
// 		name           string
// 		transfer       *Transfer
// 		stateContext   state.StateContext
// 		expectedOutput bool
// 	}{
// 		{
// 			name:           "ok",
// 			transfer:       transfer,
// 			stateContext:   *stateContext,
// 			expectedOutput: true,
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			output := tc.transfer.Validate(&tc.stateContext)
// 			if output != tc.expectedOutput {
// 				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
// 			}
// 		})
// 	}
// }

// func TestApplyStateChanges(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	masterXmss := xmss.NewXMSSFromHeight(4, 0)
// 	masterXmssPK := masterXmss.GetPK()
// 	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((masterXmssPK[:])))

// 	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
// 	slaveXmss1PK := slaveXmss1.GetPK()
// 	slaveXmss1Addr := slaveXmss1.GetAddress()
// 	slaveXmss2 := xmss.NewXMSSFromHeight(8, 0)
// 	slaveXmss2PK := slaveXmss2.GetPK()
// 	var slavePKs [][]byte
// 	slavePKs = append(slavePKs, slaveXmss2PK[:])
// 	validatorDilithium := dilithium.New()
// 	validatorDilithiumPK := validatorDilithium.GetPK()
// 	validatorDilithium2 := dilithium.New()
// 	validatorDilithium2PK := validatorDilithium2.GetPK()

// 	addrTo1 := xmss.NewXMSSFromHeight(10, 0).GetAddress()
// 	addrTo2 := xmss.NewXMSSFromHeight(12, 0).GetAddress()
// 	var addrsTo [][]byte
// 	addrsTo = append(addrsTo, addrTo1[:])
// 	addrsTo = append(addrsTo, addrTo2[:])

// 	var amounts []uint64
// 	amounts = append(amounts, 10)
// 	amounts = append(amounts, 20)
// 	fee := uint64(1)
// 	message := []byte("message")
// 	nonce := uint64(10)
// 	networkID := uint64(1)
// 	transfer := NewTransfer(networkID, addrsTo, amounts, fee, slavePKs, message, nonce, slaveXmss1PK[:], masterAddr[:])

// 	var validators [][]byte
// 	validators = append(validators, validatorDilithiumPK[:])
// 	validators = append(validators, validatorDilithium2PK[:])
// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
// 	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
// 	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
// 	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
// 		parentBlockHeaderHash, 0)
// 	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
// 	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
// 	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
// 	addressState := address.NewAddressState(masterAddr[:], nonce, 100)
// 	slaveAddressState := address.NewAddressState(slaveXmss1Addr[:], nonce, 20)
// 	addrTo1AddressState := address.NewAddressState(addrTo1[:], nonce, 0)
// 	addrTo2AddressState := address.NewAddressState(addrTo2[:], nonce, 0)
// 	addressesState := make(map[string]*address.AddressState)
// 	addressesState[hex.EncodeToString(masterAddr[:])] = addressState
// 	addressesState[hex.EncodeToString(slaveXmss1Addr[:])] = slaveAddressState
// 	addressesState[hex.EncodeToString(addrTo1[:])] = addrTo1AddressState
// 	addressesState[hex.EncodeToString(addrTo2[:])] = addrTo2AddressState
// 	masterAddressStateSerialized, _ := addressState.Serialize()
// 	slaveAddressStateSerialized, _ := slaveAddressState.Serialize()
// 	addrTo1AddressStateSerialized, _ := addrTo1AddressState.Serialize()
// 	addrTo2AddressStateSerialized, _ := addrTo2AddressState.Serialize()
// 	masterAddressStateKey := address.GetAddressStateKey(masterAddr[:])
// 	slaveAddressStateKey := address.GetAddressStateKey(slaveXmss1Addr[:])
// 	addrTo1AddressStateKey := address.GetAddressStateKey(addrTo1[:])
// 	addrTo2AddressStateKey := address.GetAddressStateKey(addrTo2[:])

// 	store := mockdb.NewMockDB(ctrl)
// 	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
// 	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(masterAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(masterAddressStateSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(slaveAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(slaveAddressStateSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(addrTo1AddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(addrTo1AddressStateSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(addrTo2AddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(addrTo2AddressStateSerialized, nil).AnyTimes()

// 	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
// 		blockSigningHash, epochMetadata)
// 	stateContext.PrepareAddressState(hex.EncodeToString(masterAddr[:]))
// 	stateContext.PrepareAddressState(hex.EncodeToString(slaveXmss1Addr[:]))
// 	stateContext.PrepareAddressState(hex.EncodeToString(addrTo1[:]))
// 	stateContext.PrepareAddressState(hex.EncodeToString(addrTo2[:]))

// 	if err != nil {
// 		t.Error("unexpected error while creating new statecontext ", err)
// 	}

// 	err = transfer.ApplyStateChanges(stateContext)
// 	if err != nil {
// 		t.Error("got unexpected error while applying state changes ", err)
// 	}
// }

func TestGetSigningHash(t *testing.T) {
	//masterXmss := xmss.NewXMSSFromHeight(4, 0)
	masterXmssPK, _ := hex.DecodeString("1002001eb39fc767971a9438d78a0bccadf561b72ece6a92640f0094626d652b08bf3238f00d3c24f01080cdb800b98ac117c5b017effc07ba591bc2d842fcc61c2e1e")
	addrTo, _ := hex.DecodeString("100500cb43ef24349ef178b7c6b83a91175d2e5e")

	amount := uint64(30)
	fee := uint64(1)
	message := []byte("message")
	nonce := uint64(30)
	networkID := uint64(1)
	transfer := NewTransfer(networkID, addrTo[:], amount, fee, 0, message, nonce, masterXmssPK[:])

	expectedHash, _ := hex.DecodeString("f29eff1776c6f3ec6fd987fb1af4f814c0b9d77f8edb09a2c167937421d277e7")

	output := transfer.GetSigningHash()
	if hex.EncodeToString(output.Bytes()) != hex.EncodeToString(expectedHash) {
		t.Errorf("expected transfer signing hash to be (%v), got (%v)", hex.EncodeToString(expectedHash), hex.EncodeToString(output.Bytes()))
	}
}
