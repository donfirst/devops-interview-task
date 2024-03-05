install_tools:
	brew install protobuf
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go get github.com/prometheus/client_golang/prometheus

generate_protos:
	protoc \
		--proto_path=. \
		--go_out=. \
	  	--go_opt=paths=source_relative \
	  	--go-grpc_out=. \
	  	--go-grpc_opt=paths=source_relative \
	  	pkg/muzz/*.proto
