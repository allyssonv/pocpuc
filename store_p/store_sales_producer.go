package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"net/http"
	_ "os"
)

type Order struct {
	OrderNo  string
	Name     string
	Category string
	Price    int
	Comment  string
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func Sale(w http.ResponseWriter, r *http.Request) {

	enableCors(&w)
	//====================================
	var order Order
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	//====================================
	broker := "localhost"
	topic := "sales"

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
	}

	fmt.Printf("Created Producer %v\n", p)

	b, err := json.Marshal(order) // parse to json
	if err != nil {
		fmt.Println("FAIL json parse", err)
		return
	}

	err = p.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(b),
		}, nil)

	p.Flush(15 * 1000)
}

func main() {

	http.HandleFunc("/sale", Sale)
	fmt.Println("Server running at port: 8086")
	http.ListenAndServe(":8086", nil)

}
