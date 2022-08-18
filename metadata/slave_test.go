package metadata

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/go-qrllib/xmss"
	"github.com/theQRL/zond/common"
	mockdb "github.com/theQRL/zond/db/mock"
	"github.com/theQRL/zond/misc"
	"go.etcd.io/bbolt"
)

func TestNewSlaveMetaData(t *testing.T) {
	slaveXmss := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK := slaveXmss.GetPK()
	transactionHash := sha256.New().Sum([]byte("transactionHash"))

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((validatorXmssPK[:])))

	slaveMetadata := NewSlaveMetaData(transactionHash, address[:], slaveXmssPK[:])

	if string(slaveMetadata.Address()) != string(address[:]) {
		t.Errorf("expected address (%v) got (%v)", string(slaveMetadata.Address()), string(address[:]))
	}

	if string(slaveMetadata.SlavePK()) != string(slaveXmssPK[:]) {
		t.Errorf("expected slave key (%v) got (%v)", string(slaveMetadata.SlavePK()), string(slaveXmssPK[:]))
	}
}

func TestGetSlaveMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	slaveXmss := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK := slaveXmss.GetPK()
	transactionHash := common.Hash(sha256.Sum256([]byte("transactionHash")))
	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((validatorXmssPK[:])))

	slaveMetadata := NewSlaveMetaData(transactionHash.Bytes(), address[:], slaveXmssPK[:])
	slaveMetadataSerialized, _ := slaveMetadata.Serialize()

	finalizedHeaderHash := common.Hash(sha256.Sum256([]byte("finalizedHeaderHash")))
	blockHeaderHash := common.Hash(sha256.Sum256([]byte("blockHeaderHash")))

	store := mockdb.NewMockDB(ctrl)

	store.EXPECT().
		GetFromBucket(gomock.Eq(GetSlaveMetaDataKey(address[:], slaveXmssPK[:])), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", (common.Hash(sha256.Sum256([]byte("blockHeaderHash"))).String()))))).
		Return(slaveMetadataSerialized, nil).AnyTimes()

	output, err := GetSlaveMetaData(store, address[:], slaveXmssPK[:], blockHeaderHash, finalizedHeaderHash)

	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if string(output.Address()) != string(address[:]) {
		t.Errorf("expected address (%v) got (%v)", string(slaveMetadata.Address()), string(address[:]))
	}

	if string(output.SlavePK()) != string(slaveXmssPK[:]) {
		t.Errorf("expected slave key (%v) got (%v)", string(output.SlavePK()), string(slaveXmssPK[:]))
	}
}

func TestSlaveCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	db, err := bbolt.Open("./testdb.txt", 0600, &bbolt.Options{Timeout: 1 * time.Second, InitialMmapSize: 10e6})
	if err != nil {
		t.Error(err.Error())
	}

	slaveXmss := xmss.NewXMSSFromHeight(4, 0)
	slaveXmssPK := slaveXmss.GetPK()
	transactionHash := sha256.New().Sum([]byte("transactionHash"))

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((validatorXmssPK[:])))

	slaveMetadata := NewSlaveMetaData(transactionHash, address[:], slaveXmssPK[:])

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
