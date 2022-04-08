package psql

import (
	"L0/model"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func OrderPutDb(order *model.Order, cache map[string]model.Order) {
	db := NewPostgresDB()
	defer db.Close()
	_, exists := (cache)[(*order).OrderUid]
	if exists {
		log.Println("OrderUid is exist in database")
		return
	}
	_, err := db.Exec("INSERT INTO public.delivery (name,phone,zip,city,address,region,email) VALUES ($1,$2,$3,$4,$5,$6,$7)",
		order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = db.Exec("INSERT INTO public.payment (transaction,request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)",
		order.Payment.Transaction, order.Payment.RequestId, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = db.Exec("INSERT INTO public.order (order_uid, track_number, entry,delivery_id,payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) VALUES ($1, $2, $3, (SELECT MAX(id_delivery) FROM public.delivery), (SELECT MAX(id_payment) FROM public.payment), $4, $5, $6, $7, $8, $9, $10, $11)",
		order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
	if err != nil {
		log.Println(err)
		return
	}
	for _, item := range order.Items {
		_, err = db.Exec("INSERT INTO public.item (chrt_id,track_number,price,rid,name,sale,size,total_price,nm_id,brand,status) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
			item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId,
			item.Brand, item.Status)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = db.Exec("INSERT INTO public.order_items (order_id, item_id) VALUES ((SELECT MAX(id_order) FROM public.order), (SELECT MAX(id_item) FROM public.item))")
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func FromDbToCash(cache map[string]model.Order) {
	db := NewPostgresDB()
	defer db.Close()

	var order model.Order
	var id_order, delivery_id, payment_id, item_id int
	var item model.Items

	rows, err := db.Query("select * from public.order")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id_order, &order.OrderUid, &order.TrackNumber, &order.Entry, &delivery_id, &payment_id, &order.Locale, &order.InternalSignature,
			&order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard)
		if err != nil {
			fmt.Println(err)
			continue
		}

		rowsDelivery, err := db.Query("select * from public.delivery where id_delivery = $1", delivery_id)
		if err != nil {
			panic(err)
		}
		rowsDelivery.Next()
		defer rowsDelivery.Close()
		err = rowsDelivery.Scan(&delivery_id, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip, &order.Delivery.City,
			&order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email)
		if err != nil {
			fmt.Println(err)
			continue
		}
		rowsPayment, err := db.Query("select * from public.payment where id_payment = $1", payment_id)
		if err != nil {
			panic(err)
		}
		rowsPayment.Next()
		defer rowsPayment.Close()
		err = rowsPayment.Scan(&payment_id, &order.Payment.Transaction, &order.Payment.RequestId, &order.Payment.Currency, &order.Payment.Provider,
			&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal, &order.Payment.CustomFee)
		if err != nil {
			fmt.Println(err)
			continue
		}

		order.Items = []model.Items{}
		rowsItems, err := db.Query("select * from public.order_items where order_id = $1", id_order)
		if err != nil {
			panic(err)
		}
		defer rowsItems.Close()

		for rowsItems.Next() {
			err := rowsItems.Scan(&id_order, &item_id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			rowsItem, err := db.Query("select * from public.item where id_item = $1", item_id)
			if err != nil {
				panic(err)
			}
			defer rowsItem.Close()
			rowsItem.Next()
			err = rowsItem.Scan(&item_id, &item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
				&item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
			if err != nil {
				fmt.Println(err)
				continue
			}
			order.Items = append(order.Items, item)
		}
		(cache)[order.OrderUid] = order
	}
}
