test: set_gopath
	go test -v ./src/...

set_gopath:
	GOPATH=$GOPATH:$PWD:$PWD/vendor

lint:
	golint ./src/...
	#test -z "$$(golint ./src/...)"

vet:
	go vet ./src/...
