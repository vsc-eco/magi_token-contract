package contract_test

import (
	"testing"
)

// ===================================
// Edge Cases & Negative Tests
// ===================================

// --- Invalid Payload Tests ---

func TestMintInvalidJSON(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Malformed JSON
	CallContract(t, ct, "mint", []byte(`{invalid`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferInvalidJSON(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Malformed JSON
	CallContract(t, ct, "transfer", []byte(`{"to":`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferMissingTo(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Missing 'to' field
	CallContract(t, ct, "transfer", []byte(`{"amount":100}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferEmptyTo(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Empty 'to' field
	CallContract(t, ct, "transfer", []byte(`{"to":"","amount":100}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestBalanceOfEmptyAccount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Empty account field
	CallContract(t, ct, "balanceOf", []byte(`{"account":""}`), nil, "hive:anyone", false, uint(100_000_000), "")
}

func TestAllowanceEmptyOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Empty owner
	CallContract(t, ct, "allowance", []byte(`{"owner":"","spender":"hive:spender"}`), nil, "hive:anyone", false, uint(100_000_000), "")
}

func TestAllowanceEmptySpender(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Empty spender
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:owner","spender":""}`), nil, "hive:anyone", false, uint(100_000_000), "")
}

func TestApproveEmptySpender(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Empty spender
	CallContract(t, ct, "approve", []byte(`{"spender":"","amount":100}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestChangeOwnerEmptyAddress(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Empty new owner
	CallContract(t, ct, "changeOwner", []byte(`{"newOwner":""}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferFromEmptyFrom(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Empty from
	CallContract(t, ct, "transferFrom", []byte(`{"from":"","to":"hive:recipient","amount":100}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestTransferFromEmptyTo(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Empty to
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"","amount":100}`), nil, "hive:spender", false, uint(100_000_000), "")
}

// --- Boundary Condition Tests ---

func TestMintExactlyMaxSupply(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Mint exactly max supply (1 billion)
	CallContract(t, ct, "mint", []byte(`{"amount":1000000000}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Verify supply
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":1000000000}`)
	// Cannot mint any more
	CallContract(t, ct, "mint", []byte(`{"amount":1}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestMintApproachMaxSupply(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Mint 999,999,999
	CallContract(t, ct, "mint", []byte(`{"amount":999999999}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Mint 1 more (exactly at limit)
	CallContract(t, ct, "mint", []byte(`{"amount":1}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Verify at max
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":1000000000}`)
	// Cannot mint any more
	CallContract(t, ct, "mint", []byte(`{"amount":1}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferEntireBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Transfer all tokens
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":1000}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Verify balances
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":0}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:someone"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":1000}`)
}

func TestBurnEntireBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Burn all tokens
	CallContract(t, ct, "burn", []byte(`{"amount":1000}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Verify balance and supply
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":0}`)
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":0}`)
}

func TestTransferFromNoApproval(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, ownerAddress, true, uint(100_000_000), "")
	// TransferFrom with no prior approval (allowance = 0)
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"hive:recipient","amount":100}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestTransferFromEntireAllowance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":1000}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Use entire allowance
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"hive:recipient","amount":500}`), nil, "hive:spender", true, uint(100_000_000), "")
	// Allowance should be 0
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":0}`)
	// Cannot transfer more
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"hive:recipient","amount":1}`), nil, "hive:spender", false, uint(100_000_000), "")
}

// --- Zero Amount Edge Cases ---

func TestApproveZero(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Approve 500 first
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Approve 0 to revoke (should succeed per ERC-20)
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":0}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Verify allowance is 0
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":0}`)
}

func TestIncreaseAllowanceZero(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Increase by 0 (no-op, should succeed)
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:spender","amount":0}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Allowance unchanged
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":500}`)
}

func TestDecreaseAllowanceZero(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":500}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Decrease by 0 (no-op, should succeed)
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:spender","amount":0}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Allowance unchanged
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":500}`)
}

// --- Query Non-Existent Accounts ---

func TestBalanceOfNonExistentAccount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Query balance of account that never held tokens (should return 0)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:nonexistent"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":0}`)
}

func TestAllowanceNonExistentPair(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Query allowance for pair that was never set (should return 0)
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:owner","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":0}`)
}

// --- Operations Not Initialized ---

func TestTransferNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// Transfer before init
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":100}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestBurnNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// Burn before init
	CallContract(t, ct, "burn", []byte(`{"amount":100}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestApproveNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// Approve before init
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":100}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferFromNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// TransferFrom before init
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:owner","to":"hive:recipient","amount":100}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestPauseNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// Pause before init
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestChangeOwnerNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// ChangeOwner before init
	CallContract(t, ct, "changeOwner", []byte(`{"newOwner":"hive:newowner"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestIncreaseAllowanceNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// IncreaseAllowance before init
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:spender","amount":100}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestDecreaseAllowanceNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// DecreaseAllowance before init
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:spender","amount":100}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestBalanceOfNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// BalanceOf before init
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:anyone"}`), nil, "hive:anyone", false, uint(100_000_000), "")
}

func TestTotalSupplyNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// TotalSupply before init
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", false, uint(100_000_000), "")
}

func TestGetOwnerNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// GetOwner before init
	CallContract(t, ct, "getOwner", []byte(""), nil, "hive:anyone", false, uint(100_000_000), "")
}

func TestAllowanceNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// Allowance before init
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:owner","spender":"hive:spender"}`), nil, "hive:anyone", false, uint(100_000_000), "")
}

func TestIsPausedNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// IsPaused before init
	CallContract(t, ct, "isPaused", []byte(""), nil, "hive:anyone", false, uint(100_000_000), "")
}

func TestGetInfoNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// GetInfo before init
	CallContract(t, ct, "getInfo", []byte(""), nil, "hive:anyone", false, uint(100_000_000), "")
}

// --- ChangeOwner Edge Cases ---

func TestChangeOwnerToSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Change owner to self (same address) - should succeed
	CallContract(t, ct, "changeOwner", []byte(`{"newOwner":"hive:tibfox"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Owner unchanged
	CallContract(t, ct, "getOwner", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"owner":"hive:tibfox"}`)
}
