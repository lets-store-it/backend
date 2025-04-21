package audit

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Brokers []string
	conn    *kafka.Conn
	writer  *kafka.Writer
}

func NewKafkaConfig(brokers []string) *KafkaConfig {
	return &KafkaConfig{
		Brokers: brokers,
	}
}

func (k *KafkaConfig) Connect(ctx context.Context, topic string) error {
	if k.writer != nil {
		return fmt.Errorf("connection already established")
	}

	k.writer = &kafka.Writer{
		Addr:     kafka.TCP(k.Brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}

	conn, err := kafka.DialLeader(ctx, "tcp", k.Brokers[0], topic, 0)
	if err != nil {
		return fmt.Errorf("failed to dial leader: %w", err)
	}
	defer conn.Close()

	k.conn = conn

	return nil
}

func (k *KafkaConfig) Close() error {
	if k.writer != nil {
		if err := k.writer.Close(); err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}
		k.writer = nil
	}
	if k.conn != nil {
		if err := k.conn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
		k.conn = nil
	}
	return nil
}

func (k *KafkaConfig) SendMessage(ctx context.Context, key, message []byte) error {
	if k.writer == nil {
		return fmt.Errorf("connection not established, call Connect first")
	}

	err := k.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   key,
			Value: message,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
