package types

import (
	"github.com/ksedleni/bookstore/protocol"
)

// Re-export locked protocol constants so chain code cannot drift.
const (
	RoyaltyBPS     = protocol.RoyaltyBPS
	ProtocolFeeBPS = protocol.ProtocolFeeBPS
)

// Msg type URLs (amino / proto-ready stubs for the MVP scaffold).
const (
	TypeMsgMintOriginal = "bookstore/MintOriginal"
	TypeMsgTransfer     = "bookstore/Transfer"
	TypeMsgP2PSale      = "bookstore/P2PSale"
)

// MsgMintOriginal creates a unique ownership record from a content hash.
type MsgMintOriginal struct {
	Creator      string `json:"creator"`
	ContentHash  string `json:"content_hash"`
	MetadataHash string `json:"metadata_hash"`
	MetadataURI  string `json:"metadata_uri"`
}

func (m MsgMintOriginal) Route() string { return RouterKey }
func (m MsgMintOriginal) Type() string  { return TypeMsgMintOriginal }

// MsgTransfer is an unpaid P2P ownership move (no royalty).
type MsgTransfer struct {
	From       string `json:"from"`
	To         string `json:"to"`
	OriginalID string `json:"original_id"`
}

func (m MsgTransfer) Route() string { return RouterKey }
func (m MsgTransfer) Type() string  { return TypeMsgTransfer }

// MsgP2PSale is an atomic paid sale: payment + 0.1% royalty + ownership transfer.
// Protocol fee is always 0%. No bookstore UI required.
type MsgP2PSale struct {
	Buyer       string `json:"buyer"`
	Seller      string `json:"seller"`
	OriginalID  string `json:"original_id"`
	PriceAmount string `json:"price_amount"`
	Denom       string `json:"denom"`
}

func (m MsgP2PSale) Route() string { return RouterKey }
func (m MsgP2PSale) Type() string  { return TypeMsgP2PSale }
