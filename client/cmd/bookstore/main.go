package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ksedleni/bookstore/protocol"
)

// Reference CLI: mint / transfer / P2P sale / list / buy without any storefront.
// Persists a local ledger under --home (default ~/.bookstore-demo).
func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	home := defaultHome()
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "help", "-h", "--help":
		usage()
	case "hash":
		runHash(args)
	case "faucet":
		runFaucet(home, args)
	case "mint":
		runMint(home, args)
	case "transfer":
		runTransfer(home, args)
	case "sell":
		runSell(home, args)
	case "list":
		runList(home, args)
	case "buy":
		runBuy(home, args)
	case "show":
		runShow(home, args)
	case "events":
		runEvents(home)
	case "balance":
		runBalance(home, args)
	case "pause-market":
		runPause(home, args)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `bookstore — reference CLI (no storefront required)

BookStore is a fee-free digital-book ownership protocol. Anyone can mint,
transfer, or P2P-sell with keys + this client. Optional listings exist for
discovery; third-party bookstores are optional UX over the same ledger.

Commands:
  hash <file>                         SHA-256 content commitment
  faucet --to ADDR --amount N         Demo native balance (local ledger only)
  mint --from ADDR --file PATH [--title T] [--author A]
  transfer --from ADDR --to ADDR --id ORIGINAL_ID
  sell --buyer ADDR --seller ADDR --id ORIGINAL_ID --price N
  list --from ADDR --id ORIGINAL_ID --price N [--denom ustore]
  buy --from ADDR --listing LISTING_ID
  show --id ORIGINAL_ID
  balance --addr ADDR
  events
  pause-market --paused true|false    Narrow marketplace pause only

Royalty: fixed 0.1%% on paid sales. Protocol fee: 0%%.
No Ethereum / Solana / Base RPCs are used.
`)
}

func defaultHome() string {
	if v := os.Getenv("BOOKSTORE_HOME"); v != "" {
		return v
	}
	h, _ := os.UserHomeDir()
	return filepath.Join(h, ".bookstore-demo")
}

func ledgerPath(home string) string {
	return filepath.Join(home, "ledger.json")
}

func load(home string) *protocol.LocalLedger {
	led := protocol.NewLocalLedger()
	path := ledgerPath(home)
	data, err := os.ReadFile(path)
	if err != nil {
		return led
	}
	var p protocol.Snapshot
	if err := json.Unmarshal(data, &p); err != nil {
		fmt.Fprintf(os.Stderr, "warn: corrupt ledger, starting fresh: %v\n", err)
		return led
	}
	led.Restore(p.Balances, p.Originals, p.Listings, p.Events, p.MarketPaused, p.NextOrig, p.NextList)
	return led
}

func save(home string, led *protocol.LocalLedger) {
	if err := os.MkdirAll(home, 0o755); err != nil {
		fatal(err)
	}
	snap := led.Snapshot()
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		fatal(err)
	}
	if err := os.WriteFile(ledgerPath(home), data, 0o644); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func runHash(args []string) {
	fs := flag.NewFlagSet("hash", flag.ExitOnError)
	_ = fs.Parse(args)
	if fs.NArg() != 1 {
		fatal(fmt.Errorf("usage: hash <file>"))
	}
	b, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		fatal(err)
	}
	fmt.Println(protocol.ContentHashSHA256(b))
}

func runFaucet(home string, args []string) {
	fs := flag.NewFlagSet("faucet", flag.ExitOnError)
	to := fs.String("to", "", "address")
	amount := fs.String("amount", "", "amount")
	_ = fs.Parse(args)
	led := load(home)
	if err := led.Faucet(*to, *amount); err != nil {
		fatal(err)
	}
	save(home, led)
	fmt.Printf("credited %s with %s\n", *to, *amount)
}

