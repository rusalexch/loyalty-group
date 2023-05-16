package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
	"github.com/rusalexch/loyalty-group/internal/validator"
)

func (h *handlers) getOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	orderID := chi.URLParam(r, "ID")
	ID, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		log.Println("handlers > getOrder > can't convert orderID")
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ip := r.RemoteAddr
	if h.isLimitRequest(ip) {
		w.Header().Add(contentType, text)
		w.Header().Add(retryAfter, "60")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(fmt.Sprintf("No more than %d requests per minute allowed", h.maxRequest)))
		return
	}

	if !validator.IsValid(ID) {
		log.Printf("handler > getOrder > order ID isn't valid: %d", ID)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrder(ctx, ID)
	if err != nil {
		if errors.Is(err, app.ErrOrderNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		log.Println("handleree > getOrder > can't get order")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(contentType, appJSON)
	w.Write(body)
}

func (h *handlers) isLimitRequest(ip string) bool {
	h.mx.Lock()
	defer h.mx.Unlock()

	h.reqPerMinute[ip] += 1

	return h.reqPerMinute[ip] > h.maxRequest
}
