package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/rusalexch/loyalty-group/internal/accrual/internal/app"
)

func (h *handlers) addReward(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	var reward app.Reward

	err = json.Unmarshal(body, &reward)
	if err != nil {
		log.Println("handler > addReward > can't unmarshal body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = h.service.AddReward(ctx, reward); err != nil {
		if errors.Is(err, app.ErrRewardAlreadyExist) {
			log.Println("handlre > addReward > reward already exist")
			w.WriteHeader(http.StatusConflict)
			return
		}
		log.Println("handlre > addReward > can't add reward")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
