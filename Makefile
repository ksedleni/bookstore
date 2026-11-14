.PHONY: tidy test build demo content-demo storefront node node-init node-start smoke-node e2e

export PATH := /usr/local/go-new/go/bin:$(PATH)
export GOTOOLCHAIN ?= local

BIN_DIR ?= .bin
HOME_DIR ?= $(HOME)/.bookstore
CHAIN_ID ?= bookstore-1
DENOM ?= stake
MONIKER ?= localvalidator
KEY ?= alice

tidy:
	go mod tidy

test:
	go test ./protocol/... ./content/... ./x/bookstore/... ./x/marketplace/...

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/bookstored ./cmd/bookstored
	go build -o $(BIN_DIR)/bookstore ./client/cmd/bookstore

demo: build
	rm -rf /tmp/bookstore-demo
	@echo "== offline CLI demo (no CometBFT) =="
	printf 'Once upon a chain\n' > /tmp/demo-book.txt
	BOOKSTORE_HOME=/tmp/bookstore-demo $(BIN_DIR)/bookstore faucet --to bob --amount 5000000
	BOOKSTORE_HOME=/tmp/bookstore-demo $(BIN_DIR)/bookstore mint --from alice --file /tmp/demo-book.txt --title Demo --author A
	BOOKSTORE_HOME=/tmp/bookstore-demo $(BIN_DIR)/bookstore sell --buyer bob --seller alice --id book-1 --price 1000000
	BOOKSTORE_HOME=/tmp/bookstore-demo $(BIN_DIR)/bookstore show --id book-1

# Offline content put/get/rewrap against the local envelope store (no node).
content-demo: build
	./scripts/content-demo.sh

# Optional reference storefront (Vite). Requires a running node with LCD :1317.
storefront:
	cd web && npm install && ulimit -n 65536 2>/dev/null; npm run dev

# Init a single-validator localnet (destructive for HOME_DIR).
node-init: build
	./scripts/localnet-init.sh

node-start: build
	$(BIN_DIR)/bookstored start --home $(HOME_DIR) --minimum-gas-prices=0$(DENOM)

# Convenience: wipe, init, then start in foreground.
node: node-init
	$(BIN_DIR)/bookstored start --home $(HOME_DIR) --minimum-gas-prices=0$(DENOM)

# Non-interactive smoke: init + start until a block commits, then exit.
smoke-node: build
	./scripts/smoke-node.sh

# Live e2e: mint → P2P sale → assert owner + 0.1% royalty / 0% protocol fee.
e2e: build
	./scripts/e2e-localnet.sh
