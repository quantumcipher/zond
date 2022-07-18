package state

import (
	"testing"
	"errors"
	"fmt"
	"reflect"
	"encoding/hex"
	"crypto/sha256"
	mockdb "github.com/theQRL/zond/db/mock"
	"github.com/golang/mock/gomock"
	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/go-qrllib/xmss"
	// "github.com/theQRL/zond/protos"
	"github.com/theQRL/zond/misc"
	"github.com/theQRL/zond/metadata"
	"github.com/theQRL/zond/address"
)

func TestProcessValidatorStakeAmount(t *testing.T) {
	dilithium_ := dilithium.New()
	pk := dilithium_.GetPK()

	validator := dilithium.New()
	validatorPK := validator.GetPK()

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	transaction := dilithium.New()
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
		mainChainMetaData: &metadata.MainChainMetaData {},
	}


	testCases := []struct {
		name string
		dilithiumPK []byte
		stateContext StateContext
		expectedError error
	} {
		{
			name: "ok",
			dilithiumPK: validatorPK[:],
			stateContext: *stateContext,
			expectedError: nil,
		},
		{
			name: "non-validator dilithium key",
			dilithiumPK: pk[:],
			stateContext: *stateContext,
			expectedError: errors.New(fmt.Sprintf("validator dilithium state not found for %s", pk)),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			initialCurrentStakeAmount := tc.stateContext.currentBlockTotalStakeAmount
			err := tc.stateContext.processValidatorStakeAmount(tc.dilithiumPK)
			
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}

			strKey := hex.EncodeToString(metadata.GetDilithiumMetaDataKey(tc.dilithiumPK))
			slotLeaderDilithiumMetaData := tc.stateContext.dilithiumState[strKey]
			// if !slotLeaderDilithiumMetaData {
			// 	return errors.New(fmt.Sprintf("validator dilithium state not found for %s", tc.dilithiumPK))
			// }
			if initialCurrentStakeAmount + slotLeaderDilithiumMetaData.Balance() != tc.stateContext.currentBlockTotalStakeAmount {
				t.Errorf("expected total block stake amount (%d), got (%d)", initialCurrentStakeAmount + slotLeaderDilithiumMetaData.Balance(), tc.stateContext.currentBlockTotalStakeAmount)
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
		mainChainMetaData: &metadata.MainChainMetaData {},
	}
	
	testCases := []struct {
		name string
		dilithiumPK []byte
		stateContext StateContext
		expectedError error
	} {
		{
			name: "ok",
			dilithiumPK: attestorDilithiumPK1[:],
			stateContext: *stateContext,
			expectedError: nil,
		},
		{
			name: "already attested dilithium key",
			dilithiumPK: attestorDilithiumPK2[:],
			stateContext: *stateContext,
			expectedError: errors.New("attestor already attested for this slot number"),
		},
		{
			name: "non-attestor dilithium key",
			dilithiumPK: pk[:],
			stateContext: *stateContext,
			expectedError: errors.New("attestor is not assigned to attest at this slot number"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := tc.stateContext.ProcessAttestorsFlag(tc.dilithiumPK)
			
			if !errors.Is(err, tc.expectedError) {
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
		mainChainMetaData: &metadata.MainChainMetaData {},
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
		mainChainMetaData: &metadata.MainChainMetaData {},
	}
	
	testCases := []struct {
		name string
		dilithiumPK []byte
		stateContext StateContext
		expectedError error
	} {
		{
			name: "ok",
			dilithiumPK: blockProposerPK[:],
			stateContext: *stateContext,
			expectedError: nil,
		},
		{
			name: "non block-proposer dilithium key",
			dilithiumPK: pk[:],
			stateContext: *stateContext,
			expectedError: errors.New("unexpected block proposer"),
		},
		{
			name: "block proposer already proposed",
			dilithiumPK: blockProposerPK[:],
			stateContext: *stateContext2,
			expectedError: errors.New("block proposer has already been processed"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := tc.stateContext.ProcessBlockProposerFlag(tc.dilithiumPK)
		
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
		})
	}
}

func TestPrepareAddressState(t *testing.T) {
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
		mainChainMetaData: &metadata.MainChainMetaData {},
	}
	
	testCases := []struct {
		name string
		dilithiumPK []byte
		stateContext StateContext
		expectedError error
	} {
		{
			name: "ok",
			dilithiumPK: pk[:],
			stateContext: *stateContext,
			expectedError: nil,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			xmssAddress := tc.stateContext.GetXMSSAddressByDilithiumPK(tc.dilithiumPK)
			err := tc.stateContext.PrepareAddressState(hex.EncodeToString(xmssAddress))
			
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}

			strKey := hex.EncodeToString(address.GetAddressStateKey(xmssAddress[:]))
			addressesState, ok := tc.stateContext.addressesState[strKey]
			if ok {
				t.Errorf("Address state not set for key %s", strKey)
			}

			expectedAddressState, err := address.GetAddressState(tc.stateContext.db, xmssAddress[:],
				tc.stateContext.parentBlockHeaderHash, tc.stateContext.mainChainMetaData.FinalizedBlockHeaderHash())
			if err!=nil {
				t.Errorf("Address state not set for address %s", hex.EncodeToString(xmssAddress))
			}

			if !reflect.DeepEqual(addressesState, expectedAddressState) {
				t.Errorf("expected addressstate does not match")
			}
		})
	}
}



