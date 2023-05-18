package auth

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
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

	content := r.Header.Get(contentType)
	if content != appJSON {
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

	w.Header().Add(authHeader, token)
	w.WriteHeader(http.StatusOK)
}
