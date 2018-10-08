package main

// consumer_example implements a consumer using the non-channel Poll() API
// to retrieventTypee messages and eventTypeents.
import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"os/signal"
	"syscall"
)

type Product struct {
	Id        int64
	ImagePath string
	Name      string
	Category  string
	Price     int
}

var db *gorm.DB
var err error

func consume() {

	broker := "localhost"
	group := "store.products"
	topic := "products"
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    broker,
		"group.id":             group,
		"session.timeout.ms":   6000,
		"default.topic.config": kafka.ConfigMap{"auto.offset.reset": "earliest"}})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.Subscribe(topic, nil)

	run := true

	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			eventType := c.Poll(100)
			if eventType == nil {
				continue
			}

			switch message := eventType.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n", message.TopicPartition, string(message.Value))
				p := Product{}
				err := json.Unmarshal(message.Value, &p)
				if err != nil {
					fmt.Println("===============Unmarshal Error===============")
					fmt.Println(err)

				}
				//fmt.Println("StrucT ---->>>> ", p.Name)

				//======================================

				if err := db.Create(&p).Error; err != nil {
					fmt.Println("Erro no INSERT", err)
				} else {
					fmt.Println("===============Stored in Postgres===============")
				}
				//======================================
				if message.Headers != nil {
					fmt.Printf("%% Headers: %v\n", message.Headers)
				}
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", message)
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% error: %v\n", message)
				run = false
			default:
				fmt.Printf("Ignored %v\n", message)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()
}

func main() {

	fmt.Println("Up!")

	db, err = gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=products password=jarvis sslmode=disable")

	if err != nil {
		panic("failed to connect database")
		panic(err)
	} else {
		fmt.Println("Succeed")
	}

	defer db.Close()

	db.AutoMigrate(&Product{})

	consume()
}
