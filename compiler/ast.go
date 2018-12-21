package compiler

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	chainjson "github.com/bytom/encoding/json"
)

// Contract is a compiled Equity contract.
type Contract struct {
	// Name is the contract name.
	Name string `json:"name"`

	// Params is the list of contract parameters.
	Params []*Param `json:"params,omitempty"`

	// Clauses is the list of contract clauses.
	Clauses []*Clause `json:"clauses"`

	// Value is the name of the value locked by the contract.
	Value ValueInfo `json:"value"`

	// Body is the optimized bytecode of the contract body. This is not
	// a complete program!  Use instantiate to turn this (plus some
	// arguments) into a program.
	Body chainjson.HexBytes `json:"body_bytecode"`

	// Opcodes is the human-readable string of opcodes corresponding to
	// Body.
	Opcodes string `json:"body_opcodes,omitempty"`

	// Recursive tells whether this contract calls itself.  (This is
	// used to select between two possible instantiation options.)
	Recursive bool `json:"recursive"`

	// Pre-optimized list of instruction steps, with stack snapshots.
	Steps []Step `json:"-"`
}

// Param is a contract or clause parameter.
type Param struct {
	// Name is the parameter name.
	Name string `json:"name"`

	// Type is the declared parameter type.
	Type typeDesc `json:"type"`

	// InferredType, if available, is a more-specific type than Type,
	// inferred from the logic of the contract.
	InferredType typeDesc `json:"inferred_type,omitempty"`
}

// Clause is a compiled contract clause.
type Clause struct {
	// Name is the clause name.
	Name string `json:"name"`

	// Params is the list of clause parameters.
	Params []*Param `json:"params,omitempty"`

	statements []statement

	// BlockHeight is the list of expressions passed to greater()/less() in this
	// clause.
	BlockHeight []string `json:"blockheight,omitempty"`

	// HashCalls is the list of hash functions and their arguments used
	// in this clause.
	HashCalls []HashCall `json:"hash_calls,omitempty"`

	// Values is the list of values unlocked or relocked in this clause.
	Values []ValueInfo `json:"values"`

	// Conditions is the list of condition for if-else statements which body contains
	// the lock or unlock statement in this clause.
	Conditions map[string]Condition `json:"conditions"`

	// CondValues is the map of values unlocked or relocked in this clause's
	// if-else statements which body contains the lock or unlock statement.
	CondValues map[string][]ValueInfo `json:"cond_values"`

	// Contracts is the list of contracts called by this clause.
	Contracts []string `json:"contracts,omitempty"`
}

// ValueInfo describes how a blockchain value is used in a contract clause.
type ValueInfo struct {
	// Name is the clause's name for this value.
	Name string `json:"name"`

	// Program is the program expression used to the lock the value, if
	// the value is locked with "lock." If it's unlocked with "unlock"
	// instead, this is empty.
	Program string `json:"program,omitempty"`

	// Asset is the expression describing the asset type the value must
	// have, as it appears in a clause's "requires" section. If this is
	// the contract value instead, this is empty.
	Asset string `json:"asset,omitempty"`

	// Amount is the expression describing the amount the value must
	// have, as it appears in a clause's "requires" section. If this is
	// the contract value instead, this is empty.
	Amount string `json:"amount,omitempty"`

	// ContractCalls is the list of arguments for program which is a contract.
	ContractCalls []CallArgs `json:"contract_calls,omitempty"`

	// AmountParams is the list of parameters for Amount expression.
	// If the value of amount is a variable, this is empty.
	AmountParams []*Param `json:"amount_params,omitempty"`
}

// HashCall describes a call to a hash function.
type HashCall struct {
	// HashType is "sha3" or "sha256".
	HashType string `json:"hash_type"`

	// Arg is the expression passed to the hash function.
	Arg string `json:"arg"`

	// ArgType is the type of Arg.
	ArgType string `json:"arg_type"`
}

// Condition describes a condition expression.
type Condition struct {
	// Source is the string format of condition expression.
	Source string `json:"source"`

	// Params is the list of parameters for condition expression.
	Params []*Param `json:"params,omitempty"`
}

// CallArgs describes a argument expression for function call.
type CallArgs struct {
	// Source is the string format of argument expression.
	Source string `json:"source"`

	// Position is the position of argument expression.
	Position int `json:"position"`

	// Params is the list of parameters for argument expression.
	Params []*Param `json:"params,omitempty"`
}

// ContractArg is an argument with which to instantiate a contract as
// a program. Exactly one of B, I, and S should be supplied.
type ContractArg struct {
	B *bool               `json:"boolean,omitempty"`
	I *int64              `json:"integer,omitempty"`
	S *chainjson.HexBytes `json:"string,omitempty"`
}

type statement interface {
	countVarRefs(map[string]int)
}

type defineStatement struct {
	variable *Param
	expr     expression
}

func (s defineStatement) countVarRefs(counts map[string]int) {
	s.expr.countVarRefs(counts)
}

