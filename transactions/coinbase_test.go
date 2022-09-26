package transactions

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/theQRL/go-qrllib/dilithium"
	"github.com/theQRL/go-qrllib/xmss"
	"github.com/theQRL/zond/address"
	"github.com/theQRL/zond/db"
	"github.com/theQRL/zond/metadata"
	"github.com/theQRL/zond/misc"
	"github.com/theQRL/zond/state"
	"go.etcd.io/bbolt"
)

func TestNewCoinBase(t *testing.T) {
	networkId := uint64(1)
	blockProposer := dilithium.New()
	blockProposerPK := blockProposer.GetPK()
	blockProposerReward := uint64(10)
	attestorReward := uint64(10)
	feeReward := uint64(10)
	lastCoinBaseNonce := uint64(10)

	coinbase := NewCoinBase(networkId, blockProposerPK[:], blockProposerReward, attestorReward, feeReward, lastCoinBaseNonce)

	if coinbase.BlockProposerReward() != blockProposerReward {
		t.Errorf("blockproposer reward incorrectly set in new coinbase transaction, expected (%v), got (%v)", blockProposerReward, coinbase.BlockProposerReward())
	}
}

// func TestValidateData(t *testing.T) {
// 	networkId := uint64(1)
// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()
// 	blockProposerXmss := xmss.NewXMSSFromHeight(4, 0)
// 	blockProposerXmssPK := blockProposerXmss.GetPK()
// 	blockProposerAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((blockProposerXmssPK[:])))

// 	blockProposer_case3 := dilithium.New()
// 	blockProposerPK_case3 := blockProposer_case3.GetPK()

// 	blockProposerReward := uint64(5000000000)
// 	attestorReward := uint64(2000000000)
// 	feeReward := uint64(0)
// 	lastCoinBaseNonce := uint64(10)

// 	coinbase := NewCoinBase(networkId, blockProposerPK[:], blockProposerReward, attestorReward, feeReward, lastCoinBaseNonce)
// 	coinbase_case2 := NewCoinBase(networkId, blockProposerPK[:], blockProposerReward, attestorReward, feeReward, 9)
// 	coinbase_case3 := NewCoinBase(networkId, blockProposerPK_case3[:], blockProposerReward, attestorReward, feeReward, lastCoinBaseNonce)
// 	coinbase_case4 := NewCoinBase(networkId, blockProposerPK[:], uint64(500000000), attestorReward, feeReward, lastCoinBaseNonce)
// 	coinbase_case5 := NewCoinBase(networkId, blockProposerPK[:], blockProposerReward, uint64(200000000), feeReward, lastCoinBaseNonce)
// 	coinbase_case6 := NewCoinBase(networkId, blockProposerPK[:], blockProposerReward, attestorReward, 1, lastCoinBaseNonce)

// 	binCoinBaseAddress, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")

// 	var validators [][]byte
// 	validators = append(validators, blockProposerPK[:])
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.Sum256([]byte("blockSigningHash"))
// 	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
// 	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)

// 	addressState := address.NewAddressState(binCoinBaseAddress[:], 10, 20000000000001)
// 	blockproposerAddressState := address.NewAddressState(blockProposerAddr[:], 10, 20001)
// 	addressesState := make(map[string]*address.AddressState)
// 	addressesState[hex.EncodeToString(binCoinBaseAddress[:])] = addressState
// 	addressesState[hex.EncodeToString(blockProposerAddr[:])] = blockproposerAddressState

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
// 		coinbase       *CoinBase
// 		stateContext   state.StateContext
// 		expectedOutput bool
// 	}{
// 		{
// 			name:           "ok",
// 			coinbase:       coinbase,
// 			stateContext:   *stateContext,
// 			expectedOutput: true,
// 		},
// 		{
// 			name:           "incorrect nonce",
// 			coinbase:       coinbase_case2,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "failed to process blockproposer",
// 			coinbase:       coinbase_case3,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "incorrect block proposer reward",
// 			coinbase:       coinbase_case4,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "incorrect attestor reward",
// 			coinbase:       coinbase_case5,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 		{
// 			name:           "incorrect fee reward",
// 			coinbase:       coinbase_case6,
// 			stateContext:   *stateContext,
// 			expectedOutput: false,
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			output := tc.coinbase.validateData(&tc.stateContext)
// 			if output != tc.expectedOutput {
// 				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
// 			}

