package main

import (
	"encoding/binary"
	"magi_token/sdk"

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
	sdk.StateSetObject(balanceKey(account), string(u64ToBytes(newBal)))
}

// decBalance decrements token balance of an address. Aborts if insufficient balance.
func decBalance(account string, amount uint64) {
	oldBal := getBalanceInternal(account)
	if oldBal < amount {
		sdk.Abort("Insufficient balance")
	}
	newBal := safeSub(oldBal, amount)
	sdk.StateSetObject(balanceKey(account), string(u64ToBytes(newBal)))
}

// getBalanceInternal retrieves token balance of an address.
func getBalanceInternal(account string) uint64 {
	bal := sdk.StateGetObject(balanceKey(account))
	if bal == nil || *bal == "" {
		return 0
	}
	amt := bytesToU64([]byte(*bal))
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
	if alw == nil || *alw == "" {
		return 0
	}
	amt := bytesToU64([]byte(*alw))
	return amt
}

// setAllowanceInternal sets the allowance for a spender on owner's tokens.
func setAllowanceInternal(owner, spender string, amount uint64) {
	sdk.StateSetObject(allowanceKey(owner, spender), string(u64ToBytes(amount)))
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
	if d == nil || *d == "" {
		return 0
	}
	return (*d)[0]
}

// getMaxSupply retrieves the max supply from state.
func getMaxSupply() uint64 {
	m := sdk.StateGetObject("token_max_supply")
	if m == nil || *m == "" {
		return 0
	}
	return bytesToU64([]byte(*m))
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
	return bytesToU64([]byte(*s))
}

// setSupply sets the total supply.
func setSupply(amount uint64) {
	sdk.StateSetObject("supply", string(u64ToBytes(amount)))
}

// ===================================
// uint64 <-> []byte Helper
// ===================================

func u64ToBytes(val uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, val)

	// In Little Endian, leading zeros (most significant bytes) are at the end of the slice.
	// We iterate backwards to find the last non-zero byte.
	lastNonZeroIndex := len(b) - 1
	for lastNonZeroIndex >= 0 {
		if b[lastNonZeroIndex] != 0 {
			break
		}
		lastNonZeroIndex--
	}

	// If the value was 0, ensure we return at least one byte (0x00) instead of an empty slice.
	if lastNonZeroIndex < 0 {
		return []byte{0}
	}

	return b[:lastNonZeroIndex+1]
}

func bytesToU64(b []byte) uint64 {
	if len(b) > 8 {
		sdk.Abort("byte length less than or equal to 8")
	}

	// Create an 8-byte buffer initialized to zeros.
	buf := make([]byte, 8)

	// In Little Endian, the existing bytes are the least significant and go at the start.
	// Copy the input slice into the beginning of the buffer.
	copy(buf, b)

	return binary.LittleEndian.Uint64(buf)
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
