package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

func (h *handlers) getOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	orderID := chi.URLParam(r, "orderID")
	ID, err := strconv.ParseInt(orderID, 10, 64)
	if err != nil {
		log.Println("handlers > getORder > can't convert orderID")
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

	order, err := h.service.GetOrder(ctx, ID)
	if err != nil {
		if errors.Is(err, app.ErrOrderNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return
		}
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
