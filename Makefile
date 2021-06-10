run:
	go run cmd/lecture-7-demo/main.go

lint:
	golint ./...

test:
	go test -v ./...

.PHONY: build
build: vendor-proto .generate .build

PHONY: .generate
.generate:
		mkdir -p swagger
		mkdir -p pkg/lecture-7-demo
		protoc -I vendor.protogen \
				--go_out=pkg/lecture-7-demo --go_opt=paths=import \
				--go-grpc_out=pkg/lecture-7-demo --go-grpc_opt=paths=import \
				--grpc-gateway_out=pkg/lecture-7-demo \
				--grpc-gateway_opt=logtostderr=true \
				--grpc-gateway_opt=paths=import \
				--swagger_out=allow_merge=true,merge_file_name=api:swagger \
				api/lecture-7-demo/lecture-7-demo.proto
		mv pkg/lecture-7-demo/gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo/* pkg/lecture-7-demo/
		rm -rf pkg/lecture-7-demo/gitlab.com
		mkdir -p cmd/lecture-7-demo
		cd pkg/lecture-7-demo && ls go.mod || go mod init gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo && go mod tidy

PHONY: generate
generate: .vendor-proto .generate

PHONY: .build
.build:
		go build -o cmd/lecture-7-demo cmd/lecture-7-demo/main.go

PHONY: install
install: build .install

PHONY: .install
install:
		go install cmd/grpc-server/main.go

PHONY: vendor-proto
vendor-proto: .vendor-proto

PHONY: .vendor-proto
.vendor-proto:
		mkdir -p vendor.protogen
		mkdir -p vendor.protogen/api/lecture-7-demo
		cp api/lecture-7-demo/lecture-7-demo.proto vendor.protogen/api/lecture-7-demo/lecture-7-demo.proto
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi


.PHONY: deps
deps: install-go-deps

.PHONY: install-go-deps
install-go-deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
		ls go.mod || go mod init gitlab.com/siriusfreak/lecture-7-demo
		go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
		go get -u github.com/golang/protobuf/proto
		go get -u github.com/golang/protobuf/protoc-gen-go
		go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc