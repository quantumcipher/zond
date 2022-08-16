package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/go-qrllib/xmss"
	mockdb "github.com/theQRL/zond/db/mock"
	"github.com/theQRL/zond/misc"
	"go.etcd.io/bbolt"
)

func TestNewDilithiumMetaData(t *testing.T) {
	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	transactionHash := sha256.New().Sum([]byte("transactionHash"))

	dilithiumMetadata := NewDilithiumMetaData(transactionHash, validatorDilithiumPK[:], address[:], true)

	if string(dilithiumMetadata.Address()) != string(address[:]) {
		t.Errorf("expected address (%v) got (%v)", string(dilithiumMetadata.Address()), string(address[:]))
	}

	if string(dilithiumMetadata.DilithiumPK()) != string(validatorDilithiumPK[:]) {
		t.Errorf("expected dilithium key (%v) got (%v)", string(dilithiumMetadata.DilithiumPK()), string(validatorDilithiumPK[:]))
	}
}

func TestGetDilithiumMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	transactionHash := sha256.New().Sum([]byte("transactionHash"))

	dilithiumMetadata := NewDilithiumMetaData(transactionHash, validatorDilithiumPK[:], address[:], true)
	dilithiumMetadataSerialized, _ := dilithiumMetadata.Serialize()

	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))

	store := mockdb.NewMockDB(ctrl)

	store.EXPECT().
		GetFromBucket(gomock.Eq([]byte(fmt.Sprintf("DILITHIUM-META-DATA-%s", hex.EncodeToString(validatorDilithiumPK[:])))), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("blockHeaderHash"))))))).
		Return(dilithiumMetadataSerialized, nil).AnyTimes()

	output, err := GetDilithiumMetaData(store, validatorDilithiumPK[:], blockHeaderHash, finalizedHeaderHash)

	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if string(output.Address()) != string(address[:]) {
		t.Errorf("expected address (%v) got (%v)", string(dilithiumMetadata.Address()), string(address[:]))
	}

	if string(output.DilithiumPK()) != string(validatorDilithiumPK[:]) {
		t.Errorf("expected slave key (%v) got (%v)", string(output.DilithiumPK()), string(validatorDilithiumPK[:]))
	}
}

func TestGetXMSSAddressFromDilithiumPK(t *testing.T) {
	ctrl := gomock.NewController(t)

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	transactionHash := sha256.New().Sum([]byte("transactionHash"))

	dilithiumMetadata := NewDilithiumMetaData(transactionHash, validatorDilithiumPK[:], address[:], true)
	dilithiumMetadataSerialized, _ := dilithiumMetadata.Serialize()

	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))

	store := mockdb.NewMockDB(ctrl)

	store.EXPECT().
		GetFromBucket(gomock.Eq([]byte(fmt.Sprintf("DILITHIUM-META-DATA-%s", hex.EncodeToString(validatorDilithiumPK[:])))), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("blockHeaderHash"))))))).
		Return(dilithiumMetadataSerialized, nil).AnyTimes()

	output, err := GetXMSSAddressFromDilithiumPK(store, validatorDilithiumPK[:],
		blockHeaderHash, finalizedHeaderHash)
	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if string(output) != string(address[:]) {
		t.Errorf("expected address (%v) got (%v)", string(address[:]), string(output))
	}
}

func TestDilithiumCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	db, err := bbolt.Open("./testdb.txt", 0600, &bbolt.Options{Timeout: 1 * time.Second, InitialMmapSize: 10e6})
	if err != nil {
		t.Error(err.Error())
	}

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	transactionHash := sha256.New().Sum([]byte("transactionHash"))

	dilithiumMetadata := NewDilithiumMetaData(transactionHash, validatorDilithiumPK[:], address[:], true)

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

		err = dilithiumMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}

		data := mainBucket.Get(GetDilithiumMetaDataKey(validatorDilithiumPK[:]))
		if data == nil {
			return fmt.Errorf("metadata not saved in db, got (%s)", data)
		}
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}
}
