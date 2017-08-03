TEST := $(shell go list ./... | grep -v vendor)
NOW := $(shell date +'%Y-%m-%d_%T')
SHA1 := $(shell git rev-parse HEAD)

default : test vet

vet :
	go vet $(go list ./...)

test :
	go test -timeout 30s $(TEST)

build :
	go build -ldflags "-X main.sha1ver=$(SHA1) -X main.buildTime=$(NOW)" -o pg-seer

.PHONY : test vet build