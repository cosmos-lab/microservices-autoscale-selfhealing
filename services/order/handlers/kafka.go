package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

var (
	kafkaBroker   = getEnv("KAFKA_BROKER", "kafka-broker.default.svc.cluster.local:9092")
	pendingOrders = sync.Map{}
)

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func PublishToKafka(topic, key string, payload []byte) error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBroker},
		Topic:   topic,
	})
	defer w.Close()

	return w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: payload,
	})
}

func WaitForInventoryReply(orderID string, timeout time.Duration) (*InventoryEvent, error) {
	ch := make(chan *InventoryEvent, 1)
	pendingOrders.Store(orderID, ch)
	defer pendingOrders.Delete(orderID)

	select {
	case result := <-ch:
		return result, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for inventory reply")
	}
}

func StartInventoryReplyConsumer() {
	go func() {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{kafkaBroker},
			GroupID:        "order-service",
			GroupTopics:    []string{"inventory-reserved", "inventory-failed"},
			StartOffset:    kafka.FirstOffset,
			CommitInterval: time.Second,
			MinBytes:       1,
			MaxBytes:       10e6,
		})
		defer r.Close()

		for {
			msg, err := r.FetchMessage(context.Background())
			if err != nil {
				continue
			}

			var event InventoryEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				continue
			}

			if val, ok := pendingOrders.Load(event.OrderID); ok {
				ch := val.(chan *InventoryEvent)
				ch <- &event
			}

			r.CommitMessages(context.Background(), msg)
		}
	}()
}

func consumeInventoryReplies(r *kafka.Reader) {
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		var event InventoryEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			continue
		}

		if val, ok := pendingOrders.Load(event.OrderID); ok {
			ch := val.(chan *InventoryEvent)
			ch <- &event
		}
	}
}
