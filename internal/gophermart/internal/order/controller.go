package order

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/rusalexch/loyalty-group/internal/gophermart/internal/app"
	"github.com/rusalexch/loyalty-group/internal/validator"
)

func (om *orderModule) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	isAuth := ctx.Value(app.UserKey)
	if isAuth == nil {
		log.Println("order > get > unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := isAuth.(app.User)

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

	isAuth := ctx.Value(app.UserKey)
	if isAuth == nil {
		log.Println("order > get > unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user := isAuth.(app.User)

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
	orderID := string(body)

	if !validator.IsValid(orderID) {
		log.Printf("order > create > invalid order ID: %s\n", orderID)
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
	if errors.Is(err, app.ErrNotFound) {
		err = om.add(ctx, user.ID, orderID)
		if err != nil {
			log.Println("order > create > can't create order")
			log.Println(err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
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

}
