package vm

import (
	"fmt"
)

const (
	//base runlimit gas
	GasBaseStep     int64 = 1
	GasLowSlowStep  int64 = 2
	GasSlowStep     int64 = 10
	GasHighSlowStep int64 = 17
	GasLowMidStep   int64 = 32
	GasMidStep      int64 = 33
	GasHighMidStep  int64 = 41
	GasLowFastStep  int64 = 296
	GasFastStep     int64 = 881

	GasData      int64 = 9
	GasPushData1 int64 = 138
	GasPushData2 int64 = 32779
	GasPushData4 int64 = 2147483661

	GasUnderBaseStep     int64 = -4
	GasUnderLowSlowStep  int64 = -5
	GasUnderSlowStep     int64 = -7
	GasUnderHighSlowStep int64 = -8
	GasUnderLowMidStep   int64 = -14
	GasUnderMidStep      int64 = -15
	GasUnderHighMidStep  int64 = -21
	GasUnderLowFastStep  int64 = -30
	GasUnderFastStep     int64 = -37

	//mutable factor
	ConfigLength   int64 = 64
	ConfigDistance int64 = 8
	ConfigNumPub   int64 = 3
	ConfigNumSig   int64 = 2

	//mutable runlimit gas
	GasMutBaseStep     int64 = 1 + ConfigLength
	GasMutLowSlowStep  int64 = 9 + ConfigLength
	GasMutSlowStep     int64 = -6 + ConfigLength
	GasMutHighSlowStep int64 = -6 - ConfigLength
	GasMutLowMidStep   int64 = -7 - ConfigLength
	GasMutMidStep      int64 = -71 - ConfigLength
	GasMutHighMidStep  int64 = 18 + 2*ConfigLength
	GasMutLowFastStep  int64 = 26 + 3*ConfigLength
	GasMutFastStep     int64 = -14 - 2*ConfigLength

	GasDifferBaseStep int64 = -7 + ConfigDistance
	GasDifferSlowStep int64 = -7 - ConfigDistance
	GasDifferMidStep  int64 = -12 - ConfigDistance
	GasDifferFastStep int64 = -28 - ConfigDistance

	GasMultiSig int64 = -63 + 984*ConfigNumPub - 72*ConfigNumSig
)

type gasInfo struct {
	Op   Op
	name string
	Gas  int64
}

