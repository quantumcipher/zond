package metadata

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/go-qrllib/xmss"
	"github.com/theQRL/zond/config"
	mockdb "github.com/theQRL/zond/db/mock"
	"github.com/theQRL/zond/misc"
	"go.etcd.io/bbolt"
)

func TestNewOTSIndexMetaData(t *testing.T) {
	otsIndex := uint64(8192)

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	pageNumber := otsIndex / config.GetDevConfig().OTSBitFieldPerPage
	otsIndexMetadata := NewOTSIndexMetaData(address[:], pageNumber)

	if string(otsIndexMetadata.Address()) != string(address[:]) {
		t.Errorf("expected otsIndexMetadata address (%v), got (%v)", string(address[:]), string(otsIndexMetadata.Address()))
	}

	if otsIndexMetadata.PageNumber() != pageNumber {
		t.Errorf("expected otsindex page number (%v), got (%v)", otsIndexMetadata.PageNumber(), pageNumber)
	}
}

func TestGetOTSIndexMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mockdb.NewMockDB(ctrl)

	otsIndex := uint64(8192)

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	blockHeaderHash := sha256.New().Sum([]byte("blockHeaderHash"))

	pageNumber := otsIndex / config.GetDevConfig().OTSBitFieldPerPage
	otsIndexMetadata := NewOTSIndexMetaData(address[:], pageNumber)
	otsIndexMetadataSerialized, _ := otsIndexMetadata.Serialize()

	store.EXPECT().
		GetFromBucket(gomock.Eq(GetOTSIndexMetaDataKeyByOTSIndex(address[:], otsIndex)), gomock.Eq([]byte(fmt.Sprintf("BLOCK-BUCKET-%s", hex.EncodeToString(sha256.New().Sum([]byte("blockHeaderHash"))))))).
		Return(otsIndexMetadataSerialized, nil).AnyTimes()

	output, err := GetOTSIndexMetaData(store, address[:], otsIndex,
		blockHeaderHash, finalizedHeaderHash)

	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if string(output.Address()) != string(address[:]) {
		t.Errorf("expected otsIndexMetadata address (%v), got (%v)", string(address[:]), string(output.Address()))
	}

	if output.PageNumber() != pageNumber {
		t.Errorf("expected otsindex page number (%v), got (%v)", output.PageNumber(), pageNumber)
	}
}

func TestIsOTSIndexUsed(t *testing.T) {
	otsIndex := uint64(81)

	unusedOtsIndex := uint64(82)

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	pageNumber := otsIndex / config.GetDevConfig().OTSBitFieldPerPage
	otsIndexMetadata := NewOTSIndexMetaData(address[:], pageNumber)

	isUsed1 := otsIndexMetadata.IsOTSIndexUsed(otsIndex)

	isUsed2 := otsIndexMetadata.IsOTSIndexUsed(unusedOtsIndex)

	if !isUsed1 && (isUsed1 != isUsed2) {
		t.Errorf("not able to detect if ots index used")
	}
}

func TestSetOTSIndex(t *testing.T) {
	otsIndex := uint64(81)

	unusedOtsIndex := uint64(82)

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	pageNumber := otsIndex / config.GetDevConfig().OTSBitFieldPerPage
	otsIndexMetadata := NewOTSIndexMetaData(address[:], pageNumber)

	isSet1 := otsIndexMetadata.SetOTSIndex(otsIndex)

	isSet2 := otsIndexMetadata.SetOTSIndex(unusedOtsIndex)

	if isSet1 && (isSet2 != isSet1) {
		t.Errorf("not able to detect if ots index used when setting otsIndex")
	}
}

func TestOtsIndexCommit(t *testing.T) {
	ctrl := gomock.NewController(t)
	db, err := bbolt.Open("./testdb.txt", 0600, &bbolt.Options{Timeout: 1 * time.Second, InitialMmapSize: 10e6})
	if err != nil {
		t.Error(err.Error())
	}

	otsIndex := uint64(8192)

	validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	validatorXmssPK := validatorXmss.GetPK()
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((validatorXmssPK[:])))

	pageNumber := otsIndex / config.GetDevConfig().OTSBitFieldPerPage
	otsIndexMetadata := NewOTSIndexMetaData(address[:], pageNumber)

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

		err = otsIndexMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}

		data := mainBucket.Get(GetOTSIndexMetaDataKeyByPageNumber(otsIndexMetadata.Address(), otsIndexMetadata.PageNumber()))
		if data == nil {
			return fmt.Errorf("metadata not saved in db, got (%s)", data)
		}
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}
}
