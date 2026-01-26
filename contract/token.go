package main

import (
	"magi_token/sdk"
	"strconv"

	"github.com/CosmWasm/tinyjson/jlexer"
	"github.com/CosmWasm/tinyjson/jwriter"
)

// ===================================
// MAGI Token - Core Token Functions
// ===================================

// ===================================
// Safe Math Utilities
// ===================================

// safeAdd performs a + b addition. Aborts execution if an overflow is detected.
func safeAdd(a, b uint64) uint64 {
	sum := a + b
	if sum < a {
		sdk.Abort("safeAdd overflow")
	}
	return sum
}

// safeSub performs a - b subtraction. Aborts execution if an underflow is detected.
func safeSub(a, b uint64) uint64 {
	if b > a {
		sdk.Abort("safeSub underflow")
	}
	return a - b
}

// ===================================
// Balance Management
// ===================================

// balanceKey returns the state key for an account's balance.
func balanceKey(account string) string {
	return "bal_" + account
}

// incBalance increments token balance of an address.
func incBalance(account string, amount uint64) {
	oldBal := getBalanceInternal(account)
	newBal := safeAdd(oldBal, amount)
	sdk.StateSetObject(balanceKey(account), strconv.FormatUint(newBal, 10))
}

// decBalance decrements token balance of an address. Aborts if insufficient balance.
func decBalance(account string, amount uint64) {
	oldBal := getBalanceInternal(account)
	if oldBal < amount {
		sdk.Abort("Insufficient balance")
	}
	newBal := safeSub(oldBal, amount)
	sdk.StateSetObject(balanceKey(account), strconv.FormatUint(newBal, 10))
}

// getBalanceInternal retrieves token balance of an address (internal use).
func getBalanceInternal(account string) uint64 {
	bal := sdk.StateGetObject(balanceKey(account))
	if bal == nil {
		return 0
	}
	amt, _ := strconv.ParseUint(*bal, 10, 64)
	return amt
}

// ===================================
// Allowance Management
// ===================================

// allowanceKey returns the state key for an allowance (owner approves spender).
func allowanceKey(owner, spender string) string {
	return "alw_" + owner + "_" + spender
}

// getAllowanceInternal retrieves the allowance for a spender on owner's tokens.
func getAllowanceInternal(owner, spender string) uint64 {
	alw := sdk.StateGetObject(allowanceKey(owner, spender))
	if alw == nil {
		return 0
	}
	amt, _ := strconv.ParseUint(*alw, 10, 64)
	return amt
}

// setAllowanceInternal sets the allowance for a spender on owner's tokens.
func setAllowanceInternal(owner, spender string, amount uint64) {
	sdk.StateSetObject(allowanceKey(owner, spender), strconv.FormatUint(amount, 10))
}

// ===================================
// Supply Management
// ===================================

// getSupply retrieves the current total supply.
func getSupply() uint64 {
	s := sdk.StateGetObject("supply")
	if s == nil {
		return 0
	}
	if *s == "" {
		return 0
	}
	supply, _ := strconv.ParseUint(*s, 10, 64)
	return supply
}

// setSupply sets the total supply.
func setSupply(amount uint64) {
	sdk.StateSetObject("supply", strconv.FormatUint(amount, 10))
}

// ===================================
// JSON Response Helper
// ===================================

func jsonResponse(marshaler interface{ MarshalTinyJSON(*jwriter.Writer) }) *string {
	w := jwriter.Writer{}
	marshaler.MarshalTinyJSON(&w)
	result := string(w.Buffer.BuildBytes())
	return &result
}

// ===================================
// Public Contract Entry Points (WASM)
// ===================================

// Init initializes the token contract.
// Can only be called once by the Creator.
//
//go:wasmexport init
func Init(_ *string) *string {
	if isInit() {
		sdk.Abort("Already initialized")
	}
	caller := sdk.GetEnvKey("msg.sender")
	if caller == nil {
		sdk.Abort("Caller must be creator to initialize")
	}
	if *caller != Creator {
		sdk.Abort("Caller must be creator to initialize")
	}
	sdk.StateSetObject("isInit", "1")
	sdk.StateSetObject("supply", "0")
	sdk.StateSetObject("owner", Creator)
	emitInit(Creator)
	return jsonResponse(SuccessResponse{Success: true})
}

