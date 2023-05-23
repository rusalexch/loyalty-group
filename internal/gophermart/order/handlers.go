package order

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/rusalexch/loyalty-group/internal/gophermart/app"
	"github.com/rusalexch/loyalty-group/internal/validator"
)

func (om *orderModule) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.Header.Get(app.AuthHeader)
	if token == "" {
		log.Println("order > get > empty auth header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := om.auth.CheckToken(ctx, token)
	if err != nil {
		log.Println("order > get > can't check token")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	orders, err := om.findByUserID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		log.Println("order > get > can't get users order")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonOrders := make([]app.Order, 0, len(orders))
	for _, ord := range orders {
		jsonOrders = append(jsonOrders, dbToJSON(ord))
	}

	body, err := json.Marshal(jsonOrders)
	if err != nil {
		log.Println("order > get > can't unmarshal request")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(app.ContentType, app.AppJSON)
	w.Write(body)
}

func (om *orderModule) create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.Header.Get(app.AuthHeader)
	if token == "" {
		log.Println("order > create > empty auth header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := om.auth.CheckToken(ctx, token)
	if err != nil {
		log.Println("order > create > can't check token")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Header.Get(app.ContentType) != app.Text {
		log.Println("order > create > unsupported content type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("order > create > can't read body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	orderID, err := strconv.ParseInt(string(body), 10, 0)
	if err != nil {
		log.Printf("order > create > can't read order ID: %s\n", string(body))
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !validator.IsValid(orderID) {
		log.Printf("order > create > invalid order ID: %d\n", orderID)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	order, err := om.findByID(ctx, orderID)
	if err != nil && !errors.Is(err, app.ErrNotFound) {
		log.Println("order > create > can't get order by ID")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if order.UserID == user.ID {
		log.Println("order > create > already created by user")
		w.WriteHeader(http.StatusOK)
		return
	}
	if order.UserID != user.ID {
		log.Println("order > create > already created another user")
		w.WriteHeader(http.StatusConflict)
		return
	}

	err = om.add(ctx, user.ID, orderID)
	if err != nil {
		log.Println("order > create > can't create order")
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
