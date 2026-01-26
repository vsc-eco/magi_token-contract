package main

import (
	"magi_token/sdk"
	"strconv"

	"github.com/CosmWasm/tinyjson/jwriter"
)

// ===================================
// Internal Helper Functions
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
	return "bal|" + account
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

// getBalanceInternal retrieves token balance of an address.
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
	return "alw|" + owner + "|" + spender
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
// Token Properties (from state)
// ===================================

// getTokenName retrieves the token name from state.
func getTokenName() string {
	n := sdk.StateGetObject("token_name")
	if n == nil {
		return ""
	}
	return *n
}

// getTokenSymbol retrieves the token symbol from state.
func getTokenSymbol() string {
	s := sdk.StateGetObject("token_symbol")
	if s == nil {
		return ""
	}
	return *s
}

// getTokenDecimals retrieves the token decimals from state.
func getTokenDecimals() uint8 {
	d := sdk.StateGetObject("token_decimals")
	if d == nil {
		return 0
	}
	decimals, _ := strconv.ParseUint(*d, 10, 8)
	return uint8(decimals)
}

// getMaxSupply retrieves the max supply from state.
func getMaxSupply() uint64 {
	m := sdk.StateGetObject("token_max_supply")
	if m == nil {
		return 0
	}
	maxSupply, _ := strconv.ParseUint(*m, 10, 64)
	return maxSupply
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
