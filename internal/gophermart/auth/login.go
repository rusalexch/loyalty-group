package auth

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func (am *authModule) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("auth > login > can't read body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	content := r.Header.Get(contentType)
	if content != appJSON {
		log.Println("auth > login > wrong content type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var login login
	err = json.Unmarshal(body, &login)
	if err != nil {
		log.Println("auth > login > can't unmarshal body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if login.Login == "" && login.Password == "" {
		log.Println("auth > register > login and/or password is empty")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := am.signin(ctx, login)
	if err != nil {
		log.Println("auth > login > can't signin")
		log.Println(err)
		if errors.Is(err, errUnauthorized) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(authHeader, token)
	w.WriteHeader(http.StatusOK)
}