// Mint creates new tokens and assigns them to the owner.
// Payload: {"amount": 1000}
// Only the owner can mint tokens.
//
//go:wasmexport mint
func Mint(payload *string) *string {
	assertInit()
	assertNotPaused()
	owner, isOwner := getOwner()
	if !isOwner {
		sdk.Abort("Must be owner to mint")
	}
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p MintPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Amount == 0 {
		sdk.Abort("Amount must be greater than 0")
	}
	supply := getSupply()
	newSupply := safeAdd(supply, p.Amount)
	if newSupply > MaxSupply {
		sdk.Abort("Exceeded max supply")
	}
	setSupply(newSupply)
	incBalance(owner, p.Amount)
	emitMint(owner, p.Amount)
	return jsonResponse(SuccessResponse{Success: true})
}

// Burn destroys tokens from the caller's balance.
// Payload: {"amount": 500}
//
//go:wasmexport burn
func Burn(payload *string) *string {
	assertInit()
	assertNotPaused()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p BurnPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Amount == 0 {
		sdk.Abort("Amount must be greater than 0")
	}
	caller := sdk.GetEnvKey("msg.sender")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	decBalance(*caller, p.Amount)
	supply := getSupply()
	newSupply := safeSub(supply, p.Amount)
	setSupply(newSupply)
	emitBurn(*caller, p.Amount)
	return jsonResponse(SuccessResponse{Success: true})
}

// Transfer moves tokens from caller to recipient.
// Payload: {"to": "hive:recipient", "amount": 100}
//
//go:wasmexport transfer
func Transfer(payload *string) *string {
	assertInit()
	assertNotPaused()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p TransferPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.To == "" {
		sdk.Abort("Recipient required")
	}
	if p.Amount == 0 {
		sdk.Abort("Amount must be greater than 0")
	}

	caller := sdk.GetEnvKey("msg.sender")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	from := *caller

	if from == p.To {
		sdk.Abort("Cannot transfer to self")
	}

	decBalance(from, p.Amount)
	incBalance(p.To, p.Amount)
	emitTransfer(from, p.To, p.Amount)
	return jsonResponse(SuccessResponse{Success: true})
}

// Approve sets the allowance for a spender to spend caller's tokens.
// Payload: {"spender": "hive:spender", "amount": 100}
//
//go:wasmexport approve
func Approve(payload *string) *string {
	assertInit()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p ApprovePayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Spender == "" {
		sdk.Abort("Spender required")
	}

	caller := sdk.GetEnvKey("msg.sender")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	owner := *caller

	if owner == p.Spender {
		sdk.Abort("Cannot approve self")
	}

	setAllowanceInternal(owner, p.Spender, p.Amount)
	emitApproval(owner, p.Spender, p.Amount)
	return jsonResponse(SuccessResponse{Success: true})
}

// TransferFrom moves tokens from one address to another using allowance.
// Payload: {"from": "hive:owner", "to": "hive:recipient", "amount": 100}
// Caller must have sufficient allowance from 'from' address.
//
//go:wasmexport transferFrom
func TransferFrom(payload *string) *string {
	assertInit()
	assertNotPaused()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p TransferFromPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.From == "" {
		sdk.Abort("From address required")
	}
	if p.To == "" {
		sdk.Abort("To address required")
	}
	if p.Amount == 0 {
		sdk.Abort("Amount must be greater than 0")
	}

	caller := sdk.GetEnvKey("msg.sender")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	spender := *caller

	if p.From == p.To {
		sdk.Abort("Cannot transfer to same address")
	}

	// Check and deduct allowance
	allowance := getAllowanceInternal(p.From, spender)
	if allowance < p.Amount {
		sdk.Abort("Insufficient allowance")
	}
	setAllowanceInternal(p.From, spender, safeSub(allowance, p.Amount))

	// Transfer tokens
	decBalance(p.From, p.Amount)
	incBalance(p.To, p.Amount)
	emitTransfer(p.From, p.To, p.Amount)
	return jsonResponse(SuccessResponse{Success: true})
}

// IncreaseAllowance atomically increases the allowance for a spender.
// Payload: {"spender": "hive:spender", "amount": 100}
//
//go:wasmexport increaseAllowance
func IncreaseAllowance(payload *string) *string {
	assertInit()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p AllowancePayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Spender == "" {
		sdk.Abort("Spender required")
	}

	caller := sdk.GetEnvKey("msg.sender")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	owner := *caller

	if owner == p.Spender {
		sdk.Abort("Cannot approve self")
	}

	currentAllowance := getAllowanceInternal(owner, p.Spender)
	newAllowance := safeAdd(currentAllowance, p.Amount)
	setAllowanceInternal(owner, p.Spender, newAllowance)
	emitApproval(owner, p.Spender, newAllowance)
	return jsonResponse(SuccessResponse{Success: true})
}

