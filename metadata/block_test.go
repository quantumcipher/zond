package metadata

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/theQRL/zond/db/mock"
	"go.etcd.io/bbolt"
)

func TestNewBlockMetaData(t *testing.T) {
	parentHeaderHash := sha256.New().Sum([]byte("parentHeaderHash"))
	headerHash := sha256.New().Sum([]byte("headerHash"))
	slotNumber := uint64(178)
	totalStakeAmount := []byte("100")

	blockMetadata := NewBlockMetaData(parentHeaderHash, headerHash, slotNumber, totalStakeAmount)

	if string(blockMetadata.ParentHeaderHash()) != string(parentHeaderHash) {
		t.Errorf("expected parent headerhash (%v), got (%v)", string(parentHeaderHash), string(blockMetadata.ParentHeaderHash()))
	}

	if blockMetadata.SlotNumber() != slotNumber {
		t.Errorf("expected slotnumber (%v) got (%v)", slotNumber, blockMetadata.SlotNumber())
	}
}

func TestGetBlockMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	parentHeaderHash := sha256.New().Sum([]byte("parentHeaderHash"))
	headerHash := sha256.New().Sum([]byte("headerHash"))
	slotNumber := uint64(178)
	totalStakeAmount := []byte("100")

	blockMetadata := NewBlockMetaData(parentHeaderHash, headerHash, slotNumber, totalStakeAmount)
	blockMetadataSerialized, _ := blockMetadata.Serialize()

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(GetBlockMetaDataKey(headerHash))).Return(blockMetadataSerialized, nil).AnyTimes()

	output, err := GetBlockMetaData(store, headerHash)
	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if string(output.ParentHeaderHash()) != string(parentHeaderHash) {
		t.Errorf("expected parent headerhash (%v), got (%v)", string(parentHeaderHash), string(output.ParentHeaderHash()))
	}

	if output.SlotNumber() != slotNumber {
		t.Errorf("expected slotnumber (%v) got (%v)", slotNumber, output.SlotNumber())
	}
}

func TestBlockCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	db, err := bbolt.Open("./testdb5.txt", 0600, &bbolt.Options{Timeout: 1 * time.Second, InitialMmapSize: 10e6})
	if err != nil {
		t.Error(err.Error())
	}

	parentHeaderHash := sha256.New().Sum([]byte("parentHeaderHash"))
	headerHash := sha256.New().Sum([]byte("headerHash"))
	slotNumber := uint64(178)
	totalStakeAmount := []byte("100")

	blockMetadata := NewBlockMetaData(parentHeaderHash, headerHash, slotNumber, totalStakeAmount)

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

		err = blockMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}

		data := mainBucket.Get(GetBlockMetaDataKey(headerHash))
		if data == nil {
			return fmt.Errorf("metadata not saved in db, got (%s)", data)
		}
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}
}
