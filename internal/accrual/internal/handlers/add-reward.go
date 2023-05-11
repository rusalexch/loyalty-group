package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/common"
)

func (h *handlers) addReward(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("handler > addReward > can't read body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	content := r.Header.Get(contentType)
	if content != appJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var reward common.Reward

	err = json.Unmarshal(body, &reward)
	if err != nil {
		log.Println("handler > addReward > can't unmarshal body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isExist, err := h.service.IsRewardExist(ctx, reward.ID)
	if err != nil {
		log.Println("handler > addReward > can't get reward by ID")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if isExist {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err = h.service.AddReward(ctx, reward); err != nil {
		log.Println("handlre > addReward > can't add reward")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
