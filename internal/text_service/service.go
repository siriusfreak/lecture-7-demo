package text_service

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"gitlab.com/siriusfreak/lecture-7-demo/common"
)



type Service interface {
	AddV1(context context.Context, id int64, text string,  result bool, callbackUrl string) (err error)
	CallbackFirstV1(context context.Context, id int64, result bool) (err error)
	CallbackSecondV1(context context.Context, id int64, result bool) (err error)
	StatusV1(context.Context) (map[int64]bool, error)
}

type TextService struct {
	producer 		sarama.SyncProducer
	resultInput	 	map[int64]bool
	resultOutput	map[int64]bool
}


var brokers = []string{"127.0.0.1:9094"}

func InitService() (*TextService, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)

	if err != nil {
		return nil, err
	}

	return &TextService{
		producer: producer,
		resultInput: make(map[int64]bool),
		resultOutput: make(map[int64]bool),
	}, nil
}

func prepareMessage(topic string, textMessage *common.TextMessage) (*sarama.ProducerMessage, error) {
	b, err := json.Marshal(textMessage)
	if err != nil {
		return nil, err
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.StringEncoder(string(b)),
	}

	return msg, nil
}

func (t *TextService)AddV1(context context.Context, id int64, text string,  result bool, callbackUrl string) (err error) {
	msg, err := prepareMessage("text", &common.TextMessage{
		ID:	id,
		Text: text,
		CallbackUrl: callbackUrl,
	})

	if err != nil {
		return nil
	}

	t.resultInput[id] = result

	_, _, err = t.producer.SendMessage(msg)

	return err
}

func (t *TextService)CallbackFirstV1(context context.Context, id int64, result bool) (err error) {
	t.resultOutput[id] = result

	return nil
}

func (t *TextService)CallbackSecondV1(context context.Context, id int64, result bool) (err error) {
	t.resultOutput[id] = result

	return nil
}

func (t *TextService)StatusV1(context.Context) (res map[int64]bool, err error) {
	res = make(map[int64]bool)
	for key := range t.resultInput {
		if val, ok := t.resultOutput[key]; ok {
			res[key] = t.resultInput[key] == val
		}
	}

	return res, nil
}