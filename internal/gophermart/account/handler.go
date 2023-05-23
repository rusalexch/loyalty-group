package account

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/rusalexch/loyalty-group/internal/gophermart/app"
	"github.com/rusalexch/loyalty-group/internal/validator"
)

func (am *accountModule) initHandler() {
	am.mux.Get("/api/user/balance", am.balance)
	am.mux.Post("/api/user/balance/withdraw", am.withdraw)
	am.mux.Get("/api/user/withdrawals", am.withdrawals)
}

func (am *accountModule) balance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.Header.Get(app.AuthHeader)
	if token == "" {
		log.Println("account > balance > token is empty")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := am.auth.CheckToken(ctx, token)
	if err != nil {
		log.Println("account > balance > can't check token")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	balance, err := am.currentBalance(ctx, user.ID)
	if err != nil {
		log.Println("account > balance > can't get current balance")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(balance)
	if err != nil {
		log.Println("account > balance > can't marchal response")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add(app.ContentType, app.AppJSON)
	w.Write(body)
}

func (am *accountModule) withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.Header.Get(app.AuthHeader)
	if token == "" {
		log.Println("account > balance > token is empty")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := am.auth.CheckToken(ctx, token)
	if err != nil {
		log.Println("account > balance > can't check token")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Header.Get(app.ContentType) != app.AppJSON {
		log.Println("account > balance > incorrect content type")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println("account > balance > can't read body")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var withdraw withdrawRequest
	err = json.Unmarshal(body, &withdraw)
	if err != nil {
		log.Println("account > balance > can't unmarshal request")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !validator.IsValid(withdraw.OrderID) {
		log.Println("account > balance > invalid order ID")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	balance, err := am.currentBalance(ctx, user.ID)
	if err != nil && !errors.Is(err, app.ErrNotFound) {
		log.Println("account > balance > can't get current balance")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if balance.Current < withdraw.Amount {
		log.Println("account > balance > not enough funds")
		w.WriteHeader(http.StatusPaymentRequired)
		return
	}

	am.addCredit(ctx, withdraw.OrderID, withdraw.Amount)
	w.WriteHeader(http.StatusOK)
}

func (am *accountModule) withdrawals(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.Header.Get(app.AuthHeader)
	if token == "" {
		log.Println("account > balance > token is empty")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := am.auth.CheckToken(ctx, token)
	if err != nil {
		log.Println("account > balance > can't check token")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tr, err := am.userCredit(ctx, user.ID)
	if err != nil {
		if errors.Is(err, app.ErrNotFound) {
			log.Println("account > withdrawals > credits not found")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		log.Println("account > withdrawals > can't get users credit")
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(tr)
	if err != nil {
		log.Println("account > withdrawals > can't marshal response body")
		log.Println(err)
		return
	}

	w.Header().Add(app.ContentType, app.AppJSON)
	w.Write(body)
}