// DecreaseAllowance atomically decreases the allowance for a spender.
// Payload: {"spender": "hive:spender", "amount": 100}
//
//go:wasmexport decreaseAllowance
func DecreaseAllowance(payload *string) *string {
	assertInit()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p AllowancePayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Spender == "" {
		sdk.Abort("Spender required")
	}

	caller := sdk.GetEnvKey("msg.sender")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	owner := *caller

	if owner == p.Spender {
		sdk.Abort("Cannot approve self")
	}

	currentAllowance := getAllowanceInternal(owner, p.Spender)
	if currentAllowance < p.Amount {
		sdk.Abort("Decreased allowance below zero")
	}
	newAllowance := safeSub(currentAllowance, p.Amount)
	setAllowanceInternal(owner, p.Spender, newAllowance)
	emitApproval(owner, p.Spender, newAllowance)
	return jsonResponse(SuccessResponse{Success: true})
}

// ChangeOwner transfers contract ownership to a new address.
// Payload: {"newOwner": "hive:newowner"}
//
//go:wasmexport changeOwner
func ChangeOwner(payload *string) *string {
	assertInit()
	_, isOwner := getOwner()
	if !isOwner {
		sdk.Abort("Not owner")
	}
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p ChangeOwnerPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.NewOwner == "" {
		sdk.Abort("New owner required")
	}

	sdk.StateSetObject("owner", p.NewOwner)
	emitOwnerChange(p.NewOwner)
	return jsonResponse(SuccessResponse{Success: true})
}

// Pause pauses all token transfers. Only owner can pause.
//
//go:wasmexport pause
func Pause(_ *string) *string {
	assertInit()
	owner, isOwner := getOwner()
	if !isOwner {
		sdk.Abort("Not owner")
	}
	if isPaused() {
		sdk.Abort("Already paused")
	}
	sdk.StateSetObject("paused", "1")
	emitPaused(owner)
	return jsonResponse(SuccessResponse{Success: true})
}

// Unpause unpauses all token transfers. Only owner can unpause.
//
//go:wasmexport unpause
func Unpause(_ *string) *string {
	assertInit()
	owner, isOwner := getOwner()
	if !isOwner {
		sdk.Abort("Not owner")
	}
	if !isPaused() {
		sdk.Abort("Not paused")
	}
	sdk.StateSetObject("paused", "0")
	emitUnpaused(owner)
	return jsonResponse(SuccessResponse{Success: true})
}

// ===================================
// Read-Only Getters (WASM)
// ===================================

// BalanceOf returns the token balance of an address.
// Payload: {"account": "hive:user"}
//
//go:wasmexport balanceOf
func BalanceOf(payload *string) *string {
	assertInit()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p BalancePayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Account == "" {
		sdk.Abort("Account required")
	}

	bal := getBalanceInternal(p.Account)
	return jsonResponse(BalanceResponse{Balance: bal})
}

// TotalSupply returns the current total supply.
//
//go:wasmexport totalSupply
func TotalSupply(_ *string) *string {
	assertInit()
	supply := getSupply()
	return jsonResponse(SupplyResponse{TotalSupply: supply})
}

// GetOwner returns the current contract owner.
//
//go:wasmexport getOwner
func GetOwnerExport(_ *string) *string {
	assertInit()
	owner, _ := getOwner()
	return jsonResponse(OwnerResponse{Owner: owner})
}

// Allowance returns the allowance for a spender on owner's tokens.
// Payload: {"owner": "hive:owner", "spender": "hive:spender"}
//
//go:wasmexport allowance
func Allowance(payload *string) *string {
	assertInit()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p GetAllowancePayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Owner == "" {
		sdk.Abort("Owner required")
	}
	if p.Spender == "" {
		sdk.Abort("Spender required")
	}

	alw := getAllowanceInternal(p.Owner, p.Spender)
	return jsonResponse(AllowanceResponse{Allowance: alw})
}

// GetInfo returns token metadata.
//
//go:wasmexport getInfo
func GetInfo(_ *string) *string {
	return jsonResponse(InfoResponse{
		Name:      Name,
		Symbol:    Symbol,
		Decimals:  Precision,
		MaxSupply: MaxSupply,
	})
}

// IsPausedExport returns whether the contract is paused.
//
//go:wasmexport isPaused
func IsPausedExport(_ *string) *string {
	assertInit()
	return jsonResponse(PausedResponse{Paused: isPaused()})
}
