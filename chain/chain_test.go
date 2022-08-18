package chain

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/zond/block"
	"github.com/theQRL/zond/common"
	"github.com/theQRL/zond/metadata"
	"github.com/theQRL/zond/ntp"
	"github.com/theQRL/zond/protos"
	"github.com/theQRL/zond/state"
	"go.etcd.io/bbolt"
)

func TestNewChain(t *testing.T) {
	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	_ = NewChain(state)
}

/*
func TestChainLoad(t *testing.T) {
	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	// blockProposerAddr, _ := hex.DecodeString("0005004010c7bca4f580b0369fa1c4fe62a44f719b7b6b88c23dcb467b881e650c8dcc15a7ca24")

	chain := NewChain(state)
	err := chain.Load()
	if err != nil {
		t.Error("got unexpected error while loading chain")
	}
}

func TestGetStateContext2(t *testing.T) {
	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	chain := NewChain(state)
	err := chain.Load()
	if err != nil {
		t.Error("got unexpected error while loading chain")
	}
	epoch := uint64(1)
	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	slotNumber := uint64(100)
	parentBlockHeaderHash := sha256.New().Sum([]byte("parentBlockHeaderHash"))
	partialBlockSigningHash := sha256.New().Sum([]byte("partialBlockSigningHash"))

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()
	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()
	var validators [][]byte
	validators = append(validators, validatorDilithiumPK[:])
	validators = append(validators, validatorDilithium2PK[:])
	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)

	stateContext, err := chain.GetStateContext2(slotNumber, blockProposerPK[:], parentBlockHeaderHash, partialBlockSigningHash)

	if err != nil {
		t.Error("got unexpected error while calling GetStateContext2: ", err)
	}

	if string(stateContext.PartialBlockSigningHash()) != string(partialBlockSigningHash) {
		t.Errorf("expected partial block signing hash to be (%v), got (%v)", hex.EncodeToString(partialBlockSigningHash), hex.EncodeToString(stateContext.PartialBlockSigningHash()))
	}
}

func TestGetStateContext(t *testing.T) {
	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	chain := NewChain(state)

	_, err := chain.GetStateContext()
	if err != nil {
		t.Error("got unexpected error while calling GetStateContext2: ", err)
	}
}
*/

