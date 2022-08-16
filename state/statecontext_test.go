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
	"github.com/theQRL/zond/config"
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

	transaction := xmss.NewXMSSFromHeight(4, 0)
	transactionPK := transaction.GetPK()

	ctrl := gomock.NewController(t)

	store := mockdb.NewMockDB(ctrl)

	txaddress := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK(transactionPK[:]))

	dilithiumState := make(map[string]*metadata.DilithiumMetaData)
	dilithiumState[hex.EncodeToString(metadata.GetDilithiumMetaDataKey(validatorPK[:]))] = metadata.NewDilithiumMetaData(
		sha256.New().Sum([]byte("transactionHash")),
		validatorPK[:],
		txaddress[:],
		true,
	)

	dilithiumState[hex.EncodeToString(metadata.GetDilithiumMetaDataKey(validatorPK[:]))].AddBalance(20)

	stateContext := &StateContext{
		db:                           store,
		addressesState:               make(map[string]*address.AddressState),
		dilithiumState:               dilithiumState,
		slaveState:                   make(map[string]*metadata.SlaveMetaData),
		otsIndexState:                make(map[string]*metadata.OTSIndexMetaData),
		currentBlockTotalStakeAmount: 10,
		slotNumber:                   0,
		blockProposer:                blockProposerPK[:],
		finalizedHeaderHash:          sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:        sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:              sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash:      sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:             sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress:      make(map[string][]byte),
		attestorsFlag:                make(map[string]bool),
		blockProposerFlag:            false,

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
			initialCurrentStakeAmount := tc.stateContext.currentBlockTotalStakeAmount
			err := tc.stateContext.processValidatorStakeAmount(tc.dilithiumPK[:])

			if err != nil && (err.Error() != tc.expectedError.Error()) {
				fmt.Printf("error is %s", err.Error())
				x := tc.expectedError.Error()
				t.Errorf("expected error (%s), got error (%s)", x, err.Error())
			}

			strKey := hex.EncodeToString(metadata.GetDilithiumMetaDataKey(tc.dilithiumPK))
			slotLeaderDilithiumMetaData, ok := tc.stateContext.dilithiumState[strKey]

			if ok && (initialCurrentStakeAmount+slotLeaderDilithiumMetaData.Balance() != tc.stateContext.currentBlockTotalStakeAmount) {
				t.Errorf("expected total block stake amount (%d), got (%d)", initialCurrentStakeAmount+slotLeaderDilithiumMetaData.Balance(), tc.stateContext.currentBlockTotalStakeAmount)
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
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           attestorsFlag,
		blockProposerFlag:       false,

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
			err := tc.stateContext.ProcessAttestorsFlag(tc.dilithiumPK)

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
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: &metadata.MainChainMetaData{},
	}

	stateContext2 := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       true,

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
		{
			name:          "block proposer already proposed",
			dilithiumPK:   blockProposerPK[:],
			stateContext:  *stateContext2,
			expectedError: errors.New("block proposer has already been processed"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := tc.stateContext.ProcessBlockProposerFlag(tc.dilithiumPK)

			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
		})
	}
}

