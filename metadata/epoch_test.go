package metadata

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/zond/db"
	"github.com/theQRL/zond/misc"
	"go.etcd.io/bbolt"
)

func TestNewEpochMetaData(t *testing.T) {
	epoch := uint64(1)
	prevSlotLastBlockHeaderHash := sha256.Sum256([]byte("prevSlotLastBlockHeaderHash"))

	epochMetadata := NewEpochMetaData(epoch, prevSlotLastBlockHeaderHash, nil)

	if epochMetadata.PrevSlotLastBlockHeaderHash().String() != misc.BytesToHexStr(prevSlotLastBlockHeaderHash[:]) {
		t.Errorf("expected previous slot last block headerhash (%v), got (%v)", misc.BytesToHexStr(prevSlotLastBlockHeaderHash[:]), epochMetadata.PrevSlotLastBlockHeaderHash().String())
	}

	if epochMetadata.Epoch() != epoch {
		t.Errorf("epoch not set correctly, expected (%v) got (%v)", epoch, epochMetadata.Epoch())
	}
}

func TestGetEpochMetaData(t *testing.T) {
	epoch := uint64(1)
	currentBlockSlotNumber := uint64(178)
	parentHeaderHash := sha256.Sum256([]byte("parentHeaderHash"))
	parentHeaderHash2 := sha256.Sum256([]byte("parentHeaderHash2"))
	headerHash := sha256.Sum256([]byte("headerHash"))
	slotNumber := uint64(178)
	trieRoot := sha256.Sum256([]byte("trieRoot"))
	blockMetadata := NewBlockMetaData(parentHeaderHash, headerHash, slotNumber, []byte("100"), trieRoot)

	slotNumber2 := uint64(50)
	blockMetadata2 := NewBlockMetaData(parentHeaderHash2, parentHeaderHash, slotNumber2, []byte("100"), trieRoot)

	epochMetadata := NewEpochMetaData(epoch, parentHeaderHash, nil)

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

		err = blockMetadata2.Commit(mainBucket)
		if err != nil {
			return err
		}

		err = epochMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}
		return nil
	})

	output, err := GetEpochMetaData(store, currentBlockSlotNumber, parentHeaderHash)
	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if output.PrevSlotLastBlockHeaderHash().String() != misc.BytesToHexStr(parentHeaderHash[:]) {
		t.Errorf("expected previous slot last block headerhash (%v), got (%v)", misc.BytesToHexStr(parentHeaderHash[:]), output.PrevSlotLastBlockHeaderHash().String())
	}

	if output.Epoch() != epoch {
		t.Errorf("epoch not set correctly, expected (%v) got (%v)", epoch, output.Epoch())
	}
}

func TestAllotSlots(t *testing.T) {
	epoch := uint64(1)
	headerHash := sha256.Sum256([]byte("headerHash"))

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()

	validators := make([][]byte, 2)
	validators = append(validators, validatorDilithiumPK[:])
	validators = append(validators, validatorDilithium2PK[:])

	epochMetadata := NewEpochMetaData(epoch, headerHash, validators)

	epochMetadata.AllotSlots(1, epoch, headerHash)
}

func TestEpochCommit(t *testing.T) {
	epoch := uint64(1)
	headerHash := sha256.Sum256([]byte("headerHash"))

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()

	validators := make([][]byte, 2)
	validators = append(validators, validatorDilithiumPK[:])
	validators = append(validators, validatorDilithium2PK[:])

	epochMetadata := NewEpochMetaData(epoch, headerHash, validators)

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

		err = epochMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}

		data := mainBucket.Get(GetEpochMetaDataKey(epoch, headerHash))
		if data == nil {
			return fmt.Errorf("metadata not saved in db, got (%s)", data)
		}
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}
}
