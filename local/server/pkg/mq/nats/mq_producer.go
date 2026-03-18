package nats

import (
	"context"
	"strings"

	"github.com/nats-io/nats.go"
)

func NewNatsProducerV2(config *Config, addr []string, subject string) (*MqProducer, error) {
	opts, err := BuildNatsOptions(config)
	if err != nil {
		return nil, err
	}

	// Connect to NATS server
	var conn *nats.Conn
	if len(addr) > 0 {
		// Remove nats:// prefix if present, nats.Connect handles it
		url := addr[0]
		if strings.HasPrefix(url, "nats://") {
			url = strings.TrimPrefix(url, "nats://")
		}
		conn, err = nats.Connect(url, opts...)
	} else {
		conn, err = nats.Connect(nats.DefaultURL, opts...)
	}
	if err != nil {
		return nil, err
	}

	return &MqProducer{
		subject: subject,
		conn:    conn,
	}, nil
}

type MqProducer struct {
	subject string
	conn    *nats.Conn
}

func (x *MqProducer) SendMessage(ctx context.Context, key string, value []byte) error {
	header, err := GetMQHeaderWithContext(ctx)
	if err != nil {
		return err
	}

	msg := &nats.Msg{
		Subject: x.subject,
		Data:    value,
		Header:  header,
	}

	// Use key as reply subject if provided
	if key != "" {
		msg.Reply = key
	}

	return x.conn.PublishMsg(msg)
}

func (x *MqProducer) Close() error {
	x.conn.Close()
	return nil
}
