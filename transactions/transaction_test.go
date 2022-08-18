package transactions

import (
	"encoding/hex"
	"testing"

	// "github.com/golang/mock/gomock"
	// "github.com/theQRL/go-qrllib/dilithium"
	// "github.com/theQRL/go-qrllib/xmss"
	// "github.com/theQRL/zond/address"
	// mockdb "github.com/theQRL/zond/db/mock"
	// "github.com/theQRL/zond/metadata"
	// "github.com/theQRL/zond/misc"

	"github.com/theQRL/zond/protos"
	// "github.com/theQRL/zond/state"
)

func TestProtoToTransaction(t *testing.T) {
	slaveXmss1PK, _ := hex.DecodeString("000300252a2b71c81fde22419c528ba09bf025f145335c307e09dec6eeed5cff391b0ded8dd1b67f234394255f5dea219289f334f49e2f237c37e78097765796064aac")

	networkID := uint64(1)
	gas := uint64(1)
	gasPrice := uint64(1)
	nonce := uint64(10)
	protoTx := &protos.Transaction{
		NetworkId: networkID,
		Gas:       gas,
		GasPrice:  gasPrice,
		Nonce:     nonce,
		Pk:        slaveXmss1PK,
	}

	_ = ProtoToTransaction(protoTx)
}

func TestGenerateTxHash(t *testing.T) {
	masterXmssPK, _ := hex.DecodeString("000300252a2b71c81fde22419c528ba09bf025f145335c307e09dec6eeed5cff391b0ded8dd1b67f234394255f5dea219289f334f49e2f237c37e78097765796064aac")

	addrTo := "000500f162cf94b30c2a6b96f46d8fba2b361fd7"

	amount := uint64(10)
	gas := uint64(1)
	gasPrice := uint64(1)
	nonce := uint64(10)
	networkID := uint64(1)
	data := []byte("data")
	transfer := NewTransfer(networkID, []byte(addrTo), amount, gas, gasPrice, data, nonce, masterXmssPK[:])

	expectedHash := "0xb483dda8b28c6510695ce6856fd934ae0b4921f6b3ed22664ca779635cbaf149"

	output := transfer.GenerateTxHash()
	if output.String() != expectedHash {
		t.Errorf("expected transaction hash (%v), got (%v)", expectedHash, output.String())
	}
}

/*
func TestValidateSlave(t *testing.T) {
	ctrl := gomock.NewController(t)

	masterXmss := xmss.NewXMSSFromHeight(4, 0)
	masterXmssPK := masterXmss.GetPK()
	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((masterXmssPK[:])))

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
	staking_case2 := NewStake(networkID, dilithiumPKs, stake, fee, nonce, masterXmssPK[:], masterAddr[:])
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
		{
			name:           "matching master and slave address",
			staking:        staking_case2,
			stateContext:   *stateContext,
			expectedOutput: false,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			output := tc.staking.ValidateSlave(&tc.stateContext)
			if output != tc.expectedOutput {
				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
			}

		})
	}
}
*/
