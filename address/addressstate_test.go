package address

import (
	"crypto/sha256"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/go-qrllib/xmss"
	mockdb "github.com/theQRL/zond/db/mock"
	"github.com/theQRL/zond/metadata"
)

func TestNewAddressState(t *testing.T) {
	newXmss := xmss.NewXMSSFromHeight(8, 0)
	newXmssAddr := newXmss.GetAddress()
	nonce := uint64(10)
	balance := uint64(100)
	newAddressState := NewAddressState(newXmssAddr[:], nonce, balance)

	if newAddressState.Nonce() != nonce {
		t.Errorf("nonce not set correctly in address stage, expected (%v) got (%v)", nonce, newAddressState.Nonce())
	}

	if newAddressState.Balance() != balance {
		t.Errorf("balance not set correctly in address state, expected (%v) got (%v)", balance, newAddressState.Balance())
	}
}

func TestGetAddressState(t *testing.T) {
	ctrl := gomock.NewController(t)

	newXmss := xmss.NewXMSSFromHeight(8, 0)
	newXmssAddr := newXmss.GetAddress()
	nonce := uint64(10)
	balance := uint64(100)
	newAddressState := NewAddressState(newXmssAddr[:], nonce, balance)
	newAddressStateSerialized, _ := newAddressState.Serialize()
	newAddressStateKey := GetAddressStateKey(newXmssAddr[:])

	finalizedHeaderHash := sha256.New().Sum([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))

	store := mockdb.NewMockDB(ctrl)
	store.EXPECT().GetFromBucket(gomock.Eq(newAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(newAddressStateSerialized, nil).AnyTimes()

	output, _ := GetAddressState(store, newXmssAddr[:], parentBlockHeaderHash, finalizedHeaderHash)

	if output.Nonce() != nonce {
		t.Errorf("nonce not set correctly in address stage, expected (%v) got (%v)", nonce, newAddressState.Nonce())
	}

	if output.Balance() != balance {
		t.Errorf("balance not set correctly in address state, expected (%v) got (%v)", balance, newAddressState.Balance())
	}
}
