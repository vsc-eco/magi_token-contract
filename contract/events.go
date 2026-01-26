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

func emitInit(owner, name, symbol string, decimals int, maxSupply uint64) {
	txID := sdk.GetEnvKey("tx.id")
	event := InitEvent{
		Type: "init_magi_token",
		Attributes: InitAttributes{
			Owner:     owner,
			Name:      name,
			Symbol:    symbol,
			Decimals:  decimals,
			MaxSupply: maxSupply,
		},
		Tx: *txID,
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

func emitOwnerChange(previousOwner, newOwner string) {
	txID := sdk.GetEnvKey("tx.id")
	event := OwnerChangeEvent{
		Type:       "ownerChange",
		Attributes: OwnerChangeAttributes{PreviousOwner: previousOwner, NewOwner: newOwner},
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