var (
	gasop = [256]gasInfo{
		// data pushing
		OP_FALSE: {OP_FALSE, "FALSE", GasSlowStep},
		OP_1:     {OP_1, "1", GasSlowStep},

		// sic: the PUSHDATA ops all share an implementation
		OP_PUSHDATA1: {OP_PUSHDATA1, "PUSHDATA1", GasPushData1},
		OP_PUSHDATA2: {OP_PUSHDATA2, "PUSHDATA2", GasPushData2},
		OP_PUSHDATA4: {OP_PUSHDATA4, "PUSHDATA4", GasPushData4},

		OP_1NEGATE: {OP_1NEGATE, "1NEGATE", GasHighSlowStep},

		OP_NOP: {OP_NOP, "NOP", GasBaseStep},

		// control flow
		OP_JUMP:   {OP_JUMP, "JUMP", GasBaseStep},
		OP_JUMPIF: {OP_JUMPIF, "JUMPIF", GasUnderMidStep},

		OP_VERIFY: {OP_VERIFY, "VERIFY", GasUnderHighSlowStep},
		OP_FAIL:   {OP_FAIL, "FAIL", GasBaseStep},

		OP_TOALTSTACK:   {OP_TOALTSTACK, "TOALTSTACK", GasLowSlowStep},
		OP_FROMALTSTACK: {OP_FROMALTSTACK, "FROMALTSTACK", GasLowSlowStep},
		OP_2DROP:        {OP_2DROP, "2DROP", GasMutFastStep},
		OP_2DUP:         {OP_2DUP, "2DUP", GasMutHighMidStep},
		OP_3DUP:         {OP_3DUP, "3DUP", GasMutLowFastStep},
		OP_2OVER:        {OP_2OVER, "2OVER", GasMutHighMidStep},
		OP_2ROT:         {OP_2ROT, "2ROT", GasLowSlowStep},
		OP_2SWAP:        {OP_2SWAP, "2SWAP", GasLowSlowStep},
		OP_IFDUP:        {OP_IFDUP, "IFDUP", GasBaseStep},
		OP_DEPTH:        {OP_DEPTH, "DEPTH", GasHighSlowStep},
		OP_DROP:         {OP_DROP, "DROP", GasMutLowMidStep},
		OP_DUP:          {OP_DUP, "DUP", GasMutLowSlowStep},
		OP_NIP:          {OP_NIP, "NIP", GasMutLowMidStep},
		OP_OVER:         {OP_OVER, "OVER", GasMutLowSlowStep},
		OP_PICK:         {OP_PICK, "PICK", GasMutSlowStep},
		OP_ROLL:         {OP_ROLL, "ROLL", GasUnderLowMidStep},
		OP_ROT:          {OP_ROT, "ROT", GasLowSlowStep},
		OP_SWAP:         {OP_SWAP, "SWAP", GasBaseStep},
		OP_TUCK:         {OP_TUCK, "TUCK", GasMutLowSlowStep},

		OP_CAT:         {OP_CAT, "CAT", GasUnderBaseStep},
		OP_SUBSTR:      {OP_SUBSTR, "SUBSTR", GasDifferFastStep},
		OP_LEFT:        {OP_LEFT, "LEFT", GasDifferMidStep},
		OP_RIGHT:       {OP_RIGHT, "RIGHT", GasDifferMidStep},
		OP_SIZE:        {OP_SIZE, "SIZE", GasHighSlowStep},
		OP_CATPUSHDATA: {OP_CATPUSHDATA, "CATPUSHDATA", GasUnderBaseStep},

		OP_INVERT:      {OP_INVERT, "INVERT", GasMutBaseStep},
		OP_AND:         {OP_AND, "AND", GasDifferSlowStep},
		OP_OR:          {OP_OR, "OR", GasDifferBaseStep},
		OP_XOR:         {OP_XOR, "XOR", GasDifferBaseStep},
		OP_EQUAL:       {OP_EQUAL, "EQUAL", GasMutHighSlowStep},
		OP_EQUALVERIFY: {OP_EQUALVERIFY, "EQUALVERIFY", GasMutLowMidStep},

		OP_1ADD:               {OP_1ADD, "1ADD", GasLowSlowStep},
		OP_1SUB:               {OP_1SUB, "1SUB", GasLowSlowStep},
		OP_2MUL:               {OP_2MUL, "2MUL", GasLowSlowStep},
		OP_2DIV:               {OP_2DIV, "2DIV", GasLowSlowStep},
		OP_NEGATE:             {OP_NEGATE, "NEGATE", GasLowSlowStep},
		OP_ABS:                {OP_ABS, "ABS", GasLowSlowStep},
		OP_NOT:                {OP_NOT, "NOT", GasUnderLowSlowStep},
		OP_0NOTEQUAL:          {OP_0NOTEQUAL, "0NOTEQUAL", GasUnderLowSlowStep},
		OP_ADD:                {OP_ADD, "ADD", GasUnderLowMidStep},
		OP_SUB:                {OP_SUB, "SUB", GasUnderLowMidStep},
		OP_MUL:                {OP_MUL, "MUL", GasUnderLowMidStep},
		OP_DIV:                {OP_DIV, "DIV", GasUnderLowMidStep},
		OP_MOD:                {OP_MOD, "MOD", GasUnderLowMidStep},
		OP_LSHIFT:             {OP_LSHIFT, "LSHIFT", GasUnderLowMidStep},
		OP_RSHIFT:             {OP_RSHIFT, "RSHIFT", GasUnderLowMidStep},
		OP_BOOLAND:            {OP_BOOLAND, "BOOLAND", GasUnderSlowStep},
		OP_BOOLOR:             {OP_BOOLOR, "BOOLOR", GasUnderSlowStep},
		OP_NUMEQUAL:           {OP_NUMEQUAL, "NUMEQUAL", GasUnderHighMidStep},
		OP_NUMEQUALVERIFY:     {OP_NUMEQUALVERIFY, "NUMEQUALVERIFY", GasUnderLowFastStep},
		OP_NUMNOTEQUAL:        {OP_NUMNOTEQUAL, "NUMNOTEQUAL", GasUnderHighMidStep},
		OP_LESSTHAN:           {OP_LESSTHAN, "LESSTHAN", GasUnderHighMidStep},
		OP_GREATERTHAN:        {OP_GREATERTHAN, "GREATERTHAN", GasUnderHighMidStep},
		OP_LESSTHANOREQUAL:    {OP_LESSTHANOREQUAL, "LESSTHANOREQUAL", GasUnderHighMidStep},
		OP_GREATERTHANOREQUAL: {OP_GREATERTHANOREQUAL, "GREATERTHANOREQUAL", GasUnderHighMidStep},
		OP_MIN:                {OP_MIN, "MIN", GasUnderLowMidStep},
		OP_MAX:                {OP_MAX, "MAX", GasUnderLowMidStep},
		OP_WITHIN:             {OP_WITHIN, "WITHIN", GasUnderFastStep},

		OP_SHA256:        {OP_SHA256, "SHA256", GasLowMidStep},
		OP_SHA3:          {OP_SHA3, "SHA3", GasLowMidStep},
		OP_CHECKSIG:      {OP_CHECKSIG, "CHECKSIG", GasFastStep},
		OP_CHECKMULTISIG: {OP_CHECKMULTISIG, "CHECKMULTISIG", GasMultiSig},
		OP_TXSIGHASH:     {OP_TXSIGHASH, "TXSIGHASH", GasLowFastStep},

		OP_CHECKOUTPUT: {OP_CHECKOUTPUT, "CHECKOUTPUT", GasMutMidStep},
		OP_ASSET:       {OP_ASSET, "ASSET", GasHighMidStep},
		OP_AMOUNT:      {OP_AMOUNT, "AMOUNT", GasHighSlowStep},
		OP_PROGRAM:     {OP_PROGRAM, "PROGRAM", GasMutLowSlowStep},
		//OP_MINTIME:     {OP_MINTIME, "MINTIME", GasHighSlowStep},
		//OP_MAXTIME:     {OP_MAXTIME, "MAXTIME", GasHighSlowStep},
		//OP_TXDATA:      {OP_TXDATA, "TXDATA", GasHighMidStep},
		//OP_ENTRYDATA:   {OP_ENTRYDATA, "ENTRYDATA", GasHighMidStep},
		OP_INDEX:    {OP_INDEX, "INDEX", GasHighSlowStep},
		OP_ENTRYID:  {OP_ENTRYID, "ENTRYID", GasHighMidStep},
		OP_OUTPUTID: {OP_OUTPUTID, "OUTPUTID", GasHighMidStep},
		//OP_NONCE:       {OP_NONCE, "NONCE", GasHighMidStep},
		OP_BLOCKHEIGHT: {OP_BLOCKHEIGHT, "BLOCKHEIGHT", GasHighSlowStep},

		OP_CHECKPREDICATE: {OP_CHECKPREDICATE, "CHECKPREDICATE", GasMidStep},
	}

	gasByName map[string]gasInfo
)

