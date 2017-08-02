TEST := ./...

default : test vet

vet :
	go vet $(go list ./...)

test :
	go test -timeout 30s $(TEST)

build :
	go build -o pg-seer

.PHONY : test vet