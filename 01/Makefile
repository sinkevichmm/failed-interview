gen:
	protoc --go_out=:. --go_opt=paths=source_relative \
    --go-grpc_out=:. --go-grpc_opt=paths=source_relative \
    internal/proto/*.proto

build:
	go build -o build/manager cmd/manager/manager.go
	go build -o build/server cmd/server/server.go

manager1:
	./build/manager/manager --meta=build/manager/meta1.json --limit=2 --port=8080 --auth="qweRTY"

manager2:
	./build/manager/manager --meta=build/manager/meta2.json --limit=2 --port=8081 --auth="qweRTY"

manager3:
	./build/manager/manager --meta=build/manager/meta3.json --limit=8 --port=8082 --auth="qweRTY"

server:
	./build/server/server --httpport=1234 --grpcaddress="0.0.0.0:8080 0.0.0.0:8081 0.0.0.0:8082" --auth="qweRTY"

test:
	
	go test -cover -race ./...

.PHONY: manager build server