// 		})
// 	}
// }

// func TestCoinbaseValidate(t *testing.T) {
// 	networkId := uint64(1)
// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()
// 	blockProposerXmss := xmss.NewXMSSFromHeight(4, 0)
// 	blockProposerXmssPK := blockProposerXmss.GetPK()
// 	blockProposerAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((blockProposerXmssPK[:])))

// 	blockProposerReward := uint64(5000000000)
// 	attestorReward := uint64(2000000000)
// 	feeReward := uint64(0)
// 	lastCoinBaseNonce := uint64(10)

// 	coinbase := NewCoinBase(networkId, blockProposerPK[:], blockProposerReward, attestorReward, feeReward, lastCoinBaseNonce)

// 	binCoinBaseAddress, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")

// 	var validators [][]byte
// 	validators = append(validators, blockProposerPK[:])
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.Sum256([]byte("blockSigningHash"))
// 	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
// 	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
// 	addressState := address.NewAddressState(binCoinBaseAddress[:], 10, 20000000000001)
// 	blockproposerAddressState := address.NewAddressState(blockProposerAddr[:], 10, 20001)
// 	addressesState := make(map[string]*address.AddressState)
// 	addressesState[hex.EncodeToString(binCoinBaseAddress[:])] = addressState
// 	addressesState[hex.EncodeToString(blockProposerAddr[:])] = blockproposerAddressState

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

// 	coinbase.Sign(blockProposer, coinbase.GetSigningHash(stateContext.BlockSigningHash()).Bytes())

// 	testCases := []struct {
// 		name           string
// 		coinbase       *CoinBase
// 		stateContext   state.StateContext
// 		expectedOutput bool
// 	}{
// 		{
// 			name:           "ok",
// 			coinbase:       coinbase,
// 			stateContext:   *stateContext,
// 			expectedOutput: true,
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			output := tc.coinbase.Validate(&tc.stateContext)
// 			if output != tc.expectedOutput {
// 				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
// 			}

// 		})
// 	}
// }

// func TestCoinbaseApplyStateChanges(t *testing.T) {
// 	networkId := uint64(1)
// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()
// 	blockProposerXmss := xmss.NewXMSSFromHeight(4, 0)
// 	blockProposerXmssPK := blockProposerXmss.GetPK()
// 	blockProposerAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((blockProposerXmssPK[:])))

// 	blockProposerReward := uint64(5000000000)
// 	attestorReward := uint64(2000000000)
// 	feeReward := uint64(0)
// 	lastCoinBaseNonce := uint64(10)

// 	coinbase := NewCoinBase(networkId, blockProposerPK[:], blockProposerReward, attestorReward, feeReward, lastCoinBaseNonce)

// 	binCoinBaseAddress, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")

// 	var validators [][]byte
// 	validators = append(validators, blockProposerPK[:])
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.Sum256([]byte("blockSigningHash"))
// 	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
// 	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)

// 	addressState := address.NewAddressState(binCoinBaseAddress[:], 10, 20000000000001)
// 	blockproposerAddressState := address.NewAddressState(blockProposerAddr[:], 10, 20001)
// 	addressesState := make(map[string]*address.AddressState)
// 	addressesState[hex.EncodeToString(binCoinBaseAddress[:])] = addressState
// 	addressesState[hex.EncodeToString(blockProposerAddr[:])] = blockproposerAddressState

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

// 	err = coinbase.ApplyStateChanges(stateContext)
// 	if err != nil {
// 		t.Error("got unexpected error while applying state changes in coinbase transaction ", err)
// 	}
// }

