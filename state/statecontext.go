package state

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	log "github.com/sirupsen/logrus"
	"github.com/theQRL/zond/common"
	"github.com/theQRL/zond/config"
	"github.com/theQRL/zond/db"
	"github.com/theQRL/zond/metadata"
	"go.etcd.io/bbolt"
)

type StateContext struct {
<<<<<<< HEAD
	db *db.DB
=======
	db             db.DB
	addressesState map[string]*address.AddressState
	dilithiumState map[string]*metadata.DilithiumMetaData
	slaveState     map[string]*metadata.SlaveMetaData
	otsIndexState  map[string]*metadata.OTSIndexMetaData
>>>>>>> qcfork/unit-test

	slotNumber                   uint64
	blockProposer                []byte
	finalizedHeaderHash          common.Hash
	parentBlockHeaderHash        common.Hash
	blockHeaderHash              common.Hash
	partialBlockSigningHash      common.Hash
	blockSigningHash             common.Hash
	currentBlockTotalStakeAmount *big.Int

	validatorsFlag map[string]bool // Flag just to mark if attestor or block proposer has been processed, need not to be stored in state

	epochMetaData     *metadata.EpochMetaData
	epochBlockHashes  *metadata.EpochBlockHashes
	mainChainMetaData *metadata.MainChainMetaData

	totalTransactionFee uint64
}

func (s *StateContext) GetTotalTransactionFee() uint64 {
	return s.totalTransactionFee
}

func (s *StateContext) AddTransactionFee(fee uint64) {
	s.totalTransactionFee += fee
}

func (s *StateContext) GetEpochMetaData() *metadata.EpochMetaData {
	return s.epochMetaData
}

func (s *StateContext) GetSlotNumber() uint64 {
	return s.slotNumber
}

func (s *StateContext) GetMainChainMetaData() *metadata.MainChainMetaData {
	return s.mainChainMetaData
}

func (s *StateContext) PartialBlockSigningHash() common.Hash {
	return s.partialBlockSigningHash
}

func (s *StateContext) SetPartialBlockSigningHash(p common.Hash) {
	s.partialBlockSigningHash = p
}

func (s *StateContext) BlockSigningHash() common.Hash {
	return s.blockSigningHash
}

func (s *StateContext) BlockProposer() []byte {
	return s.blockProposer
}

func (s *StateContext) ValidatorsFlag() map[string]bool {
	return s.validatorsFlag
}

