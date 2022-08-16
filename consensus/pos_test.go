package consensus

import (
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/theQRL/zond/chain"
	"github.com/theQRL/zond/config"
	mockdb "github.com/theQRL/zond/db/mock"
	"github.com/theQRL/zond/p2p"
	"github.com/theQRL/zond/state"
)

func TestNewPOS(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mockdb.NewMockDB(ctrl)

	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	chain := chain.NewChain(state)

	server, err := p2p.NewServer(chain)
	if err != nil {
		t.Error("got unexpected error while creating new server instance")
	}

	_ = NewPOS(server, chain, store)
}

func TestGetCurrentSlot(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mockdb.NewMockDB(ctrl)

	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	chain := chain.NewChain(state)

	server, err := p2p.NewServer(chain)
	if err != nil {
		t.Error("got unexpected error while creating new server instance")
	}

	pos := NewPOS(server, chain, store)

	currentTime := uint64(time.Now().Unix())
	config := config.GetConfig()
	genesisTimestamp := config.Dev.Genesis.GenesisTimestamp
	blockTiming := config.Dev.BlockTime

	expectedOutput := (currentTime - genesisTimestamp) / blockTiming

	output := pos.GetCurrentSlot()

	if expectedOutput != output {
		t.Errorf("expected output (%v), got (%v)", expectedOutput, output)
	}
}

func TestTimeRemainingForNextAction(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mockdb.NewMockDB(ctrl)

	state, _ := state.NewState("./", "testStateDb.txt")
	defer os.Remove("testStateDb.txt")
	chain := chain.NewChain(state)

	server, err := p2p.NewServer(chain)
	if err != nil {
		t.Error("got unexpected error while creating new server instance")
	}

	pos := NewPOS(server, chain, store)
	currentTime := uint64(time.Now().Unix())
	config := config.GetConfig()
	genesisTimestamp := config.Dev.Genesis.GenesisTimestamp
	blockTiming := config.Dev.BlockTime
	currentSlot := (currentTime - genesisTimestamp) / blockTiming
	nextSlotTime := genesisTimestamp + (currentSlot+1)*blockTiming

	timeRemainingForNextSlot := nextSlotTime - currentTime

	time_ := pos.TimeRemainingForNextAction()
	expectedTime := time.Duration(timeRemainingForNextSlot) * time.Second

	if time_.String() != expectedTime.String() {
		t.Errorf("expected time remaining for next action to be (%v), got (%v)", expectedTime.String(), time_.String())
	}
}
