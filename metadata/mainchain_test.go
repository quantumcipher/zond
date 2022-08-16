package metadata

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/theQRL/zond/db/mock"
	"go.etcd.io/bbolt"
)

func TestNewMainChainMetaData(t *testing.T) {
	finalizedBlockHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.New().Sum([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)

	mainChainMetaData := NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)

	if string(mainChainMetaData.FinalizedBlockHeaderHash()) != string(finalizedBlockHeaderHash) {
		t.Errorf("expected finalized block header hash (%v), got (%v)", string(mainChainMetaData.FinalizedBlockHeaderHash()), string(finalizedBlockHeaderHash))
	}

	if string(mainChainMetaData.LastBlockHeaderHash()) != string(lastBlockHeaderHash) {
		t.Errorf("expected last block header hash (%v), got (%v)", string(mainChainMetaData.LastBlockHeaderHash()), string(lastBlockHeaderHash))
	}
}

func TestGetMainChainMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	finalizedBlockHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.New().Sum([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)

	mainChainMetaData := NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)

	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(GetMainChainMetaDataKey()).Return(mainChainMetaDataSerialized, nil)

	output, err := GetMainChainMetaData(store)

	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if string(output.FinalizedBlockHeaderHash()) != string(finalizedBlockHeaderHash) {
		t.Errorf("expected finalized block header hash (%v), got (%v)", string(output.FinalizedBlockHeaderHash()), string(finalizedBlockHeaderHash))
	}

	if string(output.LastBlockHeaderHash()) != string(lastBlockHeaderHash) {
		t.Errorf("expected last block header hash (%v), got (%v)", string(output.LastBlockHeaderHash()), string(lastBlockHeaderHash))
	}
}

func TestUpdateFinalizedBlockData(t *testing.T) {
	finalizedBlockHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.New().Sum([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)

	finalizedBlockHeaderHash2 := sha256.New().Sum([]byte("finalizedHeaderHash2"))
	finalizedBlockSlotNumber2 := uint64(11)

	mainChainMetaData := NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)

	mainChainMetaData.UpdateFinalizedBlockData(finalizedBlockHeaderHash2, finalizedBlockSlotNumber2)

	if string(mainChainMetaData.FinalizedBlockHeaderHash()) != string(finalizedBlockHeaderHash2) {
		t.Errorf("the finalized block header hash not able to update")
	}

	if mainChainMetaData.FinalizedBlockSlotNumber() != finalizedBlockSlotNumber2 {
		t.Errorf("the finalized block slot number not able to update")
	}
}

func TestUpdateLastBlockData(t *testing.T) {
	finalizedBlockHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.New().Sum([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)

	lastBlockHeaderHash2 := sha256.New().Sum([]byte("lastHeaderHash2"))
	lastBlockSlotNumber2 := uint64(11)

	mainChainMetaData := NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)

	mainChainMetaData.UpdateLastBlockData(lastBlockHeaderHash2, lastBlockSlotNumber2)

	if string(mainChainMetaData.LastBlockHeaderHash()) != string(lastBlockHeaderHash2) {
		t.Errorf("the finalized block header hash not able to update")
	}

	if mainChainMetaData.LastBlockSlotNumber() != lastBlockSlotNumber2 {
		t.Errorf("the finalized block slot number not able to update")
	}
}

func TestCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	db, err := bbolt.Open("./testdb2.txt", 0600, &bbolt.Options{InitialMmapSize: 10e6})
	if err != nil {
		t.Error(err.Error())
	}

	finalizedBlockHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.New().Sum([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)

	mainChainMetaData := NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)
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

		err = mainChainMetaData.Commit(mainBucket)
		if err != nil {
			return err
		}

		data := mainBucket.Get(GetMainChainMetaDataKey())
		if data == nil {
			return fmt.Errorf("metadata not saved in db, got (%s)", data)
		}
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}
}
