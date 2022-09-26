package wallet

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/theQRL/zond/misc"
	"github.com/theQRL/zond/protos"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestWalletXmssCreation(t *testing.T) {
	outFileName := "wallet_new.txt"
	defer os.Remove(outFileName)
	wallet := NewWallet(outFileName)
	wallet.AddXMSS(4, 0)

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
}

func TestWalletDilithiumCreation(t *testing.T) {
	outFileName := "wallet_new.txt"
	defer os.Remove(outFileName)
	wallet := NewWallet(outFileName)
	wallet.AddDilithium()

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
}

// func TestWalletBalance(t *testing.T) {
// 	subtests := []struct {
// 		name          string
// 		input         string
// 		w             *Wallet
// 		expectedData  uint64
// 		expectedError error
// 	}{
// 		{
// 			name:  "Correct address format",
// 			input: "0x100200cd5efc5ab7f490a26e84adc8e97a592b1c",
// 			w: &Wallet{
// 				outFileName: "outfile",
// 				pbData: &protos.Wallet{
// 					Version: "1",
// 					Info: []*protos.Info{{
// 						Address:  "0x100200cd5efc5ab7f490a26e84adc8e97a592b1c",
// 						HexSeed:  "0x10020012af6fcebaec02890830fb753fbf382c08b2304faa7c71685328c2c7f3dbca8ec081db194ae770a2fc7a9380480d1230",
// 						Mnemonic: "badge bunny bass watery soil quite added motive couch rice draft villa cloak meter corps windy leg bicker flora milan siege divide slabby tunis loop reach expert khaki pick shut nation agree atlas campus",
// 						Type:     uint32(common2.XMSSSig),
// 					}},
// 				},
// 			},
// 			expectedData:  0,
// 			expectedError: nil,
// 		},
// 		{
// 			name:  "Incorrect address format",
// 			input: "0x10000cd5efc5ab7f490a26e84adc8e97a592b1c",
// 			w: &Wallet{
// 				outFileName: "outfile",
// 				pbData: &protos.Wallet{
// 					Version: "1",
// 					Info: []*protos.Info{
// 						{
// 							Address:  "0x100200cd5efc5ab7f490a26e84adc8e97a592b1c",
// 							HexSeed:  "0x10020012af6fcebaec02890830fb753fbf382c08b2304faa7c71685328c2c7f3dbca8ec081db194ae770a2fc7a9380480d1230",
// 							Mnemonic: "badge bunny bass watery soil quite added motive couch rice draft villa cloak meter corps windy leg bicker flora milan siege divide slabby tunis loop reach expert khaki pick shut nation agree atlas campus",
// 							Type:     uint32(common2.XMSSSig),
// 						},
// 					},
// 				},
// 			},
// 			expectedData:  0,
// 			expectedError: errors.New("Error Decoding address 0x100200cd5efc5ab7f490a26e84adc8e97a592b1c\n encoding/hex: odd length hex string"),
// 		},
// 	}

// 	for _, subtest := range subtests {
// 		t.Run(subtest.name, func(t *testing.T) {
// 			data, err := subtest.w.reqBalance(subtest.input)
// 			if !reflect.DeepEqual(data, subtest.expectedData) {
// 				t.Errorf("expected (%b), got (%b)", subtest.expectedData, data)
// 			}
// 			if !errors.Is(err, subtest.expectedError) {
// 				t.Errorf("expected error (%v), got error (%v)", subtest.expectedError, err)
// 			}
// 		})
// 	}
// }

func TestXmssAccountByIndex(t *testing.T) {
	outFileName := "wallet_new.txt"
	defer os.Remove(outFileName)
	wallet := NewWallet(outFileName)
	wallet.AddXMSS(4, 0)

	xmss, _ := wallet.GetXMSSAccountByIndex(1)
	if reflect.TypeOf(xmss).String() != "*xmss.XMSS" {
		t.Errorf("expected xmss of type *xmss.XMSS, got (%T)", xmss)
	}

	address := xmss.GetAddress()
	address2 := misc.BytesToHexStr([]byte(address[:]))

	if len(wallet.pbData.Info[0].Address) != len(address2) {
		t.Errorf("expected address length (%d), got (%d)", len(wallet.pbData.Info[0].Address), len(address2))
	}
}

func TestDilithiumAccountByIndex(t *testing.T) {
	outFileName := "wallet_new.txt"
	defer os.Remove(outFileName)
	wallet := NewWallet(outFileName)
	wallet.AddDilithium()

	dilithium, _ := wallet.GetDilithiumAccountByIndex(1)
	if reflect.TypeOf(dilithium).String() != "*dilithium.Dilithium" {
		t.Errorf("expected xmss of type *xmss.XMSS, got (%T)", dilithium)
	}

	address := dilithium.GetAddress()
	address2 := misc.BytesToHexStr([]byte(address[:]))

	if len(wallet.pbData.Info[0].Address) != len(address2) {
		t.Errorf("expected address length (%d), got (%d)", len(wallet.pbData.Info[0].Address), len(address2))
	}
}