<<<<<<< HEAD
func (s *StateContext) processValidatorStakeAmount(dilithiumPK []byte, stakeBalance *big.Int) error {
	//addr := dilithium.GetDilithiumAddressFromPK(misc.UnSizedDilithiumPKToSizedPK(dilithiumPK))
	//stakeBalance := s.db2.GetStakeBalance(addr)

	if stakeBalance.Uint64() == 0 {
		return errors.New(fmt.Sprintf("Invalid stake balance %d for pk %s", stakeBalance, dilithiumPK))
=======
func (s *StateContext) processValidatorStakeAmount(dilithiumPK []byte) error {
	strKey := hex.EncodeToString(metadata.GetDilithiumMetaDataKey(dilithiumPK))
	slotLeaderDilithiumMetaData, ok := s.dilithiumState[strKey]

	hexPK := hex.EncodeToString(dilithiumPK)
	if !ok {
		return fmt.Errorf("validator dilithium state not found for %s", hexPK)
>>>>>>> qcfork/unit-test
	}
	s.currentBlockTotalStakeAmount = s.currentBlockTotalStakeAmount.Add(s.currentBlockTotalStakeAmount, stakeBalance)
	return nil
}

func (s *StateContext) ProcessAttestorsFlag(attestorDilithiumPK []byte, stakeBalance *big.Int) error {
	if s.slotNumber == 0 {
		return nil
	}
	strAttestorDilithiumPK := hex.EncodeToString(attestorDilithiumPK)
	result, ok := s.validatorsFlag[strAttestorDilithiumPK]
	if !ok {
		return errors.New("attestor is not assigned to attest at this slot number")
	}

	if result {
		return errors.New("attestor already attested for this slot number")
	}

	err := s.processValidatorStakeAmount(attestorDilithiumPK, stakeBalance)
	if err != nil {
		return err
	}
	s.validatorsFlag[strAttestorDilithiumPK] = true
	return nil
}

func (s *StateContext) ProcessBlockProposerFlag(blockProposerDilithiumPK []byte, stakeBalance *big.Int) error {
	if s.slotNumber == 0 {
		return nil
	}
	slotInfo := s.epochMetaData.SlotInfo()[s.slotNumber%config.GetDevConfig().BlocksPerEpoch]
	slotLeader := s.epochMetaData.Validators()[slotInfo.SlotLeader]
	if !reflect.DeepEqual(slotLeader, blockProposerDilithiumPK) {
		return errors.New("unexpected block proposer")
	}
	if s.validatorsFlag[hex.EncodeToString(blockProposerDilithiumPK)] {
		return errors.New("block proposer has already been processed")
	}

	err := s.processValidatorStakeAmount(blockProposerDilithiumPK, stakeBalance)
	if err != nil {
		return err
	}
	s.validatorsFlag[hex.EncodeToString(blockProposerDilithiumPK)] = true
	return nil
}

<<<<<<< HEAD
func (s *StateContext) PrepareValidators(dilithiumPK []byte) error {
	s.validatorsFlag[hex.EncodeToString(dilithiumPK)] = false
	return nil
}

func (s *StateContext) Commit(blockStorageKey []byte, bytesBlock []byte, trieRoot common.Hash, isFinalizedState bool) error {
=======
func (s *StateContext) PrepareAddressState(addr string) error {
	binAddr, err := hex.DecodeString(addr)
	if err != nil {
		return err
	}
	strKey := hex.EncodeToString(address.GetAddressStateKey(binAddr))
	_, ok := s.addressesState[strKey]
	if ok {
		return nil
	}

	addressState, err := address.GetAddressState(s.db, binAddr,
		s.parentBlockHeaderHash, s.mainChainMetaData.FinalizedBlockHeaderHash())
	if addressState == nil {
		return err
	}
	s.addressesState[strKey] = addressState

	return err
}

func (s *StateContext) GetAddressState(addr string) (*address.AddressState, error) {
	binAddr, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}
	strKey := hex.EncodeToString(address.GetAddressStateKey(binAddr))
	addressState, ok := s.addressesState[strKey]
	if !ok {
		return nil, fmt.Errorf("address %s not found in addressesState", addr)
	}
	return addressState, nil
}

func (s *StateContext) GetAddressStateByPK(pk []byte) (*address.AddressState, error) {
	address := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK(pk))
	addr := hex.EncodeToString(address[:])
	return s.GetAddressState(addr)
}

func (s *StateContext) PrepareValidatorsToXMSSAddress(dilithiumPK []byte) error {
	xmssAddress, err := metadata.GetXMSSAddressFromDilithiumPK(s.db,
		dilithiumPK, s.parentBlockHeaderHash, s.finalizedHeaderHash)
	if err != nil {
		log.Error("Failed to PrepareValidatorsToXMSSAddress for ",
			hex.EncodeToString(dilithiumPK))
		return err
	}
	s.validatorsToXMSSAddress[hex.EncodeToString(dilithiumPK)] = xmssAddress
	return nil
}

func (s *StateContext) GetXMSSAddressByDilithiumPK(dilithiumPK []byte) []byte {
	return s.validatorsToXMSSAddress[hex.EncodeToString(dilithiumPK)]
}

func (s *StateContext) PrepareDilithiumMetaData(dilithiumPK string) error {
	binDilithiumPK, err := hex.DecodeString(dilithiumPK)
	if err != nil {
		return err
	}
	strKey := hex.EncodeToString(metadata.GetDilithiumMetaDataKey(binDilithiumPK))
	_, ok := s.dilithiumState[strKey]
	if ok {
		return nil
	}

	dilithiumMetaData, err := metadata.GetDilithiumMetaData(s.db, binDilithiumPK,
		s.parentBlockHeaderHash, s.mainChainMetaData.FinalizedBlockHeaderHash())
	if err != nil {
		return err
	}
	s.dilithiumState[strKey] = dilithiumMetaData
	return err
}

func (s *StateContext) AddDilithiumMetaData(dilithiumPK string, dilithiumMetaData *metadata.DilithiumMetaData) error {
	binDilithiumPK, err := hex.DecodeString(dilithiumPK)
	if err != nil {
		return err
	}
	strKey := hex.EncodeToString(metadata.GetDilithiumMetaDataKey(binDilithiumPK))
	_, ok := s.dilithiumState[strKey]
	if ok {
		return errors.New("dilithiumPK already exists")
	}
	s.dilithiumState[strKey] = dilithiumMetaData
	return nil
}

func (s *StateContext) GetDilithiumState(dilithiumPK string) *metadata.DilithiumMetaData {
	binDilithiumPk, err := hex.DecodeString(dilithiumPK)
	if err != nil {
		log.Error("Error decoding dilithium PK")
		return nil
	}
	strKey := hex.EncodeToString(metadata.GetDilithiumMetaDataKey(binDilithiumPk))
	dilithiumState := s.dilithiumState[strKey]
	return dilithiumState
}

func (s *StateContext) PrepareSlaveMetaData(masterAddr string, slavePK string) error {
	binMasterAddress, err := hex.DecodeString(masterAddr)
	if err != nil {
		return err
	}
	binSlavePk, err := hex.DecodeString(slavePK)
	if err != nil {
		return err
	}
	strKey := hex.EncodeToString(metadata.GetSlaveMetaDataKey(binMasterAddress, binSlavePk))
	_, ok := s.slaveState[strKey]
	if ok {
		return nil
	}

	slaveMetaData, err := metadata.GetSlaveMetaData(s.db, binMasterAddress, binSlavePk,
		s.parentBlockHeaderHash, s.mainChainMetaData.FinalizedBlockHeaderHash())
	if slaveMetaData == nil {
		return err
	}
	s.slaveState[strKey] = slaveMetaData
	return err
}

func (s *StateContext) AddSlaveMetaData(masterAddr string, slavePK string,
	slaveMetaData *metadata.SlaveMetaData) error {
	binMasterAddress, err := hex.DecodeString(masterAddr)
	if err != nil {
		return err
	}
	binSlavePk, err := hex.DecodeString(slavePK)
	if err != nil {
		return err
	}
	strKey := hex.EncodeToString(metadata.GetSlaveMetaDataKey(binMasterAddress, binSlavePk))
	_, ok := s.slaveState[strKey]
	if ok {
		return errors.New("SlaveMetaData already exists")
	}
	s.slaveState[strKey] = slaveMetaData
	return nil
}

func (s *StateContext) GetSlaveState(masterAddr string, slavePK string) *metadata.SlaveMetaData {
	binMasterAddress, err := hex.DecodeString(masterAddr)
	if err != nil {
		log.Error("Error decoding masterAddr ", err.Error())
		return nil
	}
	binSlavePk, err := hex.DecodeString(slavePK)
	if err != nil {
		log.Error("Error decoding slavePK ", err.Error())
		return nil
	}
	strKey := hex.EncodeToString(metadata.GetSlaveMetaDataKey(binMasterAddress, binSlavePk))
	slaveMetaData, _ := s.slaveState[strKey]
	return slaveMetaData
}

func (s *StateContext) PrepareOTSIndexMetaData(address string, otsIndex uint64) error {
	binAddress, err := hex.DecodeString(address)
	if err != nil {
		return err
	}
	key := metadata.GetOTSIndexMetaDataKeyByOTSIndex(binAddress, otsIndex)
	strKey := hex.EncodeToString(key)
	_, ok := s.otsIndexState[strKey]
	if ok {
		return nil
	}

	otsIndexMetaData, err := metadata.GetOTSIndexMetaData(s.db, binAddress, otsIndex,
		s.parentBlockHeaderHash, s.mainChainMetaData.FinalizedBlockHeaderHash())
	if otsIndexMetaData == nil {
		return err
	}
	s.otsIndexState[strKey] = otsIndexMetaData
	return err
}

func (s *StateContext) AddOTSIndexMetaData(address string, otsIndex uint64,
	otsIndexMetaData *metadata.OTSIndexMetaData) error {
	binAddress, err := hex.DecodeString(address)
	if err != nil {
		return err
	}
	strKey := hex.EncodeToString(metadata.GetOTSIndexMetaDataKeyByOTSIndex(binAddress, otsIndex))
	_, ok := s.otsIndexState[strKey]
	if ok {
		return errors.New("OTSIndexMetaData already exists")
	}
	s.otsIndexState[strKey] = otsIndexMetaData
	return nil
}

func (s *StateContext) GetOTSIndexState(address string, otsIndex uint64) *metadata.OTSIndexMetaData {
	binAddress, err := hex.DecodeString(address)
	if err != nil {
		log.Error("Error decoding address ", err.Error())
		return nil
	}
	strKey := hex.EncodeToString(metadata.GetOTSIndexMetaDataKeyByOTSIndex(binAddress, otsIndex))
	otsIndexMetaData, _ := s.otsIndexState[strKey]
	return otsIndexMetaData
}

func (s *StateContext) Commit(blockStorageKey []byte, bytesBlock []byte, isFinalizedState bool) error {
>>>>>>> qcfork/unit-test
	var parentBlockMetaData *metadata.BlockMetaData
	var err error
	totalStakeAmount := big.NewInt(0)
	lastBlockTotalStakeAmount := big.NewInt(0)

	if s.slotNumber != 0 {
		parentBlockMetaData, err = metadata.GetBlockMetaData(s.db, s.parentBlockHeaderHash)
		if err != nil {
			log.Error("[Commit] Failed to load Parent BlockMetaData")
			return err
		}
		parentBlockMetaData.AddChildHeaderHash(s.blockHeaderHash)

		err = totalStakeAmount.UnmarshalText(parentBlockMetaData.TotalStakeAmount())
		if err != nil {
			log.Error("[Commit] Unable to unmarshal total stake amount of parent block metadata")
			return err
		}

		lastBlockMetaData, err := metadata.GetBlockMetaData(s.db, s.mainChainMetaData.LastBlockHeaderHash())
		lastBlockHash := s.mainChainMetaData.LastBlockHeaderHash()
		if err != nil {
			log.Error("[Commit] Failed to load last block meta data ",
				hex.EncodeToString(lastBlockHash[:]))
			return err
		}
		err = lastBlockTotalStakeAmount.UnmarshalText(lastBlockMetaData.TotalStakeAmount())
		if err != nil {
			log.Error("[Commit] Unable to Unmarshal Text for lastblockmetadata total stake amount ",
				hex.EncodeToString(lastBlockHash[:]))
			return err
		}
	}

	totalStakeAmount = totalStakeAmount.Add(totalStakeAmount, s.currentBlockTotalStakeAmount)
	bytesTotalStakeAmount, err := totalStakeAmount.MarshalText()
	if err != nil {
		log.Error("[Commit] Unable to marshal total stake amount")
		return err
	}

	blockMetaData := metadata.NewBlockMetaData(s.parentBlockHeaderHash, s.blockHeaderHash,
		s.slotNumber, bytesTotalStakeAmount, trieRoot)
	return s.db.DB().Update(func(tx *bbolt.Tx) error {
		var err error
		b := tx.Bucket([]byte("DB"))
		if err := blockMetaData.Commit(b); err != nil {
			log.Error("[Commit] Failed to commit BlockMetaData")
			return err
		}
		fmt.Print("reached here")
		err = s.epochBlockHashes.AddHeaderHashBySlotNumber(s.blockHeaderHash, s.slotNumber)
		if err != nil {
			log.Error("[Commit] Failed to Add Hash into EpochBlockHashes")
			return err
		}
		if err := s.epochBlockHashes.Commit(b); err != nil {
			log.Error("[Commit] Failed to commit EpochBlockHashes")
			return err
		}

		if s.slotNumber != 0 {
			if err := parentBlockMetaData.Commit(b); err != nil {
				log.Error("[Commit] Failed to commit ParentBlockMetaData")
				return err
			}
		}

		if s.slotNumber == 0 || blockMetaData.Epoch() != parentBlockMetaData.Epoch() {
			if err := s.epochMetaData.Commit(b); err != nil {
				log.Error("[Commit] Failed to commit EpochMetaData")
				return err
			}
		}

		err = b.Put(blockStorageKey, bytesBlock)
		if err != nil {
			log.Error("[Commit] Failed to commit block")
			return err
		}

		if isFinalizedState {
			// Update Main Chain Finalized Block Data
			s.mainChainMetaData.UpdateFinalizedBlockData(s.blockHeaderHash, s.slotNumber)
			s.mainChainMetaData.UpdateLastBlockData(s.blockHeaderHash, s.slotNumber)
			if err := s.mainChainMetaData.Commit(b); err != nil {
				log.Error("[Commit] Failed to commit MainChainMetaData")
				return err
			}
		}

		if totalStakeAmount.Cmp(lastBlockTotalStakeAmount) == 1 {
			// Update Main Chain Last Block Data
			s.mainChainMetaData.UpdateLastBlockData(s.blockHeaderHash, s.slotNumber)
			if err := s.mainChainMetaData.Commit(b); err != nil {
				log.Error("[Commit] Failed to commit MainChainMetaData")
				return err
			}

			err = b.Put([]byte("mainchain-head-trie-root"), trieRoot[:])
			if err != nil {
				log.Error("[Commit] Failed to commit state trie root")
				return err
			}
		}

		if !isFinalizedState {
			b, err = tx.CreateBucketIfNotExists(metadata.GetBlockBucketName(s.blockHeaderHash))
			if err != nil {
				log.Error("[Commit] Failed to create bucket")
				return err
			}
		}

		return nil
	})
}

func (s *StateContext) Finalize(blockMetaDataPathForFinalization []*metadata.BlockMetaData) error {
	bm := blockMetaDataPathForFinalization[len(blockMetaDataPathForFinalization)-1]
	parentBlockMetaData, err := metadata.GetBlockMetaData(s.db, bm.ParentHeaderHash())
	pHash := bm.ParentHeaderHash()
	if err != nil {
		log.Error("[Finalize] Failed to load ParentBlockMetaData ",
			hex.EncodeToString(pHash[:]))
		return err
	}
	return s.db.DB().Update(func(tx *bbolt.Tx) error {
		var err error
		mainBucket := tx.Bucket([]byte("DB"))
		for i := len(blockMetaDataPathForFinalization) - 1; i >= 0; i-- {
			bm := blockMetaDataPathForFinalization[i]
			blockBucket := tx.Bucket(metadata.GetBlockBucketName(bm.HeaderHash()))
			c := blockBucket.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				err = mainBucket.Put(k, v)
				if err != nil {
					log.Error("[Finalize] Finalization failed for key = ", k,
						"value ", v)
					return err
				}
			}

			expectedPHash := bm.ParentHeaderHash()
			foundPHash := parentBlockMetaData.HeaderHash()
			if !reflect.DeepEqual(parentBlockMetaData.HeaderHash(), bm.ParentHeaderHash()) {
				log.Error("[Finalize] Unexpected error parent block header hash not matching")
				log.Error("Expected ParentBlockHeaderHash ",
					hex.EncodeToString(expectedPHash[:]))
				log.Error("ParentBlockHeaderHash found ",
					hex.EncodeToString(foundPHash[:]))
				return errors.New("unexpected error parent block header hash not matching")
			}

			parentBlockMetaData.UpdateFinalizedChildHeaderHash(bm.HeaderHash())
			err := parentBlockMetaData.Commit(mainBucket)
			if err != nil {
				log.Error("[Finalize] Failed to Commit ParentBlockMetaData ",
					hex.EncodeToString(foundPHash[:]))
				return err
			}
			parentBlockMetaData = bm
			log.Info("Finalized Block #", bm.SlotNumber())
		}

		bm = blockMetaDataPathForFinalization[0]
		finalizedBlockHeaderHash := bm.HeaderHash()
		finalizedSlotNumber := bm.SlotNumber()
		s.mainChainMetaData.UpdateFinalizedBlockData(finalizedBlockHeaderHash, finalizedSlotNumber)
		return s.mainChainMetaData.Commit(mainBucket)
	})
}

