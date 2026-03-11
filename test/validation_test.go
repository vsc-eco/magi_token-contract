package contract_test

import (
	"fmt"
	"strings"
	"testing"
)

// ===================================
// Input Validation Tests (Negative)
// ===================================

// helpers

func longString(n int) string { return strings.Repeat("a", n) }

func pipeAddress() string { return "hive:user|evil" }

// ===================================
// Address Validation — Transfer
// ===================================

func TestTransferFailsAddressTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"to":%q,"amount":"100"}`, addr))
	CallContract(t, ct, "transfer", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferFailsAddressWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	payload := []byte(fmt.Sprintf(`{"to":%q,"amount":"100"}`, pipeAddress()))
	CallContract(t, ct, "transfer", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

// ===================================
// Address Validation — TransferFrom
// ===================================

func TestTransferFromFailsFromAddressTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"from":%q,"to":"hive:other","amount":"100"}`, addr))
	CallContract(t, ct, "transferFrom", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferFromFailsFromAddressWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	payload := []byte(fmt.Sprintf(`{"from":%q,"to":"hive:other","amount":"100"}`, pipeAddress()))
	CallContract(t, ct, "transferFrom", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferFromFailsToAddressTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"from":"hive:tibfox","to":%q,"amount":"100"}`, addr))
	CallContract(t, ct, "transferFrom", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferFromFailsToAddressWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	payload := []byte(fmt.Sprintf(`{"from":"hive:tibfox","to":%q,"amount":"100"}`, pipeAddress()))
	CallContract(t, ct, "transferFrom", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

// ===================================
// Address Validation — Approve
// ===================================

func TestApproveFailsSpenderAddressTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"spender":%q,"amount":"100"}`, addr))
	CallContract(t, ct, "approve", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

func TestApproveFailsSpenderAddressWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	payload := []byte(fmt.Sprintf(`{"spender":%q,"amount":"100"}`, pipeAddress()))
	CallContract(t, ct, "approve", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

// ===================================
// Address Validation — ChangeOwner
// ===================================

func TestChangeOwnerFailsNewOwnerTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"newOwner":%q}`, addr))
	CallContract(t, ct, "changeOwner", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

func TestChangeOwnerFailsNewOwnerWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	payload := []byte(fmt.Sprintf(`{"newOwner":%q}`, pipeAddress()))
	CallContract(t, ct, "changeOwner", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

// ===================================
// Init: Name / Symbol Length
// ===================================

func TestInitFailsNameTooLong(t *testing.T) {
	ct := SetupContractTest()
	name := longString(65)
	payload := []byte(fmt.Sprintf(`{"name":%q,"symbol":"MAGI","decimals":3,"maxSupply":"1000000000"}`, name))
	CallContract(t, ct, "init", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

func TestInitFailsSymbolTooLong(t *testing.T) {
	ct := SetupContractTest()
	symbol := longString(17)
	payload := []byte(fmt.Sprintf(`{"name":"Magi Token","symbol":%q,"decimals":3,"maxSupply":"1000000000"}`, symbol))
	CallContract(t, ct, "init", payload, nil, ownerAddress, false, uint(100_000_000), "")
}
