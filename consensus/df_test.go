package consensus

import (
	"encoding/hex"
	"testing"
)

func TestComputeDelayedFunction(t *testing.T) {
	input := []byte("random byte input")
	numberOfElements := 10
	expectedOutput := "72616e646f6d206279746520696e707574"

	output := ComputeDelayedFunction(input, numberOfElements)

	if hex.EncodeToString(output[0]) != expectedOutput {
		t.Errorf("expected output (%v) got (%v)", expectedOutput, hex.EncodeToString(output[0]))
	}
}
