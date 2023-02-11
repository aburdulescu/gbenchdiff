all: build vet

build:
	go build

vet:
	go vet

fieldalignment:
	@which fieldalignment || go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	fieldalignment -test=false ./...
