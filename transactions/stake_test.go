package transactions

import (
	"encoding/hex"
	"testing"

	"github.com/theQRL/go-qrllib/dilithium"
)

func TestNewStake(t *testing.T) {
	masterDilithium := dilithium.New()
	masterDilithiumPK := masterDilithium.GetPK()
	networkID := uint64(1)
	stake := uint64(10)
	fee := uint64(1)
	nonce := uint64(10)

	staking := NewStake(networkID, stake, fee, uint64(0), nonce, masterDilithiumPK[:])
	if staking.Amount() != stake {
		t.Error("the stake is incorrectly set")
	}
}

// func TestValidateStakeData(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	masterXmss := xmss.NewXMSSFromHeight(4, 0)
// 	masterXmssPK := masterXmss.GetPK()
// 	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((masterXmssPK[:])))

// 	masterXmss2 := xmss.NewXMSSFromHeight(4, 0)
// 	masterXmss2PK := masterXmss2.GetPK()
// 	masterAddr2 := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((masterXmss2PK[:])))

// 	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
// 	slaveXmss1PK := slaveXmss1.GetPK()

// 	networkID := uint64(1)
// 	validatorDilithium := dilithium.New()
// 	validatorDilithiumPK := validatorDilithium.GetPK()
// 	validatorDilithium2 := dilithium.New()
// 	validatorDilithium2PK := validatorDilithium2.GetPK()
// 	validatorDilithium3 := dilithium.New()
// 	validatorDilithium3PK := validatorDilithium3.GetPK()
// 	validatorXmss3 := xmss.NewXMSSFromHeight(4, 0)
// 	validatorXmss3Addr := validatorXmss3.GetAddress()

// 	var dilithiumPKs [][]byte
// 	dilithiumPKs = append(dilithiumPKs, validatorDilithiumPK[:])
// 	dilithiumPKs = append(dilithiumPKs, validatorDilithium2PK[:])

// 	var dilithiumPKs_case5 [][]byte
// 	dilithiumPKs_case5 = append(dilithiumPKs_case5, validatorDilithiumPK[:])
// 	dilithiumPKs_case5 = append(dilithiumPKs_case5, validatorDilithium2PK[:])
// 	dilithiumPKs_case5 = append(dilithiumPKs_case5, validatorDilithium3PK[:])
// 	dilithiumPKs_case7 := make([][]byte, 101)

// 	var dilithiumPKs_case8 [][]byte
// 	dilithiumPKs_case8 = append(dilithiumPKs_case8, validatorDilithium2PK[:])
// 	dilithiumPKs_case8 = append(dilithiumPKs_case8, validatorDilithium3PK[:])

// 	stake := true
// 	fee := uint64(1)
// 	nonce := uint64(10)

// 	staking := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])
// 	staking_case2 := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr2[:])
// 	staking_case3 := NewStake(networkID, dilithiumPKs, stake, fee, 20, slaveXmss1PK[:], masterAddr[:])
// 	staking_case4 := NewStake(networkID, dilithiumPKs, stake, 20000000000002, nonce, slaveXmss1PK[:], masterAddr[:])
// 	staking_case5 := NewStake(networkID, dilithiumPKs_case5, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])
// 	staking_case6 := NewStake(networkID, dilithiumPKs, false, fee, nonce, slaveXmss1PK[:], masterAddr[:])
// 	staking_case7 := NewStake(networkID, dilithiumPKs_case7, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])
// 	staking_case8 := NewStake(networkID, dilithiumPKs_case8, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])

// 	var validators [][]byte
// 	validators = append(validators, validatorDilithiumPK[:])
// 	validators = append(validators, validatorDilithium2PK[:])
// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.Sum256([]byte("blockSigningHash"))
// 	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
// 	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)

// 	addressState := address.NewAddressState(masterAddr[:], nonce, 20000000000001)
// 	addressesState := make(map[string]*address.AddressState)
// 	addressesState[hex.EncodeToString(masterAddr[:])] = addressState

// 	dir, err := os.MkdirTemp("", "tempdir")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer os.RemoveAll(dir) // clean up

