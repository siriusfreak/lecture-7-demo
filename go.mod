module gitlab.com/siriusfreak/lecture-7-demo

go 1.15

require (
	github.com/Shopify/sarama v1.29.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/prometheus/client_golang v1.11.0
	gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo v0.0.0-00010101000000-000000000000
	google.golang.org/genproto v0.0.0-20210608205507-b6d2f5bf0d7d
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

replace gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo => ./pkg/lecture-7-demo
