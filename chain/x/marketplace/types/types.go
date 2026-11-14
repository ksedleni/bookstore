package types

const (
	TypeMsgList   = "marketplace/List"
	TypeMsgDelist = "marketplace/Delist"
	TypeMsgBuy    = "marketplace/Buy"
	TypeMsgPause  = "marketplace/Pause"
)

// MsgList creates an on-ledger listing. Optional convenience over P2P sale.
type MsgList struct {
	Seller      string `json:"seller"`
	OriginalID  string `json:"original_id"`
	PriceAmount string `json:"price_amount"`
	Denom       string `json:"denom"`
}

func (m MsgList) Route() string { return RouterKey }
func (m MsgList) Type() string  { return TypeMsgList }

// MsgDelist cancels a listing.
type MsgDelist struct {
	Seller    string `json:"seller"`
	ListingID string `json:"listing_id"`
}

func (m MsgDelist) Route() string { return RouterKey }
func (m MsgDelist) Type() string  { return TypeMsgDelist }

// MsgBuy fulfills a listing (0.1% royalty, 0% protocol fee).
type MsgBuy struct {
	Buyer     string `json:"buyer"`
	ListingID string `json:"listing_id"`
}

func (m MsgBuy) Route() string { return RouterKey }
func (m MsgBuy) Type() string  { return TypeMsgBuy }

// MsgSetPaused is the narrow marketplace-only pause (does not block P2P/mint).
type MsgSetPaused struct {
	Authority string `json:"authority"`
	Paused    bool   `json:"paused"`
}

func (m MsgSetPaused) Route() string { return RouterKey }
func (m MsgSetPaused) Type() string  { return TypeMsgPause }
