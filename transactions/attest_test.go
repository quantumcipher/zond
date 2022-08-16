package transactions

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/go-qrllib/xmss"
	"github.com/theQRL/zond/address"
	mockdb "github.com/theQRL/zond/db/mock"
	"github.com/theQRL/zond/metadata"
	"github.com/theQRL/zond/misc"
	"github.com/theQRL/zond/protos"
	"github.com/theQRL/zond/state"
	"google.golang.org/protobuf/proto"
)

func TestNewAttest(t *testing.T) {
	networkID := uint64(1)
	coinBaseNonce := uint64(10)
	attest := NewAttest(networkID, coinBaseNonce)

	if attest.Nonce() != coinBaseNonce {
		t.Error("the nonce is not set correctly while attesting transaction")
	}
}

func TestAttestValidateData(t *testing.T) {
	ctrl := gomock.NewController(t)

	networkID := uint64(1)
	coinBaseNonce := uint64(10)
	attest := NewAttest(networkID, coinBaseNonce)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	blockProposerXmss := xmss.NewXMSSFromHeight(4, 0)
	blockProposerXmssPK := blockProposerXmss.GetPK()
	blockProposerAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((blockProposerXmssPK[:])))
	blockProposerDilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), blockProposerPK[:], blockProposerAddr[:], true)
	blockProposerDilithiumMetadataSerialized, _ := blockProposerDilithiumMetadata.Serialize()

	attestorXmss := xmss.NewXMSSFromHeight(8, 0)
	attestorXmssAddr := attestorXmss.GetAddress()
	attestorDilithium := dilithium.New()
	attestorDilithiumPK := attestorDilithium.GetPK()
	attestorDilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), attestorDilithiumPK[:], attestorXmssAddr[:], true)
	attestorDilithiumMetadataSerialized, _ := attestorDilithiumMetadata.Serialize()

	var validators [][]byte
	validators = append(validators, blockProposerPK[:])
	validators = append(validators, attestorDilithiumPK[:])
	epoch := uint64(1)
	slotNumber := uint64(100)
	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
	epochMetadaPBData := &protos.EpochMetaData{
		SlotInfo: []*protos.SlotInfo{{
			SlotLeader: 0,
			Attestors:  []uint64{uint64(1)},
		}},
		Validators: validators,
	}
	epochMetadata := metadata.NewEpochMetaData(epoch, nil, nil)
	epochMetadataSerialized, _ := proto.Marshal(epochMetadaPBData)
	epochMetadata.DeSerialize(epochMetadataSerialized)
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
		parentBlockHeaderHash, 0)
	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
	blockproposerMetadataKey := metadata.GetDilithiumMetaDataKey(blockProposerPK[:])

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(blockproposerMetadataKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(blockProposerDilithiumMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(metadata.GetDilithiumMetaDataKey(attestorDilithiumPK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(attestorDilithiumMetadataSerialized, nil).AnyTimes()

	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
		blockSigningHash, epochMetadata)
	// stateContext.PrepareAddressState(hex.EncodeToString(binCoinBaseAddress[:]))
	stateContext.PrepareDilithiumMetaData(hex.EncodeToString(blockProposerPK[:]))
	stateContext.PrepareDilithiumMetaData(hex.EncodeToString(attestorDilithiumPK[:]))
	if err != nil {
		t.Error("unexpected error while creating new statecontext ", err)
	}

	testCases := []struct {
		name           string
		attest         *Attest
		stateContext   state.StateContext
		dilithiumKey   *dilithium.Dilithium
		expectedOutput bool
	}{
		{
			name:           "ok",
			attest:         attest,
			stateContext:   *stateContext,
			dilithiumKey:   attestorDilithium,
			expectedOutput: true,
		},
		{
			name:           "attestor key not present",
			attest:         attest,
			stateContext:   *stateContext,
			dilithiumKey:   blockProposer,
			expectedOutput: false,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.attest.Sign(tc.dilithiumKey, attest.GetSigningHash(stateContext.PartialBlockSigningHash()))

			output := tc.attest.validateData(stateContext)
			if output != tc.expectedOutput {
				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
			}

		})
	}
}

func TestAttestValidate(t *testing.T) {
	ctrl := gomock.NewController(t)

	networkID := uint64(1)
	coinBaseNonce := uint64(10)
	attest := NewAttest(networkID, coinBaseNonce)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	blockProposerXmss := xmss.NewXMSSFromHeight(4, 0)
	blockProposerXmssPK := blockProposerXmss.GetPK()
	blockProposerAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((blockProposerXmssPK[:])))
	blockProposerDilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), blockProposerPK[:], blockProposerAddr[:], true)
	blockProposerDilithiumMetadataSerialized, _ := blockProposerDilithiumMetadata.Serialize()

	attestorXmss := xmss.NewXMSSFromHeight(8, 0)
	attestorXmssAddr := attestorXmss.GetAddress()
	attestorDilithium := dilithium.New()
	attestorDilithiumPK := attestorDilithium.GetPK()
	attestorDilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), attestorDilithiumPK[:], attestorXmssAddr[:], true)
	attestorDilithiumMetadataSerialized, _ := attestorDilithiumMetadata.Serialize()

	var validators [][]byte
	validators = append(validators, blockProposerPK[:])
	validators = append(validators, attestorDilithiumPK[:])
	epoch := uint64(1)
	slotNumber := uint64(100)
	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
	epochMetadaPBData := &protos.EpochMetaData{
		SlotInfo: []*protos.SlotInfo{{
			SlotLeader: 0,
			Attestors:  []uint64{uint64(1)},
		}},
		Validators: validators,
	}
	epochMetadata := metadata.NewEpochMetaData(epoch, nil, nil)
	epochMetadataSerialized, _ := proto.Marshal(epochMetadaPBData)
	epochMetadata.DeSerialize(epochMetadataSerialized)
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
		parentBlockHeaderHash, 0)
	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
	blockproposerMetadataKey := metadata.GetDilithiumMetaDataKey(blockProposerPK[:])

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(blockproposerMetadataKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(blockProposerDilithiumMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(metadata.GetDilithiumMetaDataKey(attestorDilithiumPK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(attestorDilithiumMetadataSerialized, nil).AnyTimes()

	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
		blockSigningHash, epochMetadata)
	// stateContext.PrepareAddressState(hex.EncodeToString(binCoinBaseAddress[:]))
	stateContext.PrepareDilithiumMetaData(hex.EncodeToString(blockProposerPK[:]))
	stateContext.PrepareDilithiumMetaData(hex.EncodeToString(attestorDilithiumPK[:]))
	if err != nil {
		t.Error("unexpected error while creating new statecontext ", err)
	}

	testCases := []struct {
		name           string
		attest         *Attest
		stateContext   state.StateContext
		dilithiumKey   *dilithium.Dilithium
		expectedOutput bool
	}{
		{
			name:           "ok",
			attest:         attest,
			stateContext:   *stateContext,
			dilithiumKey:   attestorDilithium,
			expectedOutput: true,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.attest.Sign(tc.dilithiumKey, attest.GetSigningHash(stateContext.PartialBlockSigningHash()))

			output := tc.attest.Validate(stateContext)
			if output != tc.expectedOutput {
				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
			}

		})
	}
}

func TestAttestSetAffectedAddress(t *testing.T) {
	ctrl := gomock.NewController(t)

	networkID := uint64(1)
	coinBaseNonce := uint64(10)
	attest := NewAttest(networkID, coinBaseNonce)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	blockProposerXmss := xmss.NewXMSSFromHeight(4, 0)
	blockProposerXmssPK := blockProposerXmss.GetPK()
	blockProposerAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((blockProposerXmssPK[:])))
	blockProposerDilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), blockProposerPK[:], blockProposerAddr[:], true)
	blockProposerDilithiumMetadataSerialized, _ := blockProposerDilithiumMetadata.Serialize()
	blockProposerAddressState := address.NewAddressState(blockProposerAddr[:], coinBaseNonce, 100)
	blockProposerAddressesStateSerialized, _ := blockProposerAddressState.Serialize()
	blockProposerAddressStateKey := address.GetAddressStateKey(blockProposerAddr[:])

	attestorXmss := xmss.NewXMSSFromHeight(8, 0)
	attestorXmssAddr := attestorXmss.GetAddress()
	attestorDilithium := dilithium.New()
	attestorDilithiumPK := attestorDilithium.GetPK()
	attestorDilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), attestorDilithiumPK[:], attestorXmssAddr[:], true)
	attestorDilithiumMetadataSerialized, _ := attestorDilithiumMetadata.Serialize()
	attestorAddressState := address.NewAddressState(attestorXmssAddr[:], coinBaseNonce, 100)
	attestorAddressesStateSerialized, _ := attestorAddressState.Serialize()
	attestorAddressStateKey := address.GetAddressStateKey(attestorXmssAddr[:])

	var validators [][]byte
	validators = append(validators, blockProposerPK[:])
	validators = append(validators, attestorDilithiumPK[:])
	epoch := uint64(1)
	slotNumber := uint64(100)
	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
	epochMetadaPBData := &protos.EpochMetaData{
		SlotInfo: []*protos.SlotInfo{{
			SlotLeader: 0,
			Attestors:  []uint64{uint64(1)},
		}},
		Validators: validators,
	}
	epochMetadata := metadata.NewEpochMetaData(epoch, nil, nil)
	epochMetadataSerialized, _ := proto.Marshal(epochMetadaPBData)
	epochMetadata.DeSerialize(epochMetadataSerialized)
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
		parentBlockHeaderHash, 0)
	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
	blockproposerMetadataKey := metadata.GetDilithiumMetaDataKey(blockProposerPK[:])

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(blockproposerMetadataKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(blockProposerDilithiumMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(metadata.GetDilithiumMetaDataKey(attestorDilithiumPK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(attestorDilithiumMetadataSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(blockProposerAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(blockProposerAddressesStateSerialized, nil).AnyTimes()
	store.EXPECT().GetFromBucket(gomock.Eq(attestorAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(attestorAddressesStateSerialized, nil).AnyTimes()

	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
		blockSigningHash, epochMetadata)
	if err != nil {
		t.Error("unexpected error while creating new statecontext ", err)
	}
	stateContext.PrepareValidatorsToXMSSAddress(attestorDilithiumPK[:])
	stateContext.PrepareValidatorsToXMSSAddress(blockProposerPK[:])

	attest.Sign(attestorDilithium, attest.GetSigningHash(stateContext.PartialBlockSigningHash()))

	err = attest.SetAffectedAddress(stateContext)
	if err != nil {
		t.Error("got unexpected error while setting affected state of attestor transaction ", err)
	}
}
