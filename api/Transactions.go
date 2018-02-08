package api

import (
	"encoding/json"
	"github.com/vahriin/MT/db"
	"github.com/vahriin/MT/model"
	"net/http"
)

func TransactionsHandler(cdb *db.CacheDB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			first, amount, group, err := getTransactionsForm(req)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest)+"\n"+
					err.Error(), http.StatusBadRequest)
				return
			}

			transactions, _ := cdb.GetTransactions(amount, first, group)

			encoder := json.NewEncoder(rw)
			err = encoder.Encode(transactions)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)

		} else if req.Method == http.MethodPost {
			/*if err := blockNoJSON(req); err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest)+
					"\n"+err.Error(), http.StatusBadRequest)
				return
			}*/

			inputTransaction := new(model.InputTransaction)

			decoder := json.NewDecoder(req.Body)
			if err := decoder.Decode(inputTransaction); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := inputTransactionValidation(inputTransaction); err != nil {
				http.Error(rw, err.Error(),	http.StatusBadRequest)
				return
			}

			if err := cdb.AddTransaction(inputTransaction); err != nil {
				http.Error(rw, err.Error(),	http.StatusInternalServerError)
				return
			}

			rw.WriteHeader(http.StatusCreated)

		} else {
			http.Error(rw, http.StatusText(http.StatusBadRequest)+
				"\nSupported only POST and GET", http.StatusBadRequest)
		}
		return
	})
}


