package transactions

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/go-qrllib/xmss"
	"github.com/theQRL/zond/address"
	"github.com/theQRL/zond/config"
	mockdb "github.com/theQRL/zond/db/mock"
	"github.com/theQRL/zond/metadata"
	"github.com/theQRL/zond/misc"
	"github.com/theQRL/zond/state"
)

func TestNewStake(t *testing.T) {
	masterXmss := xmss.NewXMSSFromHeight(4, 0)
	masterXmssPK := masterXmss.GetPK()
	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((masterXmssPK[:])))

	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
	slaveXmss1PK := slaveXmss1.GetPK()

	networkID := uint64(1)
	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()
	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()
	var dilithiumPKs [][]byte
	dilithiumPKs = append(dilithiumPKs, validatorDilithiumPK[:])
	dilithiumPKs = append(dilithiumPKs, validatorDilithium2PK[:])

	stake := true
	fee := uint64(1)
	nonce := uint64(10)

	staking := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])
	if staking.Stake() != stake {
		t.Error("the stake is incorrectly set")
	}
}

func TestValidateStakeData(t *testing.T) {
	ctrl := gomock.NewController(t)

	masterXmss := xmss.NewXMSSFromHeight(4, 0)
	masterXmssPK := masterXmss.GetPK()
	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((masterXmssPK[:])))

	masterXmss2 := xmss.NewXMSSFromHeight(4, 0)
	masterXmss2PK := masterXmss2.GetPK()
	masterAddr2 := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((masterXmss2PK[:])))

	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
	slaveXmss1PK := slaveXmss1.GetPK()

	networkID := uint64(1)
	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()
	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()
	validatorDilithium3 := dilithium.New()
	validatorDilithium3PK := validatorDilithium3.GetPK()
	validatorXmss3 := xmss.NewXMSSFromHeight(4, 0)
	validatorXmss3Addr := validatorXmss3.GetAddress()

	var dilithiumPKs [][]byte
	dilithiumPKs = append(dilithiumPKs, validatorDilithiumPK[:])
	dilithiumPKs = append(dilithiumPKs, validatorDilithium2PK[:])

	var dilithiumPKs_case5 [][]byte
	dilithiumPKs_case5 = append(dilithiumPKs_case5, validatorDilithiumPK[:])
	dilithiumPKs_case5 = append(dilithiumPKs_case5, validatorDilithium2PK[:])
	dilithiumPKs_case5 = append(dilithiumPKs_case5, validatorDilithium3PK[:])
	dilithiumPKs_case7 := make([][]byte, 101)

	var dilithiumPKs_case8 [][]byte
	dilithiumPKs_case8 = append(dilithiumPKs_case8, validatorDilithium2PK[:])
	dilithiumPKs_case8 = append(dilithiumPKs_case8, validatorDilithium3PK[:])
	validator3DilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithium3PK[:], validatorXmss3Addr[:], true)
	dilithiumMetadataSerialized, _ := validator3DilithiumMetadata.Serialize()

	stake := true
	fee := uint64(1)
	nonce := uint64(10)

	staking := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])
	staking_case2 := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr2[:])
	staking_case3 := NewStake(networkID, dilithiumPKs, stake, fee, 20, slaveXmss1PK[:], masterAddr[:])
	staking_case4 := NewStake(networkID, dilithiumPKs, stake, 20000000000002, nonce, slaveXmss1PK[:], masterAddr[:])
	staking_case5 := NewStake(networkID, dilithiumPKs_case5, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])
	staking_case6 := NewStake(networkID, dilithiumPKs, false, fee, nonce, slaveXmss1PK[:], masterAddr[:])
	staking_case7 := NewStake(networkID, dilithiumPKs_case7, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])
	staking_case8 := NewStake(networkID, dilithiumPKs_case8, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])

	var validators [][]byte
	validators = append(validators, validatorDilithiumPK[:])
	validators = append(validators, validatorDilithium2PK[:])
	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	epoch := uint64(1)
	slotNumber := uint64(100)
	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
		parentBlockHeaderHash, 0)
	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
	addressState := address.NewAddressState(masterAddr[:], nonce, 20000000000001)
	addressesState := make(map[string]*address.AddressState)
	addressesState[hex.EncodeToString(masterAddr[:])] = addressState
	addressesStateSerialized, _ := addressState.Serialize()
	dbAddressStateKey := address.GetAddressStateKey(masterAddr[:])
	dilithiumMetadataKey := metadata.GetDilithiumMetaDataKey(validatorDilithium3PK[:])

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(dbAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(addressesStateSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(dilithiumMetadataKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(dilithiumMetadataSerialized, nil).AnyTimes()

	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
		blockSigningHash, epochMetadata)
	stateContext.PrepareAddressState(hex.EncodeToString(masterAddr[:]))
	stateContext.PrepareDilithiumMetaData(hex.EncodeToString(validatorDilithium3PK[:]))
	if err != nil {
		t.Error("unexpected error while creating new statecontext ", err)
	}

	testCases := []struct {
		name           string
		staking        *Stake
		stateContext   state.StateContext
		expectedOutput bool
	}{
		{
			name:           "ok",
			staking:        staking,
			stateContext:   *stateContext,
			expectedOutput: true,
		},
		{
			name:           "from address missing from statecontext",
			staking:        staking_case2,
			stateContext:   *stateContext,
			expectedOutput: false,
		},
		{
			name:           "incorrect nonce",
			staking:        staking_case3,
			stateContext:   *stateContext,
			expectedOutput: false,
		},
		{
			name:           "insufficient balance",
			staking:        staking_case4,
			stateContext:   *stateContext,
			expectedOutput: false,
		},
		{
			name:           "insufficient staking balance",
			staking:        staking_case5,
			stateContext:   *stateContext,
			expectedOutput: false,
		},
		{
			name:           "dilithium metadata not found",
			staking:        staking_case6,
			stateContext:   *stateContext,
			expectedOutput: false,
		},
		{
			name:           "dilithium PK length beyond limit",
			staking:        staking_case7,
			stateContext:   *stateContext,
			expectedOutput: false,
		},
		{
			name:           "dilithium key associated with another QRL address",
			staking:        staking_case8,
			stateContext:   *stateContext,
			expectedOutput: false,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			output := tc.staking.validateData(&tc.stateContext)
			if output != tc.expectedOutput {
				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
			}

		})
	}
}

