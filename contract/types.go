package main

// ===================================
// MAGI Token - JSON Types (tinyjson)
// ===================================

// ===================================
// Payload Types (Input)
// ===================================

// InitPayload for init action
type InitPayload struct {
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	Decimals  uint8  `json:"decimals"`
	MaxSupply uint64 `json:"maxSupply"`
}

// TransferPayload for transfer action
type TransferPayload struct {
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

// TransferFromPayload for transferFrom action
type TransferFromPayload struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

// ApprovePayload for approve action
type ApprovePayload struct {
	Spender string `json:"spender"`
	Amount  uint64 `json:"amount"`
}

// AllowancePayload for increaseAllowance/decreaseAllowance actions
type AllowancePayload struct {
	Spender string `json:"spender"`
	Amount  uint64 `json:"amount"`
}

// AllowanceQueryPayload for getAllowance query
type AllowanceQueryPayload struct {
	Owner   string `json:"owner"`
	Spender string `json:"spender"`
}

// BalancePayload for getBalance query
type BalancePayload struct {
	Account string `json:"account"`
}

// MintPayload for mint action
type MintPayload struct {
	Amount uint64 `json:"amount"`
}

// BurnPayload for burn action
type BurnPayload struct {
	Amount uint64 `json:"amount"`
}

// ChangeOwnerPayload for changeOwner action
type ChangeOwnerPayload struct {
	NewOwner string `json:"newOwner"`
}

// ===================================
// Response Types (Output)
// ===================================

// BalanceResponse for balance queries
type BalanceResponse struct {
	Balance uint64 `json:"balance"`
}

// SupplyResponse for supply queries
type SupplyResponse struct {
	TotalSupply uint64 `json:"totalSupply"`
}

// AllowanceResponse for allowance queries
type AllowanceResponse struct {
	Allowance uint64 `json:"allowance"`
}

// OwnerResponse for owner queries
type OwnerResponse struct {
	Owner string `json:"owner"`
}

// InfoResponse for token info queries
type InfoResponse struct {
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	Decimals  int    `json:"decimals"`
	MaxSupply uint64 `json:"maxSupply"`
}

// PausedResponse for isPaused queries
type PausedResponse struct {
	Paused bool `json:"paused"`
}

// SuccessResponse for mutation operations
type SuccessResponse struct {
	Success bool `json:"success"`
}

// ===================================
// Event Types
// ===================================

// InitEvent for contract initialization
type InitEvent struct {
	Type       string         `json:"type"`
	Attributes InitAttributes `json:"attributes"`
}

type InitAttributes struct {
	Owner     string `json:"owner"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	Decimals  int    `json:"decimals"`
	MaxSupply uint64 `json:"maxSupply"`
}

// TransferEvent for token transfers
type TransferEvent struct {
	Type       string             `json:"type"`
	Attributes TransferAttributes `json:"attributes"`
}

type TransferAttributes struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
}

// ApprovalEvent for allowance approvals
type ApprovalEvent struct {
	Type       string             `json:"type"`
	Attributes ApprovalAttributes `json:"attributes"`
}

type ApprovalAttributes struct {
	Owner   string `json:"owner"`
	Spender string `json:"spender"`
	Amount  uint64 `json:"amount"`
}

// OwnerChangeEvent for ownership transfers
type OwnerChangeEvent struct {
	Type       string                `json:"type"`
	Attributes OwnerChangeAttributes `json:"attributes"`
}

type OwnerChangeAttributes struct {
	PreviousOwner string `json:"previousOwner"`
	NewOwner      string `json:"newOwner"`
}

// PausedEvent for pause action
type PausedEvent struct {
	Type       string           `json:"type"`
	Attributes PausedAttributes `json:"attributes"`
}

type PausedAttributes struct {
	By string `json:"by"`
}

// UnpausedEvent for unpause action
type UnpausedEvent struct {
	Type       string             `json:"type"`
	Attributes UnpausedAttributes `json:"attributes"`
}

type UnpausedAttributes struct {
	By string `json:"by"`
}