type assignStatement struct {
	variable *Param
	expr     expression
}

func (s assignStatement) countVarRefs(counts map[string]int) {
	s.expr.countVarRefs(counts)
}

// IfStatmentBody describes the content of if-else structure
type IfStatmentBody struct {
	// if body statements
	trueBody []statement

	// else body statements
	falseBody []statement
}

type ifStatement struct {
	condition expression
	body      *IfStatmentBody
}

func (s ifStatement) countVarRefs(counts map[string]int) {
	s.condition.countVarRefs(counts)
}

type verifyStatement struct {
	expr expression
}

func (s verifyStatement) countVarRefs(counts map[string]int) {
	s.expr.countVarRefs(counts)
}

type lockStatement struct {
	lockedAmount expression
	lockedAsset  expression
	program      expression

	// Added as a decoration, used by CHECKOUTPUT
	index int64
}

func (s lockStatement) countVarRefs(counts map[string]int) {
	s.lockedAmount.countVarRefs(counts)
	s.lockedAsset.countVarRefs(counts)
	s.program.countVarRefs(counts)
}

type unlockStatement struct {
	unlockedAmount expression
	unlockedAsset  expression
}

func (s unlockStatement) countVarRefs(counts map[string]int) {
	s.unlockedAmount.countVarRefs(counts)
	s.unlockedAsset.countVarRefs(counts)
}

type expression interface {
	String() string
	typ(*environ) typeDesc
	countVarRefs(map[string]int)
}

type binaryExpr struct {
	left, right expression
	op          *binaryOp
}

func (e binaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.left, e.op.op, e.right)
}

func (e binaryExpr) typ(*environ) typeDesc {
	return e.op.result
}

func (e binaryExpr) countVarRefs(counts map[string]int) {
	e.left.countVarRefs(counts)
	e.right.countVarRefs(counts)
}

type unaryExpr struct {
	op   *unaryOp
	expr expression
}

func (e unaryExpr) String() string {
	return fmt.Sprintf("%s%s", e.op.op, e.expr)
}

func (e unaryExpr) typ(*environ) typeDesc {
	return e.op.result
}

func (e unaryExpr) countVarRefs(counts map[string]int) {
	e.expr.countVarRefs(counts)
}

type callExpr struct {
	fn   expression
	args []expression
}

func (e callExpr) String() string {
	var argStrs []string
	for _, a := range e.args {
		argStrs = append(argStrs, a.String())
	}
	return fmt.Sprintf("%s(%s)", e.fn, strings.Join(argStrs, ", "))
}

func (e callExpr) typ(env *environ) typeDesc {
	if b := referencedBuiltin(e.fn); b != nil {
		switch b.name {
		case "sha3":
			if len(e.args) == 1 {
				switch e.args[0].typ(env) {
				case strType:
					return sha3StrType
				case pubkeyType:
					return sha3PubkeyType
				}
			}

		case "sha256":
			if len(e.args) == 1 {
				switch e.args[0].typ(env) {
				case strType:
					return sha256StrType
				case pubkeyType:
					return sha256PubkeyType
				}
			}
		}

		return b.result
	}
	if e.fn.typ(env) == predType {
		return boolType
	}
	if e.fn.typ(env) == contractType {
		return progType
	}
	return nilType
}

func (e callExpr) countVarRefs(counts map[string]int) {
	e.fn.countVarRefs(counts)
	for _, a := range e.args {
		a.countVarRefs(counts)
	}
}

type varRef string

func (v varRef) String() string {
	return string(v)
}

func (v varRef) typ(env *environ) typeDesc {
	if entry := env.lookup(string(v)); entry != nil {
		return entry.t
	}
	return nilType
}

func (v varRef) countVarRefs(counts map[string]int) {
	counts[string(v)]++
}

type bytesLiteral []byte

func (e bytesLiteral) String() string {
	return "0x" + hex.EncodeToString([]byte(e))
}

func (bytesLiteral) typ(*environ) typeDesc {
	return "String"
}

func (bytesLiteral) countVarRefs(map[string]int) {}

type integerLiteral int64

func (e integerLiteral) String() string {
	return strconv.FormatInt(int64(e), 10)
}

func (integerLiteral) typ(*environ) typeDesc {
	return "Integer"
}

func (integerLiteral) countVarRefs(map[string]int) {}

type booleanLiteral bool

func (e booleanLiteral) String() string {
	if e {
		return "true"
	}
	return "false"
}

func (booleanLiteral) typ(*environ) typeDesc {
	return "Boolean"
}

func (booleanLiteral) countVarRefs(map[string]int) {}

type listExpr []expression

func (e listExpr) String() string {
	var elts []string
	for _, elt := range e {
		elts = append(elts, elt.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elts, ", "))
}

func (listExpr) typ(*environ) typeDesc {
	return "List"
}

func (e listExpr) countVarRefs(counts map[string]int) {
	for _, elt := range e {
		elt.countVarRefs(counts)
	}
}
