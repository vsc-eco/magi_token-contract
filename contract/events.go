package main

import (
	"magi_token/sdk"
	"math/big"

	"github.com/CosmWasm/tinyjson/jwriter"
)

// ==========================================
// MAGI Token - Event Emission (tinyjson)
// ==========================================

// ==============
// Init Event
// ==============

func emitInit(owner, name, symbol string, decimals int, maxSupply *big.Int) {
	event := InitEvent{
		Type: "init_magi_token",
		Attributes: InitAttributes{
			Owner:     owner,
			Name:      name,
			Symbol:    symbol,
			Decimals:  decimals,
			MaxSupply: maxSupply,
		},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ==================
// Transfer Event
// ==================

func emitTransfer(from string, to string, amount *big.Int) {
	event := TransferEvent{
		Type:       "transfer",
		Attributes: TransferAttributes{From: from, To: to, Amount: amount},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// Approval Event
// ======================

func emitApproval(owner string, spender string, amount *big.Int) {
	event := ApprovalEvent{
		Type:       "approval",
		Attributes: ApprovalAttributes{Owner: owner, Spender: spender, Amount: amount},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// Owner Change Event
// ======================

func emitOwnerChange(previousOwner, newOwner string) {
	event := OwnerChangeEvent{
		Type:       "ownerChange",
		Attributes: OwnerChangeAttributes{PreviousOwner: previousOwner, NewOwner: newOwner},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

// ======================
// Pause Events
// ======================

func emitPaused(by string) {
	event := PausedEvent{
		Type:       "paused",
		Attributes: PausedAttributes{By: by},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}

func emitUnpaused(by string) {
	event := UnpausedEvent{
		Type:       "unpaused",
		Attributes: UnpausedAttributes{By: by},
	}
	w := jwriter.Writer{}
	event.MarshalTinyJSON(&w)
	sdk.Log(string(w.Buffer.BuildBytes()))
}
