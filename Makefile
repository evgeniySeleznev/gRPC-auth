generate-auth-api:
	mkdir -p pkg/authV1
	protoc --proto_path api/gRPC-auth-1 \
	--go_out=pkg/authV1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go.exe \
	--go-grpc_out=pkg/authV1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc.exe \
	api/gRPC-auth-1/auth.proto

migrate:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations
