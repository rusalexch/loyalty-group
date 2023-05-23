package auth

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/rusalexch/loyalty-group/internal/gophermart/app"
)

func (am *authModule) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("auth > register > can't read body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	content := r.Header.Get(app.ContentType)
	if content != app.AppJSON {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var signup signup

	if err := json.Unmarshal(body, &signup); err != nil {
		log.Println("auth > register > can't unmarshal body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if signup.Login == "" && signup.Password == "" {
		log.Println("auth > register > login and/or password is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := am.signup(ctx, signup)
	if err != nil {
		log.Println("auth > register > can't signup")
		log.Println(err)
		if errors.Is(err, errLoginAlreadyExist) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(app.AuthHeader, token)
	w.WriteHeader(http.StatusOK)
}
