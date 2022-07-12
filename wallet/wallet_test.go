package wallet

import (
	"testing"
	"encoding/hex"
	"errors"
	"reflect"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"github.com/theQRL/zond/misc"
	"github.com/theQRL/zond/protos"
)


func TestWalletCreation(t *testing.T) {
	outFileName := "wallet_new.txt"
	wallet := NewWallet(outFileName)
	wallet.Add(4, 0)

	wallet.pbData = &protos.Wallet{}
	if !misc.FileExists(outFileName) {
		t.Error("Unable to create private key file")
	}
	data, err := ioutil.ReadFile(outFileName)
	if err != nil {
		t.Error("Unable to read private key file")
	}
	err = protojson.Unmarshal(data, wallet.pbData)
	if err != nil {
		t.Error("Unable to decode private key file", err)
	}

	expectedAddress := "000200c12f67d1941f95b23825fcc88c7be6e9a56f1953b2b50daec862896f45d90dd058a2ef79"
	expectedHexSeed := "00020052fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939cb"

	if(wallet.pbData.XmssInfo[0].Address != expectedAddress) {
		t.Fatalf("Address Mismatch\nExpected: %s\nFound: %s", expectedAddress, wallet.pbData.XmssInfo[0].Address)
	}

	if(wallet.pbData.XmssInfo[0].HexSeed != expectedHexSeed) {
		t.Fatalf("HexSeed Mismatch\nExpected: %s\nFound: %s", expectedHexSeed, wallet.pbData.XmssInfo[0].HexSeed)
	}
}


func TestWalletBalance(t *testing.T) {
	subtests := []struct {
		name          string
		input         string
		w             *Wallet
		expectedData  uint64
		expectedError error
	}{
		{
			name:  "Correct address format",
			input: "000200c12f67d1941f95b23825fcc88c7be6e9a56f1953b2b50daec862896f45d90dd058a2ef79",
			w: &Wallet{
				outFileName: "outfile",
				pbData: &protos.Wallet{
					Version: "1",
					XmssInfo: []*protos.XMSSInfo{&protos.XMSSInfo{
						Address:  "000200c12f67d1941f95b23825fcc88c7be6e9a56f1953b2b50daec862896f45d90dd058a2ef79",
						HexSeed:  "00020052fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e00167939cb",
						Mnemonic: "aback bunny flood taxi also blew harp vein dough vase onset cache isle found shrimp sparse affect sham leak ramp amidst bronze shirt nylon maid pushy heil span many brute broke accuse laity over",
					}},
				},
			},
			expectedData:  0,
			expectedError: nil,
		},
		{
			name:  "Incorrect address format",
			input: "000200c12f67d1941f95b23825fcc88c7be6e9a56f1953b2b50daec862896f45d90dd058a2ef79",
			w: &Wallet{
				outFileName: "outfile",
				pbData: &protos.Wallet{
					Version: "1",
					XmssInfo: []*protos.XMSSInfo{
						&protos.XMSSInfo{
							Address:  "000200c12f67d1941f95b23825fcc88c7be6e9a56f1953b2b50daec862896f45d90dd058a2ef79",
							HexSeed:  "00020052fdfc072182654f163f5f0f9a621d729566c74d10037c4d7bbb0407d1e2c64981855ad8681d0d86d1e91e0016793",
							Mnemonic: "aback bunny flood taxi also blew harp vein dough vase onset cache isle found shrimp sparse affect sham leak ramp amidst bronze shirt nylon maid pushy heil span many brute broke accuse laity over",
						},
					},
				},
			},
			expectedData:  0,
			expectedError: errors.New("Error Decoding address 000200c12f67d1941f95b23825fcc88c7be6e9a56f1953b2b50daec862896f45d90dd058a2ef79\n encoding/hex: odd length hex string"),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			data, err := subtest.w.reqBalance(subtest.input)
			if !reflect.DeepEqual(data, subtest.expectedData) {
				t.Errorf("expected (%b), got (%b)", subtest.expectedData, data)
			}
			if !errors.Is(err, subtest.expectedError) {
				t.Errorf("expected error (%v), got error (%v)", subtest.expectedError, err)
			}
		})
	}
}


func TestXmssByIndex(t *testing.T) {
	outFileName := "wallet_new.txt"
	wallet := NewWallet(outFileName)

	xmss, _ := wallet.GetXMSSByIndex(1)
	if reflect.TypeOf(xmss).String() != "*xmss.XMSS" {
		t.Errorf("expected xmss of type *xmss.XMSS, got (%T)", xmss)
	}

	address := xmss.GetAddress()
	decoded := hex.EncodeToString(address[:])
	if decoded!= wallet.pbData.XmssInfo[0].Address {
		t.Errorf("expected address (%s), got (%s)", wallet.pbData.XmssInfo[0].Address, decoded)
	}
}