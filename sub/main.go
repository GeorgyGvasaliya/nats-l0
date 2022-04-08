package main

import (
	"L0/model"
	"L0/psql"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
)

const (
	clusterName = "test-cluster"
	clientName  = "subscriber"
	channel     = "orders"
)

func main() {
	cache := make(map[string]model.Order)
	psql.FromDbToCash(cache)
	handler := model.NewHandler(cache)
	sc, _ := stan.Connect(clusterName, clientName)
	sub, _ := sc.Subscribe(channel, func(msg *stan.Msg) {
		var order model.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Print(err)
			return
		}
		if err != nil {
			log.Print(err)
			return
		}
		psql.OrderPutDb(&order, cache)
		cache[order.OrderUid] = order
		handler = model.NewHandler(cache)

	})
	router := httprouter.New()
	handler.Register(router)
	log.Fatal(http.ListenAndServe(":3000", router))
	sub.Unsubscribe()

}
