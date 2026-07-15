.PHONY: test race vet run compare testbench fmt

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

run:
	go run ./cmd/sim

compare:
	go run ./cmd/compare -trace all

testbench:
	go run ./cmd/testbench

fmt:
	gofmt -w ./cmd ./internal
