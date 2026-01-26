package contract_test

import (
	"testing"
)

// ===================================
// MAGI Token Contract Tests
// ===================================

// Init tests
func TestInitSuccess(t *testing.T) {
	ct := SetupContractTest()
	// Init should succeed when called by creator
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Verify owner is set
	CallContract(t, ct, "getOwner", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"owner":"hive:vaultec"}`)
	// Verify supply is 0
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":0}`)
}

func TestInitFailsNonCreator(t *testing.T) {
	ct := SetupContractTest()
	// Init should fail when called by non-creator
	CallContract(t, ct, "init", []byte(""), nil, "hive:someone", false, uint(100_000_000), "")
}

func TestInitFailsDouble(t *testing.T) {
	ct := SetupContractTest()
	// First init succeeds
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Second init fails
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", false, uint(100_000_000), "")
}

// Mint tests
func TestMintSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Mint 1000 tokens
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check supply
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":1000}`)
	// Check balance
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:vaultec"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":1000}`)
}

func TestMintFailsNonOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Non-owner cannot mint
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:someone", false, uint(100_000_000), "")
}

func TestMintFailsNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// Cannot mint before init
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestMintFailsZeroAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot mint 0 tokens
	CallContract(t, ct, "mint", []byte(`{"amount":0}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestMintFailsExceedsMaxSupply(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot mint more than max supply (1 billion)
	CallContract(t, ct, "mint", []byte(`{"amount":1000000001}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

// Transfer tests
func TestTransferSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Transfer 500 to someone
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check balances
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:vaultec"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":500}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:someone"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":500}`)
}

func TestTransferFailsInsufficientBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":100}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot transfer more than balance
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":500}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestTransferFailsToSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot transfer to self
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:vaultec","amount":500}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestTransferFailsZeroAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot transfer 0 tokens
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":0}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

// Burn tests
func TestBurnSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Burn 300 tokens
	CallContract(t, ct, "burn", []byte(`{"amount":300}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check supply and balance
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":700}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:vaultec"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":700}`)
}

func TestBurnFailsInsufficientBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":100}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot burn more than balance
	CallContract(t, ct, "burn", []byte(`{"amount":500}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestBurnFailsZeroAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot burn 0 tokens
	CallContract(t, ct, "burn", []byte(`{"amount":0}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

// ChangeOwner tests
func TestChangeOwnerSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Change owner
	CallContract(t, ct, "changeOwner", []byte(`{"newOwner":"hive:newowner"}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Verify new owner
	CallContract(t, ct, "getOwner", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"owner":"hive:newowner"}`)
	// Old owner cannot mint
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", false, uint(100_000_000), "")
	// New owner can mint
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:newowner", true, uint(100_000_000), "")
}

func TestChangeOwnerFailsNonOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Non-owner cannot change owner
	CallContract(t, ct, "changeOwner", []byte(`{"newOwner":"hive:newowner"}`), nil, "hive:someone", false, uint(100_000_000), "")
}

// GetInfo test
func TestGetInfo(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Get token info
	CallContract(t, ct, "getInfo", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"name":"Magi Token","symbol":"MAGI","decimals":3,"maxSupply":1000000000}`)
}

// Integration test - full flow
func TestFullFlow(t *testing.T) {
	ct := SetupContractTest()
	// Initialize
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Mint tokens
	CallContract(t, ct, "mint", []byte(`{"amount":10000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Transfer to user A
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:userA","amount":3000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Transfer to user B
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:userB","amount":2000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// User A transfers to user B
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:userB","amount":1000}`), nil, "hive:userA", true, uint(100_000_000), "")
	// User B burns some
	CallContract(t, ct, "burn", []byte(`{"amount":500}`), nil, "hive:userB", true, uint(100_000_000), "")
	// Check final balances
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:vaultec"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":5000}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:userA"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":2000}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:userB"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":2500}`)
	// Check final supply
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":9500}`)
}

// ===================================
// Allowance Tests (ERC-20 style)
// ===================================

func TestApproveSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Approve spender for 500 tokens
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check allowance
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:vaultec","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":500}`)
}

func TestApproveFailsSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot approve self
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:vaultec","amount":500}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestApproveOverwrite(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Approve 500
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Overwrite with 1000
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check new allowance
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:vaultec","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":1000}`)
}

func TestTransferFromSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Owner approves spender
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Spender transfers from owner to recipient
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:vaultec","to":"hive:recipient","amount":300}`), nil, "hive:spender", true, uint(100_000_000), "")
	// Check balances
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:vaultec"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":700}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:recipient"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":300}`)
	// Check remaining allowance
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:vaultec","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":200}`)
}

func TestTransferFromFailsInsufficientAllowance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Owner approves spender for only 100
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":100}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Spender tries to transfer 500 (more than allowance)
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:vaultec","to":"hive:recipient","amount":500}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestTransferFromFailsInsufficientBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":100}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Owner approves spender for 500 (more than balance)
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Spender tries to transfer 500 (allowance ok, but balance insufficient)
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:vaultec","to":"hive:recipient","amount":500}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestTransferFromFailsZeroAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot transfer 0
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:vaultec","to":"hive:recipient","amount":0}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestTransferFromFailsSameAddress(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot transfer to same address
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:vaultec","to":"hive:vaultec","amount":100}`), nil, "hive:spender", false, uint(100_000_000), "")
}

// Integration test with allowance - DEX-like flow
func TestAllowanceFlow(t *testing.T) {
	ct := SetupContractTest()
	// Initialize and mint
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":10000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Transfer to user
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:user","amount":5000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// User approves DEX contract
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:dex","amount":2000}`), nil, "hive:user", true, uint(100_000_000), "")
	// DEX swaps user's tokens to another user (simulating a trade)
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:user","to":"hive:buyer","amount":1500}`), nil, "hive:dex", true, uint(100_000_000), "")
	// Check final state
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:user"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":3500}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:buyer"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":1500}`)
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:user","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":500}`)
}

// ===================================
// IncreaseAllowance / DecreaseAllowance Tests
// ===================================

func TestIncreaseAllowanceSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Approve 500 initially
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Increase by 300
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:spender","amount":300}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check allowance is now 800
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:vaultec","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":800}`)
}

func TestIncreaseAllowanceFromZero(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Increase from zero
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:spender","amount":100}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check allowance is 100
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:vaultec","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":100}`)
}

func TestIncreaseAllowanceFailsSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot increase allowance for self
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:vaultec","amount":100}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestDecreaseAllowanceSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Approve 500 initially
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Decrease by 200
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:spender","amount":200}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check allowance is now 300
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:vaultec","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":300}`)
}

func TestDecreaseAllowanceToZero(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Approve 500 initially
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Decrease by 500 (to zero)
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check allowance is 0
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:vaultec","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":0}`)
}

func TestDecreaseAllowanceFailsBelowZero(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Approve 500 initially
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot decrease below zero
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:spender","amount":600}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestDecreaseAllowanceFailsSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot decrease allowance for self
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:vaultec","amount":100}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

// ===================================
// Pausable Tests
// ===================================

func TestPauseSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Pause contract
	CallContract(t, ct, "pause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check paused state
	CallContract(t, ct, "isPaused", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"paused":true}`)
}

func TestPauseFailsNonOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Non-owner cannot pause
	CallContract(t, ct, "pause", []byte(""), nil, "hive:someone", false, uint(100_000_000), "")
}

func TestPauseFailsAlreadyPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot pause when already paused
	CallContract(t, ct, "pause", []byte(""), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestUnpauseSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Unpause contract
	CallContract(t, ct, "unpause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check not paused
	CallContract(t, ct, "isPaused", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"paused":false}`)
}

func TestUnpauseFailsNonOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Non-owner cannot unpause
	CallContract(t, ct, "unpause", []byte(""), nil, "hive:someone", false, uint(100_000_000), "")
}

func TestUnpauseFailsNotPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Cannot unpause when not paused
	CallContract(t, ct, "unpause", []byte(""), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestTransferFailsWhenPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Transfer fails when paused
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":100}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestTransferFromFailsWhenPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// TransferFrom fails when paused
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:vaultec","to":"hive:recipient","amount":100}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestMintFailsWhenPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Mint fails when paused
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestBurnFailsWhenPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Burn fails when paused
	CallContract(t, ct, "burn", []byte(`{"amount":100}`), nil, "hive:vaultec", false, uint(100_000_000), "")
}

func TestOperationsResumeAfterUnpause(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Pause
	CallContract(t, ct, "pause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Transfer fails
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":100}`), nil, "hive:vaultec", false, uint(100_000_000), "")
	// Unpause
	CallContract(t, ct, "unpause", []byte(""), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Transfer works again
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":100}`), nil, "hive:vaultec", true, uint(100_000_000), "")
	// Check balance
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:someone"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":100}`)
}
