package compiler

import (
	"encoding/hex"
	"strings"
	"testing"
)

const TrivialLock = `
contract TrivialLock() locks amount of asset {
  clause trivialUnlock() {
    unlock amount of asset
  }
}
`

const LockWithPublicKey = `
contract LockWithPublicKey(publicKey: PublicKey) locks amount of asset {
  clause unlockWithSig(sig: Signature) {
    verify checkTxSig(publicKey, sig)
    unlock amount of asset
  }
}
`

const LockWithPKHash = `
contract LockWithPublicKeyHash(pubKeyHash: Hash) locks amount of asset {
  clause spend(pubKey: PublicKey, sig: Signature) {
    verify sha3(pubKey) == pubKeyHash
    verify checkTxSig(pubKey, sig)
    unlock amount of asset
  }
}
`

const LockWith2of3Keys = `
contract LockWith3Keys(pubkey1, pubkey2, pubkey3: PublicKey) locks amount of asset {
  clause unlockWith2Sigs(sig1, sig2: Signature) {
    verify checkTxMultiSig([pubkey1, pubkey2, pubkey3], [sig1, sig2])
    unlock amount of asset
  }
}
`

const LockToOutput = `
contract LockToOutput(address: Program) locks amount of asset {
  clause relock() {
    lock amount of asset with address
  }
}
`

const TradeOffer = `
contract TradeOffer(requestedAsset: Asset, requestedAmount: Amount, sellerProgram: Program, sellerKey: PublicKey) locks amount of asset {
  clause trade() {
    lock requestedAmount of requestedAsset with sellerProgram
    unlock amount of asset
  }
  clause cancel(sellerSig: Signature) {
    verify checkTxSig(sellerKey, sellerSig)
    unlock amount of asset
  }
}
`

const EscrowedTransfer = `
contract EscrowedTransfer(agent: PublicKey, sender: Program, recipient: Program) locks amount of asset {
  clause approve(sig: Signature) {
    verify checkTxSig(agent, sig)
    lock amount of asset with recipient
  }
  clause reject(sig: Signature) {
    verify checkTxSig(agent, sig)
    lock amount of asset with sender
  }
}
`

const CollateralizedLoan = `
contract CollateralizedLoan(balanceAsset: Asset, balanceAmount: Amount, finalHeight: Integer, lender: Program, borrower: Program) locks valueAmount of valueAsset {
  clause repay() {
    lock balanceAmount of balanceAsset with lender
    lock valueAmount of valueAsset with borrower
  }
  clause default() {
    verify above(finalHeight)
    lock valueAmount of valueAsset with lender
  }
}
`

const RevealPreimage = `
contract RevealPreimage(hash: Hash) locks amount of asset {
  clause reveal(string: String) {
    verify sha3(string) == hash
    unlock amount of asset
  }
}
`

const PriceChanger = `
contract PriceChanger(askAmount: Amount, askAsset: Asset, sellerKey: PublicKey, sellerProg: Program) locks valueAmount of valueAsset {
  clause changePrice(newAmount: Amount, newAsset: Asset, sig: Signature) {
    verify checkTxSig(sellerKey, sig)
    lock valueAmount of valueAsset with PriceChanger(newAmount, newAsset, sellerKey, sellerProg)
  }
  clause redeem() {
    lock askAmount of askAsset with sellerProg
    unlock valueAmount of valueAsset
  }
}
`

const CallOptionWithSettlement = `
contract CallOptionWithSettlement(strikePrice: Amount,
                    strikeCurrency: Asset,
                    sellerProgram: Program,
                    sellerKey: PublicKey,
                    buyerKey: PublicKey,
                    finalHeight: Integer) locks valueAmount of valueAsset {
  clause exercise(buyerSig: Signature) {
    verify below(finalHeight)
    verify checkTxSig(buyerKey, buyerSig)
    lock strikePrice of strikeCurrency with sellerProgram
    unlock valueAmount of valueAsset
  }
  clause expire() {
    verify above(finalHeight)
    lock valueAmount of valueAsset with sellerProgram
  }
  clause settle(sellerSig: Signature, buyerSig: Signature) {
    verify checkTxSig(sellerKey, sellerSig)
    verify checkTxSig(buyerKey, buyerSig)
    unlock valueAmount of valueAsset
  }
}
`

const TestDefineVar = `
contract TestDefineVar(result: Integer) locks valueAmount of valueAsset {
  clause LockWithMath(left: Integer, right: Integer) {
    define calculate: Integer = left + right
    verify left != calculate
    verify result == calculate
    unlock valueAmount of valueAsset
  }
}
`

const TestAssignVar = `
contract TestAssignVar(result: Integer) locks valueAmount of valueAsset {
  clause LockWithMath(first: Integer, second: Integer) {
    define calculate: Integer = first
    assign calculate = calculate + second
    verify result == calculate
    unlock valueAmount of valueAsset
  }
}
`

