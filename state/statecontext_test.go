package state

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/go-qrllib/xmss"
	"github.com/theQRL/zond/address"
	"github.com/theQRL/zond/common"
	mockdb "github.com/theQRL/zond/db/mock"
	"github.com/theQRL/zond/metadata"
	"github.com/theQRL/zond/misc"
	"go.etcd.io/bbolt"
)

func TestProcessValidatorStakeAmount(t *testing.T) {
	dilithium_ := dilithium.New()
	pk := dilithium_.GetPK()

	validator := dilithium.New()
	validatorPK := validator.GetPK()

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	// transaction := xmss.NewXMSSFromHeight(4, 0)
	// transactionPK := transaction.GetPK()

	ctrl := gomock.NewController(t)

	store := mockdb.NewMockDB(ctrl)

	// txaddress := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK(transactionPK[:]))

	// dilithiumState := make(map[string]*metadata.DilithiumMetaData)
	// dilithiumState[hex.EncodeToString(metadata.GetDilithiumMetaDataKey(validatorPK[:]))] = metadata.NewDilithiumMetaData(
	// 	sha256.New().Sum([]byte("transactionHash")),
	// 	validatorPK[:],
	// 	txaddress[:],
	// 	true,
	// )

	// dilithiumState[hex.EncodeToString(metadata.GetDilithiumMetaDataKey(validatorPK[:]))].AddBalance(20)

	stateContext := &StateContext{
		db:                           store,
		currentBlockTotalStakeAmount: big.NewInt(10),
		slotNumber:                   0,
		blockProposer:                blockProposerPK[:],
		finalizedHeaderHash:          common.Hash(sha256.Sum256([]byte("finalizedHeaderHash"))),
		parentBlockHeaderHash:        common.Hash(sha256.Sum256([]byte("parentBlockHeaderHash"))),
		blockHeaderHash:              common.Hash(sha256.Sum256([]byte("blockHeaderHash"))),
		partialBlockSigningHash:      common.Hash(sha256.Sum256([]byte("partialBlockSigningHash"))),
		blockSigningHash:             common.Hash(sha256.Sum256([]byte("blockSigningHash"))),

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: &metadata.MainChainMetaData{},
	}

	testCases := []struct {
		name          string
		dilithiumPK   []byte
		stateContext  StateContext
		expectedError error
	}{
		{
			name:          "ok",
			dilithiumPK:   validatorPK[:],
			stateContext:  *stateContext,
			expectedError: nil,
		},
		{
			name:          "non-validator dilithium key",
			dilithiumPK:   pk[:],
			stateContext:  *stateContext,
			expectedError: fmt.Errorf("validator dilithium state not found for %s", hex.EncodeToString(pk[:])),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := tc.stateContext.processValidatorStakeAmount(tc.dilithiumPK[:], big.NewInt(10))

			if err != nil && (err.Error() != tc.expectedError.Error()) {
				fmt.Printf("error is %s", err.Error())
				x := tc.expectedError.Error()
				t.Errorf("expected error (%s), got error (%s)", x, err.Error())
			}

		})
	}
}

func TestProcessAttestorsFlag(t *testing.T) {
	dilithium_ := dilithium.New()
	pk := dilithium_.GetPK()

	ctrl := gomock.NewController(t)

	store := mockdb.NewMockDB(ctrl)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	attestor1 := dilithium.New()
	attestorDilithiumPK1 := attestor1.GetPK()

	attestor2 := dilithium.New()
	attestorDilithiumPK2 := attestor2.GetPK()

	attestorsFlag := make(map[string]bool)
	strAttestorDilithiumPK1 := hex.EncodeToString(attestorDilithiumPK1[:])
	attestorsFlag[strAttestorDilithiumPK1] = false

	strAttestorDilithiumPK2 := hex.EncodeToString(attestorDilithiumPK2[:])
	attestorsFlag[strAttestorDilithiumPK2] = true

	stateContext := &StateContext{
		db: store,

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     common.Hash(sha256.Sum256([]byte("finalizedHeaderHash"))),
		parentBlockHeaderHash:   common.Hash(sha256.Sum256([]byte("parentBlockHeaderHash"))),
		blockHeaderHash:         common.Hash(sha256.Sum256([]byte("blockHeaderHash"))),
		partialBlockSigningHash: common.Hash(sha256.Sum256([]byte("partialBlockSigningHash"))),
		blockSigningHash:        common.Hash(sha256.Sum256([]byte("blockSigningHash"))),

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: &metadata.MainChainMetaData{},
	}

	testCases := []struct {
		name          string
		dilithiumPK   []byte
		stateContext  StateContext
		expectedError error
	}{
		{
			name:          "ok",
			dilithiumPK:   attestorDilithiumPK1[:],
			stateContext:  *stateContext,
			expectedError: nil,
		},
		{
			name:          "already attested dilithium key",
			dilithiumPK:   attestorDilithiumPK2[:],
			stateContext:  *stateContext,
			expectedError: errors.New("attestor already attested for this slot number"),
		},
		{
			name:          "non-attestor dilithium key",
			dilithiumPK:   pk[:],
			stateContext:  *stateContext,
			expectedError: errors.New("attestor is not assigned to attest at this slot number"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := tc.stateContext.ProcessAttestorsFlag(tc.dilithiumPK, big.NewInt(10))

			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
		})
	}
}

