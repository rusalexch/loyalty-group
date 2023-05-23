package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

func (h *handlers) addOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("handler > addORder > can't read body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	content := r.Header.Get(contentType)
	if content != appJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var order app.OrderGoods

	if err := json.Unmarshal(body, &order); err != nil {
		log.Println("handler > addOrder > can't unmarshal body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !order.IsValid() {
		log.Printf("handler > getOrder > order ID isn't valid: %s", order.ID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isExist, err := h.service.GetOrder(ctx, order.ID)
	if err == nil && isExist.ID == order.ID {
		w.WriteHeader(http.StatusConflict)
		return
	} else if !errors.Is(err, app.ErrOrderNotFound) {
		log.Println("handler > addOrder > can't get order by ID")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.service.AddOrder(ctx, order); err != nil {
		log.Println("handler > addOrder > can't add new order")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
