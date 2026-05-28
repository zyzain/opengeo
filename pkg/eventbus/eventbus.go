package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// EventBus 事件总线
type EventBus struct {
	producer  *kafka.Writer
	consumers []*kafka.Reader
	brokers   []string
}

// NewEventBus 创建事件总线
func NewEventBus(brokers []string) *EventBus {
	producer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &EventBus{
		producer:  producer,
		consumers: make([]*kafka.Reader, 0),
		brokers:   brokers,
	}
}

// Publish 发布事件
func (eb *EventBus) Publish(ctx context.Context, topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Value: data,
		Headers: []kafka.Header{
			{Key: "timestamp", Value: []byte(time.Now().Format(time.RFC3339))},
		},
	}

	if err := eb.producer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Subscribe 订阅事件
func (eb *EventBus) Subscribe(ctx context.Context, topic, groupID string, handler EventHandler) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  eb.brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	eb.consumers = append(eb.consumers, reader)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					log.Printf("Failed to read message: %v", err)
					continue
				}

				if err := handler(ctx, msg.Value); err != nil {
					log.Printf("Failed to handle message: %v", err)
					// TODO: 发送到死信队列
				}
			}
		}
	}()

	return nil
}

// Close 关闭事件总线
func (eb *EventBus) Close() error {
	if err := eb.producer.Close(); err != nil {
		return err
	}

	for _, consumer := range eb.consumers {
		if err := consumer.Close(); err != nil {
			return err
		}
	}

	return nil
}

// EventHandler 事件处理器
type EventHandler func(ctx context.Context, data []byte) error