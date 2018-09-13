package compiler

import (
	"encoding/hex"
	"strings"
	"testing"

	chainjson "github.com/bytom/encoding/json"

	"github.com/equity/compiler/equitytest"
)

func TestCompile(t *testing.T) {
	cases := []struct {
		name     string
		contract string
		want     string
	}{
		{
			"TrivialLock",
			equitytest.TrivialLock,
			"51",
		},
		{
			"LockWithPublicKey",
			equitytest.LockWithPublicKey,
			"ae7cac",
		},
		{
			"LockWithPublicKeyHash",
			equitytest.LockWithPKHash,
			"5279aa887cae7cac",
		},
		{
			"LockWith2of3Keys",
			equitytest.LockWith2of3Keys,
			"537a547a526bae71557a536c7cad",
		},
		{
			"LockToOutput",
			equitytest.LockToOutput,
			"00c3c251547ac1",
		},
		{
			"TradeOffer",
			equitytest.TradeOffer,
			"547a6413000000007b7b51547ac1631a000000547a547aae7cac",
		},
		{
			"EscrowedTransfer",
			equitytest.EscrowedTransfer,
			"537a641a000000537a7cae7cac6900c3c251557ac16328000000537a7cae7cac6900c3c251547ac1",
		},
		{
			"RevealPreimage",
			equitytest.RevealPreimage,
			"7caa87",
		},
		{
			"PriceChanger",
			equitytest.PriceChanger,
			"557a6432000000557a5479ae7cac6900c3c25100597a89587a89587a89587a89557a890274787e008901c07ec1633a000000007b537a51567ac1",
		},
		{
			"TestDefineVar",
			equitytest.TestDefineVar,
			"52797b937b7887916987",
		},
		{
			"TestSigIf",
			equitytest.TestSigIf,
			"53797b879169765379a00087641c00000052795279a0696321000000765279a069",
		},
		{
			"TestIfAndMultiClause",
			equitytest.TestIfAndMultiClause,
			"7b641f0000007087916976547aa00087641a000000765379a06963240000007b7bae7cac",
		},
		{
			"TestIfNesting",
			equitytest.TestIfNesting,
			"7b644400000054795279879169765579a00087643500000052795479a000876429000000765379a06952795579879169633a000000765479a06953797b8791635c0000007654798791695279a000876459000000527978a0697d8791",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := strings.NewReader(c.contract)
			compiled, err := Compile(r)
			if err != nil {
				t.Fatal(err)
			}

			contract := compiled[len(compiled)-1]
			got := []byte(contract.Body)

			want, err := hex.DecodeString(c.want)
			if err != nil {
				t.Fatal(err)
			}

			if string(got) != string(want) {
				t.Errorf("%s got  %s\nwant %s", c.name, hex.EncodeToString(got), hex.EncodeToString(want))
			}
		})
	}
}

func mustDecodeHex(h string) *chainjson.HexBytes {
	bits, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	result := chainjson.HexBytes(bits)
	return &result
}
