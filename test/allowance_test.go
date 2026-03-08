package contract_test

import (
	"testing"
)

// ===================================
// Allowance Tests (ERC-20 style)
// ===================================

// --- Approve Tests ---

func TestApproveSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Approve spender for 500 tokens
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Check allowance
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"500"}`)
}

func TestApproveFailsSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot approve self
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:tibfox","amount":"500"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestApproveOverwrite(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Approve 500
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Overwrite with 1000
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Check new allowance
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"1000"}`)
}

// --- TransferFrom Tests ---

func TestTransferFromSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Owner approves spender
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Spender transfers from owner to recipient
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"hive:recipient","amount":"300"}`), nil, "hive:spender", true, uint(100_000_000), "")
	// Check balances
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"700"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:recipient"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"300"}`)
	// Check remaining allowance
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"200"}`)
}

func TestTransferFromFailsInsufficientAllowance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Owner approves spender for only 100
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"100"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Spender tries to transfer 500 (more than allowance)
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"hive:recipient","amount":"500"}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestTransferFromFailsInsufficientBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"100"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Owner approves spender for 500 (more than balance)
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Spender tries to transfer 500 (allowance ok, but balance insufficient)
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"hive:recipient","amount":"500"}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestTransferFromFailsZeroAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot transfer 0
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"hive:recipient","amount":"0"}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestTransferFromFailsSameAddress(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot transfer to same address
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"hive:tibfox","amount":"100"}`), nil, "hive:spender", false, uint(100_000_000), "")
}

// --- IncreaseAllowance Tests ---

func TestIncreaseAllowanceSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Approve 500 initially
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Increase by 300
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:spender","amount":"300"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Check allowance is now 800
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"800"}`)
}

func TestIncreaseAllowanceFromZero(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Increase from zero
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:spender","amount":"100"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Check allowance is 100
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"100"}`)
}

func TestIncreaseAllowanceFailsSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot increase allowance for self
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:tibfox","amount":"100"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

// --- DecreaseAllowance Tests ---

func TestDecreaseAllowanceSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Approve 500 initially
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Decrease by 200
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:spender","amount":"200"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Check allowance is now 300
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"300"}`)
}

func TestDecreaseAllowanceToZero(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Approve 500 initially
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Decrease by 500 (to zero)
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Check allowance is 0
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"0"}`)
}

func TestDecreaseAllowanceFailsBelowZero(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Approve 500 initially
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot decrease below zero
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:spender","amount":"600"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestDecreaseAllowanceFailsSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot decrease allowance for self
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:tibfox","amount":"100"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

// --- Allowance Integration Test - DEX-like Flow ---

func TestAllowanceFlow(t *testing.T) {
	ct := SetupContractTest()
	// Initialize and mint
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"10000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Transfer to user
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:user","amount":"5000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// User approves DEX contract
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:dex","amount":"2000"}`), nil, "hive:user", true, uint(100_000_000), "")
	// DEX swaps user's tokens to another user (simulating a trade)
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:user","to":"hive:buyer","amount":"1500"}`), nil, "hive:dex", true, uint(100_000_000), "")
	// Check final state
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:user"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"3500"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:buyer"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"1500"}`)
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:user","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"500"}`)
}
