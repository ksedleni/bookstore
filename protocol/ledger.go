package protocol

import (
	"fmt"
	"math/big"
	"sync"
	"time"
)

// LocalLedger is an in-process sovereign ledger for CLI demos and tests.
// It implements mint, P2P transfer/sale, and optional listings with the
// locked fee model — no bookstore UI and no incumbent-chain RPCs.
type LocalLedger struct {
	mu         sync.RWMutex
	originals  map[string]*Original // id -> original
	byHash     map[string]string    // contentHash -> id
	listings   map[string]*Listing  // listing id -> listing
	balances   map[string]*big.Int  // address -> native balance
	events     []Event
	nextOrig   uint64
	nextList   uint64
	height     int64
	marketPaused bool
}

// NewLocalLedger creates an empty demo ledger.
func NewLocalLedger() *LocalLedger {
	return &LocalLedger{
		originals: make(map[string]*Original),
		byHash:    make(map[string]string),
		listings:  make(map[string]*Listing),
		balances:  make(map[string]*big.Int),
		events:    nil,
	}
}

func (l *LocalLedger) credit(addr string, amt *big.Int) {
	if l.balances[addr] == nil {
		l.balances[addr] = big.NewInt(0)
	}
	l.balances[addr].Add(l.balances[addr], amt)
}

func (l *LocalLedger) debit(addr string, amt *big.Int) error {
	bal := l.balances[addr]
	if bal == nil || bal.Cmp(amt) < 0 {
		return fmt.Errorf("insufficient balance for %s", addr)
	}
	bal.Sub(bal, amt)
	return nil
}

// Faucet credits native units for local demos (not a mainnet concept).
func (l *LocalLedger) Faucet(addr string, amount string) error {
	amt, ok := new(big.Int).SetString(amount, 10)
	if !ok || amt.Sign() <= 0 {
		return fmt.Errorf("invalid amount")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.credit(addr, amt)
	return nil
}

// Balance returns the native balance string for addr.
func (l *LocalLedger) Balance(addr string) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.balances[addr] == nil {
		return "0"
	}
	return l.balances[addr].String()
}

// SetMarketplacePaused toggles the narrow marketplace pause.
// P2P transfer/sale and mint remain available when paused.
func (l *LocalLedger) SetMarketplacePaused(paused bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.marketPaused = paused
}

// MarketplacePaused reports listing-module pause state.
func (l *LocalLedger) MarketplacePaused() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.marketPaused
}

func (l *LocalLedger) emit(ev Event) {
	l.height++
	ev.Height = l.height
	if ev.Timestamp.IsZero() {
		ev.Timestamp = time.Now().UTC()
	}
	l.events = append(l.events, ev)
}