const (
	GasAmount    int64 = 8
	GasAsset     int64 = 32
	GasBoolean   int64 = 1
	GasHash      int64 = 32
	GasInteger   int64 = 8
	GasProgram   int64 = 64
	GasPublicKey int64 = 32
	GasSignature int64 = 64
	GasString    int64 = 32
	GasTime      int64 = 8

	//push parament factor
	ContractParams int64 = 9
	ClauseParams   int64 = 8
)

var paramGasMap = map[string]int64{
	"Amount":    GasAmount,
	"Asset":     GasAsset,
	"Boolean":   GasBoolean,
	"Hash":      GasHash,
	"Integer":   GasInteger,
	"Program":   GasProgram,
	"PublicKey": GasPublicKey,
	"Signature": GasSignature,
	"String":    GasString,
	"Time":      GasTime,
}

func InitGas() {
	for i := 1; i <= 75; i++ {
		gasop[i] = gasInfo{Op(i), fmt.Sprintf("DATA_%d", i), GasData + int64(i)}
	}

	for i := uint8(0); i <= 15; i++ {
		op := uint8(OP_1) + i
		gasop[op] = gasInfo{Op(op), fmt.Sprintf("%d", i+1), GasSlowStep}
	}

	gasByName = make(map[string]gasInfo)
	for _, info := range gasop {
		gasByName[info.name] = info
	}
	gasByName["0"] = gasop[OP_FALSE]
	gasByName["TRUE"] = gasop[OP_1]

	gasop[OP_0] = gasop[OP_FALSE]
	gasop[OP_TRUE] = gasop[OP_1]

	for i := 0; i <= 255; i++ {
		if gasop[i].name == "" {
			gasop[i] = gasInfo{Op(i), fmt.Sprintf("NOPx%02x", i), GasBaseStep}
		}
	}
}

func GetGas(op Op) int64 {
	return gasop[op].Gas
}

func GetContractParamGas(typename string) int64 {
	paramGas, ok := paramGasMap[typename]
	if !ok {
		return -1
	}

	return paramGas + ContractParams
}

func GetClauseParamGas(typename string) int64 {
	paramGas, ok := paramGasMap[typename]
	if !ok {
		return -1
	}

	return paramGas + ClauseParams
}
