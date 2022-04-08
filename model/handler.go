package model

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
)

type Handler interface {
	Register(router *httprouter.Router)
}

type handler struct {
	cache map[string]Order
}

func NewHandler(cache map[string]Order) Handler {
	return &handler{
		cache: cache,
	}
}
func (h *handler) Register(router *httprouter.Router) {
	router.GET("/", h.GetAllId)
	router.GET("/orders", h.GetById)
}

func (h *handler) GetAllId(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	tmpl, err := template.ParseFiles("../static/orders.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Page not found", 404)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Page not found", 404)
		return
	}
}

func (h *handler) GetById(w http.ResponseWriter, r *http.Request, param httprouter.Params) {

	id := r.URL.Query().Get("id")
	val, exists := h.cache[id]
	if exists {
		jsonData, _ := json.Marshal(val)
		w.Write(jsonData)
	} else {
		w.Write([]byte("Заказ не найден"))
	}

}