func TestProcessBlockProposerFlag(t *testing.T) {
	dilithium_ := dilithium.New()
	pk := dilithium_.GetPK()

	ctrl := gomock.NewController(t)

	store := mockdb.NewMockDB(ctrl)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	stateContext := &StateContext{
		db:                      store,
		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     common.Hash(sha256.Sum256([]byte("finalizedHeaderHash"))),
		parentBlockHeaderHash:   common.Hash(sha256.Sum256([]byte("parentBlockHeaderHash"))),
		blockHeaderHash:         common.Hash(sha256.Sum256([]byte("blockHeaderHash"))),
		partialBlockSigningHash: common.Hash(sha256.Sum256([]byte("partialBlockSigningHash"))),
		blockSigningHash:        common.Hash(sha256.Sum256([]byte("blockSigningHash"))),

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: &metadata.MainChainMetaData{},
	}

	testCases := []struct {
		name          string
		dilithiumPK   []byte
		stateContext  StateContext
		expectedError error
	}{
		{
			name:          "ok",
			dilithiumPK:   blockProposerPK[:],
			stateContext:  *stateContext,
			expectedError: nil,
		},
		{
			name:          "non block-proposer dilithium key",
			dilithiumPK:   pk[:],
			stateContext:  *stateContext,
			expectedError: errors.New("unexpected block proposer"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := tc.stateContext.ProcessBlockProposerFlag(tc.dilithiumPK, big.NewInt(10))

			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
		})
	}
}

func TestCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	db, err := bbolt.Open("./testdb.txt", 0600, &bbolt.Options{Timeout: 1 * time.Second, InitialMmapSize: 10e6})
	if err != nil {
		t.Error(err.Error())
	}

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	bytesBlock := []byte("byteblock")
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
	blockStorageKey := []byte(fmt.Sprintf("BLOCK-%s", blockHeaderHash))

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	strKey := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((validatorXmssPK[:])))

	validatorsToXMSSAddress := make(map[string][]byte)
	validatorsToXMSSAddress[hex.EncodeToString(validatorDilithiumPK[:])] = validatorXmssPK[:]

	attestor1 := dilithium.New()
	attestorDilithiumPK1 := attestor1.GetPK()

	attestor2 := dilithium.New()
	attestorDilithiumPK2 := attestor2.GetPK()

	attestorsFlag := make(map[string]bool)
	strAttestorDilithiumPK1 := hex.EncodeToString(attestorDilithiumPK1[:])
	attestorsFlag[strAttestorDilithiumPK1] = false

	strAttestorDilithiumPK2 := hex.EncodeToString(attestorDilithiumPK2[:])
	attestorsFlag[strAttestorDilithiumPK2] = true

	slaveXmss := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK := slaveXmss.GetPK()

	mainChainMetaData := metadata.NewMainChainMetaData(common.BytesToHash(sha256.New().Sum([]byte("finalizedBlockHeaderHash"))), 1,
		common.BytesToHash(sha256.New().Sum([]byte("lastBlockHeaderHash"))), 0)
	addressesState := make(map[string]*address.AddressState)
	addressesState[(hex.EncodeToString(address.GetAddressStateKey(strKey[:])))] = address.NewAddressState(strKey[:], 0, 10)
	// dilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithiumPK[:], validatorXmssPK[:], false)
	// dilithiumMetadataSerialized, _ := dilithiumMetadata.Serialize()
	// dilithiumState := make(map[string]*metadata.DilithiumMetaData)
	// dilithiumState[hex.EncodeToString(metadata.GetDilithiumMetaDataKey(validatorDilithiumPK[:]))] = dilithiumMetadata

	slaveMetadata := metadata.NewSlaveMetaData(sha256.New().Sum([]byte("transactionHash")), strKey[:], slaveXmssPK[:])
	slaveMetadataSerialized, _ := slaveMetadata.Serialize()
	slaveState := make(map[string]*metadata.SlaveMetaData)
	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK[:]))] = slaveMetadata

	// otsIndex := uint64(8192)

	// otsIndexMetadata := metadata.NewOTSIndexMetaData(strKey[:], otsIndex/config.GetDevConfig().OTSBitFieldPerPage)
	// otsIndexMetadataSerialized, _ := otsIndexMetadata.Serialize()
	// otsIndexState := make(map[string]*metadata.OTSIndexMetaData)
	// otsIndexState[hex.EncodeToString(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex))] = otsIndexMetadata

	totalStakeAmount, _ := big.NewInt(10).MarshalText()
	parentBlockMetadata := metadata.NewBlockMetaData(common.BytesToHash(sha256.New().Sum([]byte("parentsparentBlockHeaderHash"))), common.BytesToHash(sha256.New().Sum([]byte("parentBlockHeaderHash"))), 0, totalStakeAmount, common.Hash{})
	parentBlockMetadataSerialized, _ := parentBlockMetadata.Serialize()

	lastBlockMetadata := metadata.NewBlockMetaData(common.BytesToHash(sha256.New().Sum([]byte("parentBlockHeaderHash"))), common.BytesToHash(sha256.New().Sum([]byte("lastBlockHeaderHash"))), 0, totalStakeAmount, common.Hash{})
	lastBlockMetadataSerialized, _ := lastBlockMetadata.Serialize()

	store := mockdb.NewMockDB(ctrl)

	stateContext := &StateContext{
		db: store,

		slotNumber:                   1,
		blockProposer:                blockProposerPK[:],
		finalizedHeaderHash:          common.Hash(sha256.Sum256([]byte("finalizedHeaderHash"))),
		parentBlockHeaderHash:        common.Hash(sha256.Sum256([]byte("parentBlockHeaderHash"))),
		blockHeaderHash:              common.Hash(sha256.Sum256([]byte("blockHeaderHash"))),
		partialBlockSigningHash:      common.Hash(sha256.Sum256([]byte("partialBlockSigningHash"))),
		blockSigningHash:             common.Hash(sha256.Sum256([]byte("blockSigningHash"))),
		currentBlockTotalStakeAmount: big.NewInt(10),

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name             string
		blockStorageKey  []byte
		bytesBlock       []byte
		isFinalizedState bool
		stateContext     StateContext
		buildStubs       func(store *mockdb.MockDB)
		expectedError    error
	}{
		{
			name:             "ok",
			blockStorageKey:  blockStorageKey,
			bytesBlock:       bytesBlock,
			isFinalizedState: false,
			stateContext:     *stateContext,
			buildStubs: func(store *mockdb.MockDB) {
				// store.EXPECT().
				// 	GetFromBucket(gomock.Eq([]byte(fmt.Sprintf("DILITHIUM-META-DATA-%s", hex.EncodeToString(validatorDilithiumPK[:])))), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
				// 	Return(dilithiumMetadataSerialized, nil).AnyTimes()
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK[:])), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(slaveMetadataSerialized, nil).AnyTimes()
				// store.EXPECT().
				// 	GetFromBucket(gomock.Eq(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex)), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
				// 	Return(otsIndexMetadataSerialized, nil).AnyTimes()
				store.EXPECT().DB().Return(db).AnyTimes()
				store.EXPECT().Get(gomock.Eq(metadata.GetBlockMetaDataKey(common.BytesToHash(sha256.New().Sum([]byte("parentBlockHeaderHash")))))).Return(parentBlockMetadataSerialized, nil).AnyTimes()
				store.EXPECT().Get(gomock.Eq(metadata.GetBlockMetaDataKey(common.BytesToHash(sha256.New().Sum([]byte("lastBlockHeaderHash")))))).Return(lastBlockMetadataSerialized, nil).AnyTimes()
			},
			expectedError: nil,
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)

			err := store.DB().Update(func(tx *bbolt.Tx) error {
				mainBucket := tx.Bucket([]byte("DB"))
				if mainBucket == nil {
					_, err := tx.CreateBucket([]byte("DB"))
					if err != nil {
						return fmt.Errorf("create bucket: %s", err)
					}
					return nil
				}
				return nil
			})
			if err != nil {
				t.Error("error creating bucket", err)
			}
			err = tc.stateContext.Commit(tc.blockStorageKey, tc.bytesBlock, common.Hash{}, tc.isFinalizedState)
			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
		})
	}
}