// 	file := filepath.Join(dir, "tmpfile.txt")
// 	if err := os.WriteFile(file, []byte("content"), 0666); err != nil {
// 		t.Error(err)
// 	}

// 	store, err := db.NewDB(dir, file)
// 	if err != nil {
// 		t.Error("unexpected error while creating new db ", err)
// 	}

// 	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
// 		blockSigningHash, epochMetadata)
// 	if err != nil {
// 		t.Error("unexpected error while creating new statecontext ", err)
// 	}

// 	testCases := []struct {
// 		name           string
// 		staking        *Stake
// 		stateContext   state.StateContext
// 		expectedOutput bool
// 	}{
// 		{
// 			name:           "ok",
// 			staking:        staking,
// 			stateContext:   *stateContext,
// 			expectedOutput: true,
// 		},
// 		{
// 			name:           "from address missing from statecontext",
// 			staking:        staking_case2,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "incorrect nonce",
// 			staking:        staking_case3,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "insufficient balance",
// 			staking:        staking_case4,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "insufficient staking balance",
// 			staking:        staking_case5,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "dilithium metadata not found",
// 			staking:        staking_case6,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "dilithium PK length beyond limit",
// 			staking:        staking_case7,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "dilithium key associated with another QRL address",
// 			staking:        staking_case8,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			output := tc.staking.validateData(&tc.stateContext)
// 			if output != tc.expectedOutput {
// 				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
// 			}

// 		})
// 	}
// }

// func TestStakeValidate(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	masterXmss := xmss.NewXMSSFromHeight(4, 0)
// 	masterXmssPK := masterXmss.GetPK()
// 	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((masterXmssPK[:])))

// 	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
// 	slaveXmss1PK := slaveXmss1.GetPK()

// 	networkID := uint64(1)
// 	validatorDilithium := dilithium.New()
// 	validatorDilithiumPK := validatorDilithium.GetPK()
// 	validatorDilithium2 := dilithium.New()
// 	validatorDilithium2PK := validatorDilithium2.GetPK()
// 	var dilithiumPKs [][]byte
// 	dilithiumPKs = append(dilithiumPKs, validatorDilithiumPK[:])
// 	dilithiumPKs = append(dilithiumPKs, validatorDilithium2PK[:])

// 	stake := true
// 	fee := uint64(1)
// 	nonce := uint64(10)

// 	staking := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])
// 	if staking.Stake() != stake {
// 		t.Error("the stake is incorrectly set")
// 	}
// 	staking.Sign(slaveXmss1, staking.GetSigningHash())

// 	var validators [][]byte
// 	validators = append(validators, validatorDilithiumPK[:])
// 	validators = append(validators, validatorDilithium2PK[:])
// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.Sum256([]byte("blockSigningHash"))
// 	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
// 	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
// 	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
// 		parentBlockHeaderHash, 0)
// 	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
// 	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
// 	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
// 	addressState := address.NewAddressState(masterAddr[:], nonce, 20000000000001)
// 	addressesState := make(map[string]*address.AddressState)
// 	addressesState[hex.EncodeToString(masterAddr[:])] = addressState
// 	addressesStateSerialized, _ := addressState.Serialize()
// 	dbAddressStateKey := address.GetAddressStateKey(masterAddr[:])
// 	slaveMetadata1 := metadata.NewSlaveMetaData(sha256.Sum256([]byte("transactionHash")), masterAddr[:], slaveXmss1PK[:])
// 	slaveMetadataSerialized, _ := slaveMetadata1.Serialize()
// 	slaveState := make(map[string]*metadata.SlaveMetaData)
// 	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:]))] = slaveMetadata1

// 	store := mockdb.NewMockDB(ctrl)
// 	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
// 	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(dbAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(addressesStateSerialized, nil).AnyTimes()
// 	store.EXPECT().
// 		GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).
// 		Return(slaveMetadataSerialized, nil).AnyTimes()

