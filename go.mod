module gitlab.com/siriusfreak/lecture-7-demo

go 1.15

require (
	github.com/Shopify/sarama v1.29.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.11.0
	github.com/uber/jaeger-client-go v2.29.1+incompatible
	github.com/uber/jaeger-lib v2.4.1+incompatible
	gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo v0.0.0-00010101000000-000000000000
	go.uber.org/atomic v1.8.0 // indirect
	google.golang.org/genproto v0.0.0-20210608205507-b6d2f5bf0d7d
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
)

replace gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo => ./pkg/lecture-7-demo
