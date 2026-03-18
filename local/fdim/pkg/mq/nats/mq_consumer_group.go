package nats

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"fdim/pkg/log"
	"fdim/pkg/mcontext"
	"github.com/nats-io/nats.go"
)

func NewMConsumerGroupV2(ctx context.Context, conf *Config, groupID string, subjects []string, autoCommitEnable bool) (*MqConsumerGroup, error) {
	opts, err := BuildNatsOptions(conf)
	if err != nil {
		return nil, err
	}

	// Connect to NATS server
	var conn *nats.Conn
	if len(conf.Addr) > 0 {
		// Remove nats:// prefix if present, nats.Connect handles it
		url := conf.Addr[0]
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

	// Try to get JetStream context (for stream-based consumption)
	js, _ := conn.JetStream()

	mcg := &MqConsumerGroup{
		subjects: subjects,
		groupID:  groupID,
		conn:     conn,
		js:       js,
		msg:      make(chan *consumerMessage, 64),
	}
	mcg.ctx, mcg.cancel = context.WithCancel(ctx)
	go mcg.loopConsume()
	return mcg, nil
}

type consumerMessage struct {
	Msg     *nats.Msg
	Subject string
	Sub     *nats.Subscription
}

type MqConsumerGroup struct {
	subjects []string
	groupID  string
	conn     *nats.Conn
	js       nats.JetStreamContext
	ctx      context.Context
	cancel   context.CancelFunc
	msg      chan *consumerMessage
	once     sync.Once
	subs     []*nats.Subscription
}

func (x *MqConsumerGroup) closeMsgChan() {
	x.once.Do(func() {
		x.cancel()
		close(x.msg)
	})
}

func (x *MqConsumerGroup) loopConsume() {
	defer x.closeMsgChan()
	ctx := mcontext.SetOperationID(x.ctx, fmt.Sprintf("consumer_group_%s_%s_%d", strings.Join(x.subjects, "_"), x.groupID, rand.Uint32()))

	// If JetStream is available, use it for queue groups
	if x.js != nil {
		for _, subject := range x.subjects {
			// Use queue group for load balancing
			var sub *nats.Subscription
			sub, err := x.js.QueueSubscribe(subject, x.groupID, func(msg *nats.Msg) {
				select {
				case <-x.ctx.Done():
					return
				case x.msg <- &consumerMessage{Msg: msg, Subject: subject}:
				}
			}, nats.Durable(x.groupID))
			if err != nil {
				log.ZWarn(ctx, "QueueSubscribe failed, falling back to regular subscribe", err, "subject", subject, "groupID", x.groupID)
				// Fall back to regular subscription
				sub, err = x.conn.QueueSubscribe(subject, x.groupID, func(msg *nats.Msg) {
					select {
					case <-x.ctx.Done():
						return
					case x.msg <- &consumerMessage{Msg: msg, Subject: subject, Sub: sub}:
					}
				})
				if err != nil {
					log.ZWarn(ctx, "QueueSubscribe failed", err, "subject", subject)
					continue
				}
			}
			x.subs = append(x.subs, sub)
		}
	} else {
		// Use regular queue subscriptions
		for _, subject := range x.subjects {
			sub, err := x.conn.QueueSubscribe(subject, x.groupID, func(msg *nats.Msg) {
				select {
				case <-x.ctx.Done():
					return
				case x.msg <- &consumerMessage{Msg: msg, Subject: subject}:
				}
			})
			if err != nil {
				log.ZWarn(ctx, "QueueSubscribe failed", err, "subject", subject)
				continue
			}
			x.subs = append(x.subs, sub)
		}
	}

	// Keep running until context is cancelled
	<-x.ctx.Done()
}

func (x *MqConsumerGroup) Subscribe(ctx context.Context, fn func(msg Message) error) error {
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case msg, ok := <-x.msg:
		if !ok {
			return errors.New("consumer closed")
		}
		msgCtx := GetContextWithMQHeader(msg.Msg.Header)
		if err := fn(Message{ctx: msgCtx, msg: msg.Msg, js: x.js}); err != nil {
			return err
		}
		return nil
	}
}

func (x *MqConsumerGroup) Close() error {
	x.cancel()
	// Unsubscribe all subscriptions
	for _, sub := range x.subs {
		_ = sub.Unsubscribe()
	}
	x.conn.Close()
	return nil
}