func TestStakeValidate(t *testing.T) {
	ctrl := gomock.NewController(t)

	masterXmss := xmss.NewXMSSFromHeight(4, 0)
	masterXmssPK := masterXmss.GetPK()
	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((masterXmssPK[:])))

	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
	slaveXmss1PK := slaveXmss1.GetPK()

	networkID := uint64(1)
	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()
	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()
	var dilithiumPKs [][]byte
	dilithiumPKs = append(dilithiumPKs, validatorDilithiumPK[:])
	dilithiumPKs = append(dilithiumPKs, validatorDilithium2PK[:])

	stake := true
	fee := uint64(1)
	nonce := uint64(10)

	staking := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])
	if staking.Stake() != stake {
		t.Error("the stake is incorrectly set")
	}
	staking.Sign(slaveXmss1, staking.GetSigningHash())

	var validators [][]byte
	validators = append(validators, validatorDilithiumPK[:])
	validators = append(validators, validatorDilithium2PK[:])
	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	epoch := uint64(1)
	slotNumber := uint64(100)
	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
		parentBlockHeaderHash, 0)
	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
	addressState := address.NewAddressState(masterAddr[:], nonce, 20000000000001)
	addressesState := make(map[string]*address.AddressState)
	addressesState[hex.EncodeToString(masterAddr[:])] = addressState
	addressesStateSerialized, _ := addressState.Serialize()
	dbAddressStateKey := address.GetAddressStateKey(masterAddr[:])
	slaveMetadata1 := metadata.NewSlaveMetaData(sha256.New().Sum([]byte("transactionHash")), masterAddr[:], slaveXmss1PK[:])
	slaveMetadataSerialized, _ := slaveMetadata1.Serialize()
	slaveState := make(map[string]*metadata.SlaveMetaData)
	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:]))] = slaveMetadata1

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(dbAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(addressesStateSerialized, nil).AnyTimes()
	store.EXPECT().
		GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).
		Return(slaveMetadataSerialized, nil).AnyTimes()

	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
		blockSigningHash, epochMetadata)
	stateContext.PrepareAddressState(hex.EncodeToString(masterAddr[:]))
	stateContext.PrepareSlaveMetaData(hex.EncodeToString(masterAddr[:]), hex.EncodeToString(slaveXmss1PK[:]))
	if err != nil {
		t.Error("unexpected error while creating new statecontext ", err)
	}

	testCases := []struct {
		name           string
		staking        *Stake
		stateContext   state.StateContext
		expectedOutput bool
	}{
		{
			name:           "ok",
			staking:        staking,
			stateContext:   *stateContext,
			expectedOutput: true,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			output := tc.staking.Validate(&tc.stateContext)
			if output != tc.expectedOutput {
				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
			}

		})
	}
}

