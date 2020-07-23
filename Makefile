default: build test testacc

build:
	go install

test:
	go test -v ./...

testacc:
	TF_ACC=1 go test -v ./... $(TESTARGS) -timeout 120m