func TestFinalize(t *testing.T) {
	ctrl := gomock.NewController(t)
	db, err := bbolt.Open("./testdb2.txt", 0600, &bbolt.Options{Timeout: 1 * time.Second, InitialMmapSize: 10e6})
	if err != nil {
		t.Error(err.Error())
	}

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	totalStakeAmount, _ := big.NewInt(10).MarshalText()
	parentBlockMetadata := metadata.NewBlockMetaData(common.BytesToHash(sha256.New().Sum([]byte("parentsparentBlockHeaderHash"))), common.BytesToHash(sha256.New().Sum([]byte("parentBlockHeaderHash"))), 0, totalStakeAmount, common.Hash{})
	parentBlockMetadataSerialized, _ := parentBlockMetadata.Serialize()

	lastBlockMetadata := metadata.NewBlockMetaData(common.BytesToHash(sha256.New().Sum([]byte("parentBlockHeaderHash"))), common.BytesToHash(sha256.New().Sum([]byte("lastBlockHeaderHash"))), 0, totalStakeAmount, common.Hash{})
	lastBlockMetadataSerialized, _ := lastBlockMetadata.Serialize()

	blockMetaDataPathForFinalization := make([]*metadata.BlockMetaData, 0)
	blockMetaDataPathForFinalization = append(blockMetaDataPathForFinalization, lastBlockMetadata)

	mainChainMetaData := metadata.NewMainChainMetaData(common.BytesToHash(sha256.New().Sum([]byte("finalizedBlockHeaderHash"))), 1,
		common.BytesToHash(sha256.New().Sum([]byte("parentBlockHeaderHash"))), 0)
	store := mockdb.NewMockDB(ctrl)

	stateContext := &StateContext{
		db: store,

		slotNumber:                   1,
		blockProposer:                blockProposerPK[:],
		finalizedHeaderHash:          common.Hash(sha256.Sum256([]byte("finalizedHeaderHash"))),
		parentBlockHeaderHash:        common.Hash(sha256.Sum256([]byte("parentBlockHeaderHash"))),
		blockHeaderHash:              common.Hash(sha256.Sum256([]byte("blockHeaderHash"))),
		partialBlockSigningHash:      common.Hash(sha256.Sum256([]byte("partialBlockSigningHash"))),
		blockSigningHash:             common.Hash(sha256.Sum256([]byte("blockSigningHash"))),
		currentBlockTotalStakeAmount: big.NewInt(10),

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name                   string
		blockMetaDataPathArray []*metadata.BlockMetaData
		stateContext           StateContext
		buildStubs             func(store *mockdb.MockDB)
		expectedError          error
	}{
		{
			name:                   "ok",
			blockMetaDataPathArray: blockMetaDataPathForFinalization,
			stateContext:           *stateContext,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().Get(gomock.Eq(metadata.GetBlockMetaDataKey(common.BytesToHash(sha256.New().Sum([]byte("parentBlockHeaderHash")))))).Return(parentBlockMetadataSerialized, nil).AnyTimes()
				store.EXPECT().Get(gomock.Eq(metadata.GetBlockMetaDataKey(common.BytesToHash(sha256.New().Sum([]byte("lastBlockHeaderHash")))))).Return(lastBlockMetadataSerialized, nil).AnyTimes()
				store.EXPECT().DB().Return(db).AnyTimes()
			},
			expectedError: nil,
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)

			err := store.DB().Update(func(tx *bbolt.Tx) error {
				mainBucket := tx.Bucket([]byte("DB"))
				if mainBucket == nil {
					_, err := tx.CreateBucket([]byte("DB"))
					if err != nil {
						return fmt.Errorf("create bucket: %s", err)
					}
					return nil
				}
				for i := len(tc.blockMetaDataPathArray) - 1; i >= 0; i-- {
					bm := tc.blockMetaDataPathArray[i]
					blockBucket := tx.Bucket(metadata.GetBlockBucketName(bm.HeaderHash()))
					if blockBucket == nil {
						_, err := tx.CreateBucket(metadata.GetBlockBucketName(bm.HeaderHash()))
						if err != nil {
							return fmt.Errorf("error to create bucket: %s", err)
						}
						return nil
					}
				}
				return nil
			})
			if err != nil {
				t.Error("error creating bucket", err)
			}
			err = tc.stateContext.Finalize(tc.blockMetaDataPathArray)
			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
		})
	}
}

func TestNewStateContext(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	slotNumber := uint64(0)
	finalizedHeaderHash := common.Hash(sha256.Sum256([]byte("finalizedHeaderHash")))
	parentBlockHeaderHash := common.Hash(sha256.Sum256([]byte("parentBlockHeaderHash")))
	blockHeaderHash := common.Hash(sha256.Sum256([]byte("blockHeaderHash")))
	partialBlockSigningHash := common.Hash(sha256.Sum256([]byte("partialBlockSigningHash")))
	blockSigningHash := common.Hash(sha256.Sum256([]byte("blockSigningHash")))
	epochMetaData := &metadata.EpochMetaData{}
	mainChainMetaData := metadata.NewMainChainMetaData(common.Hash(sha256.Sum256([]byte("finalizedBlockHeaderHash"))), 1,
		common.Hash(sha256.Sum256([]byte("parentBlockHeaderHash"))), 0)
	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()

	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(0)
	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(0))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()

	newStateContext, err := NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash,
		partialBlockSigningHash, blockSigningHash, epochMetaData)

	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if newStateContext.GetSlotNumber() != slotNumber {
		t.Errorf("expected slotnumber (%v) got (%v)", slotNumber, newStateContext.GetSlotNumber())
	}
}
