mock:
	go generate -v ./...

test:
	go test -short -coverprofile=cov.out ./...

gen:
	protoc --proto_path=./proto ./proto/*.proto --go_opt=paths=source_relative --go_out=./pb --go-grpc_opt=paths=source_relative --go-grpc_out=./pb 