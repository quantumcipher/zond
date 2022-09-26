package metadata

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/theQRL/zond/db"
	"github.com/theQRL/zond/misc"
	"go.etcd.io/bbolt"
)

func TestNewMainChainMetaData(t *testing.T) {
	finalizedBlockHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.Sum256([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)

	mainChainMetaData := NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)

	if mainChainMetaData.FinalizedBlockHeaderHash().String() != misc.BytesToHexStr(finalizedBlockHeaderHash[:]) {
		t.Errorf("expected finalized block header hash (%v), got (%v)", misc.BytesToHexStr(finalizedBlockHeaderHash[:]), mainChainMetaData.FinalizedBlockHeaderHash().String())
	}

	if mainChainMetaData.LastBlockHeaderHash().String() != misc.BytesToHexStr(lastBlockHeaderHash[:]) {
		t.Errorf("expected last block header hash (%v), got (%v)", mainChainMetaData.LastBlockHeaderHash().String(), misc.BytesToHexStr(lastBlockHeaderHash[:]))
	}
}

func TestGetMainChainMetaData(t *testing.T) {
	finalizedBlockHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.Sum256([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)

	mainChainMetaData := NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)

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

		err = mainChainMetaData.Commit(mainBucket)
		if err != nil {
			return err
		}
		return nil
	})

	output, err := GetMainChainMetaData(store)

	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if output.FinalizedBlockHeaderHash().String() != misc.BytesToHexStr(finalizedBlockHeaderHash[:]) {
		t.Errorf("expected finalized block header hash (%v), got (%v)", misc.BytesToHexStr(finalizedBlockHeaderHash[:]), output.FinalizedBlockHeaderHash().String())
	}

	if output.LastBlockHeaderHash().String() != misc.BytesToHexStr(lastBlockHeaderHash[:]) {
		t.Errorf("expected last block header hash (%v), got (%v)", misc.BytesToHexStr(lastBlockHeaderHash[:]), output.LastBlockHeaderHash().String())
	}
}

func TestUpdateFinalizedBlockData(t *testing.T) {
	finalizedBlockHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.Sum256([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)

	finalizedBlockHeaderHash2 := sha256.Sum256([]byte("finalizedHeaderHash2"))
	finalizedBlockSlotNumber2 := uint64(11)

	mainChainMetaData := NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)

	mainChainMetaData.UpdateFinalizedBlockData(finalizedBlockHeaderHash2, finalizedBlockSlotNumber2)

	if mainChainMetaData.FinalizedBlockHeaderHash().String() != misc.BytesToHexStr(finalizedBlockHeaderHash2[:]) {
		t.Errorf("the finalized block header hash not able to update")
	}

	if mainChainMetaData.FinalizedBlockSlotNumber() != finalizedBlockSlotNumber2 {
		t.Errorf("the finalized block slot number not able to update")
	}
}

func TestUpdateLastBlockData(t *testing.T) {
	finalizedBlockHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.Sum256([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)

	lastBlockHeaderHash2 := sha256.Sum256([]byte("lastHeaderHash2"))
	lastBlockSlotNumber2 := uint64(11)

	mainChainMetaData := NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)

	mainChainMetaData.UpdateLastBlockData(lastBlockHeaderHash2, lastBlockSlotNumber2)

	if mainChainMetaData.LastBlockHeaderHash().String() != misc.BytesToHexStr(lastBlockHeaderHash2[:]) {
		t.Errorf("the finalized block header hash not able to update")
	}

	if mainChainMetaData.LastBlockSlotNumber() != lastBlockSlotNumber2 {
		t.Errorf("the finalized block slot number not able to update")
	}
}

func TestCommit(t *testing.T) {
	finalizedBlockHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.Sum256([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)

	mainChainMetaData := NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)
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
