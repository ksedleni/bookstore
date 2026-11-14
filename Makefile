.PHONY: test build demo tidy

test:
	go test ./...

tidy:
	go mod tidy

build:
	mkdir -p bin
	go build -o bin/bookstore ./client/cmd/bookstore
	go build -o bin/bookstored ./chain/cmd/bookstored

demo: build
	rm -rf /tmp/bookstore-demo
	@echo "== hash =="
	printf 'Once upon a chain\n' > /tmp/demo-book.txt
	./bin/bookstore hash /tmp/demo-book.txt
	@echo "== mint (alice) =="
	BOOKSTORE_HOME=/tmp/bookstore-demo ./bin/bookstore faucet --to bob --amount 5000000
	BOOKSTORE_HOME=/tmp/bookstore-demo ./bin/bookstore mint --from alice --file /tmp/demo-book.txt --title "Demo" --author "A"
	@echo "== P2P sale alice -> bob (no storefront) =="
	BOOKSTORE_HOME=/tmp/bookstore-demo ./bin/bookstore sell --buyer bob --seller alice --id book-1 --price 1000000
	BOOKSTORE_HOME=/tmp/bookstore-demo ./bin/bookstore show --id book-1
	@echo "== events =="
	BOOKSTORE_HOME=/tmp/bookstore-demo ./bin/bookstore events
