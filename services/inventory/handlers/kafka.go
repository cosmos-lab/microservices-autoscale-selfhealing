package handlers

import (
	"context"
	"encoding/json"
	"log"
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

var kafkaBroker = getEnv("KAFKA_BROKER", "kafka-broker.default.svc.cluster.local:9092")

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func StartOrderConsumer() {

	log.Println("🚀 Inventory Kafka Consumer Starting...")
	log.Println("Connecting to broker:", kafkaBroker)

	for {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{kafkaBroker},
			GroupID:        "inventory-service",
			Topic:          "orders",
			StartOffset:    kafka.FirstOffset,
			CommitInterval: time.Second,
			MinBytes:       1,
			MaxBytes:       10e6,
		})

		for {
			msg, err := r.FetchMessage(context.Background())
			if err != nil {
				log.Println("❌ Kafka fetch error:", err)
				time.Sleep(2 * time.Second)
				continue
			}

			log.Println("📩 Message received from Kafka")

			var event OrderEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Println("❌ JSON unmarshal error:", err)
				continue
			}

			processOrder(event)

			if err := r.CommitMessages(context.Background(), msg); err != nil {
				log.Println("❌ Commit failed:", err)
			}
		}
	}
}

func processOrder(event OrderEvent) {
	log.Println("Processing inventory for Product:", event.ProductID)

	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()

	stock, ok := inventory[event.ProductID]

	if !ok || stock < event.Quantity {
		log.Println("❌ Insufficient stock for:", event.ProductID)

		publishReply("inventory-failed", InventoryEvent{
			OrderID:   event.OrderID,
			ProductID: event.ProductID,
			Reason:    "insufficient stock",
			Timestamp: time.Now().UnixMilli(),
		})
		return
	}

	inventory[event.ProductID] = stock - event.Quantity

	log.Println("✅ Stock updated for:", event.ProductID)

	publishReply("inventory-reserved", InventoryEvent{
		OrderID:   event.OrderID,
		ProductID: event.ProductID,
		NewStock:  inventory[event.ProductID],
		Timestamp: time.Now().UnixMilli(),
	})
}

func publishReply(topic string, event InventoryEvent) {
	log.Println("Publishing reply to topic:", topic)

	payload, _ := json.Marshal(event)
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBroker},
		Topic:   topic,
	})
	defer w.Close()

	err := w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(event.OrderID),
		Value: payload,
	})

	if err != nil {
		log.Println("❌ Failed to publish reply:", err)
	} else {
		log.Println("Reply published for Order:", event.OrderID)
	}
}