func TestGetTotalStakeAmount(t *testing.T) {
	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	chain := NewChain(state)

	finalizedBlockHeaderHash := common.BytesToHash(sha256.New().Sum([]byte("finalizedHeaderHash")))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := common.BytesToHash(sha256.New().Sum([]byte("lastblockHeaderhash")))
	lastBlockSlotNumber := uint64(9)
	totalStakeAmount := []byte("100")

	mainChainMetaData := metadata.NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)
	parentBlockMetadata := metadata.NewBlockMetaData(lastBlockHeaderHash, lastBlockHeaderHash, lastBlockSlotNumber, totalStakeAmount, common.Hash{})

	err := state.DB().DB().Update(func(tx *bbolt.Tx) error {
		mainBucket := tx.Bucket([]byte("DB"))
		if mainBucket == nil {
			_, err := tx.CreateBucket([]byte("DB"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}

		err := mainChainMetaData.Commit(mainBucket)
		if err != nil {
			return err
		}
		err = parentBlockMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}

	amt, err := chain.GetTotalStakeAmount()
	if err != nil {
		t.Error("got unexpected error while getting total stake amount")
	}

	if amt.Cmp(big.NewInt(100)) != 0 {
		t.Errorf("expected total stake amount (%v), got (%v)", big.NewInt(100), amt)
	}
}

func TestGetStartingNonFinalizedEpoch(t *testing.T) {
	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	chain := NewChain(state)

	finalizedBlockHeaderHash := common.BytesToHash(sha256.New().Sum([]byte("finalizedHeaderHash")))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := common.BytesToHash(sha256.New().Sum([]byte("lastblockHeaderhash")))
	lastBlockSlotNumber := uint64(9)

	mainChainMetaData := metadata.NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)

	err := state.DB().DB().Update(func(tx *bbolt.Tx) error {
		mainBucket := tx.Bucket([]byte("DB"))
		if mainBucket == nil {
			_, err := tx.CreateBucket([]byte("DB"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}

		err := mainChainMetaData.Commit(mainBucket)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}

	epoch, err := chain.GetStartingNonFinalizedEpoch()
	if err != nil {
		t.Error("got unexpected error getting starting non finalized epoch: ", err)
	}

	if epoch != uint64(0) {
		t.Errorf("expected epoch (%v), got (%v)", 0, epoch)
	}
}

func TestGetSlotLeaderDilithiumPKBySlotNumber(t *testing.T) {
	slotNumber := uint64(201)
	parentSlotNumber := uint64(150)
	parentSlotNumber2 := uint64(50)
	parentHeaderHash := common.BytesToHash(sha256.New().Sum([]byte("parentHeaderHash")))
	parentHeaderHash2 := common.BytesToHash(sha256.New().Sum([]byte("parentHeaderHash2")))
	totalStakeAmount := []byte("100")

	parentBlockMetadata := metadata.NewBlockMetaData(parentHeaderHash2, parentHeaderHash, parentSlotNumber, totalStakeAmount, common.Hash{})
	parentBlockMetadata2 := metadata.NewBlockMetaData(parentHeaderHash2, parentHeaderHash2, parentSlotNumber2, totalStakeAmount, common.Hash{})

	// validatorXmss := xmss.NewXMSSFromHeight(4, 0)
	// validatorXmssPK := validatorXmss.GetPK()
	// address_ := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((validatorXmssPK[:])))
	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()
	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()
	var validators [][]byte
	validators = append(validators, validatorDilithiumPK[:])
	validators = append(validators, validatorDilithium2PK[:])
	epochMetadata := metadata.NewEpochMetaData(1, parentHeaderHash2, validators)

	finalizedBlockHeaderHash := common.BytesToHash(sha256.New().Sum([]byte("finalizedHeaderHash")))
	finalizedBlockSlotNumber := uint64(202)
	lastBlockHeaderHash := common.BytesToHash(sha256.New().Sum([]byte("lastBlockHeaderhash")))
	lastBlockSlotNumber := uint64(155)
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)
	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	chain := NewChain(state)

	networkId := uint64(1)
	timestamp := ntp.GetNTP().Time()

	var txs []*protos.Transaction
	var protocolTxs []*protos.ProtocolTransaction
	lastCoinBaseNonce := uint64(10)

	_ = block.NewBlock(networkId, timestamp, validatorDilithiumPK[:], parentSlotNumber, parentHeaderHash2, txs, protocolTxs, lastCoinBaseNonce)

	err := state.DB().DB().Update(func(tx *bbolt.Tx) error {
		mainBucket := tx.Bucket([]byte("DB"))
		if mainBucket == nil {
			_, err := tx.CreateBucket([]byte("DB"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}

		err := mainChainMetaData.Commit(mainBucket)
		if err != nil {
			return err
		}
		err = parentBlockMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}
		err = parentBlockMetadata2.Commit(mainBucket)
		if err != nil {
			return err
		}
		err = epochMetadata.Commit(mainBucket)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}

	// err = newBlock.Commit(state.DB(), parentHeaderHash, true)
	// if err != nil {
	// 	t.Error("unexpected error committing block to database: ", err)
	// }

	err = state.DB().DB().Update(func(tx *bbolt.Tx) error {
		bucketName := metadata.GetBlockBucketName(parentHeaderHash)
		blockBucket := tx.Bucket([]byte(bucketName))
		if blockBucket == nil {
			_, err := tx.CreateBucket(bucketName)
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}

		// err = dilithiumMetadata.Commit(blockBucket)
		// if err != nil {
		// 	return err
		// }
		// err = dilithiumMetadata2.Commit(blockBucket)
		// if err != nil {
		// 	return err
		// }

		return nil
	})

	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}

	_, err = chain.GetSlotLeaderDilithiumPKBySlotNumber(common.Hash{}, slotNumber, parentHeaderHash)
	if err != nil {
		t.Error("got unexpected error while fetching slot leader: ", err)
	}
}
