package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	// "errors"
	"fmt"

	"github.com/theQRL/zond/common"
	"github.com/theQRL/zond/config"
	"github.com/theQRL/zond/db"

	// "github.com/theQRL/zond/protos"
	"go.etcd.io/bbolt"
)

func TestNewEpochBlockHashes(t *testing.T) {
	epoch := uint64(1)
	epochBlockHashes := NewEpochBlockHashes(epoch)

	slotNumber := epoch*config.GetDevConfig().BlocksPerEpoch + 1

	if epochBlockHashes.Epoch() != epoch {
		t.Error("epoch not correctly set in epoch block hashes")
	}

	if epochBlockHashes.BlockHashesBySlotNumber()[1].GetSlotNumber() != slotNumber {
		t.Error("slot number not correctly set")
	}
}

func TestGetEpochBlockHashes(t *testing.T) {
	epoch := uint64(1)

	slotNumber := epoch*config.GetDevConfig().BlocksPerEpoch + 1

	epochBlockHashesMetadata := NewEpochBlockHashes(epoch)

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
			_, err = tx.CreateBucket([]byte("DB"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}

		err = epochBlockHashesMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}
		return nil
	})

	output, err := GetEpochBlockHashes(store, epoch)
	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if output.Epoch() != epoch {
		t.Error("epoch not correctly set in epoch block hashes")
	}

	if output.BlockHashesBySlotNumber()[1].GetSlotNumber() != slotNumber {
		t.Error("slot number not correctly set")
	}
}

func TestAddHeaderHashBySlotNumber(t *testing.T) {
	epoch := uint64(1)

	slotNumber := epoch*config.GetDevConfig().BlocksPerEpoch + 1

	headerHash := sha256.Sum256([]byte("headerHash"))

	epochBlockHashesMetadata := NewEpochBlockHashes(epoch)

	epochBlockHashesMetadata2 := NewEpochBlockHashes(uint64(2))
	epochBlockHashesMetadata2.pbData.BlockHashesBySlotNumber[0].SlotNumber = 2000

	epochBlockHashesMetadata3 := NewEpochBlockHashes(uint64(3))
	epochBlockHashesMetadata3.pbData.BlockHashesBySlotNumber[0].HeaderHashes = append(epochBlockHashesMetadata3.pbData.BlockHashesBySlotNumber[0].HeaderHashes, headerHash[:])

	testCases := []struct {
		name             string
		headerHash       common.Hash
		slotNumber       uint64
		epochBlockHashes *EpochBlockHashes
		expectedError    error
	}{
		{
			name:             "ok",
			headerHash:       headerHash,
			slotNumber:       slotNumber,
			epochBlockHashes: epochBlockHashesMetadata,
			expectedError:    nil,
		},
		{
			name:             "slotNumber out of range",
			headerHash:       headerHash,
			slotNumber:       2000,
			epochBlockHashes: epochBlockHashesMetadata,
			expectedError:    fmt.Errorf("SlotNumber %d doesn't belong to epoch %d", 2000, 1),
		},
		{
			name:             "unexpected slotNumber",
			headerHash:       headerHash,
			slotNumber:       2*config.GetDevConfig().BlocksPerEpoch + 0,
			epochBlockHashes: epochBlockHashesMetadata2,
			expectedError:    fmt.Errorf("Unexpected slot number %d at index %d", 2000, 0),
		},
		{
			name:             "already existing headerHash",
			headerHash:       headerHash,
			slotNumber:       3*config.GetDevConfig().BlocksPerEpoch + 0,
			epochBlockHashes: epochBlockHashesMetadata3,
			expectedError:    fmt.Errorf("Headerhash %s already exists", hex.EncodeToString(headerHash[:])),
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			err := tc.epochBlockHashes.AddHeaderHashBySlotNumber(tc.headerHash, tc.slotNumber)
			if err != nil && (err.Error() != tc.expectedError.Error()) {
				t.Errorf("expected error (%v), got error (%v)", tc.expectedError, err)
			}
		})
	}
}

func TestEpochBlockHashesCommit(t *testing.T) {
	epoch := uint64(1)

	epochBlockHashesMetadata := NewEpochBlockHashes(epoch)

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
			_, err = tx.CreateBucket([]byte("DB"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}

		err = epochBlockHashesMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}

		data := mainBucket.Get(GetEpochBlockHashesKey(epoch))
		if data == nil {
			return fmt.Errorf("metadata not saved in db, got (%s)", data)
		}
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}
}