// Mint creates a unique Original from contentHash. Permissionless; duplicate hash rejected.
func (l *LocalLedger) Mint(creator, contentHash, metadataHash, metadataURI string) (*Original, error) {
	if err := ValidateContentHash(contentHash); err != nil {
		return nil, err
	}
	if creator == "" {
		return nil, fmt.Errorf("creator required")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if id, exists := l.byHash[contentHash]; exists {
		return nil, fmt.Errorf("content hash already minted as %s", id)
	}

	l.nextOrig++
	id := fmt.Sprintf("book-%d", l.nextOrig)
	o := &Original{
		ID:           id,
		ContentHash:  contentHash,
		MetadataHash: metadataHash,
		MetadataURI:  metadataURI,
		Creator:      creator,
		Owner:        creator,
		RoyaltyPayee: creator,
		CreatedAt:    time.Now().UTC(),
	}
	l.originals[id] = o
	l.byHash[contentHash] = id
	l.emit(Event{
		Type:       EventMinted,
		OriginalID: id,
		Attrs: map[string]string{
			"content_hash": contentHash,
			"creator":      creator,
			"owner":        creator,
		},
	})
	return cloneOriginal(o), nil
}

// Transfer moves ownership without payment (no royalty).
func (l *LocalLedger) Transfer(from, to, originalID string) error {
	if to == "" {
		return fmt.Errorf("recipient required")
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	o, ok := l.originals[originalID]
	if !ok {
		return fmt.Errorf("original not found: %s", originalID)
	}
	if o.Owner != from {
		return fmt.Errorf("not owner")
	}
	o.Owner = to
	l.emit(Event{
		Type:       EventTransfer,
		OriginalID: originalID,
		Attrs: map[string]string{
			"from": from,
			"to":   to,
		},
	})
	return nil
}

// P2PSale settles a paid peer-to-peer sale: pay + 0.1% royalty + ownership move.
// Works with marketplace paused. No protocol fee. No bookstore required.
func (l *LocalLedger) P2PSale(buyer, seller, originalID, priceAmount string) (*SaleSplit, error) {
	price, ok := new(big.Int).SetString(priceAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid price")
	}
	split, err := SplitSale(price)
	if err != nil {
		return nil, err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	o, found := l.originals[originalID]
	if !found {
		return nil, fmt.Errorf("original not found: %s", originalID)
	}
	if o.Owner != seller {
		return nil, fmt.Errorf("seller is not owner")
	}
	if buyer == seller {
		return nil, fmt.Errorf("buyer and seller must differ")
	}
	if err := l.debit(buyer, price); err != nil {
		return nil, err
	}
	l.credit(seller, split.SellerProceeds)
	l.credit(o.RoyaltyPayee, split.Royalty)
	o.Owner = buyer

	// Clear any active listing by former seller.
	for _, listing := range l.listings {
		if listing.OriginalID == originalID && listing.Active {
			listing.Active = false
		}
	}

	l.emit(Event{
		Type:       EventSale,
		OriginalID: originalID,
		Attrs: map[string]string{
			"buyer":           buyer,
			"seller":          seller,
			"price":           price.String(),
			"royalty":         split.Royalty.String(),
			"seller_proceeds": split.SellerProceeds.String(),
			"protocol_fee":    "0",
			"path":            "p2p",
		},
	})
	return &split, nil
}

// List creates an on-ledger listing. Blocked only by marketplace pause.
func (l *LocalLedger) List(seller, originalID, priceAmount, denom string) (*Listing, error) {
	if priceAmount == "" || denom == "" {
		return nil, fmt.Errorf("price and denom required")
	}
	if _, ok := new(big.Int).SetString(priceAmount, 10); !ok {
		return nil, fmt.Errorf("invalid price")
	}

	l.mu.Lock()
	defer l.mu.Unlock()
	if l.marketPaused {
		return nil, fmt.Errorf("marketplace module paused (P2P sale/transfer still allowed)")
	}
	o, ok := l.originals[originalID]
	if !ok {
		return nil, fmt.Errorf("original not found: %s", originalID)
	}
	if o.Owner != seller {
		return nil, fmt.Errorf("not owner")
	}

	l.nextList++
	id := fmt.Sprintf("list-%d", l.nextList)
	listing := &Listing{
		ID:          id,
		OriginalID:  originalID,
		Seller:      seller,
		PriceAmount: priceAmount,
		Denom:       denom,
		Active:      true,
		CreatedAt:   time.Now().UTC(),
	}
	l.listings[id] = listing
	l.emit(Event{
		Type:       EventListed,
		OriginalID: originalID,
		Attrs: map[string]string{
			"listing_id": id,
			"seller":     seller,
			"price":      priceAmount,
			"denom":      denom,
		},
	})
	return cloneListing(listing), nil
}

// Delist cancels an active listing.
func (l *LocalLedger) Delist(seller, listingID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	listing, ok := l.listings[listingID]
	if !ok {
		return fmt.Errorf("listing not found")
	}
	if listing.Seller != seller {
		return fmt.Errorf("not listing seller")
	}
	if !listing.Active {
		return fmt.Errorf("listing already inactive")
	}
	listing.Active = false
	l.emit(Event{
		Type:       EventDelisted,
		OriginalID: listing.OriginalID,
		Attrs:      map[string]string{"listing_id": listingID},
	})
	return nil
}

// BuyListed fulfills an active listing (same royalty rules; 0% protocol fee).
func (l *LocalLedger) BuyListed(buyer, listingID string) (*SaleSplit, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.marketPaused {
		return nil, fmt.Errorf("marketplace module paused (use P2P sale instead)")
	}
	listing, ok := l.listings[listingID]
	if !ok || !listing.Active {
		return nil, fmt.Errorf("listing not active")
	}
	o, found := l.originals[listing.OriginalID]
	if !found {
		return nil, fmt.Errorf("original missing")
	}
	if o.Owner != listing.Seller {
		return nil, fmt.Errorf("seller no longer owns original")
	}
	price, ok := new(big.Int).SetString(listing.PriceAmount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid listing price")
	}
	split, err := SplitSale(price)
	if err != nil {
		return nil, err
	}
	if buyer == listing.Seller {
		return nil, fmt.Errorf("buyer and seller must differ")
	}
	if err := l.debit(buyer, price); err != nil {
		return nil, err
	}
	l.credit(listing.Seller, split.SellerProceeds)
	l.credit(o.RoyaltyPayee, split.Royalty)
	o.Owner = buyer
	listing.Active = false

	l.emit(Event{
		Type:       EventSale,
		OriginalID: listing.OriginalID,
		Attrs: map[string]string{
			"buyer":           buyer,
			"seller":          listing.Seller,
			"price":           price.String(),
			"royalty":         split.Royalty.String(),
			"seller_proceeds": split.SellerProceeds.String(),
			"protocol_fee":    "0",
			"path":            "listed",
			"listing_id":      listingID,
		},
	})
	return &split, nil
}

// GetOriginal returns a copy of an ownership record.
func (l *LocalLedger) GetOriginal(id string) (*Original, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	o, ok := l.originals[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return cloneOriginal(o), nil
}

// Events returns a copy of emitted ledger events.
func (l *LocalLedger) Events() []Event {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Event, len(l.events))
	copy(out, l.events)
	return out
}

func cloneOriginal(o *Original) *Original {
	c := *o
	return &c
}

func cloneListing(l *Listing) *Listing {
	c := *l
	return &c
}

// Snapshot is a JSON-friendly dump for the local demo CLI.
type Snapshot struct {
	Balances     map[string]string `json:"balances"`
	Originals    []*Original       `json:"originals"`
	Listings     []*Listing        `json:"listings"`
	Events       []Event           `json:"events"`
	MarketPaused bool              `json:"market_paused"`
	NextOrig     uint64            `json:"next_orig"`
	NextList     uint64            `json:"next_list"`
}

// Snapshot exports ledger state for persistence.
func (l *LocalLedger) Snapshot() Snapshot {
	l.mu.RLock()
	defer l.mu.RUnlock()
	s := Snapshot{
		Balances:     make(map[string]string, len(l.balances)),
		Originals:    make([]*Original, 0, len(l.originals)),
		Listings:     make([]*Listing, 0, len(l.listings)),
		Events:       append([]Event(nil), l.events...),
		MarketPaused: l.marketPaused,
		NextOrig:     l.nextOrig,
		NextList:     l.nextList,
	}
	for k, v := range l.balances {
		s.Balances[k] = v.String()
	}
	for _, o := range l.originals {
		s.Originals = append(s.Originals, cloneOriginal(o))
	}
	for _, listing := range l.listings {
		s.Listings = append(s.Listings, cloneListing(listing))
	}
	return s
}

// Restore loads a previously snapshotted ledger.
func (l *LocalLedger) Restore(
	balances map[string]string,
	originals []*Original,
	listings []*Listing,
	events []Event,
	marketPaused bool,
	nextOrig, nextList uint64,
) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.balances = make(map[string]*big.Int)
	for addr, amt := range balances {
		n, ok := new(big.Int).SetString(amt, 10)
		if ok {
			l.balances[addr] = n
		}
	}
	l.originals = make(map[string]*Original)
	l.byHash = make(map[string]string)
	for _, o := range originals {
		c := cloneOriginal(o)
		l.originals[c.ID] = c
		l.byHash[c.ContentHash] = c.ID
	}
	l.listings = make(map[string]*Listing)
	for _, listing := range listings {
		l.listings[listing.ID] = cloneListing(listing)
	}
	l.events = append([]Event(nil), events...)
	l.marketPaused = marketPaused
	l.nextOrig = nextOrig
	l.nextList = nextList
}