const TestSigIf = `
contract TestSigIf(a: Integer, count:Integer) locks valueAmount of valueAsset {
  clause check(b: Integer, c: Integer) {
    verify b != count
    if a > b {
        verify b > c
    } else {
        verify a > c
    }
    unlock valueAmount of valueAsset
  }
}
`

const TestIfAndMultiClause = `
contract TestIfAndMultiClause(a: Integer, cancelKey: PublicKey) locks valueAmount of valueAsset {
  clause check(b: Integer, c: Integer) {
    verify b != c
    if a > b {
        verify a > c
    }
    unlock valueAmount of valueAsset
  }
  clause cancel(sellerSig: Signature) {
    verify checkTxSig(cancelKey, sellerSig)
    unlock valueAmount of valueAsset
  }
}
`

const TestIfNesting = `
contract TestIfNesting(a: Integer, count:Integer) locks valueAmount of valueAsset {
  clause check(b: Integer, c: Integer, d: Integer) {
    verify b != count
    if a > b {
        if d > c {
           verify a > d
        }
        verify d != b
    } else {
        verify a > c
    }
    verify c != count
    unlock valueAmount of valueAsset
  }
  clause cancel(e: Integer, f: Integer) {
    verify a != e
    if a > f {
      verify e > count
    }
    verify f != count
    unlock valueAmount of valueAsset
  }
}
`

const TestConstantMath = `
contract TestConstantMath(result: Integer, hashByte: Hash, hashStr: Hash, outcome: Boolean) locks valueAmount of valueAsset {
  clause calculation(left: Integer, right: Integer, boolResult: Boolean) {
    verify result == left + right + 10
    verify hashByte == sha3(0x31323330)
    verify hashStr == sha3('string')
    verify !outcome
    verify boolResult && (result == left + 20)
    unlock valueAmount of valueAsset
  }
}
`

const VerifySignature = `
contract VerifySignature(sig1: Sign, sig2: Sign, msgHash: Hash) locks valueAmount of valueAsset {
  clause check(publicKey1: PublicKey, publicKey2: PublicKey) {
    verify checkMsgSig(publicKey1, msgHash, sig1)
    verify checkMsgSig(publicKey2, msgHash, sig2)
    unlock valueAmount of valueAsset
  }
}
`

func TestCompile(t *testing.T) {
	cases := []struct {
		name     string
		contract string
		want     string
	}{
		{
			"TrivialLock",
			TrivialLock,
			"51",
		},
		{
			"LockWithPublicKey",
			LockWithPublicKey,
			"ae7cac",
		},
		{
			"LockWithPublicKeyHash",
			LockWithPKHash,
			"5279aa887cae7cac",
		},
		{
			"LockWith2of3Keys",
			LockWith2of3Keys,
			"537a547a526bae71557a536c7cad",
		},
		{
			"LockToOutput",
			LockToOutput,
			"00c3c251547ac1",
		},
		{
			"TradeOffer",
			TradeOffer,
			"547a6413000000007b7b51547ac1631a000000547a547aae7cac",
		},
		{
			"EscrowedTransfer",
			EscrowedTransfer,
			"537a641a000000537a7cae7cac6900c3c251557ac16328000000537a7cae7cac6900c3c251547ac1",
		},
		{
			"CollateralizedLoan",
			CollateralizedLoan,
			"557a641b000000007b7b51557ac16951c3c251557ac163260000007bcd9f6900c3c251567ac1",
		},
		{
			"RevealPreimage",
			RevealPreimage,
			"7caa87",
		},
		{
			"PriceChanger",
			PriceChanger,
			"557a6432000000557a5479ae7cac6900c3c25100597a89587a89587a89587a89557a890274787e008901c07ec1633a000000007b537a51567ac1",
		},
		{
			"CallOptionWithSettlement",
			CallOptionWithSettlement,
			"567a76529c64360000006425000000557acda06971ae7cac69007c7b51547ac16346000000557acd9f6900c3c251567ac1634600000075577a547aae7cac69557a547aae7cac",
		},
		{
			"TestDefineVar",
			TestDefineVar,
			"52797b937b7887916987",
		},
		{
			"TestAssignVar",
			TestAssignVar,
			"7b7b9387",
		},
		{
			"TestSigIf",
			TestSigIf,
			"53797b879169765379a09161641c00000052795279a0696321000000765279a069",
		},
		{
			"TestIfAndMultiClause",
			TestIfAndMultiClause,
			"7b641f0000007087916976547aa09161641a000000765379a06963240000007b7bae7cac",
		},
		{
			"TestIfNesting",
			TestIfNesting,
			"7b644400000054795279879169765579a09161643500000052795479a091616429000000765379a06952795579879169633a000000765479a06953797b8791635c0000007654798791695279a091616459000000527978a0697d8791",
		},
		{
			"TestConstantMath",
			TestConstantMath,
			"765779577a935a93887c0431323330aa887c06737472696e67aa887c91697b011493879a",
		},
		{
			"VerifySignature",
			VerifySignature,
			"5279557aac697c7bac",
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