<<<<<<< HEAD
func NewStateContext(db *db.DB, slotNumber uint64,
	blockProposer []byte, finalizedHeaderHash common.Hash,
	parentBlockHeaderHash common.Hash, blockHeaderHash common.Hash,
	partialBlockSigningHash common.Hash, blockSigningHash common.Hash,
=======
func NewStateContext(db db.DB, slotNumber uint64,
	blockProposer []byte, finalizedHeaderHash []byte,
	parentBlockHeaderHash []byte, blockHeaderHash []byte,
	partialBlockSigningHash []byte, blockSigningHash []byte,
>>>>>>> qcfork/unit-test
	epochMetaData *metadata.EpochMetaData) (*StateContext, error) {

	mainChainMetaData, err := metadata.GetMainChainMetaData(db)
	if err != nil {
		return nil, err
	}

	epoch := slotNumber / config.GetDevConfig().BlocksPerEpoch
	epochBlockHashes, err := metadata.GetEpochBlockHashes(db, epoch)
	if err != nil {
		epochBlockHashes = metadata.NewEpochBlockHashes(epoch)
	}

	attestorsFlag := make(map[string]bool)
	if slotNumber > 0 {
		slotInfo := epochMetaData.SlotInfo()[slotNumber%config.GetDevConfig().BlocksPerEpoch]
		for _, attestorsIndex := range slotInfo.Attestors {
			attestorsFlag[hex.EncodeToString(epochMetaData.Validators()[attestorsIndex])] = false
		}
	}

	return &StateContext{
		db: db,

		slotNumber:                   slotNumber,
		blockProposer:                blockProposer,
		finalizedHeaderHash:          finalizedHeaderHash,
		parentBlockHeaderHash:        parentBlockHeaderHash,
		blockHeaderHash:              blockHeaderHash,
		partialBlockSigningHash:      partialBlockSigningHash,
		blockSigningHash:             blockSigningHash,
		currentBlockTotalStakeAmount: big.NewInt(0),

		validatorsFlag: make(map[string]bool),

		epochMetaData:     epochMetaData,
		epochBlockHashes:  epochBlockHashes,
		mainChainMetaData: mainChainMetaData,
	}, nil
}
