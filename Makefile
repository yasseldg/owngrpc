.PHONY: proto
proto: generate_grpc_code_helloworld generate_grpc_code_polygon

generate_grpc_code_helloworld:
	@echo "gRPC code generation: Helloworld"
	@protoc \
		-I ./proto/helloworld \
		--go_out ./proto/helloworld \
		--go_opt paths=source_relative \
		--go-grpc_out ./proto/helloworld \
		--go-grpc_opt paths=source_relative \
		./proto/helloworld/*.proto

generate_grpc_code_polygon:
	@echo "gRPC code generation: Polygon"
	@protoc \
		-I ./proto/polygon \
		--go_out ./proto/polygon \
		--go_opt paths=source_relative \
		--go-grpc_out ./proto/polygon \
		--go-grpc_opt paths=source_relative \
		./proto/polygon/*.proto

create-certs:
	@cd certs && ./create.sh
