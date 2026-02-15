package handlers

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type OrderEvent struct {
	OrderID   string `json:"orderId"`
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
	Timestamp int64  `json:"timestamp"`
}

type InventoryEvent struct {
	OrderID   string `json:"orderId"`
	ProductID string `json:"productId"`
	NewStock  int    `json:"newStock,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Timestamp int64  `json:"timestamp"`
}

var kafkaBroker = getEnv("KAFKA_BROKER", "kafka.default.svc.cluster.local:9092")

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func StartOrderConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		GroupID: "inventory-service",
		Topic:   "orders",
	})
	defer r.Close()

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			continue
		}

		var event OrderEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			continue
		}

		processOrder(event)
	}
}

func processOrder(event OrderEvent) {
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	stock, ok := inventory[event.ProductID]

	if !ok || stock < event.Quantity {
		publishReply("inventory-failed", InventoryEvent{
			OrderID:   event.OrderID,
			ProductID: event.ProductID,
			Reason:    "insufficient stock",
		})
		return
	}

	inventory[event.ProductID] = stock - event.Quantity

	publishReply("inventory-reserved", InventoryEvent{
		OrderID:   event.OrderID,
		ProductID: event.ProductID,
		NewStock:  inventory[event.ProductID],
	})
}

func publishReply(topic string, event InventoryEvent) {
	payload, _ := json.Marshal(event)
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBroker},
		Topic:   topic,
	})
	defer w.Close()

	w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(event.OrderID),
		Value: payload,
	})
}