func TestPrepareAddressState(t *testing.T) {
	ctrl := gomock.NewController(t)
	db, _ := bbolt.Open("./", 0600, &bbolt.Options{Timeout: 1 * time.Second, InitialMmapSize: 10e6})

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()

	validatorsToXMSSAddress := make(map[string][]byte)
	validatorsToXMSSAddress[hex.EncodeToString(validatorDilithiumPK[:])] = validatorXmssPK[:]

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	addressesState := make(map[string]*address.AddressState)
	strKey := hex.EncodeToString(address.GetAddressStateKey(validatorXmssPK[:]))

	addressesState[strKey] = address.NewAddressState(validatorXmssPK[:], 0, 10)
	serialized_address, _ := address.NewAddressState(validatorXmssPK[:], 0, 10).Serialize()
	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().GetFromBucket(gomock.Any(), gomock.Any()).Return(serialized_address, nil).AnyTimes()
	store.EXPECT().Get(gomock.Any()).AnyTimes()
	store.EXPECT().DB().Return(db).AnyTimes()

	stateContext := &StateContext{
		db:             store,
		addressesState: addressesState,
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: validatorsToXMSSAddress,
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}

	testCases := []struct {
		name          string
		dilithiumPK   []byte
		stateContext  StateContext
		expectedError error
	}{
		{
			name:          "ok",
			dilithiumPK:   validatorDilithiumPK[:],
			stateContext:  *stateContext,
			expectedError: nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := tc.stateContext.PrepareAddressState(hex.EncodeToString(validatorXmssPK[:]))

			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}

			strKey := hex.EncodeToString(address.GetAddressStateKey(validatorXmssPK[:]))
			addressesState := tc.stateContext.addressesState[strKey]

			expectedAddressState, err := address.GetAddressState(tc.stateContext.db, validatorXmssPK[:],
				tc.stateContext.parentBlockHeaderHash, tc.stateContext.mainChainMetaData.FinalizedBlockHeaderHash())
			if err != nil {
				t.Errorf("Address state not set for address %s", hex.EncodeToString(validatorXmssPK[:]))
			}

			if string(addressesState.Address()) != string(expectedAddressState.Address()) {
				t.Errorf("expected addressstate addresses does not match")
			}

			if addressesState.Balance() != expectedAddressState.Balance() {
				t.Errorf("expected addressstate balance does not match")
			}
		})
	}
}

func TestGetAddressState(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()

	randomXmss := xmss.NewXMSSFromHeight(10, 0)
	randomXmssPK := randomXmss.GetPK()

	validatorsToXMSSAddress := make(map[string][]byte)
	validatorsToXMSSAddress[hex.EncodeToString(validatorDilithiumPK[:])] = validatorXmssPK[:]

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	addressesState := make(map[string]*address.AddressState)
	strKey := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	addressesState[(hex.EncodeToString(address.GetAddressStateKey(strKey[:])))] = address.NewAddressState(strKey[:], 0, 10)

	serialized_address, _ := address.NewAddressState(strKey[:], 0, 10).Serialize()
	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().GetFromBucket(gomock.Any(), gomock.Any()).Return(serialized_address, nil).AnyTimes()
	store.EXPECT().Get(gomock.Any()).AnyTimes()

	stateContext := &StateContext{
		db:             store,
		addressesState: addressesState,
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: validatorsToXMSSAddress,
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}

	testCases := []struct {
		name                 string
		xmssPK               []byte
		stateContext         StateContext
		expectedError        error
		expectedAddressState *address.AddressState
	}{
		{
			name:                 "ok",
			xmssPK:               strKey[:],
			stateContext:         *stateContext,
			expectedAddressState: address.NewAddressState(strKey[:], 0, 10),
			expectedError:        nil,
		},
		{
			name:                 "address does not exist",
			xmssPK:               randomXmssPK[:],
			stateContext:         *stateContext,
			expectedAddressState: nil,
			expectedError:        fmt.Errorf("address %s not found in addressesState", hex.EncodeToString(randomXmssPK[:])),
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			addressState, err := tc.stateContext.GetAddressState(hex.EncodeToString(tc.xmssPK))

			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}

			if err == nil && (string(addressState.Address()) != string(tc.expectedAddressState.Address())) {
				t.Errorf("expected addressState addresses does not match")
			}

			if err == nil && (addressState.Balance() != tc.expectedAddressState.Balance()) {
				t.Errorf("expected addressState balance does not match")
			}
		})
	}
}

