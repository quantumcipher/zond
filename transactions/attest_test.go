package transactions

import (
	"testing"
)

func TestNewAttest(t *testing.T) {
	networkID := uint64(1)
	coinBaseNonce := uint64(10)
	attest := NewAttest(networkID, coinBaseNonce)

	if attest.Nonce() != coinBaseNonce {
		t.Error("the nonce is not set correctly while attesting transaction")
	}
}

// func TestAttestValidateData(t *testing.T) {
// 	networkID := uint64(1)
// 	coinBaseNonce := uint64(10)
// 	attest := NewAttest(networkID, coinBaseNonce)

// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()
// 	attestorDilithium := dilithium.New()
// 	attestorDilithiumPK := attestorDilithium.GetPK()

// 	var validators [][]byte
// 	validators = append(validators, blockProposerPK[:])
// 	validators = append(validators, attestorDilithiumPK[:])
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.Sum256([]byte("blockSigningHash"))
// 	prevSlotLastBlockHeaderHash := sha256.Sum256([]byte("prevSlotLastBlockHeaderHash"))
// 	epochMetadaPBData := &protos.EpochMetaData{
// 		SlotInfo: []*protos.SlotInfo{{
// 			SlotLeader: 0,
// 			Attestors:  []uint64{uint64(1)},
// 		}},
// 		Validators: validators,
// 	}
// 	epochMetadata := metadata.NewEpochMetaData(epoch, prevSlotLastBlockHeaderHash, nil)
// 	epochMetadataSerialized, _ := proto.Marshal(epochMetadaPBData)
// 	epochMetadata.DeSerialize(epochMetadataSerialized)

// 	dir, err := os.MkdirTemp("", "tempdir")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	defer os.RemoveAll(dir) // clean up

// 	file := filepath.Join(dir, "tmpfile.txt")
// 	if err := os.WriteFile(file, []byte("content"), 0666); err != nil {
// 		t.Error(err)
// 	}

// 	store, err := db.NewDB(dir, "tmpfile.txt")
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
// 		attest         *Attest
// 		stateContext   state.StateContext
// 		dilithiumKey   *dilithium.Dilithium
// 		expectedOutput bool
// 	}{
// 		{
// 			name:           "ok",
// 			attest:         attest,
// 			stateContext:   *stateContext,
// 			dilithiumKey:   attestorDilithium,
// 			expectedOutput: true,
// 		},
// 		{
// 			name:           "attestor key not present",
// 			attest:         attest,
// 			stateContext:   *stateContext,
// 			dilithiumKey:   blockProposer,
// 			expectedOutput: false,
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.attest.Sign(tc.dilithiumKey, attest.GetSigningHash(stateContext.PartialBlockSigningHash()).Bytes())

// 			output := tc.attest.validateData(stateContext)
// 			if output != tc.expectedOutput {
// 				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
// 			}

// 		})
// 	}
// }

// func TestAttestValidate(t *testing.T) {
// 	networkID := uint64(1)
// 	coinBaseNonce := uint64(10)
// 	attest := NewAttest(networkID, coinBaseNonce)

// 	blockProposer := dilithium.New()
// 	blockProposerPK := blockProposer.GetPK()

// 	attestorDilithium := dilithium.New()
// 	attestorDilithiumPK := attestorDilithium.GetPK()

// 	var validators [][]byte
// 	validators = append(validators, blockProposerPK[:])
// 	validators = append(validators, attestorDilithiumPK[:])
// 	epoch := uint64(1)
// 	slotNumber := uint64(100)
// 	finalizedHeaderHash := sha256.Sum256([]byte("finalizedHeaderHash"))
// 	parentBlockHeaderHash := sha256.Sum256([]byte("parentBlockHeaderHash"))
// 	blockHeaderHash := sha256.Sum256([]byte("blockHeaderHash"))
// 	partialBlockSigningHash := sha256.Sum256([]byte("partialBlockSigningHash"))
// 	blockSigningHash := sha256.Sum256([]byte("blockSigningHash"))
// 	prevSlotLastBlockHeaderHash := sha256.Sum256([]byte("prevSlotLastBlockHeaderHash"))
// 	epochMetadaPBData := &protos.EpochMetaData{
// 		SlotInfo: []*protos.SlotInfo{{
// 			SlotLeader: 0,
// 			Attestors:  []uint64{uint64(1)},
// 		}},
// 		Validators: validators,
// 	}
// 	epochMetadata := metadata.NewEpochMetaData(epoch, prevSlotLastBlockHeaderHash, nil)
// 	epochMetadataSerialized, _ := proto.Marshal(epochMetadaPBData)
// 	epochMetadata.DeSerialize(epochMetadataSerialized)

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

// 	stateContext, err := state.NewStateContext(store, slotNumber, blockProposerPK[:], finalizedHeaderHash, parentBlockHeaderHash, blockHeaderHash, partialBlockSigningHash,
// 		blockSigningHash, epochMetadata)

// 	if err != nil {
// 		t.Error("unexpected error while creating new statecontext ", err)
// 	}

// 	testCases := []struct {
// 		name           string
// 		attest         *Attest
// 		stateContext   state.StateContext
// 		dilithiumKey   *dilithium.Dilithium
// 		expectedOutput bool
// 	}{
// 		{
// 			name:           "ok",
// 			attest:         attest,
// 			stateContext:   *stateContext,
// 			dilithiumKey:   attestorDilithium,
// 			expectedOutput: true,
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.attest.Sign(tc.dilithiumKey, attest.GetSigningHash(stateContext.PartialBlockSigningHash()).Bytes())

// 			output := tc.attest.Validate(stateContext)
// 			if output != tc.expectedOutput {
// 				t.Errorf("expected output of validate data to be (%v) but returned (%v)", tc.expectedOutput, output)
// 			}

// 		})
// 	}
// }
