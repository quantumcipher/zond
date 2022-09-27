package chain

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"path"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/go-qrllib/xmss"
	"github.com/theQRL/zond/block"
	"github.com/theQRL/zond/common"
	"github.com/theQRL/zond/core/rawdb"
	state2 "github.com/theQRL/zond/core/state"
	"github.com/theQRL/zond/metadata"
	"github.com/theQRL/zond/ntp"
	"github.com/theQRL/zond/protos"
	"github.com/theQRL/zond/state"
	"github.com/theQRL/zond/transactions"
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
	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))

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
	parentHeaderHash := sha256.Sum256([]byte("parentHeaderHash"))
	finalizedBlockHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.Sum256([]byte("lastblockHeaderhash"))
	lastBlockSlotNumber := uint64(9)
	totalStakeAmount := []byte("100")
	trieRoot := sha256.Sum256([]byte("trieRoot"))
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)
	parentBlockMetadata := metadata.NewBlockMetaData(parentHeaderHash, lastBlockHeaderHash, lastBlockSlotNumber, totalStakeAmount, trieRoot)

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

	finalizedBlockHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(10)
	lastBlockHeaderHash := sha256.Sum256([]byte("lastblockHeaderhash"))
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
	parentHeaderHash := sha256.Sum256([]byte("parentHeaderHash"))
	parentHeaderHash2 := sha256.Sum256([]byte("parentHeaderHash2"))
	totalStakeAmount := []byte("100")
	trieRoot := common.Hash{}

	parentBlockMetadata := metadata.NewBlockMetaData(parentHeaderHash2, parentHeaderHash, parentSlotNumber, totalStakeAmount, trieRoot)
	parentBlockMetadata2 := metadata.NewBlockMetaData(parentHeaderHash2, parentHeaderHash2, parentSlotNumber2, totalStakeAmount, trieRoot)

	validatorDilithium := dilithium.New()
	validatorDilithiumPK := validatorDilithium.GetPK()
	validatorDilithium2 := dilithium.New()
	validatorDilithium2PK := validatorDilithium2.GetPK()
	var validators [][]byte
	validators = append(validators, validatorDilithiumPK[:])
	validators = append(validators, validatorDilithium2PK[:])
	epochMetadata := metadata.NewEpochMetaData(1, parentHeaderHash2, validators)

	finalizedBlockHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
	finalizedBlockSlotNumber := uint64(202)
	lastBlockHeaderHash := sha256.Sum256([]byte("lastBlockHeaderhash"))
	lastBlockSlotNumber := uint64(155)
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedBlockHeaderHash, finalizedBlockSlotNumber,
		lastBlockHeaderHash, lastBlockSlotNumber)

	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	chain := NewChain(state)

	err := state.DB().DB().Update(func(tx *bbolt.Tx) error {
		bucketName := metadata.GetBlockBucketName(parentHeaderHash)
		blockBucket := tx.Bucket([]byte(bucketName))
		if blockBucket == nil {
			_, err := tx.CreateBucket(bucketName)
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}

		return nil
	})

	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}

	err = state.DB().DB().Update(func(tx *bbolt.Tx) error {
		bucketName := metadata.GetBlockBucketName(parentHeaderHash2)
		blockBucket := tx.Bucket([]byte(bucketName))
		if blockBucket == nil {
			_, err := tx.CreateBucket(bucketName)
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}

		return nil
	})

	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}

	db2, err := rawdb.NewLevelDBDatabaseWithFreezer(
		path.Join(chain.config.User.DataDir(), chain.config.Dev.DB2Name), 16,
		16, path.Join(chain.config.User.DataDir(), chain.config.Dev.DB2FreezerName),
		chain.config.Dev.DB2Name, false)

	if err != nil {
		log.Error("Failed to create db2")
		t.Error(err)
	}

	chain.db2 = state2.NewDatabaseWithConfig(db2, nil)

	networkId := uint64(1)
	timestamp := ntp.GetNTP().Time()
	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	addrTo := xmss.NewXMSSFromHeight(10, 0).GetAddress()
	masterXmss := xmss.NewXMSSFromHeight(4, 0)
	masterXmssPK := masterXmss.GetPK()
	//txnHash := sha256.Sum256([]byte("txHash"))
	var txs []*protos.Transaction
	amount := uint64(30)
	fee := uint64(1)
	message := []byte("message")
	nonce := uint64(30)
	networkID := uint64(1)
	txn1 := transactions.NewTransfer(networkID, addrTo[:], amount, fee, 0, message, nonce, masterXmssPK[:])
	txs = append(txs, txn1.PBData())

	protocolTxs := make([]*protos.ProtocolTransaction, 0)
	lastCoinBaseNonce := uint64(10)

	newBlock := block.NewBlock(networkId, timestamp, blockProposerPK[:], 0, parentHeaderHash2, txs, protocolTxs, lastCoinBaseNonce)
	newBlock.PBData().Header.Hash = block.ComputeBlockHash(newBlock).Bytes()
	bytesBlock, _ := newBlock.Serialize()
	//chain.AddBlock(newBlock)
	err = state.DB().DB().Update(func(tx *bbolt.Tx) error {
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

		err = mainBucket.Put(block.GetBlockStorageKey(parentHeaderHash), bytesBlock)
		if err != nil {
			log.Error("[Commit] Failed to commit block")
			return err
		}

		return nil
	})

	if err != nil {
		t.Errorf("unexpected error committing to database (%v)", err)
	}
	_, err = block.GetBlock(state.DB(), parentHeaderHash)
	if err != nil {
		t.Error(err)
	}

	// newBlock2 := block.NewBlock(networkId, timestamp, blockProposerPK[:], parentSlotNumber, parentHeaderHash2, txs, protocolTxs, lastCoinBaseNonce)
	// newBlock2.PBData().Header.Hash = block.ComputeBlockHash(newBlock2).Bytes()
	// chain.AddBlock(newBlock2)

	_, err = chain.GetSlotLeaderDilithiumPKBySlotNumber(trieRoot, slotNumber, parentHeaderHash)
	if err != nil {
		t.Error("got unexpected error while fetching slot leader: ", err)
	}
}