func TestGetAddressStateByPK(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()

	validatorsToXMSSAddress := make(map[string][]byte)
	validatorsToXMSSAddress[hex.EncodeToString(validatorDilithiumPK[:])] = validatorXmssPK[:]

	addressesState := make(map[string]*address.AddressState)
	strKey := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	addressesState[(hex.EncodeToString(address.GetAddressStateKey(strKey[:])))] = address.NewAddressState(strKey[:], 0, 10)
	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	serialized_address, _ := address.NewAddressState(strKey[:], 0, 10).Serialize()
	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().GetFromBucket(gomock.Any(), gomock.Any()).Return(serialized_address, nil).AnyTimes()
	store.EXPECT().Get(gomock.Any()).AnyTimes()

	stateContext := &StateContext{
		db:             store,
		addressesState: addressesState,
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: validatorsToXMSSAddress,
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name                 string
		xmssPK               []byte
		stateContext         StateContext
		expectedError        error
		expectedAddressState *address.AddressState
	}{
		{
			name:                 "ok",
			xmssPK:               validatorXmssPK[:],
			stateContext:         *stateContext,
			expectedAddressState: address.NewAddressState(strKey[:], 0, 10),
			expectedError:        nil,
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			addressState, err := tc.stateContext.GetAddressStateByPK(tc.xmssPK)
			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}

			if err == nil && (string(addressState.Address()) != string(tc.expectedAddressState.Address())) {
				t.Errorf("expected addressState addresses does not match")
			}

			if err == nil && (addressState.Balance() != tc.expectedAddressState.Balance()) {
				t.Errorf("expected addressState balance does not match")
			}
		})
	}
}

func TestPrepareValidatorsToXMSSAddress(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()

	store := mockdb.NewMockDB(ctrl)

	dilithiumMetadata, _ := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithiumPK[:], validatorXmssPK[:], false).Serialize()

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: &metadata.MainChainMetaData{},
	}

	testCases := []struct {
		name          string
		dilithiumPK   []byte
		stateContext  StateContext
		buildStubs    func(store *mockdb.MockDB)
		expectedError error
	}{
		{
			name:         "ok",
			dilithiumPK:  validatorDilithiumPK[:],
			stateContext: *stateContext,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq([]byte(fmt.Sprintf("DILITHIUM-META-DATA-%s", hex.EncodeToString(validatorDilithiumPK[:])))), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(dilithiumMetadata, nil).AnyTimes()
			},
			expectedError: nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			err := tc.stateContext.PrepareValidatorsToXMSSAddress(tc.dilithiumPK)

			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}

		})
	}
}

func TestPrepareDilithiumMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()

	store := mockdb.NewMockDB(ctrl)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	dilithiumMetadata, _ := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithiumPK[:], validatorXmssPK[:], false).Serialize()

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name          string
		dilithiumPK   []byte
		stateContext  StateContext
		buildStubs    func(store *mockdb.MockDB)
		expectedError error
	}{
		{
			name:         "ok",
			dilithiumPK:  validatorDilithiumPK[:],
			stateContext: *stateContext,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq([]byte(fmt.Sprintf("DILITHIUM-META-DATA-%s", hex.EncodeToString(validatorDilithiumPK[:])))), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(dilithiumMetadata, nil).AnyTimes()
			},
			expectedError: nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			err := tc.stateContext.PrepareDilithiumMetaData(hex.EncodeToString(tc.dilithiumPK))
			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}

			if err == nil && (string(tc.stateContext.dilithiumState[hex.EncodeToString(metadata.GetDilithiumMetaDataKey(tc.dilithiumPK))].Address()) != string(validatorXmssPK[:])) {
				t.Errorf("expected address does not match")
			}
		})
	}
}

func TestAddDilithiumMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorDilithium2 := dilithium.New()
	validatorDilithiumPK2 := validatorDilithium2.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()

	store := mockdb.NewMockDB(ctrl)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	dilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithiumPK[:], validatorXmssPK[:], false)
	dilithiumMetadataSerialized, _ := dilithiumMetadata.Serialize()
	dilithiumMetadata2 := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithiumPK2[:], validatorXmssPK[:], false)
	dilithiumMetadata2Serialized, _ := dilithiumMetadata2.Serialize()
	dilithiumState := make(map[string]*metadata.DilithiumMetaData)
	dilithiumState[hex.EncodeToString(metadata.GetDilithiumMetaDataKey(validatorDilithiumPK2[:]))] = dilithiumMetadata2

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: dilithiumState,
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name              string
		dilithiumPK       []byte
		stateContext      StateContext
		dilithiumMetadata *metadata.DilithiumMetaData
		buildStubs        func(store *mockdb.MockDB)
		expectedError     error
	}{
		{
			name:              "ok",
			dilithiumPK:       validatorDilithiumPK[:],
			stateContext:      *stateContext,
			dilithiumMetadata: dilithiumMetadata,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq([]byte(fmt.Sprintf("DILITHIUM-META-DATA-%s", hex.EncodeToString(validatorDilithiumPK[:])))), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(dilithiumMetadataSerialized, nil).AnyTimes()
			},
			expectedError: nil,
		},
		{
			name:              "meta data already exists",
			dilithiumPK:       validatorDilithiumPK2[:],
			stateContext:      *stateContext,
			dilithiumMetadata: dilithiumMetadata2,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq([]byte(fmt.Sprintf("DILITHIUM-META-DATA-%s", hex.EncodeToString(validatorDilithiumPK2[:])))), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(dilithiumMetadata2Serialized, nil).AnyTimes()
			},
			expectedError: errors.New("dilithiumPK already exists"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			err := tc.stateContext.AddDilithiumMetaData(hex.EncodeToString(tc.dilithiumPK), tc.dilithiumMetadata)
			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}

			if err == nil && (string(tc.stateContext.dilithiumState[hex.EncodeToString(metadata.GetDilithiumMetaDataKey(tc.dilithiumPK))].Address()) != string(validatorXmssPK[:])) {
				t.Errorf("expected address does not match")
			}
		})
	}
}

func TestGetDilithiumState(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()

	store := mockdb.NewMockDB(ctrl)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	dilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithiumPK[:], validatorXmssPK[:], false)
	dilithiumState := make(map[string]*metadata.DilithiumMetaData)
	dilithiumState[hex.EncodeToString(metadata.GetDilithiumMetaDataKey(validatorDilithiumPK[:]))] = dilithiumMetadata

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: dilithiumState,
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name          string
		dilithiumPK   []byte
		stateContext  StateContext
		expectedError error
	}{
		{
			name:          "ok",
			dilithiumPK:   validatorDilithiumPK[:],
			stateContext:  *stateContext,
			expectedError: nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			dilithiumstate := tc.stateContext.GetDilithiumState(hex.EncodeToString(tc.dilithiumPK))

			if string(dilithiumstate.Address()) != string(validatorXmssPK[:]) {
				t.Errorf("expected address does not match")
			}
		})
	}
}

func TestPrepareSlaveMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	strKey := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	slaveXmss := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK := slaveXmss.GetPK()

	store := mockdb.NewMockDB(ctrl)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	slaveMetadata := metadata.NewSlaveMetaData(sha256.New().Sum([]byte("transactionHash")), strKey[:], slaveXmssPK[:])
	slaveMetadataSerialized, _ := slaveMetadata.Serialize()

	slaveState := make(map[string]*metadata.SlaveMetaData)
	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK[:]))] = slaveMetadata

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name          string
		masterAddr    string
		slaveXmssPK   string
		stateContext  StateContext
		slaveMetadata *metadata.SlaveMetaData
		buildStubs    func(store *mockdb.MockDB)
		expectedError error
	}{
		{
			name:          "ok",
			masterAddr:    hex.EncodeToString(strKey[:]),
			slaveXmssPK:   hex.EncodeToString(slaveXmssPK[:]),
			stateContext:  *stateContext,
			slaveMetadata: slaveMetadata,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK[:])), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(slaveMetadataSerialized, nil).AnyTimes()
			},
			expectedError: nil,
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			err := tc.stateContext.PrepareSlaveMetaData(tc.masterAddr, tc.slaveXmssPK)
			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
			slaveXmssPKBin, _ := hex.DecodeString(tc.slaveXmssPK)
			masterAddrBin, _ := hex.DecodeString(tc.masterAddr)
			if err == nil && (string(tc.stateContext.slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(masterAddrBin, slaveXmssPKBin))].SlavePK()) != string(slaveXmssPK[:])) {
				t.Errorf("expected address does not match")
			}
		})
	}
}

func TestAddSlaveMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	strKey := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	slaveXmss := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK := slaveXmss.GetPK()

	slaveXmss2 := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK2 := slaveXmss2.GetPK()

	store := mockdb.NewMockDB(ctrl)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	slaveMetadata := metadata.NewSlaveMetaData(sha256.New().Sum([]byte("transactionHash")), strKey[:], slaveXmssPK[:])
	slaveMetadataSerialized, _ := slaveMetadata.Serialize()

	slaveMetadata2 := metadata.NewSlaveMetaData(sha256.New().Sum([]byte("transactionHash")), strKey[:], slaveXmssPK[:])
	slaveMetadataSerialized2, _ := slaveMetadata2.Serialize()

	slaveState := make(map[string]*metadata.SlaveMetaData)
	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK2[:]))] = slaveMetadata2

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name          string
		masterAddr    string
		slaveXmssPK   string
		stateContext  StateContext
		slaveMetadata *metadata.SlaveMetaData
		buildStubs    func(store *mockdb.MockDB)
		expectedError error
	}{
		{
			name:          "ok",
			masterAddr:    hex.EncodeToString(strKey[:]),
			slaveXmssPK:   hex.EncodeToString(slaveXmssPK[:]),
			stateContext:  *stateContext,
			slaveMetadata: slaveMetadata,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK[:])), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(slaveMetadataSerialized, nil).AnyTimes()
			},
			expectedError: nil,
		},
		{
			name:          "slave metadata already exists",
			masterAddr:    hex.EncodeToString(strKey[:]),
			slaveXmssPK:   hex.EncodeToString(slaveXmssPK[:]),
			stateContext:  *stateContext,
			slaveMetadata: slaveMetadata,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK2[:])), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(slaveMetadataSerialized2, nil).AnyTimes()
			},
			expectedError: errors.New("SlaveMetaData already exists"),
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			err := tc.stateContext.AddSlaveMetaData(tc.masterAddr, tc.slaveXmssPK, tc.slaveMetadata)
			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
			slaveXmssPKBin, _ := hex.DecodeString(tc.slaveXmssPK)
			masterAddrBin, _ := hex.DecodeString(tc.masterAddr)
			if err == nil && (string(tc.stateContext.slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(masterAddrBin, slaveXmssPKBin))].SlavePK()) != string(slaveXmssPKBin[:])) {
				t.Errorf("expected address does not match")
			}
		})
	}
}

func TestGetSlaveState(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	strKey := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	slaveXmss := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK := slaveXmss.GetPK()

	store := mockdb.NewMockDB(ctrl)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	slaveMetadata := metadata.NewSlaveMetaData(sha256.New().Sum([]byte("transactionHash")), strKey[:], slaveXmssPK[:])
	slaveMetadataSerialized, _ := slaveMetadata.Serialize()

	slaveState := make(map[string]*metadata.SlaveMetaData)
	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK[:]))] = slaveMetadata

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     slaveState,
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}

	testCases := []struct {
		name          string
		masterAddr    string
		slaveXmssPK   string
		stateContext  StateContext
		slaveMetadata *metadata.SlaveMetaData
		buildStubs    func(store *mockdb.MockDB)
		expectedError error
	}{
		{
			name:          "ok",
			masterAddr:    hex.EncodeToString(strKey[:]),
			slaveXmssPK:   hex.EncodeToString(slaveXmssPK[:]),
			stateContext:  *stateContext,
			slaveMetadata: slaveMetadata,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK[:])), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(slaveMetadataSerialized, nil).AnyTimes()
			},
			expectedError: nil,
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			slavemetadata := tc.stateContext.GetSlaveState(tc.masterAddr, tc.slaveXmssPK)

			slaveXmssPKBin, _ := hex.DecodeString(tc.slaveXmssPK)

			if string(slavemetadata.SlavePK()) != string(slaveXmssPKBin[:]) {
				t.Errorf("expected address does not match")
			}
		})
	}
}

func TestPrepareOTSIndexMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	strKey := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	store := mockdb.NewMockDB(ctrl)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	otsIndex := uint64(8192)
	otsIndexMetadata := metadata.NewOTSIndexMetaData(strKey[:], otsIndex/config.GetDevConfig().OTSBitFieldPerPage)
	otsIndexMetadataSerialized, _ := otsIndexMetadata.Serialize()
	otsIndexState := make(map[string]*metadata.OTSIndexMetaData)
	otsIndexState[hex.EncodeToString(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex))] = otsIndexMetadata

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name             string
		address          string
		otsIndex         uint64
		stateContext     StateContext
		otsIndexMetadata *metadata.OTSIndexMetaData
		buildStubs       func(store *mockdb.MockDB)
		expectedError    error
	}{
		{
			name:             "ok",
			address:          hex.EncodeToString(strKey[:]),
			otsIndex:         uint64(8192),
			stateContext:     *stateContext,
			otsIndexMetadata: otsIndexMetadata,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex)), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(otsIndexMetadataSerialized, nil).AnyTimes()
			},
			expectedError: nil,
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			err := tc.stateContext.PrepareOTSIndexMetaData(tc.address, tc.otsIndex)
			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
			addressBin, _ := hex.DecodeString(tc.address)

			if err == nil && (string(tc.stateContext.otsIndexState[hex.EncodeToString(metadata.GetOTSIndexMetaDataKeyByOTSIndex(addressBin, tc.otsIndex))].Address()) != string(addressBin[:])) {
				t.Errorf("expected address does not match")
			}
		})
	}
}

func TestAddOTSIndexMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	strKey := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	store := mockdb.NewMockDB(ctrl)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	otsIndex := uint64(8192)
	otsIndex2 := uint64(8192 * 2)
	otsIndexMetadata := metadata.NewOTSIndexMetaData(strKey[:], otsIndex/config.GetDevConfig().OTSBitFieldPerPage)
	otsIndexMetadataSerialized, _ := otsIndexMetadata.Serialize()
	otsIndexMetadata2 := metadata.NewOTSIndexMetaData(strKey[:], otsIndex2/config.GetDevConfig().OTSBitFieldPerPage)
	otsIndexState := make(map[string]*metadata.OTSIndexMetaData)
	otsIndexState[hex.EncodeToString(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex))] = otsIndexMetadata2

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name             string
		address          string
		otsIndex         uint64
		stateContext     StateContext
		otsIndexMetadata *metadata.OTSIndexMetaData
		buildStubs       func(store *mockdb.MockDB)
		expectedError    error
	}{
		{
			name:             "ok",
			address:          hex.EncodeToString(strKey[:]),
			otsIndex:         uint64(8192),
			stateContext:     *stateContext,
			otsIndexMetadata: otsIndexMetadata,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex)), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(otsIndexMetadataSerialized, nil).AnyTimes()
			},
			expectedError: nil,
		},
		{
			name:             "ots metadata already exists",
			address:          hex.EncodeToString(strKey[:]),
			otsIndex:         uint64(8192 * 2),
			stateContext:     *stateContext,
			otsIndexMetadata: otsIndexMetadata,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex)), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(otsIndexMetadataSerialized, nil).AnyTimes()
			},
			expectedError: errors.New("OTSIndexMetaData already exists"),
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			err := tc.stateContext.AddOTSIndexMetaData(tc.address, tc.otsIndex, tc.otsIndexMetadata)
			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
			addressBin, _ := hex.DecodeString(tc.address)

			if err == nil && (string(tc.stateContext.otsIndexState[hex.EncodeToString(metadata.GetOTSIndexMetaDataKeyByOTSIndex(addressBin, tc.otsIndex))].Address()) != string(addressBin[:])) {
				t.Errorf("expected address does not match")
			}
		})
	}
}