func TestCoinbaseGetSigningHash(t *testing.T) {
	networkId := uint64(1)
	blockProposerPK, _ := hex.DecodeString("135febcea1b6c55fb951df921e40fdf445367967ff6bc7ed965b69c10940d9a8de63bcc06d5b66cf5088840fc0070c6a8c7072db68e3a5237df59be112b54e805a43c4e4bb41cfe5ae73011598bfebeaf6f81f8fdc4343767dc317d41fb7d686c07d342e77980b01fb38f6bcf8e36a7d2de9bbfc564b35638aa0d1a12d9b7ea48440117a4de9586ae0d4cd2ccd637e59c3fa75faa9bc4b65bdc078e9c6583e24503fa47af8fa8b6f24de491744eaad66ec4e072595fa8eeb688707d0f6083deedb5375952cb6b49ed5745f384827909e6b311b6ec9657211c2c611e3bfbfb613a7b4eb98b8fab70d3060716d6fa42480a8b738488110e672c00720d909d0293b286018d117a72bfa329d75cca061a415b185b581baf558ad1abf978eaf399ed2e46db8e9922a025c7ef0efe3414d0b704ab8900808d25a11aa1332a9e1c703738b86971f7f10f9d89d051238db68a0b671224ae6302769fcffabdda54d4b7c8f2bc37ed8c64e239dff3d0d1cb91dd771289a61c4c331e10433ab1db2dc5aa4dd72eaeaf02f46d4943e9c073cba6888dc114a1a880d7097c59fb7536033eb8fe33d64ffb47680a6e125a9bf85039d5087767ec1dba4db7cfbcc3897627ab17d448c915ce36da1d848e40687a6fc71c4101fd0a8c5338df962cd192264fb43c1e5ff525a139ad0433241039238927f0cac35a8e2455bcc6175fa8e91b8ffa606c252f484fbd4ea46b03685642fa24e55015d6a9eedd4f02119934aa44450dbf4f92ee86c4a6fcc7f4c138a4ec2cbef3edabd26a01b32513f728c9acb3c2859ded525bf0c717a10440c47859403e7d3519890d83b0438fb67a4dc146c07d3f2d2834062745e25562650ab039b108fc5949f07bbd289e1afb96f5c29b24f459276fc32361d98b2e3ac1e2ecc0f8d5f213145ff75fecdb2a8e20c2bfc05999c0b88669d3a6da2e9ca583562188e21f82f82bd61c4ceb73482818ae0bd40204c88c4a04aa954e2d42918bcace86f2929f92ae6d4ba8d5caba2a9210cf159f6b54076b49b63f7efde737cb3c1f5a5071458270a78655bd2872b438e6cbecb94ab208429ddbf9c4508c5b819fb278f7c85d3b2ea6b88c3b6604c9f34a95d8f823566ce3c46432bd9b4a280fe55a1ef1f250d92b5101bb649fa71cae1b1c1c048ee38d38b7a4616320cf3d4cb95d8133db2cf7ea8c2f6a48d0812c90cffa9f4ba4493ac483b1942d5690f77380351d6df48cc51dca17f1bf786f0981470036414520fc39cfb1eb3d0c12fe8f6a11fb0aef6569b2c67f567b996370b8808b90339802a89e09b0b4a95df01d7d0eb8568edafbc98b6847d4825dd8959ac99c7272c90ffe9ef88b642e60706604f48fa33b5e3a7b87a68b5f89a63e55dd453d90eb47ddf6f3bd5c397d13e4c2d59dff5a9f969f5516ee7afed194b5aba83d0f2a2a838aaddb1103e4a503ccbc100633a61076b037408994ae9e586bc307e68519c47a503358b9c388df500cf74ad104d11de30a4a83778850029eff3b29c746f9de92038a2109859cdc6e4f480ca692e845aef592801a5ec4e71433f8ee162d49877fc7009d80cb1c6cb0b3ed407c5a80a7fda255ab2a358eed06caf7d51db6d560d8e60cfe71bbe1c2ca6a7462c817ab2a71f09aff1e36688aafb82eb094622f8c7e0dc6a19833cdaa6cb421cfebd475b9ddbd8257a96a75fe6745d59fece8faa540e5204f3d2b0146e8658a2e006f751df647b86c3ea2dbe66f5bd3f36a5a51fec55f7bb2d3595d52caac3d6f7eabfda533bfd7ace942ef44ec5ce70f7b67c2a229ad0b2855be60f2fe22771cf53de469c4cea25bdaea60c46aa754d8a4cc40f8a6d87610b4c01c32fa6c404f8d571c2d76149982e752ba802930472950bc15d38e98ab865b9ebbd9bf5801df3587110285663daf311ea7e5a0444ce03b6a7e574557521c068490e774d2a331ac71176391c434c3ec65e8fae60505ba408c661998b9c1469773fbdac6e2be31747333611bbdfdbd9a49abee970b92a4f2e26fbaa1549940923c8d4d9d12a4a0bc0f09e453dc0d75124f010cd6d55a5a590b17fb94a6913b5")
	blockProposerXmssPK, _ := hex.DecodeString("00030001e9460db4c459269f5545524d78da462945821ff47e5191b9664735b071e18f5fec7dcb78e39a9ad6f075b7c6ef629325bd716f9719a452362744880b5341e7")
	blockProposerAddr := xmss.GetXMSSAddressFromPK(misc.UnSizedXMSSPKToSizedPK((blockProposerXmssPK[:])))

	blockProposerReward := uint64(5000000000)
	attestorReward := uint64(2000000000)
	feeReward := uint64(0)
	lastCoinBaseNonce := uint64(10)

	coinbase := NewCoinBase(networkId, blockProposerPK[:], blockProposerReward, attestorReward, feeReward, lastCoinBaseNonce)

	binCoinBaseAddress, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")

	var validators [][]byte
	validators = append(validators, blockProposerPK[:])
	epoch := uint64(1)
	slotNumber := uint64(100)
	finalizedHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))
	blockSigningHash := sha256.Sum256([]byte("blockSigningHash"))
	epochMetadata := metadata.NewEpochMetaData(epoch, parentBlockHeaderHash, validators)
	epochMetadata.AllotSlots(1, epoch, parentBlockHeaderHash)
	mainChainMetaData := metadata.NewMainChainMetaData(finalizedHeaderHash, 1,
		parentBlockHeaderHash, 0)
	addressState := address.NewAddressState(binCoinBaseAddress[:], 10, 20000000000001)
	blockproposerAddressState := address.NewAddressState(blockProposerAddr[:], 10, 20001)
	addressesState := make(map[string]*address.AddressState)
	addressesState[hex.EncodeToString(binCoinBaseAddress[:])] = addressState
	addressesState[hex.EncodeToString(blockProposerAddr[:])] = blockproposerAddressState

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

	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
		blockSigningHash, epochMetadata)
	if err != nil {
		t.Error("unexpected error while creating new statecontext ", err)
	}

	expectedHash, _ := hex.DecodeString("eff0390caec8827bbdd5aa41d75de330d701fc1c2741a432205bdde787cd3131")
	output := coinbase.GetSigningHash(stateContext.BlockSigningHash())

	if hex.EncodeToString(output.Bytes()) != hex.EncodeToString(expectedHash) {
		t.Errorf("expected coinbase signing hash (%v), got (%v)", hex.EncodeToString(expectedHash), hex.EncodeToString(output.Bytes()))
	}
}