func TestStakeApplyStateChanges(t *testing.T) {
	ctrl := gomock.NewController(t)

	masterXmss := xmss.NewXMSSFromHeight(4, 0)
	masterXmssPK := masterXmss.GetPK()
	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((masterXmssPK[:])))

	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
	slaveXmss1Addr := slaveXmss1.GetAddress()
	slaveXmss1PK := slaveXmss1.GetPK()

	networkID := uint64(1)
	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()
	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()
	var dilithiumPKs [][]byte
	dilithiumPKs = append(dilithiumPKs, validatorDilithiumPK[:])
	dilithiumPKs = append(dilithiumPKs, validatorDilithium2PK[:])

	stake := true
	fee := uint64(1)
	nonce := uint64(10)

	staking := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])

	var validators [][]byte
	validators = append(validators, validatorDilithiumPK[:])
	validators = append(validators, validatorDilithium2PK[:])
	validatorXmss := xmss.NewXMSSFromHeight(8, 0)
	validatorXmssAddr := validatorXmss.GetAddress()
	validatorXmss2 := xmss.NewXMSSFromHeight(10, 0)
	validatorXmss2Addr := validatorXmss2.GetAddress()
	validatorDilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithiumPK[:], validatorXmssAddr[:], true)
	validatorDilithiumMetadataSerialized, _ := validatorDilithiumMetadata.Serialize()
	validatorDilithiumMetadata2 := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithium2PK[:], validatorXmss2Addr[:], true)
	validatorDilithiumMetadata2Serialized, _ := validatorDilithiumMetadata2.Serialize()

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	epoch := uint64(1)
	slotNumber := uint64(100)
	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
		parentBlockHeaderHash, 0)
	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
	addressState := address.NewAddressState(masterAddr[:], nonce, 20000000000001)
	slaveAddressState := address.NewAddressState(masterAddr[:], nonce, 20)
	slaveAddressStateSerialized, _ := slaveAddressState.Serialize()
	addressesState := make(map[string]*address.AddressState)
	addressesState[hex.EncodeToString(masterAddr[:])] = addressState
	addressesStateSerialized, _ := addressState.Serialize()
	dbAddressStateKey := address.GetAddressStateKey(masterAddr[:])
	slaveMetadata1 := metadata.NewSlaveMetaData(sha256.New().Sum([]byte("transactionHash")), masterAddr[:], slaveXmss1PK[:])
	slaveMetadataSerialized, _ := slaveMetadata1.Serialize()
	slaveState := make(map[string]*metadata.SlaveMetaData)
	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:]))] = slaveMetadata1

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(dbAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(addressesStateSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(address.GetAddressStateKey(slaveXmss1Addr[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(slaveAddressStateSerialized, nil).AnyTimes()
	store.EXPECT().
		GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).
		Return(slaveMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(metadata.GetDilithiumMetaDataKey(validatorDilithiumPK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(validatorDilithiumMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(metadata.GetDilithiumMetaDataKey(validatorDilithium2PK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(validatorDilithiumMetadata2Serialized, nil).AnyTimes()

	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
		blockSigningHash, epochMetadata)
	stateContext.PrepareAddressState(hex.EncodeToString(masterAddr[:]))
	stateContext.PrepareAddressState(hex.EncodeToString(slaveXmss1Addr[:]))
	stateContext.PrepareSlaveMetaData(hex.EncodeToString(masterAddr[:]), hex.EncodeToString(slaveXmss1PK[:]))
	stateContext.PrepareDilithiumMetaData(hex.EncodeToString(validatorDilithiumPK[:]))
	stateContext.PrepareDilithiumMetaData(hex.EncodeToString(validatorDilithium2PK[:]))
	if err != nil {
		t.Error("unexpected error while creating new statecontext ", err)
	}

	err = staking.ApplyStateChanges(stateContext)
	if err != nil {
		t.Error("got unexpected error while applying state changes in staking transaction ", err)
	}
}

func TestStakeSetAffectedAddress(t *testing.T) {
	ctrl := gomock.NewController(t)

	masterXmss := xmss.NewXMSSFromHeight(4, 0)
	masterXmssPK := masterXmss.GetPK()
	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((masterXmssPK[:])))

	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
	slaveXmss1Addr := slaveXmss1.GetAddress()
	slaveXmss1PK := slaveXmss1.GetPK()

	networkID := uint64(1)
	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()
	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()
	var dilithiumPKs [][]byte
	dilithiumPKs = append(dilithiumPKs, validatorDilithiumPK[:])
	dilithiumPKs = append(dilithiumPKs, validatorDilithium2PK[:])

	stake := true
	fee := uint64(1)
	nonce := uint64(10)

	staking := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])
	staking.Sign(slaveXmss1, staking.GetSigningHash())

	var validators [][]byte
	validators = append(validators, validatorDilithiumPK[:])
	validators = append(validators, validatorDilithium2PK[:])
	validatorXmss := xmss.NewXMSSFromHeight(8, 0)
	validatorXmssAddr := validatorXmss.GetAddress()
	validatorXmss2 := xmss.NewXMSSFromHeight(10, 0)
	validatorXmss2Addr := validatorXmss2.GetAddress()
	validatorDilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithiumPK[:], validatorXmssAddr[:], true)
	validatorDilithiumMetadataSerialized, _ := validatorDilithiumMetadata.Serialize()
	validatorDilithiumMetadata2 := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithium2PK[:], validatorXmss2Addr[:], true)
	validatorDilithiumMetadata2Serialized, _ := validatorDilithiumMetadata2.Serialize()

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	epoch := uint64(1)
	slotNumber := uint64(100)
	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
		parentBlockHeaderHash, 0)
	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
	addressState := address.NewAddressState(masterAddr[:], nonce, 20000000000001)
	slaveAddressState := address.NewAddressState(masterAddr[:], nonce, 20)
	slaveAddressStateSerialized, _ := slaveAddressState.Serialize()
	addressesState := make(map[string]*address.AddressState)
	addressesState[hex.EncodeToString(masterAddr[:])] = addressState
	addressesStateSerialized, _ := addressState.Serialize()
	dbAddressStateKey := address.GetAddressStateKey(masterAddr[:])
	slaveMetadata1 := metadata.NewSlaveMetaData(sha256.New().Sum([]byte("transactionHash")), masterAddr[:], slaveXmss1PK[:])
	slaveMetadataSerialized, _ := slaveMetadata1.Serialize()
	slaveState := make(map[string]*metadata.SlaveMetaData)
	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:]))] = slaveMetadata1

	otsIndex := uint64(binary.BigEndian.Uint32(staking.pbData.Signature[0:4]))
	otsIndexMetadata := metadata.NewOTSIndexMetaData(slaveXmss1Addr[:], otsIndex/config.GetDevConfig().OTSBitFieldPerPage)
	otsIndexMetadataSerialized, _ := otsIndexMetadata.Serialize()
	otsStateKey := metadata.GetOTSIndexMetaDataKeyByOTSIndex(slaveXmss1Addr[:], otsIndex)

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(dbAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(addressesStateSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(address.GetAddressStateKey(slaveXmss1Addr[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(slaveAddressStateSerialized, nil).AnyTimes()
	store.EXPECT().
		GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).
		Return(slaveMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(metadata.GetDilithiumMetaDataKey(validatorDilithiumPK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(validatorDilithiumMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(metadata.GetDilithiumMetaDataKey(validatorDilithium2PK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(validatorDilithiumMetadata2Serialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(otsStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(otsIndexMetadataSerialized, nil).AnyTimes()

	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
		blockSigningHash, epochMetadata)

	if err != nil {
		t.Error("unexpected error while creating new statecontext ", err)
	}

	err = staking.SetAffectedAddress(stateContext)
	if err != nil {
		t.Error("got unexpected error while stting affected addresses ", err)
	}
}

func TestStakeGetSigningHash(t *testing.T) {
	masterAddr, _ := hex.DecodeString("0003009cf8640426af83ccd3adbdb3d290f71245")
	slaveXmss1PK, _ := hex.DecodeString("000300252a2b71c81fde22419c528ba09bf025f145335c307e09dec6eeed5cff391b0ded8dd1b67f234394255f5dea219289f334f49e2f237c37e78097765796064aac")

	networkID := uint64(1)
	validatorDilithiumPK, _ := hex.DecodeString("135febcea1b6c55fb951df921e40fdf445367967ff6bc7ed965b69c10940d9a8de63bcc06d5b66cf5088840fc0070c6a8c7072db68e3a5237df59be112b54e805a43c4e4bb41cfe5ae73011598bfebeaf6f81f8fdc4343767dc317d41fb7d686c07d342e77980b01fb38f6bcf8e36a7d2de9bbfc564b35638aa0d1a12d9b7ea48440117a4de9586ae0d4cd2ccd637e59c3fa75faa9bc4b65bdc078e9c6583e24503fa47af8fa8b6f24de491744eaad66ec4e072595fa8eeb688707d0f6083deedb5375952cb6b49ed5745f384827909e6b311b6ec9657211c2c611e3bfbfb613a7b4eb98b8fab70d3060716d6fa42480a8b738488110e672c00720d909d0293b286018d117a72bfa329d75cca061a415b185b581baf558ad1abf978eaf399ed2e46db8e9922a025c7ef0efe3414d0b704ab8900808d25a11aa1332a9e1c703738b86971f7f10f9d89d051238db68a0b671224ae6302769fcffabdda54d4b7c8f2bc37ed8c64e239dff3d0d1cb91dd771289a61c4c331e10433ab1db2dc5aa4dd72eaeaf02f46d4943e9c073cba6888dc114a1a880d7097c59fb7536033eb8fe33d64ffb47680a6e125a9bf85039d5087767ec1dba4db7cfbcc3897627ab17d448c915ce36da1d848e40687a6fc71c4101fd0a8c5338df962cd192264fb43c1e5ff525a139ad0433241039238927f0cac35a8e2455bcc6175fa8e91b8ffa606c252f484fbd4ea46b03685642fa24e55015d6a9eedd4f02119934aa44450dbf4f92ee86c4a6fcc7f4c138a4ec2cbef3edabd26a01b32513f728c9acb3c2859ded525bf0c717a10440c47859403e7d3519890d83b0438fb67a4dc146c07d3f2d2834062745e25562650ab039b108fc5949f07bbd289e1afb96f5c29b24f459276fc32361d98b2e3ac1e2ecc0f8d5f213145ff75fecdb2a8e20c2bfc05999c0b88669d3a6da2e9ca583562188e21f82f82bd61c4ceb73482818ae0bd40204c88c4a04aa954e2d42918bcace86f2929f92ae6d4ba8d5caba2a9210cf159f6b54076b49b63f7efde737cb3c1f5a5071458270a78655bd2872b438e6cbecb94ab208429ddbf9c4508c5b819fb278f7c85d3b2ea6b88c3b6604c9f34a95d8f823566ce3c46432bd9b4a280fe55a1ef1f250d92b5101bb649fa71cae1b1c1c048ee38d38b7a4616320cf3d4cb95d8133db2cf7ea8c2f6a48d0812c90cffa9f4ba4493ac483b1942d5690f77380351d6df48cc51dca17f1bf786f0981470036414520fc39cfb1eb3d0c12fe8f6a11fb0aef6569b2c67f567b996370b8808b90339802a89e09b0b4a95df01d7d0eb8568edafbc98b6847d4825dd8959ac99c7272c90ffe9ef88b642e60706604f48fa33b5e3a7b87a68b5f89a63e55dd453d90eb47ddf6f3bd5c397d13e4c2d59dff5a9f969f5516ee7afed194b5aba83d0f2a2a838aaddb1103e4a503ccbc100633a61076b037408994ae9e586bc307e68519c47a503358b9c388df500cf74ad104d11de30a4a83778850029eff3b29c746f9de92038a2109859cdc6e4f480ca692e845aef592801a5ec4e71433f8ee162d49877fc7009d80cb1c6cb0b3ed407c5a80a7fda255ab2a358eed06caf7d51db6d560d8e60cfe71bbe1c2ca6a7462c817ab2a71f09aff1e36688aafb82eb094622f8c7e0dc6a19833cdaa6cb421cfebd475b9ddbd8257a96a75fe6745d59fece8faa540e5204f3d2b0146e8658a2e006f751df647b86c3ea2dbe66f5bd3f36a5a51fec55f7bb2d3595d52caac3d6f7eabfda533bfd7ace942ef44ec5ce70f7b67c2a229ad0b2855be60f2fe22771cf53de469c4cea25bdaea60c46aa754d8a4cc40f8a6d87610b4c01c32fa6c404f8d571c2d76149982e752ba802930472950bc15d38e98ab865b9ebbd9bf5801df3587110285663daf311ea7e5a0444ce03b6a7e574557521c068490e774d2a331ac71176391c434c3ec65e8fae60505ba408c661998b9c1469773fbdac6e2be31747333611bbdfdbd9a49abee970b92a4f2e26fbaa1549940923c8d4d9d12a4a0bc0f09e453dc0d75124f010cd6d55a5a590b17fb94a6913b5")
	validatorDilithium2PK, _ := hex.DecodeString(" 1e7d06b1ddbb7175a83c6d4cdf16e6ae51be17e862c5c515bad27687dce086c3ea31864b960945b0edf368dbbd6fc33b7864a4026093d8fd640fcff125cfc205c3fb2504a45572be03e38e5f8fed757fbeed7e0d0eae5a1ce441eea3cef59552b3b514514d5247c7ef330ce882b2bed0ff10b396b50c0bb53b39912e571674d4a9ac224d03c4ca86082dc523df1afad4f013b59365ab9d72087e0564bb9cb87b896d467fc31fb8e67f975ca1427132714fda7e2069020e9aabbce0e5a6a95253cf779a1ee1feae67bde79e77248f2e462dd61285bed1aed1a3cfd0358870be3b27e29461ab8f7abbc2fa25fc1f3502437bc334217e7d1d7b9d24fe6cc89e0477afa8893eb61e05ce157a08dd6897564bcb4fb335aede1a22652ec5d19a7e22a97fd6287ce4164856c13a331ae056b1a3937a7f54a42356a24192f1386952d2cfbb85184f6c4e2bf189732e050b5baf5c82e00f27b42d5788f7e4a3140fe0b3c62930206a94f3cfa4e9720efa6644ab35d9c1d7cf88b61a0904e67d7c9e63a4aeb58e7965e25d7ac6dc33c3e9ec42da394377e760c1a1891417c728182b20a61623d5d07274d10c98d167551ce69710580ded56d115cbc7fae2f21529c075c13eed07548bdbade5be7e964bb68696f1bb8fc9334374cfe77fbab2ef840fc5833f33e09eca4598fc80e7067fdee7afe59428065a1cf72707cdc7141cff584a51c6fb11b38f762a2f0157eb3e83292140354be62d0251df646a9e0d973d7bc9dcd5349d834c62ce661c57aed0255758a6a88ed02c0e021cbc9b6fa8cc3f3f6d2b6c04ab4c04da272d70dc66364cab2cd0cf888a5a5b53bc12c1b1db5c6d77e1c5c2305c83fbb7d80aed40deab9744b08a46de4ffcaf325f02bb308c7632c04d211a5f78f91cd037f8d305ce6c2f4b8eba3a2c670884489af822432f2f03fb8beaca71e5934668749297d1b2d002df19e343a3a2cdc50f00ba3813ffcbac75ce0fd94795115f458614c85be7036b57e54638d3101b82c5130ccfc1138eff5c4c7715ac79bd19f62eb6dd13eceea1aabc3cea5790e0454c42a293c741bdb9c7f6e7a1d28bac8ec541f8171635ef0753550ee12bb4c3b655f9fe78f77a121c471e6314664d27d8851ae75b29ac5da7f866b1246f75754d2063d202e0b7eb4211b975ecabed6294c95b57fd58b6d8f10ae5dbd05becdb6434630f3c9243fd75109c52238ec46af87885beecb247a082d88854b3945feb711b09318e309b0b7a7cc3b3f18361dadbec27f05d0d1ab59020641a635178ee17ae4050efe0e0013b6a188247ed770089e18cabeee47396267ac6ec91b07f49477213e3fa60b26f4be680cfed3049bb95774b899d4cbf2e4b620fd99b4349ccbce67d92698ed2a166a0633d78a947c45aa3e78cb274027e9ffef5cb4237d0c16a74483a88a651eeb2ed453567b45612fe8bc843fb5187b2c786739274a7622dc4b1ff940f95dc217a317e3bed71662cb8cb13658af2e68a81fe8989f21957f8e1c01bf5a243e475b2b6b65a52dd8a5fcb90896dfbabbc36dd93849cf246457fd741641d09059c2f271f90daab905f350524a844d7a5715dc73f7176f96ae8349fec809886d4838789db927a898001b4c07f1aec3ab96fe652b46e65396898378dabcda24db570c1f058543b8db9c07b729daf3f322d681b1458d89f2bec9cc6607a68361f701a2da9976fcdd8bcecb1a9f36556b3f5dabebda1083557ce4874f9ec926162526ce081bf8f001d98bfecdc0542fa68a5ae16c1589c55aa3597b43bc7ac5d4593107b2d74b19c36659347e283eed7a5f6a83a688683d0509d18a116fdef4d1661f3ac96efc889f2f8e17051f8daf8b71b8294817799357ce0059f06348ac595d7423b75b340e92c2cfb1490f5442244b2c5c68c66e1c60d5b4d7024f19745e87fcba0af7601c19419b1fb3739e0c0184e074508369039a023aecb7b10a012007e7e49873b2a0ae98409b4fb65f6128a729f3bcb838e8664112bdb21abbdddb3aee9c2e62fd257046f6a742e36f7fb536b514017157aec9d58ab0c31db1a80099f1431eacebf7b7b")
	var dilithiumPKs [][]byte
	dilithiumPKs = append(dilithiumPKs, validatorDilithiumPK[:])
	dilithiumPKs = append(dilithiumPKs, validatorDilithium2PK[:])

	stake := true
	fee := uint64(1)
	nonce := uint64(10)

	staking := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])

	expectedSigningHash, _ := hex.DecodeString("00441165fa0cc969de3b28a6e8219494d5c92fe05b62368020297f08b9172d3a")
	output := staking.GetSigningHash()
	if hex.EncodeToString(output) != hex.EncodeToString(expectedSigningHash) {
		t.Errorf("expected stake signing hash (%v), got (%v)", hex.EncodeToString(output), hex.EncodeToString(expectedSigningHash))
	}
}
