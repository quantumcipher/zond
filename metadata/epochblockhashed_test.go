package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"

	// "errors"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/zond/config"
	mockdb "github.com/theQRL/zond/db/mock"

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
	ctrl := gomock.NewController(t)

	epoch := uint64(1)

	slotNumber := epoch*config.GetDevConfig().BlocksPerEpoch + 1

	epochBlockHashesMetadata := NewEpochBlockHashes(epoch)
	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()

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

	headerHash := sha256.New().Sum([]byte("headerHash"))

	epochBlockHashesMetadata := NewEpochBlockHashes(epoch)

	epochBlockHashesMetadata2 := NewEpochBlockHashes(uint64(2))
	epochBlockHashesMetadata2.pbData.BlockHashesBySlotNumber[0].SlotNumber = 2000

	epochBlockHashesMetadata3 := NewEpochBlockHashes(uint64(3))
	epochBlockHashesMetadata3.pbData.BlockHashesBySlotNumber[0].HeaderHashes = append(epochBlockHashesMetadata3.pbData.BlockHashesBySlotNumber[0].HeaderHashes, headerHash)

	testCases := []struct {
		name             string
		headerHash       []byte
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
			expectedError:    fmt.Errorf("Headerhash %s already exists", hex.EncodeToString(headerHash)),
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
	ctrl := gomock.NewController(t)
	db, err := bbolt.Open("./testdb3.txt", 0600, &bbolt.Options{Timeout: 1 * time.Second, InitialMmapSize: 10e6})
	if err != nil {
		t.Error(err.Error())
	}

	epoch := uint64(1)

	epochBlockHashesMetadata := NewEpochBlockHashes(epoch)

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().DB().Return(db).AnyTimes()

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
