package ml_service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"gitlab.com/siriusfreak/lecture-7-demo/common"
)

func subscribe(topic string, consumer sarama.Consumer) error {
	partitionList, err := consumer.Partitions(topic) //get all partitions on the given topic
	if err != nil {
		return err
	}
	initialOffset := sarama.OffsetOldest //get offset for the oldest message on the topic

	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				messageReceived(message)
			}
		}(pc)
	}

	return nil
}

func messageReceived(message *sarama.ConsumerMessage) {
	tracer := opentracing.GlobalTracer()

	span := tracer.StartSpan("ConsumeMessage")
	defer span.Finish()

	common.IncConsumedMessages()

	fmt.Printf("Analyzing message: %s\n", string(message.Value))
	var msg common.TextMessage
	err := json.Unmarshal(message.Value, &msg)
	if err != nil {
		fmt.Printf("Error unmarshalling message: %s\n", err)
	}
	span.SetTag("id", msg.ID)

	result := rand.Int() % 2

	req, _ := http.NewRequest("POST", msg.CallbackUrl, nil)
	req.Form = url.Values{"id": {strconv.Itoa(int(msg.ID))}, "result": {strconv.Itoa(int(result))}}
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, msg.CallbackUrl)
	ext.HTTPMethod.Set(span, "POST")

	tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("Error call callback: %v\n", err)
	} else if resp.StatusCode != 200 {
		fmt.Printf("Return code not 200: %d\n", resp.StatusCode)
	}
}

type MLService interface {
	StartConsuming() error
}

type Service struct {

}

func InitMLService() *Service {
	return &Service{}
}

var brokers = []string{"127.0.0.1:9094"}

func (s *Service) StartConsuming() error {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return err
	}

	return subscribe("text",  consumer)
}