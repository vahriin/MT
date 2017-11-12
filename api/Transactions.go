package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vahriin/MT/db"
	"github.com/vahriin/MT/model"
	"net/http"
	"strconv"
)

func TransactionsHandler(cdb *db.CacheDB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			first, amount, err := getTransactionsParsedForm(req)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest)+"\n"+
					err.Error(), http.StatusBadRequest)
				return
			}

			transactions, _ := cdb.GetTransactions(amount, first)

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
				http.Error(rw, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				fmt.Fprintln(rw, err)
				return
			}

			if err := inputTransactionValidation(inputTransaction); err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest),
					http.StatusBadRequest)
				fmt.Fprintln(rw, err)
				return
			}

			if err := cdb.AddTransaction(inputTransaction); err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				fmt.Fprintln(rw, err)
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

func getTransactionsParsedForm(req *http.Request) (int, int, error) {
	req.ParseForm()

	form1, ok := req.Form["first"]
	if !ok {
		return 0, 0, errors.New("\"first\" parameter not found")
	}

	form2, ok := req.Form["amount"]
	if !ok {
		return 0, 0, errors.New("\"amount\" parameter not found")
	}

	first, err := strconv.ParseInt(form1[0], 10, 32)
	if err != nil {
		return 0, 0, errors.New("No number in \"first\" ")
	}

	amount, err := strconv.ParseInt(form2[0], 10, 32)
	if err != nil {
		return 0, 0, errors.New("No number in \"amount\" ")
	}

	return int(first), int(amount), nil
}

func blockNoJSON(req *http.Request) error {
	req.ParseForm()
	if cType, ok := req.Form["Content-Type"]; ok {
		if cType[0] == "application/json" {
			return nil
		} else {
			return errors.New("Content-Type is not JSON")
		}
	} else {
		return errors.New("no Content-Type in header")
	}
}