// 	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
// 		blockSigningHash, epochMetadata)
// 	stateContext.PrepareAddressState(hex.EncodeToString(masterAddr[:]))
// 	stateContext.PrepareSlaveMetaData(hex.EncodeToString(masterAddr[:]), hex.EncodeToString(slaveXmss1PK[:]))
// 	if err != nil {
// 		t.Error("unexpected error while creating new statecontext ", err)
// 	}

// 	testCases := []struct {
// 		name           string
// 		staking        *Stake
// 		stateContext   state.StateContext
// 		expectedOutput bool
// 	}{
// 		{
// 			name:           "ok",
// 			staking:        staking,
// 			stateContext:   *stateContext,
// 			expectedOutput: true,
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			output := tc.staking.Validate(&tc.stateContext)
// 			if output != tc.expectedOutput {
// 				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
// 			}

// 		})
// 	}
// }

// func TestStakeApplyStateChanges(t *testing.T) {
// 	ctrl := gomock.NewController(t)

// 	masterXmss := xmss.NewXMSSFromHeight(4, 0)
// 	masterXmssPK := masterXmss.GetPK()
// 	masterAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedPKToSizedPK((masterXmssPK[:])))

// 	slaveXmss1 := xmss.NewXMSSFromHeight(6, 0)
// 	slaveXmss1Addr := slaveXmss1.GetAddress()
// 	slaveXmss1PK := slaveXmss1.GetPK()

// 	networkID := uint64(1)
// 	validatorDilithium := dilithium.New()
// 	validatorDilithiumPK := validatorDilithium.GetPK()
// 	validatorDilithium2 := dilithium.New()
// 	validatorDilithium2PK := validatorDilithium2.GetPK()
// 	var dilithiumPKs [][]byte
// 	dilithiumPKs = append(dilithiumPKs, validatorDilithiumPK[:])
// 	dilithiumPKs = append(dilithiumPKs, validatorDilithium2PK[:])

// 	stake := true
// 	fee := uint64(1)
// 	nonce := uint64(10)

// 	staking := NewStake(networkID, dilithiumPKs, stake, fee, nonce, slaveXmss1PK[:], masterAddr[:])

// 	var validators [][]byte
// 	validators = append(validators, validatorDilithiumPK[:])
// 	validators = append(validators, validatorDilithium2PK[:])
// 	validatorXmss := xmss.NewXMSSFromHeight(8, 0)
// 	validatorXmssAddr := validatorXmss.GetAddress()
// 	validatorXmss2 := xmss.NewXMSSFromHeight(10, 0)
// 	validatorXmss2Addr := validatorXmss2.GetAddress()
// 	validatorDilithiumMetadata := metadata.NewDilithiumMetaData(sha256.Sum256([]byte("transactionHash")), validatorDilithiumPK[:], validatorXmssAddr[:], true)
// 	validatorDilithiumMetadataSerialized, _ := validatorDilithiumMetadata.Serialize()
// 	validatorDilithiumMetadata2 := metadata.NewDilithiumMetaData(sha256.Sum256([]byte("transactionHash")), validatorDilithium2PK[:], validatorXmss2Addr[:], true)
// 	validatorDilithiumMetadata2Serialized, _ := validatorDilithiumMetadata2.Serialize()

// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.Sum256([]byte("blockSigningHash"))
// 	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
// 	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
// 	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
// 		parentBlockHeaderHash, 0)
// 	mainChainMetaDataSerialized, _ := mainChainMetaData.Serialize()
// 	epochBlockHashesMetadata := metadata.NewEpochBlockHashes(epoch)
// 	epochBlockHashesMetadataSerialized, _ := epochBlockHashesMetadata.Serialize()
// 	addressState := address.NewAddressState(masterAddr[:], nonce, 20000000000001)
// 	slaveAddressState := address.NewAddressState(masterAddr[:], nonce, 20)
// 	slaveAddressStateSerialized, _ := slaveAddressState.Serialize()
// 	addressesState := make(map[string]*address.AddressState)
// 	addressesState[hex.EncodeToString(masterAddr[:])] = addressState
// 	addressesStateSerialized, _ := addressState.Serialize()
// 	dbAddressStateKey := address.GetAddressStateKey(masterAddr[:])
// 	slaveMetadata1 := metadata.NewSlaveMetaData(sha256.Sum256([]byte("transactionHash")), masterAddr[:], slaveXmss1PK[:])
// 	slaveMetadataSerialized, _ := slaveMetadata1.Serialize()
// 	slaveState := make(map[string]*metadata.SlaveMetaData)
// 	slaveState[hex.EncodeToString(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:]))] = slaveMetadata1

