package ml_service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"

	"github.com/Shopify/sarama"
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
	fmt.Printf("Analyzing message: %s\n", string(message.Value))
	var msg common.TextMessage
	err := json.Unmarshal(message.Value, &msg)
	if err != nil {
		fmt.Printf("Error unmarshalling message: %s\n", err)
	}

	result := rand.Int() % 2

	resp, err := http.PostForm(msg.CallbackUrl,  url.Values{"id": {strconv.Itoa(int(msg.ID))}, "result": {strconv.Itoa(int(result))}})

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