package contract_test

import (
	"testing"
)

// ===================================
// Pausable Tests
// ===================================

// --- Pause Tests ---

func TestPauseSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Pause contract
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Check paused state
	CallContract(t, ct, "isPaused", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"paused":true}`)
}

func TestPauseFailsNonOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Non-owner cannot pause
	CallContract(t, ct, "pause", []byte(""), nil, "hive:someone", false, uint(100_000_000), "")
}

func TestPauseFailsAlreadyPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot pause when already paused
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, false, uint(100_000_000), "")
}

// --- Unpause Tests ---

func TestUnpauseSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Unpause contract
	CallContract(t, ct, "unpause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Check not paused
	CallContract(t, ct, "isPaused", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"paused":false}`)
}

func TestUnpauseFailsNonOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Non-owner cannot unpause
	CallContract(t, ct, "unpause", []byte(""), nil, "hive:someone", false, uint(100_000_000), "")
}

func TestUnpauseFailsNotPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot unpause when not paused
	CallContract(t, ct, "unpause", []byte(""), nil, ownerAddress, false, uint(100_000_000), "")
}

// --- Operations When Paused ---

func TestTransferFailsWhenPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Transfer fails when paused
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":"100"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferFromFailsWhenPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// TransferFrom fails when paused
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"hive:recipient","amount":"100"}`), nil, "hive:spender", false, uint(100_000_000), "")
}

func TestMintFailsWhenPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Mint fails when paused
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestBurnFailsWhenPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Burn fails when paused
	CallContract(t, ct, "burn", []byte(`{"amount":"100"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

// --- Operations That Work When Paused ---

func TestApproveWhilePaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Approve should work while paused (allowance management not blocked)
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
}

func TestIncreaseAllowanceWhilePaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// IncreaseAllowance should work while paused
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:spender","amount":"100"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"600"}`)
}

func TestDecreaseAllowanceWhilePaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:spender","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// DecreaseAllowance should work while paused
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:spender","amount":"100"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:spender"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"400"}`)
}

func TestChangeOwnerWhilePaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// ChangeOwner should work while paused
	CallContract(t, ct, "changeOwner", []byte(`{"newOwner":"hive:newowner"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "getOwner", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"owner":"hive:newowner"}`)
}

// --- Resume After Unpause ---

func TestOperationsResumeAfterUnpause(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Pause
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Transfer fails
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":"100"}`), nil, ownerAddress, false, uint(100_000_000), "")
	// Unpause
	CallContract(t, ct, "unpause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	// Transfer works again
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":"100"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Check balance
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:someone"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"100"}`)
}
