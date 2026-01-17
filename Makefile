.PHONY: fmt test build check adr-check adr-check-strict adr-diff

BASE ?= origin/main

fmt:
	gofmt -w .

test:
	go test ./...

build:
	go build ./cmd/decider

adr-check:
	./decider check adr

adr-check-strict:
	./decider check adr --strict

adr-diff:
	./decider check diff --base $(BASE)

check: fmt test build adr-check-strict
