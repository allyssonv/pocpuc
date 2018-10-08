package main

import (
	"encoding/json" //new
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka" //new
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var db *gorm.DB
var err error

type Product struct {
	//gorm.Model // isso cria Id e mais 2 colunas proprias do gorm
	Id        int64
	ImagePath string
	Name      string
	Category  string
	Price     uint
}

var templ = template.Must(template.ParseGlob("./*.html"))

func sendMessageToKafka(product *Product) {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	//===============================================================
	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	//===============================================================

	topic := "products"

	b, err := json.Marshal(product) // parse to json
	if err != nil {
		fmt.Println(err)
		return
	}

	p.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(b),
		}, nil)
	p.Flush(15 * 1000)

}

func Index(w http.ResponseWriter, r *http.Request) {
	var records []Product
	if err := db.Find(&records).Error; err != nil {
		log.Println(err)
	} else {
		templ.ExecuteTemplate(w, "Index", records)
	}
}

func New(w http.ResponseWriter, r *http.Request) {

	templ.ExecuteTemplate(w, "Create", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {

	product := Product{}
	str := r.FormValue("price")
	p, _ := strconv.ParseUint(str, 10, 64)
	price := uint(p)

	if r.Method == "POST" {
		product.ImagePath = r.FormValue("imagePath")
		product.Name = r.FormValue("name")
		product.Category = r.FormValue("category")
		product.Price = price

		if err := db.Create(&product).Error; err != nil {
			log.Println(err)
		} else {
			log.Println("Created.")
			go sendMessageToKafka(&product)
		}
	}
	http.Redirect(w, r, "/", 301)
}

func Edit(w http.ResponseWriter, r *http.Request) {

	product := Product{}
	nId := r.URL.Query().Get("id")
	db.First(&product, nId)
	templ.ExecuteTemplate(w, "Update", product)
}

func Update(w http.ResponseWriter, r *http.Request) {

	product := Product{}
	nId := r.FormValue("uid")
	key, _ := strconv.ParseInt(nId, 10, 64)
	str := r.FormValue("price")
	p, _ := strconv.ParseUint(str, 10, 64)
	pricex := uint(p)

	if r.Method == "POST" {

		product.Id = key
		product.ImagePath = r.FormValue("imagePath")
		product.Name = r.FormValue("name")
		product.Category = r.FormValue("category")
		product.Price = pricex

		db.Model(&product).Where("id = ?", product.Id).Update("name", product.Name).Update("category", product.Category).Update("price", product.Price).Update("imagePath", product.ImagePath)
	}
	http.Redirect(w, r, "/", 301)
}

func PreDelete(w http.ResponseWriter, r *http.Request) {
	product := Product{}
	nId := r.URL.Query().Get("id")
	db.First(&product, nId)
	templ.ExecuteTemplate(w, "Delete", product)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	product := Product{}
	nId := r.URL.Query().Get("id")
	db.Where("id = ?", nId).Delete(&product)
	log.Println("Deleted.")
	http.Redirect(w, r, "/", 301)
}

func main() {

	db, err = gorm.Open("mysql", "root:jarvis@tcp(127.0.0.1:3306)/crud?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Product{}) // Migrate the schema

	http.HandleFunc("/", Index)
	http.HandleFunc("/create", Create)
	http.HandleFunc("/new", New)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/pdelete", PreDelete)
	http.HandleFunc("/delete", Delete)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	log.Println("Server listening on port 8080 ")
	http.ListenAndServe(":8080", nil)

}

/**
	db.Create(&Product{Name: "baz", Category: "Shirt", Price: 49})

	var records []Product
	db.Find(&records)
	fmt.Println("Rows:")
	for _, record := range records {
		fmt.Printf("%s\t", record.Name)
		fmt.Printf("%s\t", record.Category)
		fmt.Printf("%d\t\n", record.Price)
	}

**/