func TestGetOTSIndexState(t *testing.T) {
	ctrl := gomock.NewController(t)

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	strKey := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	store := mockdb.NewMockDB(ctrl)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)

	otsIndex := uint64(8192)
	otsIndexMetadata := metadata.NewOTSIndexMetaData(strKey[:], otsIndex/config.GetDevConfig().OTSBitFieldPerPage)
	otsIndexMetadataSerialized, _ := otsIndexMetadata.Serialize()
	otsIndexState := make(map[string]*metadata.OTSIndexMetaData)
	otsIndexState[hex.EncodeToString(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex))] = otsIndexMetadata

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  otsIndexState,

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name             string
		address          string
		otsIndex         uint64
		stateContext     StateContext
		otsIndexMetadata *metadata.OTSIndexMetaData
		buildStubs       func(store *mockdb.MockDB)
		expectedError    error
	}{
		{
			name:             "ok",
			address:          hex.EncodeToString(strKey[:]),
			otsIndex:         uint64(8192),
			stateContext:     *stateContext,
			otsIndexMetadata: otsIndexMetadata,
			buildStubs: func(store *mockdb.MockDB) {
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex)), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(otsIndexMetadataSerialized, nil).AnyTimes()
			},
			expectedError: nil,
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			tc.buildStubs(store)
			otsIndexstate := tc.stateContext.GetOTSIndexState(tc.address, tc.otsIndex)
			// if err != nil && (err.Error() != tc.expectedError.Error()) {
			// 	t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			// }
			addressBin, _ := hex.DecodeString(tc.address)

			if string(otsIndexstate.Address()) != string(addressBin[:]) {
				t.Errorf("expected address does not match")
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
	strKey := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

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

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("lastBlockHeaderHash")), 0)
	addressesState := make(map[string]*address.AddressState)
	addressesState[(hex.EncodeToString(address.GetAddressStateKey(strKey[:])))] = address.NewAddressState(strKey[:], 0, 10)
	dilithiumMetadata := metadata.NewDilithiumMetaData(sha256.New().Sum([]byte("transactionHash")), validatorDilithiumPK[:], validatorXmssPK[:], false)
	dilithiumMetadataSerialized, _ := dilithiumMetadata.Serialize()
	dilithiumState := make(map[string]*metadata.DilithiumMetaData)
	dilithiumState[hex.EncodeToString(metadata.GetDilithiumMetaDataKey(validatorDilithiumPK[:]))] = dilithiumMetadata

	slaveMetadata := metadata.NewSlaveMetaData(sha256.New().Sum([]byte("transactionHash")), strKey[:], slaveXmssPK[:])
	slaveMetadataSerialized, _ := slaveMetadata.Serialize()
	slaveState := make(map[string]*metadata.SlaveMetaData)
	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK[:]))] = slaveMetadata

	otsIndex := uint64(8192)

	otsIndexMetadata := metadata.NewOTSIndexMetaData(strKey[:], otsIndex/config.GetDevConfig().OTSBitFieldPerPage)
	otsIndexMetadataSerialized, _ := otsIndexMetadata.Serialize()
	otsIndexState := make(map[string]*metadata.OTSIndexMetaData)
	otsIndexState[hex.EncodeToString(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex))] = otsIndexMetadata

	totalStakeAmount, _ := big.NewInt(10).MarshalText()
	parentBlockMetadata := metadata.NewBlockMetaData(sha256.New().Sum([]byte("parentsparentBlockHeaderHash")), sha256.New().Sum([]byte("parentBlockHeaderHash")), 0, totalStakeAmount)
	parentBlockMetadataSerialized, _ := parentBlockMetadata.Serialize()

	lastBlockMetadata := metadata.NewBlockMetaData(sha256.New().Sum([]byte("parentBlockHeaderHash")), sha256.New().Sum([]byte("lastBlockHeaderHash")), 0, totalStakeAmount)
	lastBlockMetadataSerialized, _ := lastBlockMetadata.Serialize()

	store := mockdb.NewMockDB(ctrl)

	stateContext := &StateContext{
		db:             store,
		addressesState: addressesState,
		dilithiumState: dilithiumState,
		slaveState:     slaveState,
		otsIndexState:  otsIndexState,

		slotNumber:              1,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: validatorsToXMSSAddress,
		attestorsFlag:           attestorsFlag,
		blockProposerFlag:       false,

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
				store.EXPECT().
					GetFromBucket(gomock.Eq([]byte(fmt.Sprintf("DILITHIUM-META-DATA-%s", hex.EncodeToString(validatorDilithiumPK[:])))), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(dilithiumMetadataSerialized, nil).AnyTimes()
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(strKey[:], slaveXmssPK[:])), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(slaveMetadataSerialized, nil).AnyTimes()
				store.EXPECT().
					GetFromBucket(gomock.Eq(metadata.GetOTSIndexMetaDataKeyByOTSIndex(strKey[:], otsIndex)), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("parentBlockHeaderHash"))))))).
					Return(otsIndexMetadataSerialized, nil).AnyTimes()
				store.EXPECT().DB().Return(db).AnyTimes()
				store.EXPECT().Get(gomock.Eq(metadata.GetBlockMetaDataKey(sha256.New().Sum([]byte("parentBlockHeaderHash"))))).Return(parentBlockMetadataSerialized, nil).AnyTimes()
				store.EXPECT().Get(gomock.Eq(metadata.GetBlockMetaDataKey(sha256.New().Sum([]byte("lastBlockHeaderHash"))))).Return(lastBlockMetadataSerialized, nil).AnyTimes()
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
			err = tc.stateContext.Commit(tc.blockStorageKey, tc.bytesBlock, tc.isFinalizedState)
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
	parentBlockMetadata := metadata.NewBlockMetaData(sha256.New().Sum([]byte("parentsparentBlockHeaderHash")), sha256.New().Sum([]byte("parentBlockHeaderHash")), 0, totalStakeAmount)
	parentBlockMetadataSerialized, _ := parentBlockMetadata.Serialize()

	lastBlockMetadata := metadata.NewBlockMetaData(sha256.New().Sum([]byte("parentBlockHeaderHash")), sha256.New().Sum([]byte("lastBlockHeaderHash")), 0, totalStakeAmount)
	lastBlockMetadataSerialized, _ := lastBlockMetadata.Serialize()

	blockMetaDataPathForFinalization := make([]*metadata.BlockMetaData, 0)
	blockMetaDataPathForFinalization = append(blockMetaDataPathForFinalization, lastBlockMetadata)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("parentBlockHeaderHash")), 0)
	store := mockdb.NewMockDB(ctrl)

	stateContext := &StateContext{
		db:             store,
		addressesState: make(map[string]*address.AddressState),
		dilithiumState: make(map[string]*metadata.DilithiumMetaData),
		slaveState:     make(map[string]*metadata.SlaveMetaData),
		otsIndexState:  make(map[string]*metadata.OTSIndexMetaData),

		slotNumber:              1,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.New().Sum([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.New().Sum([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.New().Sum([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.New().Sum([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.New().Sum([]byte("blockSigningHash")),
		validatorsToXMSSAddress: make(map[string][]byte),
		attestorsFlag:           make(map[string]bool),
		blockProposerFlag:       false,

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
				store.EXPECT().Get(gomock.Eq(metadata.GetBlockMetaDataKey(sha256.New().Sum([]byte("parentBlockHeaderHash"))))).Return(parentBlockMetadataSerialized, nil).AnyTimes()
				store.EXPECT().Get(gomock.Eq(metadata.GetBlockMetaDataKey(sha256.New().Sum([]byte("lastBlockHeaderHash"))))).Return(lastBlockMetadataSerialized, nil).AnyTimes()
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
	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))
	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))
	blockSigningHash := sha256.New().Sum([]byte("blockSigningHash"))
	epochMetaData := &metadata.EpochMetaData{}
	mainChainMetaData := metadata.NewMainChainMetaData(sha256.New().Sum([]byte("finalizedBlockHeaderHash")), 1,
		sha256.New().Sum([]byte("parentBlockHeaderHash")), 0)
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
