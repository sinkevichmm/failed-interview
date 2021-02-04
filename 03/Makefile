
clean:
	rm -rf bin/*

dependencies:
	go mod download

build: dependencies build-api

build-api: 
	go build  -o ./bin/api api/main.go

linux-binaries:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "$(LIBRARY_ENV) netgo" -installsuffix netgo -o $(BIN_DIR)/api api/main.go
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -tags "$(LIBRARY_ENV) netgo" -installsuffix netgo -o $(BIN_DIR)/search cmd/main.go

ci: dependencies test	

build-mocks:
	@go get github.com/golang/mock/gomock
	@go install github.com/golang/mock/mockgen
	@~/go/bin/mockgen -source=usecase/crawler/interface.go -destination=usecase/crawler/mock/crawler.go -package=mock

test:
	go test -tags testing ./...  -v -race