func TestCoinbaseUnsignedHash(t *testing.T) {
	networkId := uint64(1)
	blockProposerPK, _ := hex.DecodeString("135febcea1b6c55fb951df921e40fdf445367967ff6bc7ed965b69c10940d9a8de63bcc06d5b66cf5088840fc0070c6a8c7072db68e3a5237df59be112b54e805a43c4e4bb41cfe5ae73011598bfebeaf6f81f8fdc4343767dc317d41fb7d686c07d342e77980b01fb38f6bcf8e36a7d2de9bbfc564b35638aa0d1a12d9b7ea48440117a4de9586ae0d4cd2ccd637e59c3fa75faa9bc4b65bdc078e9c6583e24503fa47af8fa8b6f24de491744eaad66ec4e072595fa8eeb688707d0f6083deedb5375952cb6b49ed5745f384827909e6b311b6ec9657211c2c611e3bfbfb613a7b4eb98b8fab70d3060716d6fa42480a8b738488110e672c00720d909d0293b286018d117a72bfa329d75cca061a415b185b581baf558ad1abf978eaf399ed2e46db8e9922a025c7ef0efe3414d0b704ab8900808d25a11aa1332a9e1c703738b86971f7f10f9d89d051238db68a0b671224ae6302769fcffabdda54d4b7c8f2bc37ed8c64e239dff3d0d1cb91dd771289a61c4c331e10433ab1db2dc5aa4dd72eaeaf02f46d4943e9c073cba6888dc114a1a880d7097c59fb7536033eb8fe33d64ffb47680a6e125a9bf85039d5087767ec1dba4db7cfbcc3897627ab17d448c915ce36da1d848e40687a6fc71c4101fd0a8c5338df962cd192264fb43c1e5ff525a139ad0433241039238927f0cac35a8e2455bcc6175fa8e91b8ffa606c252f484fbd4ea46b03685642fa24e55015d6a9eedd4f02119934aa44450dbf4f92ee86c4a6fcc7f4c138a4ec2cbef3edabd26a01b32513f728c9acb3c2859ded525bf0c717a10440c47859403e7d3519890d83b0438fb67a4dc146c07d3f2d2834062745e25562650ab039b108fc5949f07bbd289e1afb96f5c29b24f459276fc32361d98b2e3ac1e2ecc0f8d5f213145ff75fecdb2a8e20c2bfc05999c0b88669d3a6da2e9ca583562188e21f82f82bd61c4ceb73482818ae0bd40204c88c4a04aa954e2d42918bcace86f2929f92ae6d4ba8d5caba2a9210cf159f6b54076b49b63f7efde737cb3c1f5a5071458270a78655bd2872b438e6cbecb94ab208429ddbf9c4508c5b819fb278f7c85d3b2ea6b88c3b6604c9f34a95d8f823566ce3c46432bd9b4a280fe55a1ef1f250d92b5101bb649fa71cae1b1c1c048ee38d38b7a4616320cf3d4cb95d8133db2cf7ea8c2f6a48d0812c90cffa9f4ba4493ac483b1942d5690f77380351d6df48cc51dca17f1bf786f0981470036414520fc39cfb1eb3d0c12fe8f6a11fb0aef6569b2c67f567b996370b8808b90339802a89e09b0b4a95df01d7d0eb8568edafbc98b6847d4825dd8959ac99c7272c90ffe9ef88b642e60706604f48fa33b5e3a7b87a68b5f89a63e55dd453d90eb47ddf6f3bd5c397d13e4c2d59dff5a9f969f5516ee7afed194b5aba83d0f2a2a838aaddb1103e4a503ccbc100633a61076b037408994ae9e586bc307e68519c47a503358b9c388df500cf74ad104d11de30a4a83778850029eff3b29c746f9de92038a2109859cdc6e4f480ca692e845aef592801a5ec4e71433f8ee162d49877fc7009d80cb1c6cb0b3ed407c5a80a7fda255ab2a358eed06caf7d51db6d560d8e60cfe71bbe1c2ca6a7462c817ab2a71f09aff1e36688aafb82eb094622f8c7e0dc6a19833cdaa6cb421cfebd475b9ddbd8257a96a75fe6745d59fece8faa540e5204f3d2b0146e8658a2e006f751df647b86c3ea2dbe66f5bd3f36a5a51fec55f7bb2d3595d52caac3d6f7eabfda533bfd7ace942ef44ec5ce70f7b67c2a229ad0b2855be60f2fe22771cf53de469c4cea25bdaea60c46aa754d8a4cc40f8a6d87610b4c01c32fa6c404f8d571c2d76149982e752ba802930472950bc15d38e98ab865b9ebbd9bf5801df3587110285663daf311ea7e5a0444ce03b6a7e574557521c068490e774d2a331ac71176391c434c3ec65e8fae60505ba408c661998b9c1469773fbdac6e2be31747333611bbdfdbd9a49abee970b92a4f2e26fbaa1549940923c8d4d9d12a4a0bc0f09e453dc0d75124f010cd6d55a5a590b17fb94a6913b5")
	blockProposerReward := uint64(5000000000)
	attestorReward := uint64(2000000000)
	feeReward := uint64(0)
	lastCoinBaseNonce := uint64(10)

	coinbase := NewCoinBase(networkId, blockProposerPK[:], blockProposerReward, attestorReward, feeReward, lastCoinBaseNonce)

	expectedHash, _ := hex.DecodeString("7ef0387466dc9f64955b73151c59a87d148940b18b67b515fa65691e23469a52")
	output := coinbase.GetUnsignedHash()

	if hex.EncodeToString(output.Bytes()) != hex.EncodeToString(expectedHash) {
		t.Errorf("expected coinbase unsigned hash (%v), got (%v)", hex.EncodeToString(expectedHash), hex.EncodeToString(output.Bytes()))
	}
}
