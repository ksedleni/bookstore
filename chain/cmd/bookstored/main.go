package main

import (
	"fmt"
	"os"

	"github.com/ksedleni/bookstore/chain/app"
	"github.com/ksedleni/bookstore/protocol"
)

// bookstored is the sovereign appchain binary entrypoint (MVP stub).
// Full CometBFT start/init commands land in the next Cosmos SDK spike.
func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("bookstored 0.1.0-mvp (sovereign BookStore appchain scaffold)")
		fmt.Printf("protocol royalty=%d bps (0.1%%) protocol_fee=%d bps\n", protocol.RoyaltyBPS, protocol.ProtocolFeeBPS)
		fmt.Printf("app=%s\n", app.Name)
		return
	}

	fmt.Fprintf(os.Stderr, `bookstored — BookStore sovereign Cosmos SDK appchain (MVP scaffold)

Usage:
  bookstored version

This binary will grow into the full Cosmos SDK node (init, start, genesis).
Until then, use the reference CLI for local mint / P2P sale / list / buy:

  go run ./client/cmd/bookstore --help

Invariants (locked):
  - 0%% protocol/marketplace fee
  - 0.1%% creator royalty on paid sales only
  - No Ethereum / Solana / Base dependency for core trade
  - Marketplace pause does not block P2P transfer/sale or mint
`)
	os.Exit(0)
}
