package metadata

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/theQRL/go-qrllib/xmss"
	"github.com/theQRL/zond/db"
	"github.com/theQRL/zond/misc"
	"go.etcd.io/bbolt"
)

func TestNewSlaveMetaData(t *testing.T) {
	slaveXmss := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK := slaveXmss.GetPK()
	transactionHash := sha256.Sum256([]byte("transactionHash"))

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((validatorXmssPK[:])))

	slaveMetadata := NewSlaveMetaData(transactionHash[:], address[:], slaveXmssPK[:])

	if misc.BytesToHexStr(slaveMetadata.Address()) != misc.BytesToHexStr(address[:]) {
		t.Errorf("expected address (%v) got (%v)", misc.BytesToHexStr(slaveMetadata.Address()), misc.BytesToHexStr(address[:]))
	}

	if misc.BytesToHexStr(slaveMetadata.SlavePK()) != misc.BytesToHexStr(slaveXmssPK[:]) {
		t.Errorf("expected slave key (%v) got (%v)", misc.BytesToHexStr(slaveMetadata.SlavePK()), misc.BytesToHexStr(slaveXmssPK[:]))
	}
}

func TestGetSlaveMetaData(t *testing.T) {
	slaveXmss := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK := slaveXmss.GetPK()
	transactionHash := sha256.Sum256([]byte("transactionHash"))
	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((validatorXmssPK[:])))
	trieRoot := sha256.Sum256([]byte("trieRoot"))
	slotNumber := uint64(178)
	totalStakeAmount := []byte("100")

	slaveMetadata := NewSlaveMetaData(transactionHash[:], address[:], slaveXmssPK[:])

	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
	parentHeaderHash := sha256.Sum256([]byte("parentHeaderHash"))

	blockMetadata := NewBlockMetaData(parentHeaderHash, blockHeaderHash, slotNumber, totalStakeAmount, trieRoot)

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
		err = slaveMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}
		return nil
	})

	output, err := GetSlaveMetaData(store, address[:], slaveXmssPK[:], blockHeaderHash, blockHeaderHash)

	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if misc.BytesToHexStr(output.Address()) != misc.BytesToHexStr(address[:]) {
		t.Errorf("expected address (%v) got (%v)", misc.BytesToHexStr(slaveMetadata.Address()), misc.BytesToHexStr(address[:]))
	}

	if misc.BytesToHexStr(output.SlavePK()) != misc.BytesToHexStr(slaveXmssPK[:]) {
		t.Errorf("expected slave key (%v) got (%v)", misc.BytesToHexStr(output.SlavePK()), misc.BytesToHexStr(slaveXmssPK[:]))
	}
}

func TestSlaveCommit(t *testing.T) {
	slaveXmss := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK := slaveXmss.GetPK()
	transactionHash := sha256.Sum256([]byte("transactionHash"))

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((validatorXmssPK[:])))

	slaveMetadata := NewSlaveMetaData(transactionHash[:], address[:], slaveXmssPK[:])

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

		err = slaveMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error committing slave metadata to database (%v)", err)
	}
}
