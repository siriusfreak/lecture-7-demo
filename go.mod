module gitlab.com/siriusfreak/lecture-7-demo

go 1.15

require (
	github.com/Shopify/sarama v1.29.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.56.3
	google.golang.org/protobuf v1.30.0
)

replace gitlab.com/siriusfreak/lecture-7-demo/pkg/lecture-7-demo => ./pkg/lecture-7-demo
