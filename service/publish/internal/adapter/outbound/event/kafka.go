package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"

	"opengeo/service/publish/internal/domain/model"
	"opengeo/service/publish/internal/port"
)

// KafkaConsumer Kafka事件消费者
type KafkaConsumer struct {
	brokers []string
	groupID string
	readers []*kafka.Reader
}

// NewKafkaConsumer 创建Kafka消费者
func NewKafkaConsumer(brokers []string, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		brokers: brokers,
		groupID: groupID,
		readers: make([]*kafka.Reader, 0),
	}
}

// SubscribeContentOptimized 订阅内容优化完成事件
func (c *KafkaConsumer) SubscribeContentOptimized(ctx context.Context, handler port.ContentOptimizedHandler) error {
	return c.subscribe(ctx, "content.optimized", func(msg kafka.Message) error {
		var event model.ContentOptimizedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}
		return handler(ctx, &event)
	})
}

// SubscribePublishRequested 订阅发布请求事件
func (c *KafkaConsumer) SubscribePublishRequested(ctx context.Context, handler port.PublishRequestedHandler) error {
	return c.subscribe(ctx, "publish.requested", func(msg kafka.Message) error {
		var event model.PublishRequestedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			return fmt.Errorf("failed to unmarshal event: %w", err)
		}
		return handler(ctx, &event)
	})
}

// subscribe 订阅主题
func (c *KafkaConsumer) subscribe(ctx context.Context, topic string, handler func(kafka.Message) error) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.brokers,
		GroupID:  c.groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	c.readers = append(c.readers, reader)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := reader.ReadMessage(ctx)
				if err != nil {
					log.Printf("Failed to read message from %s: %v", topic, err)
					continue
				}

				if err := handler(msg); err != nil {
					log.Printf("Failed to handle message from %s: %v", topic, err)
				}
			}
		}
	}()

	return nil
}

// Start 启动消费者
func (c *KafkaConsumer) Start(ctx context.Context) error {
	log.Println("Kafka consumer started")
	return nil
}

// Stop 停止消费者
func (c *KafkaConsumer) Stop() error {
	for _, reader := range c.readers {
		if err := reader.Close(); err != nil {
			log.Printf("Failed to close reader: %v", err)
		}
	}
	return nil
}

// KafkaProducer Kafka事件生产者
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer 创建Kafka生产者
func NewKafkaProducer(brokers []string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{writer: writer}
}

// PublishContentOptimized 发布内容优化事件
func (p *KafkaProducer) PublishContentOptimized(ctx context.Context, event *model.ContentOptimizedEvent) error {
	return p.publish(ctx, "content.optimized", event)
}

// PublishPublishSuccess 发布成功事件
func (p *KafkaProducer) PublishPublishSuccess(ctx context.Context, event *model.PublishSuccessEvent) error {
	return p.publish(ctx, "publish.success", event)
}

// PublishPublishFailed 发布失败事件
func (p *KafkaProducer) PublishPublishFailed(ctx context.Context, event *model.PublishFailedEvent) error {
	return p.publish(ctx, "publish.failed", event)
}

// publish 发布消息
func (p *KafkaProducer) publish(ctx context.Context, topic string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Topic: topic,
		Value: data,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Close 关闭生产者
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