// 	store := mockdb.NewMockDB(ctrl)
// 	store.EXPECT().Get(gomock.Eq(metadata.GetMainChainMetaDataKey())).Return(mainChainMetaDataSerialized, nil).AnyTimes()
// 	store.EXPECT().Get(gomock.Eq(metadata.GetEpochBlockHashesKey(epoch))).Return(epochBlockHashesMetadataSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(dbAddressStateKey), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(addressesStateSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(address.GetAddressStateKey(slaveXmss1Addr[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(slaveAddressStateSerialized, nil).AnyTimes()
// 	store.EXPECT().
// 		GetFromBucket(gomock.Eq(metadata.GetSlaveMetaDataKey(masterAddr[:], slaveXmss1PK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).
// 		Return(slaveMetadataSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(metadata.GetDilithiumMetaDataKey(validatorDilithiumPK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(validatorDilithiumMetadataSerialized, nil).AnyTimes()
// 	store.EXPECT().GetFromBucket(gomock.Eq(metadata.GetDilithiumMetaDataKey(validatorDilithium2PK[:])), gomock.Eq(metadata.GetBlockBucketName(parentBlockHeaderHash))).Return(validatorDilithiumMetadata2Serialized, nil).AnyTimes()

// 	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
// 		blockSigningHash, epochMetadata)
// 	stateContext.PrepareAddressState(hex.EncodeToString(masterAddr[:]))
// 	stateContext.PrepareAddressState(hex.EncodeToString(slaveXmss1Addr[:]))
// 	stateContext.PrepareSlaveMetaData(hex.EncodeToString(masterAddr[:]), hex.EncodeToString(slaveXmss1PK[:]))
// 	stateContext.PrepareDilithiumMetaData(hex.EncodeToString(validatorDilithiumPK[:]))
// 	stateContext.PrepareDilithiumMetaData(hex.EncodeToString(validatorDilithium2PK[:]))
// 	if err != nil {
// 		t.Error("unexpected error while creating new statecontext ", err)
// 	}

// 	err = staking.ApplyStateChanges(stateContext)
// 	if err != nil {
// 		t.Error("got unexpected error while applying state changes in staking transaction ", err)
// 	}
// }

func TestStakeGetSigningHash(t *testing.T) {
	masterDilithiumPK, _ := hex.DecodeString("a2391b9b2464c680ecb638d2664ed3a35b45e99d0550aa4b79218cfe24dd81e8e94991644f90eb427f53a8d64e33f9ef95272b8fb4a3a644684704c299757af6f463337c612a4c022bc259b3db34c057939745b385d919f3804ae493cfe921f64dc2faf63a8887b2ffb93e23c4ace09312dedd3fc6a5afd3788cc1896881ccc989675d8de66baead300847509b578cf657dedfd3e889c589455e4ab23eb4d52e92a66826a0596c25ce787f7ee7b5d4111ee0e328d8068e5f3f6bda92dd3a7dfa8b19aac10d8eef73cf261eaf81adaef2f437d7d121098d24ad4bd14eeb39136807e107ab9cffac8baf5ef989317d8e2c71907260692bd6f4b976ecc329a60630c5056a4fbe05a347932a4b5e15c6495ad4ee1f50565d0fae8b97b78411e99dd11a2d2666785a8074b436864175bce2d6ba33ee0a377a1f7aef572e2fd0621daec16e26dadef1e82bad9e5c0d58ca4742c45121dcd4f70f72c75a888c1ff92e9e57c6b4eda7812cee85919192dce5ffd58b109812a2cebd2e088af230204cab6f268b65d48592dd9f02f43724f9d083e1e63c3d1ababa9c814ec47ad63c3197b98c33d14a984d8dab36c228284fff02a3a79dac03355d72f93a561eaac60b8f50ed8b4e842eebef8e1fe9470f3b599566043168f62bf1e250ee6ba46ff89be16e16b9eaf47a171df43b601447e32cdce1c13df9aef825975594df0e6e97aa06afd56bff24d50364cdc5e59cc1265e5b7f6b6ed2545e110267cb91dc3d779821fe7c332468ebf5e08b6fd301a401ab0d2852fa615d7f5ff8088c5a6d385d708d3718cce8edbbe592ec712bfd3e4075c10b347c19713f8b0ff8d006ae4aea02a8143dd6f9cb1ff1e317bc10ce5fc8e2040e093b53b70bb8a1558e58682bd6515524e23ae74ede7d855d885265a9342954e01a628c6605f25a09ac10758de9a4eecf210765b42315916def2467e43157e88e620ec5dc2d8ceae64843380e9855cfe9fd0b7ee62b03e6b2b55c5e5bce411611e66613c1c155cb9ff9f2e3e850bd51ac88f68b3d5f969134bbc590ee907e952f2b950a0cd60cf1aac1e0ffa603485d44fe1f30829ece14233b4a3201445791cd0514ec47d1a112e89bcfedf413a09be64a87f29a093c1989caa8c46e54f31669732b6aeaf05f51363e6b0f93a8001144cff199483ec0efd0ed1b130759fd814b07533baf055049ac9a224cdd88b6416bab6cb8b321c108b18637e297364e4678da24d518f9236c1bc1a19388ba33b35c2dbb76f2ffd2fb8f9ed277bb0202b12ce2ca3138b3c430399e78593fea82e51a5efe9c45fc500ad73ea537db049beeb3d0a6bab135b990a825c528f1563213f0c9dc6d1be2d03913e8f5b67b171b66ff1e5ff786302a37f0d1a0913d0e45c751f8a1ff70d15506a218dcb7a6b390af0117f8cd08b169d945c2649463150bae3fe994968a7a4e18e655b9c139d74d8dac8be031ac53c18ff1931a9b9e1285846cc97945a094c2eb844f8104086fb58593838257a288482f1fa52c0eaf98040f5ba398506cc61e3d60a2f3f2c6d25b2db1a5bdf9e59aebb17241002f37df32d8f06e9bdbf07c1a6ac331e9f7f4a81ccc265218d13ea772471a3e6134d15ed05ff8e221e3fdffbdde7f4a05624b33475a87200554c85ae575173b2fd304c4831a56792a5a7b1f371eb8a3e43f14604ad7c591818d916d41588a1512f61049b8f2d491055f1834f5df18e2d08302d389c4857876016cee399e88901293e894ccf5ee20a7fc20fe7c0c131611cdd276a93ad1eff804ad7e2a4c807cd8af3651f06c26183ddee13ae44083693b69424934b535d03bf807abe47b774dcacc545869c0261b936259e569686aa99b6e388b237af5bbd4401cd90b6f3ce96223f3e9522e9c571ed118b79dd1470d62d140ab7114aa3c2efe64114c32a53af6ef0078f0d3bee84514ea9a11c0f62c184fcbd5909be9a62c51e036ddd2daa31321fad3fdaca826a4ec2ac73312a36c996f967e891d2afe3d6a53d8b71e9593ac88a922433b1da676e0f429d22a15431dfceacad1f6a516e2f76befbf3c2751109c39c8879800")

	networkID := uint64(1)
	stake := uint64(10)
	fee := uint64(1)
	nonce := uint64(10)

	staking := NewStake(networkID, stake, fee, 0, nonce, masterDilithiumPK[:])

	expectedSigningHash, _ := hex.DecodeString("4648ce1859c445e5e9ad0f64f252c5a7690b997923a0ff76006e73fd1f06d35b")
	output := staking.GetSigningHash()
	if hex.EncodeToString(output.Bytes()) != hex.EncodeToString(expectedSigningHash) {
		t.Errorf("expected stake signing hash (%v), got (%v)", hex.EncodeToString(expectedSigningHash), hex.EncodeToString(output.Bytes()))
	}
}
