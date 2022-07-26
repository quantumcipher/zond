package metadata

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/go-qrllib/dilithium"
	mockdb "github.com/theQRL/zond/db/mock"
	"go.etcd.io/bbolt"
)

func TestNewEpochMetaData(t *testing.T) {
	epoch := uint64(1)
	prevSlotLastBlockHeaderHash := sha256.New().Sum([]byte("prevSlotLastBlockHeaderHash"))

	epochMetadata := NewEpochMetaData(epoch, prevSlotLastBlockHeaderHash, nil)

	if string(epochMetadata.PrevSlotLastBlockHeaderHash()) != string(prevSlotLastBlockHeaderHash) {
		t.Errorf("expected previous slot last block headerhash (%v), got (%v)", string(prevSlotLastBlockHeaderHash), string(epochMetadata.PrevSlotLastBlockHeaderHash()))
	}

	if epochMetadata.Epoch() != epoch {
		t.Errorf("epoch not set correctly, expected (%v) got (%v)", epoch, epochMetadata.Epoch())
	}
}

func TestGetEpochMetaData(t *testing.T) {
	ctrl := gomock.NewController(t)

	epoch := uint64(1)
	currentBlockSlotNumber := uint64(178)
	parentHeaderHash := sha256.New().Sum([]byte("parentHeaderHash"))
	headerHash := sha256.New().Sum([]byte("headerHash"))
	slotNumber := uint64(178)
	blockMetadata := NewBlockMetaData(parentHeaderHash, headerHash, slotNumber, []byte("100"))
	blockMetadataSerialized, _ := blockMetadata.Serialize()

	slotNumber2 := uint64(50)
	blockMetadata2 := NewBlockMetaData(nil, parentHeaderHash, slotNumber2, []byte("100"))
	blockMetadataSerialized2, _ := blockMetadata2.Serialize()

	epochMetadata := NewEpochMetaData(epoch, parentHeaderHash, nil)
	epochMetadataSerialized, _ := epochMetadata.Serialize()

	fmt.Printf("blockheader key is %s\n", GetBlockMetaDataKey(headerHash))

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().Get(gomock.Eq(GetBlockMetaDataKey(headerHash))).Return(blockMetadataSerialized, nil).AnyTimes()
	store.EXPECT().Get(gomock.Eq(GetBlockMetaDataKey(parentHeaderHash))).Return(blockMetadataSerialized2, nil).AnyTimes()
	store.EXPECT().Get(gomock.Eq(GetEpochMetaDataKey(epoch, parentHeaderHash))).Return(epochMetadataSerialized, nil).AnyTimes()

	output, err := GetEpochMetaData(store, currentBlockSlotNumber, parentHeaderHash)
	if err != nil {
		t.Errorf("got unexpected error (%v)", err)
	}

	if string(output.PrevSlotLastBlockHeaderHash()) != string(parentHeaderHash) {
		t.Errorf("expected previous slot last block headerhash (%v), got (%v)", string(parentHeaderHash), string(output.PrevSlotLastBlockHeaderHash()))
	}

	if output.Epoch() != epoch {
		t.Errorf("epoch not set correctly, expected (%v) got (%v)", epoch, output.Epoch())
	}
}

func TestAllotSlots(t *testing.T) {
	epoch := uint64(1)
	headerHash := sha256.New().Sum([]byte("headerHash"))

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
	ctrl := gomock.NewController(t)
	db, err := bbolt.Open("./testdb4.txt", 0600, &bbolt.Options{Timeout: 1 * time.Second, InitialMmapSize: 10e6})
	if err != nil {
		t.Error(err.Error())
	}

	epoch := uint64(1)
	headerHash := sha256.New().Sum([]byte("headerHash"))

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()

	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()

	validators := make([][]byte, 2)
	validators = append(validators, validatorDilithiumPK[:])
	validators = append(validators, validatorDilithium2PK[:])

	epochMetadata := NewEpochMetaData(epoch, headerHash, validators)

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
