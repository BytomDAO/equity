package compiler

import (
	"encoding/hex"
	"strings"
	"testing"

	chainjson "github.com/bytom/encoding/json"

	"github.com/equity/compiler/equitytest"
)

func TestCompile(t *testing.T) {
	int64Pointer := func(x int64) *int64 { return &x }

	cases := []struct {
		name     string
		contract string
		args     []ContractArg
		want     string
	}{
		{
			"TrivialLock",
			equitytest.TrivialLock,
			[]ContractArg{},
			"74015100c0",
		},
		{
			"LockWithPublicKey",
			equitytest.LockWithPublicKey,
			[]ContractArg{
				ContractArg{S: mustDecodeHex("cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f013")},
			},
			"20cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f0137403ae7cac00c0",
		},
		{
			"LockWithPublicKeyHash",
			equitytest.LockWithPKHash,
			[]ContractArg{
				ContractArg{S: mustDecodeHex("5a6e7792029f8e84dd1260f694f2131e00fd1810d44696b56af062e91e667fc0")},
			},
			"205a6e7792029f8e84dd1260f694f2131e00fd1810d44696b56af062e91e667fc074085279aa887cae7cac00c0",
		},
		{
			"LockWith2of3Keys",
			equitytest.LockWith2of3Keys,
			[]ContractArg{
				ContractArg{S: mustDecodeHex("cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f011")},
				ContractArg{S: mustDecodeHex("cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f012")},
				ContractArg{S: mustDecodeHex("cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f013")},
			},
			"20cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f01320cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f01220cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f011740e537a547a526bae71557a536c7cad00c0",
		},
		{
			"LockToOutput",
			equitytest.LockToOutput,
			[]ContractArg{
				ContractArg{S: mustDecodeHex("00206c9d25c8aa63eba373a3a10c39cd4e884fd99480241e50ecef9d6efddc424382")},
			},
			"2200206c9d25c8aa63eba373a3a10c39cd4e884fd99480241e50ecef9d6efddc424382740700c3c251547ac100c0",
		},
		{
			"TradeOffer",
			equitytest.TradeOffer,
			[]ContractArg{
				ContractArg{S: mustDecodeHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
				ContractArg{I: int64Pointer(5)},
				ContractArg{S: mustDecodeHex("00208c50e01321a52d273859e6d01a65c1456d9892b188d9142152dee28593b2343b")},
				ContractArg{S: mustDecodeHex("cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f013")},
			},
			"20cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f0132200208c50e01321a52d273859e6d01a65c1456d9892b188d9142152dee28593b2343b5520ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff741a547a6413000000007b7b51547ac1631a000000547a547aae7cac00c0",
		},
		{
			"EscrowedTransfer",
			equitytest.EscrowedTransfer,
			[]ContractArg{
				ContractArg{S: mustDecodeHex("cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f013")},
				ContractArg{S: mustDecodeHex("00204e478aae9f8167d68732cb65b225837d4d1a4f4ca82ced9c931fe83f02361b7f")},
				ContractArg{S: mustDecodeHex("0020743ee85f3779bf77e4ecb2107ad43308fdf467695b00050366af97fd7e1902fe")},
			},
			"220020743ee85f3779bf77e4ecb2107ad43308fdf467695b00050366af97fd7e1902fe2200204e478aae9f8167d68732cb65b225837d4d1a4f4ca82ced9c931fe83f02361b7f20cc8cefdf1a433cd6c43718c7a2c1324aa13cb7e5b0875d6f727248a55bf7f0137428537a641a000000537a7cae7cac6900c3c251557ac16328000000537a7cae7cac6900c3c251547ac100c0",
		},
		{
			"RevealPreimage",
			equitytest.RevealPreimage,
			[]ContractArg{
				ContractArg{S: mustDecodeHex("a03ab19b866fc585b5cb1812a2f63ca861e7e7643ee5d43fd7106b623725fd67")},
			},
			"20a03ab19b866fc585b5cb1812a2f63ca861e7e7643ee5d43fd7106b623725fd6774037caa8700c0",
		},
		/*
			{
				"CollateralizedLoan",
				equitytest.CollateralizedLoan,
				`[{"name":"CollateralizedLoan","params":[{"name":"balanceAsset","type":"Asset"},{"name":"balanceAmount","type":"Amount"},{"name":"finalHeight","type":"Integer"},{"name":"lender","type":"Program"},{"name":"borrower","type":"Program"}],"clauses":[{"name":"repay","reqs":[{"name":"payment","asset":"balanceAsset","amount":"balanceAmount"}],"values":[{"name":"payment","program":"lender","asset":"balanceAsset","amount":"balanceAmount"},{"name":"collateral","program":"borrower"}]},{"name":"default","blockheight":["finalHeight"],"values":[{"name":"collateral","program":"lender"}]}],"value":"collateral","body_bytecode":"557a641b000000007b7b51557ac16951c3c251557ac163260000007bcd9f6900c3c251567ac1","body_opcodes":"5 ROLL JUMPIF:$default $repay 0 ROT ROT 1 5 ROLL CHECKOUTPUT VERIFY 1 AMOUNT ASSET 1 5 ROLL CHECKOUTPUT JUMP:$_end $default ROT BLOCKHEIGHT LESSTHAN VERIFY 0 AMOUNT ASSET 1 6 ROLL CHECKOUTPUT $_end","recursive":false}]`,
			},
			{
				"CallOptionWithSettlement",
				equitytest.CallOptionWithSettlement,
				`[{"name":"CallOptionWithSettlement","params":[{"name":"strikePrice","type":"Amount"},{"name":"strikeCurrency","type":"Asset"},{"name":"sellerProgram","type":"Program"},{"name":"sellerKey","type":"PublicKey"},{"name":"buyerKey","type":"PublicKey"},{"name":"finalHeight","type":"Integer"}],"clauses":[{"name":"exercise","params":[{"name":"buyerSig","type":"Signature"}],"reqs":[{"name":"payment","asset":"strikeCurrency","amount":"strikePrice"}],"blockheight":["finalHeight"],"values":[{"name":"payment","program":"sellerProgram","asset":"strikeCurrency","amount":"strikePrice"},{"name":"underlying"}]},{"name":"expire","blockheight":["finalHeight"],"values":[{"name":"underlying","program":"sellerProgram"}]},{"name":"settle","params":[{"name":"sellerSig","type":"Signature"},{"name":"buyerSig","type":"Signature"}],"values":[{"name":"underlying"}]}],"value":"underlying","body_bytecode":"567a76529c64360000006425000000557acda06971ae7cac69007c7b51547ac16346000000557acd9f6900c3c251567ac1634600000075577a547aae7cac69557a547aae7cac","body_opcodes":"6 ROLL DUP 2 NUMEQUAL JUMPIF:$settle JUMPIF:$expire $exercise 5 ROLL BLOCKHEIGHT GREATERTHAN VERIFY 2ROT TXSIGHASH SWAP CHECKSIG VERIFY 0 SWAP ROT 1 4 ROLL CHECKOUTPUT JUMP:$_end $expire 5 ROLL BLOCKHEIGHT LESSTHAN VERIFY 0 AMOUNT ASSET 1 6 ROLL CHECKOUTPUT JUMP:$_end $settle DROP 7 ROLL 4 ROLL TXSIGHASH SWAP CHECKSIG VERIFY 5 ROLL 4 ROLL TXSIGHASH SWAP CHECKSIG $_end","recursive":false}]`,
			},
			{
				"PriceChanger",
				equitytest.PriceChanger,
				`[{"name":"PriceChanger","params":[{"name":"askAmount","type":"Amount"},{"name":"askAsset","type":"Asset"},{"name":"sellerKey","type":"PublicKey"},{"name":"sellerProg","type":"Program"}],"clauses":[{"name":"changePrice","params":[{"name":"newAmount","type":"Amount"},{"name":"newAsset","type":"Asset"},{"name":"sig","type":"Signature"}],"values":[{"name":"offered","program":"PriceChanger(newAmount, newAsset, sellerKey, sellerProg)"}],"contracts":["PriceChanger"]},{"name":"redeem","reqs":[{"name":"payment","asset":"askAsset","amount":"askAmount"}],"values":[{"name":"payment","program":"sellerProg","asset":"askAsset","amount":"askAmount"},{"name":"offered"}]}],"value":"offered","body_bytecode":"557a6432000000557a5479ae7cac6900c3c25100597a89587a89587a89587a89557a890274787e008901c07ec1633a000000007b537a51567ac1","body_opcodes":"5 ROLL JUMPIF:$redeem $changePrice 5 ROLL 4 PICK TXSIGHASH SWAP CHECKSIG VERIFY 0 AMOUNT ASSET 1 0 9 ROLL CATPUSHDATA 8 ROLL CATPUSHDATA 8 ROLL CATPUSHDATA 8 ROLL CATPUSHDATA 5 ROLL CATPUSHDATA 0x7478 CAT 0 CATPUSHDATA 192 CAT CHECKOUTPUT JUMP:$_end $redeem 0 ROT 3 ROLL 1 6 ROLL CHECKOUTPUT $_end","recursive":true}]`,
			},
			{
				"OneTwo",
				equitytest.OneTwo,
				`[{"name":"Two","params":[{"name":"b","type":"Program"},{"name":"c","type":"Program"},{"name":"expirationHeight","type":"Integer"}],"clauses":[{"name":"redeem","blockheight":["expirationHeight"],"values":[{"name":"value","program":"b"}]},{"name":"default","blockheight":["expirationHeight"],"values":[{"name":"value","program":"c"}]}],"value":"value","body_bytecode":"537a64170000007bcda06900c3c251547ac163220000007bcd9f6900c3c251557ac1","body_opcodes":"3 ROLL JUMPIF:$default $redeem ROT BLOCKHEIGHT GREATERTHAN VERIFY 0 AMOUNT ASSET 1 4 ROLL CHECKOUTPUT JUMP:$_end $default ROT BLOCKHEIGHT LESSTHAN VERIFY 0 AMOUNT ASSET 1 5 ROLL CHECKOUTPUT $_end","recursive":false},{"name":"One","params":[{"name":"a","type":"Program"},{"name":"b","type":"Program"},{"name":"c","type":"Program"},{"name":"switchHeight","type":"Integer"},{"name":"blockHeight","type":"Integer"}],"clauses":[{"name":"redeem","blockheight":["switchHeight"],"values":[{"name":"value","program":"a"}]},{"name":"switch","blockheight":["switchHeight"],"values":[{"name":"value","program":"Two(b, c, blockHeight)"}],"contracts":["Two"]}],"value":"value","body_bytecode":"557a6418000000537acda06900c3c251547ac16358000000537acd9f6900c3c25100587a89577a89567a8901747e22537a64170000007bcda06900c3c251547ac163220000007bcd9f6900c3c251557ac189008901c07ec1","body_opcodes":"5 ROLL JUMPIF:$switch $redeem 3 ROLL BLOCKHEIGHT GREATERTHAN VERIFY 0 AMOUNT ASSET 1 4 ROLL CHECKOUTPUT JUMP:$_end $switch 3 ROLL BLOCKHEIGHT LESSTHAN VERIFY 0 AMOUNT ASSET 1 0 8 ROLL CATPUSHDATA 7 ROLL CATPUSHDATA 6 ROLL CATPUSHDATA 116 CAT 0x537a64170000007bcda06900c3c251547ac163220000007bcd9f6900c3c251557ac1 CATPUSHDATA 0 CATPUSHDATA 192 CAT CHECKOUTPUT $_end","recursive":false}]`,
			},*/
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := strings.NewReader(c.contract)
			compiled, err := Compile(r)
			if err != nil {
				t.Fatal(err)
			}

			contract := compiled[len(compiled)-1]
			got, err := Instantiate(contract.Body, contract.Params, false, c.args)
			if err != nil {
				t.Fatal(err)
			}

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
