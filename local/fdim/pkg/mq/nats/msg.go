package nats

import (
	"context"

	"github.com/nats-io/nats.go"
)

// Message implements mq.Message interface
type Message struct {
	ctx    context.Context
	msg    *nats.Msg
	sub    *nats.Subscription
	js     nats.JetStreamContext
	stream string
}

func (m Message) Context() context.Context {
	return m.ctx
}

func (m Message) Key() string {
	// NATS doesn't have a key concept like Kafka, use subject as key
	return m.msg.Subject
}

func (m Message) Value() []byte {
	return m.msg.Data
}

func (m Message) Mark() {
	// For NATS, marking is done via Ack
	if m.js != nil && m.stream != "" {
		_ = m.msg.Ack()
	}
}

func (m Message) Commit() {
	// For NATS, commit is same as ack
	if m.js != nil && m.stream != "" {
		_ = m.msg.Ack()
	}
}
