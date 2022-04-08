package main

import (
	Model "L0/model"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"strconv"
)

const (
	clusterName = "test-cluster"
	clientName  = "publisher"
	channel     = "orders"
)

func CreateStruct() Model.Order {
	var order = Model.Order{
		OrderUid:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: Model.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: Model.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmId:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          "1",
	}
	item := Model.Items{
		ChrtId:      9934930,
		TrackNumber: "WBILMTESTTRACK",
		Price:       453,
		Rid:         "ab4219087a764ae0btest",
		Name:        "Mascaras",
		Sale:        30,
		Size:        "0",
		TotalPrice:  317,
		NmId:        2389212,
		Brand:       "Vivienne Sabo",
		Status:      202,
	}
	order.Items = append(order.Items, item)
	return order
}

func main() {
	order := CreateStruct()
	sc, err := stan.Connect(clusterName, clientName)
	defer sc.Close()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		orderId := i
		order.OrderUid = fmt.Sprintf("b563feb7b2b84b6test" + strconv.Itoa(orderId))
		jsonData, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}
		err = sc.Publish(channel, jsonData)
		if err != nil {
			log.Fatalf("Error during publish: %v\n", err)
		}
		log.Println("Published successfully")
	}
}
