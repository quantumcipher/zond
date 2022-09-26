package state

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/zond/db"
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

	dir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir) // clean up

	file := filepath.Join(dir, "tmpfile")
	if err := os.WriteFile(file, []byte(""), 0666); err != nil {
		t.Error(err)
	}

	store, err := db.NewDB(dir, "tmpfile")
	if err != nil {
		t.Error("unexpected error while creating new db ", err)
	}

	//txaddress := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK(transactionPK[:]))

	stateContext := &StateContext{
		db:                           store,
		currentBlockTotalStakeAmount: big.NewInt(10000000000000),
		slotNumber:                   0,
		blockProposer:                blockProposerPK[:],
		finalizedHeaderHash:          sha256.Sum256([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:        sha256.Sum256([]byte("parentBlockHeaderHash")),
		blockHeaderHash:              sha256.Sum256([]byte("blockHeaderHash")),
		partialBlockSigningHash:      sha256.Sum256([]byte("partialBlockSigningHash")),
		blockSigningHash:             sha256.Sum256([]byte("blockSigningHash")),

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: &metadata.MainChainMetaData{},
	}

	testCases := []struct {
		name          string
		dilithiumPK   []byte
		stateContext  StateContext
		stakeBalance  *big.Int
		expectedError error
	}{
		{
			name:          "ok",
			dilithiumPK:   validatorPK[:],
			stateContext:  *stateContext,
			stakeBalance:  big.NewInt(10000000000000),
			expectedError: nil,
		},
		{
			name:          "invalid stake balance",
			dilithiumPK:   validatorPK[:],
			stateContext:  *stateContext,
			stakeBalance:  big.NewInt(10000000000),
			expectedError: fmt.Errorf("Invalid stake balance %d for address %s", big.NewInt(10000000000), misc.GetAddressFromUnSizedPK(validatorPK[:])),
		},
		{
			name:          "non-validator dilithium key",
			dilithiumPK:   pk[:],
			stateContext:  *stateContext,
			stakeBalance:  big.NewInt(10000000000000),
			expectedError: fmt.Errorf("validator dilithium state not found for %s", hex.EncodeToString(pk[:])),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			//initialCurrentStakeAmount := tc.stateContext.currentBlockTotalStakeAmount
			err := tc.stateContext.processValidatorStakeAmount(tc.dilithiumPK[:], tc.stakeBalance)

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

	dir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir) // clean up

	file := filepath.Join(dir, "tmpfile")
	if err := os.WriteFile(file, []byte(""), 0666); err != nil {
		t.Error(err)
	}

	store, err := db.NewDB(dir, "tmpfile")
	if err != nil {
		t.Error("unexpected error while creating new db ", err)
	}

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
		finalizedHeaderHash:     sha256.Sum256([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.Sum256([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.Sum256([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.Sum256([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.Sum256([]byte("blockSigningHash")),

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: &metadata.MainChainMetaData{},
	}

	testCases := []struct {
		name          string
		dilithiumPK   []byte
		stateContext  StateContext
		stakeBalance  *big.Int
		expectedError error
	}{
		{
			name:          "ok",
			dilithiumPK:   attestorDilithiumPK1[:],
			stateContext:  *stateContext,
			stakeBalance:  big.NewInt(10),
			expectedError: nil,
		},
		{
			name:          "already attested dilithium key",
			dilithiumPK:   attestorDilithiumPK2[:],
			stateContext:  *stateContext,
			stakeBalance:  big.NewInt(10),
			expectedError: errors.New("attestor already attested for this slot number"),
		},
		{
			name:          "non-attestor dilithium key",
			dilithiumPK:   pk[:],
			stateContext:  *stateContext,
			stakeBalance:  big.NewInt(10),
			expectedError: errors.New("attestor is not assigned to attest at this slot number"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := tc.stateContext.ProcessAttestorsFlag(tc.dilithiumPK, tc.stakeBalance)

			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
		})
	}
}

func TestProcessBlockProposerFlag(t *testing.T) {
	dilithium_ := dilithium.New()
	pk := dilithium_.GetPK()

	dir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir) // clean up

	file := filepath.Join(dir, "tmpfile")
	if err := os.WriteFile(file, []byte(""), 0666); err != nil {
		t.Error(err)
	}

	store, err := db.NewDB(dir, "tmpfile")
	if err != nil {
		t.Error("unexpected error while creating new db ", err)
	}

	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	stateContext := &StateContext{
		db: store,

		slotNumber:              0,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.Sum256([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.Sum256([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.Sum256([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.Sum256([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.Sum256([]byte("blockSigningHash")),
		epochMetaData:           &metadata.EpochMetaData{},
		epochBlockHashes:        metadata.NewEpochBlockHashes(0),
		mainChainMetaData:       &metadata.MainChainMetaData{},
	}

	testCases := []struct {
		name          string
		dilithiumPK   []byte
		stateContext  StateContext
		stakeBalance  *big.Int
		expectedError error
	}{
		{
			name:          "ok",
			dilithiumPK:   blockProposerPK[:],
			stateContext:  *stateContext,
			stakeBalance:  big.NewInt(10),
			expectedError: nil,
		},
		{
			name:          "non block-proposer dilithium key",
			dilithiumPK:   pk[:],
			stateContext:  *stateContext,
			stakeBalance:  big.NewInt(10),
			expectedError: errors.New("unexpected block proposer"),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := tc.stateContext.ProcessBlockProposerFlag(tc.dilithiumPK, tc.stakeBalance)

			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
		})
	}
}

func TestFinalize(t *testing.T) {
	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	totalStakeAmount, _ := big.NewInt(10).MarshalText()

	parentBlockMetadata := metadata.NewBlockMetaData(sha256.Sum256([]byte("parentsparentBlockHeaderHash")), sha256.Sum256([]byte("parentBlockHeaderHash")), 0, totalStakeAmount, sha256.Sum256([]byte("trieRoot")))
	lastBlockMetadata := metadata.NewBlockMetaData(sha256.Sum256([]byte("parentBlockHeaderHash")), sha256.Sum256([]byte("lastBlockHeaderHash")), 0, totalStakeAmount, sha256.Sum256([]byte("trieRoot")))
	blockMetaDataPathForFinalization := make([]*metadata.BlockMetaData, 0)
	blockMetaDataPathForFinalization = append(blockMetaDataPathForFinalization, lastBlockMetadata)

	mainChainMetaData := metadata.NewMainChainMetaData(sha256.Sum256([]byte("finalizedBlockHeaderHash")), 1,
		sha256.Sum256([]byte("parentBlockHeaderHash")), 0)
	dir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir) // clean up

	file := filepath.Join(dir, "tmpfile")
	if err := os.WriteFile(file, []byte(""), 0666); err != nil {
		t.Error(err)
	}

	store, err := db.NewDB(dir, "tmpfile")
	if err != nil {
		t.Error("unexpected error while creating new db ", err)
	}

	err = store.DB().Update(func(tx *bbolt.Tx) error {
		mainBucket := tx.Bucket([]byte("DB"))
		if mainBucket == nil {
			_, err := tx.CreateBucket([]byte("DB"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}

		err := mainChainMetaData.Commit(mainBucket)
		if err != nil {
			return err
		}

		err = parentBlockMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}

		err = lastBlockMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}
	_, err = metadata.GetBlockMetaData(store, sha256.Sum256([]byte("parentBlockHeaderHash")))
	if err != nil {
		t.Error(err)
	}

	stateContext := &StateContext{
		db:                      store,
		slotNumber:              1,
		blockProposer:           blockProposerPK[:],
		finalizedHeaderHash:     sha256.Sum256([]byte("finalizedHeaderHash")),
		parentBlockHeaderHash:   sha256.Sum256([]byte("parentBlockHeaderHash")),
		blockHeaderHash:         sha256.Sum256([]byte("blockHeaderHash")),
		partialBlockSigningHash: sha256.Sum256([]byte("partialBlockSigningHash")),
		blockSigningHash:        sha256.Sum256([]byte("blockSigningHash")),

		epochMetaData:     &metadata.EpochMetaData{},
		epochBlockHashes:  metadata.NewEpochBlockHashes(0),
		mainChainMetaData: mainChainMetaData,
	}
	testCases := []struct {
		name                   string
		blockMetaDataPathArray []*metadata.BlockMetaData
		stateContext           StateContext
		expectedError          error
	}{
		{
			name:                   "ok",
			blockMetaDataPathArray: blockMetaDataPathForFinalization,
			stateContext:           *stateContext,
			expectedError:          nil,
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
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
	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()

	slotNumber := uint64(0)
	finalizedHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))
	blockSigningHash := sha256.Sum256([]byte("blockSigningHash"))
	epochMetaData := &metadata.EpochMetaData{}
	mainChainMetaData := metadata.NewMainChainMetaData(sha256.Sum256([]byte("finalizedBlockHeaderHash")), 1,
		sha256.Sum256([]byte("parentBlockHeaderHash")), 0)

	dir, err := os.MkdirTemp("", "tempdir")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(dir) // clean up

	file := filepath.Join(dir, "tmpfile")
	if err := os.WriteFile(file, []byte(""), 0666); err != nil {
		t.Error(err)
	}

	store, err := db.NewDB(dir, "tmpfile")
	if err != nil {
		t.Error("unexpected error while creating new db ", err)
	}

	err = store.DB().Update(func(tx *bbolt.Tx) error {
		mainBucket := tx.Bucket([]byte("DB"))
		if mainBucket == nil {
			_, err := tx.CreateBucket([]byte("DB"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}

		err := mainChainMetaData.Commit(mainBucket)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}

	newStateContext, err := NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash,
		partialBlockSigningHash, blockSigningHash, epochMetaData)

	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if newStateContext.GetSlotNumber() != slotNumber {
		t.Errorf("expected slotnumber (%v) got (%v)", slotNumber, newStateContext.GetSlotNumber())
	}
}
