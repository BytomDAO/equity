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
			"TestAssignVar",
			equitytest.TestAssignVar,
			"7b7b9387",
		},
		{
			"TestSigIf",
			equitytest.TestSigIf,
			"53797b879169765379a091641b00000052795279a0696320000000765279a069",
		},
		{
			"TestIfAndMultiClause",
			equitytest.TestIfAndMultiClause,
			"7b641e0000007087916976547aa0916419000000765379a06963230000007b7bae7cac",
		},
		{
			"TestIfNesting",
			equitytest.TestIfNesting,
			"7b644200000054795279879169765579a091643300000052795479a0916427000000765379a069527955798791696338000000765479a06953797b879163590000007654798791695279a0916456000000527978a0697d8791",
		},
		{
			"TestConstantMath",
			equitytest.TestConstantMath,
			"765779577a935a93887c0431323330aa887c06737472696e67aa887c91697b011493879a",
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
