package order

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rusalexch/loyalty-group/internal/gophermart/app"
)

func (om *orderModule) process() {
	for range om.tick.C {
		om.getAccrual()
	}
}

func (om *orderModule) getAccrual() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	orders, err := om.findRegistered(ctx)
	if err != nil && !errors.Is(err, app.ErrNotFound) {
		log.Println("order > getAccrual > can't get registered orders")
	}

	for _, ord := range orders {
		res, err := http.Get(om.requestURL(ord.ID))
		if err != nil {
			log.Println("order > getAccrual > can't request order")
			log.Println(err)
			continue
		}

		if res.StatusCode != http.StatusOK {
			log.Printf("order > getAccrual > response status not OK: %d\n", res.StatusCode)
			continue
		}
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			log.Println("order > getAccrual > can't read response body")
			continue
		}
		var accrual accrual

		err = json.Unmarshal(body, &accrual)
		if err != nil {
			log.Println("order > getAccrual > can't unmarshal body")
			continue
		}
		if accrual.Status != ord.Status || *accrual.Accrual != ord.Accrual.Float64 {
			err = om.updateOrder(ctx, updateOrder{
				ID:     accrual.ID,
				Status: accrual.Status,
				Accrual: sql.NullFloat64{
					Float64: *accrual.Accrual,
					Valid:   true,
				},
			})
			if err != nil {
				log.Println("order > getAccrual > can't update order")
				log.Println(err)
				continue
			}
			if accrual.Status == "PROCESSED" && *accrual.Accrual > 0 {
				om.account.Add(ctx, accrual.ID, *accrual.Accrual)
			}
		}

	}
}

func (om *orderModule) requestURL(orderID string) string {
	return fmt.Sprintf("%s/api/orders/%s", om.accrualAddress, orderID)
}