func runMint(home string, args []string) {
	fs := flag.NewFlagSet("mint", flag.ExitOnError)
	from := fs.String("from", "", "creator address")
	file := fs.String("file", "", "path to file bytes")
	title := fs.String("title", "", "title")
	author := fs.String("author", "", "author")
	_ = fs.Parse(args)
	b, err := os.ReadFile(*file)
	if err != nil {
		fatal(err)
	}
	hash := protocol.ContentHashSHA256(b)
	meta := protocol.Metadata{
		Title:       *title,
		Author:      *author,
		SizeBytes:   int64(len(b)),
		ContentHash: hash,
	}
	metaBytes, _ := json.Marshal(meta)
	metaHash := protocol.ContentHashSHA256(metaBytes)

	led := load(home)
	orig, err := led.Mint(*from, hash, metaHash, "")
	if err != nil {
		fatal(err)
	}
	save(home, led)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(orig)
}

func runTransfer(home string, args []string) {
	fs := flag.NewFlagSet("transfer", flag.ExitOnError)
	from := fs.String("from", "", "")
	to := fs.String("to", "", "")
	id := fs.String("id", "", "")
	_ = fs.Parse(args)
	led := load(home)
	if err := led.Transfer(*from, *to, *id); err != nil {
		fatal(err)
	}
	save(home, led)
	fmt.Println("transferred", *id, "->", *to)
}

func runSell(home string, args []string) {
	fs := flag.NewFlagSet("sell", flag.ExitOnError)
	buyer := fs.String("buyer", "", "")
	seller := fs.String("seller", "", "")
	id := fs.String("id", "", "")
	price := fs.String("price", "", "")
	_ = fs.Parse(args)
	led := load(home)
	split, err := led.P2PSale(*buyer, *seller, *id, *price)
	if err != nil {
		fatal(err)
	}
	save(home, led)
	fmt.Printf("p2p sale ok royalty=%s seller_proceeds=%s protocol_fee=%s\n",
		split.Royalty, split.SellerProceeds, split.ProtocolFee)
}

func runList(home string, args []string) {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	from := fs.String("from", "", "")
	id := fs.String("id", "", "")
	price := fs.String("price", "", "")
	denom := fs.String("denom", "ustore", "")
	_ = fs.Parse(args)
	led := load(home)
	listing, err := led.List(*from, *id, *price, *denom)
	if err != nil {
		fatal(err)
	}
	save(home, led)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(listing)
}

func runBuy(home string, args []string) {
	fs := flag.NewFlagSet("buy", flag.ExitOnError)
	from := fs.String("from", "", "buyer")
	listing := fs.String("listing", "", "")
	_ = fs.Parse(args)
	led := load(home)
	split, err := led.BuyListed(*from, *listing)
	if err != nil {
		fatal(err)
	}
	save(home, led)
	fmt.Printf("listed buy ok royalty=%s seller_proceeds=%s protocol_fee=%s\n",
		split.Royalty, split.SellerProceeds, split.ProtocolFee)
}

func runShow(home string, args []string) {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	id := fs.String("id", "", "")
	_ = fs.Parse(args)
	led := load(home)
	o, err := led.GetOriginal(*id)
	if err != nil {
		fatal(err)
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(o)
}

func runBalance(home string, args []string) {
	fs := flag.NewFlagSet("balance", flag.ExitOnError)
	addr := fs.String("addr", "", "")
	_ = fs.Parse(args)
	led := load(home)
	fmt.Println(led.Balance(*addr))
}

func runEvents(home string) {
	led := load(home)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(led.Events())
}

func runPause(home string, args []string) {
	fs := flag.NewFlagSet("pause-market", flag.ExitOnError)
	paused := fs.String("paused", "true", "true|false")
	_ = fs.Parse(args)
	led := load(home)
	led.SetMarketplacePaused(*paused == "true")
	save(home, led)
	fmt.Printf("marketplace paused=%v (P2P/mint unaffected)\n", led.MarketplacePaused())
}
