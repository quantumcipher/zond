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

func TestNewBlockMetaData(t *testing.T) {
	parentHeaderHash := sha256.Sum256([]byte("parentHeaderHash"))
	headerHash := sha256.Sum256([]byte("headerHash"))
	trieRoot := sha256.Sum256([]byte("trieRoot"))
	slotNumber := uint64(178)
	totalStakeAmount := []byte("100")

	blockMetadata := NewBlockMetaData(parentHeaderHash, headerHash, slotNumber, totalStakeAmount, trieRoot)

	if blockMetadata.ParentHeaderHash().String() != misc.BytesToHexStr(parentHeaderHash[:]) {
		t.Errorf("expected parent headerhash (%v), got (%v)", misc.BytesToHexStr(parentHeaderHash[:]), blockMetadata.ParentHeaderHash())
	}

	if blockMetadata.SlotNumber() != slotNumber {
		t.Errorf("expected slotnumber (%v) got (%v)", slotNumber, blockMetadata.SlotNumber())
	}
}

func TestGetBlockMetaData(t *testing.T) {
	parentHeaderHash := sha256.Sum256([]byte("parentHeaderHash"))
	headerHash := sha256.Sum256([]byte("headerHash"))
	trieRoot := sha256.Sum256([]byte("trieRoot"))
	slotNumber := uint64(178)
	totalStakeAmount := []byte("100")

	blockMetadata := NewBlockMetaData(parentHeaderHash, headerHash, slotNumber, totalStakeAmount, trieRoot)

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

		err := blockMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}

	output, err := GetBlockMetaData(store, headerHash)
	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if output.ParentHeaderHash().String() != misc.BytesToHexStr(parentHeaderHash[:]) {
		t.Errorf("expected parent headerhash (%v), got (%v)", misc.BytesToHexStr(parentHeaderHash[:]), output.ParentHeaderHash().String())
	}

	if output.SlotNumber() != slotNumber {
		t.Errorf("expected slotnumber (%v) got (%v)", slotNumber, output.SlotNumber())
	}
}

func TestBlockCommit(t *testing.T) {
	parentHeaderHash := sha256.Sum256([]byte("parentHeaderHash"))
	headerHash := sha256.Sum256([]byte("headerHash"))
	slotNumber := uint64(178)
	trieRoot := sha256.Sum256([]byte("trieRoot"))
	totalStakeAmount := []byte("100")

	blockMetadata := NewBlockMetaData(parentHeaderHash, headerHash, slotNumber, totalStakeAmount, trieRoot)

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
