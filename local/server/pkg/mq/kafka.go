package mq

import (
	"context"

	"fdim/pkg/mq/nats"
)

// NewNatsConsumer creates a NATS consumer
func NewNatsConsumer(brokers []string, groupID string, topic string) (Consumer, error) {
	config := nats.DefaultConfig()
	config.Addr = brokers
	natsConsumer, err := nats.NewMConsumerGroupV2(context.Background(), config, groupID, []string{topic}, true)
	if err != nil {
		return nil, err
	}
	return &natsConsumerWrapper{consumer: natsConsumer}, nil
}

// NewNatsProducer creates a NATS producer
func NewNatsProducer(brokers []string, topic string) (Producer, error) {
	config := nats.DefaultConfig()
	config.Addr = brokers
	natsProducer, err := nats.NewNatsProducerV2(config, brokers, topic)
	if err != nil {
		return nil, err
	}
	return &natsProducerWrapper{producer: natsProducer}, nil
}

// Deprecated: Use NewNatsConsumer instead
func NewKafkaConsumer(brokers []string, groupID string, topic string) (Consumer, error) {
	return NewNatsConsumer(brokers, groupID, topic)
}

// Deprecated: Use NewNatsProducer instead
func NewKafkaProducer(brokers []string, topic string) (Producer, error) {
	return NewNatsProducer(brokers, topic)
}

type natsConsumerWrapper struct {
	consumer *nats.MqConsumerGroup
}

func (w *natsConsumerWrapper) Subscribe(ctx context.Context, fn Handler) error {
	return w.consumer.Subscribe(ctx, func(msg nats.Message) error {
		return fn(&natsMessageWrapper{msg: msg})
	})
}

func (w *natsConsumerWrapper) Close() error {
	return w.consumer.Close()
}

type natsProducerWrapper struct {
	producer *nats.MqProducer
}

func (w *natsProducerWrapper) SendMessage(ctx context.Context, key string, value []byte) error {
	return w.producer.SendMessage(ctx, key, value)
}

func (w *natsProducerWrapper) Close() error {
	return w.producer.Close()
}

type natsMessageWrapper struct {
	msg nats.Message
}

func (m *natsMessageWrapper) Context() context.Context { return m.msg.Context() }
func (m *natsMessageWrapper) Key() string               { return m.msg.Key() }
func (m *natsMessageWrapper) Value() []byte             { return m.msg.Value() }
func (m *natsMessageWrapper) Mark()                     { m.msg.Mark() }
func (m *natsMessageWrapper) Commit()                   { m.msg.Commit() }
