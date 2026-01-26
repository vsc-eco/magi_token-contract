package main

import (
	"magi_token/sdk"

	"github.com/CosmWasm/tinyjson/jwriter"
)

// ==========================================
// MAGI Token - Event Emission (tinyjson)
// ==========================================

// ==============
// Init Event
// ==============

func emitInit(owner string) {
	txID := sdk.GetEnvKey("tx.id")
	event := InitEvent{
		Type:       "init",
		Attributes: InitAttributes{Owner: owner},
		Tx:         *txID,
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ==============
// Mint Event
// ==============

func emitMint(to string, amount uint64) {
	txID := sdk.GetEnvKey("tx.id")
	event := MintEvent{
		Type:       "mint",
		Attributes: MintAttributes{To: to, Amount: amount},
		Tx:         *txID,
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ==============
// Burn Event
// ==============

func emitBurn(from string, amount uint64) {
	txID := sdk.GetEnvKey("tx.id")
	event := BurnEvent{
		Type:       "burn",
		Attributes: BurnAttributes{From: from, Amount: amount},
		Tx:         *txID,
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ==================
// Transfer Event
// ==================

func emitTransfer(from string, to string, amount uint64) {
	txID := sdk.GetEnvKey("tx.id")
	event := TransferEvent{
		Type:       "transfer",
		Attributes: TransferAttributes{From: from, To: to, Amount: amount},
		Tx:         *txID,
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// Approval Event
// ======================

func emitApproval(owner string, spender string, amount uint64) {
	txID := sdk.GetEnvKey("tx.id")
	event := ApprovalEvent{
		Type:       "approval",
		Attributes: ApprovalAttributes{Owner: owner, Spender: spender, Amount: amount},
		Tx:         *txID,
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// Owner Change Event
// ======================

func emitOwnerChange(newOwner string) {
	txID := sdk.GetEnvKey("tx.id")
	event := OwnerChangeEvent{
		Type:       "ownerChange",
		Attributes: OwnerChangeAttributes{NewOwner: newOwner},
		Tx:         *txID,
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// Pause Events
// ======================

func emitPaused(by string) {
	txID := sdk.GetEnvKey("tx.id")
	event := PausedEvent{
		Type:       "paused",
		Attributes: PausedAttributes{By: by},
		Tx:         *txID,
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

func emitUnpaused(by string) {
	txID := sdk.GetEnvKey("tx.id")
	event := UnpausedEvent{
		Type:       "unpaused",
		Attributes: UnpausedAttributes{By: by},
		Tx:         *txID,
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}
