package api

import (
	"net/http"
	"strconv"
	"github.com/vahriin/MT/db"
	"errors"
	"encoding/json"
	"github.com/vahriin/MT/model"
	"fmt"
)

func TransactionsHandler(cdb *db.CacheDB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			first, amount, err := getParseForm(req)
			if err != nil {
				http.Error(rw, "400 " +
					http.StatusText(http.StatusBadRequest) + "\n" +
						err.Error(), http.StatusBadRequest)
				return
			}

			transactions, _ := cdb.GetTransactions(amount, first)

			encoder := json.NewEncoder(rw)
			err = encoder.Encode(transactions)
			rw.Header().Set("Content-Type", "application/json")
			if err != nil {
				http.Error(rw, "500" + http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}
		} else if req.Method == http.MethodPost {
			/* TODO: Add processing of Content-Type type (block non-JSON) */

			inputTransaction := new(model.InputTransaction)

			decoder := json.NewDecoder(req.Body)
			err := decoder.Decode(inputTransaction)
			if err != nil {
				fmt.Fprintln(rw, err)
				http.Error(rw, "500" + http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}

			err = cdb.AddTransaction(inputTransaction)
			if err != nil {
				fmt.Fprintln(rw, err)
				http.Error(rw, "500" + http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}

			rw.WriteHeader(http.StatusCreated)
		} else {
			http.Error(rw, "400 " + http.StatusText(http.StatusBadRequest) +
				"\nSupported only POST and GET", http.StatusBadRequest)
				return
		}
	})
}

func getParseForm(req *http.Request) (int, int, error) {
	req.ParseForm()

	form1, ok := req.Form["first"]
	if !ok {
		return 0, 0, errors.New("\"first\" parameter not found")
	}

	form2, ok := req.Form["amount"]
	if !ok {
		return 0, 0, errors.New("\"amount\" parameter not found")
	}

	/* TODO: add error processing */
	first, _ := strconv.ParseInt(form1[0], 10, 32)
	amount, _ := strconv.ParseInt(form2[0], 10, 32)

	return int(first), int(amount), nil
}
