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
			first, amount, err := getTransactionsParsedForm(req)
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
				http.Error(rw, "500 " + http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}
		} else if req.Method == http.MethodPost {
			/* TODO: Add processing of Content-Type type (block non-JSON) */
			if err := blockNoJSON(req); err != nil {
				http.Error(rw, "400 " + http.StatusText(http.StatusBadRequest) +
					"\n" + err.Error(), http.StatusBadRequest)
					return
			}

			inputTransaction := new(model.InputTransaction)

			decoder := json.NewDecoder(req.Body)
			err := decoder.Decode(inputTransaction)
			if err != nil {
				fmt.Fprintln(rw, err)
				http.Error(rw, "500 " + http.StatusText(http.StatusInternalServerError),
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
		return 0, 0, errors.New("No number in \"first\"")
	}

	amount, err := strconv.ParseInt(form2[0], 10, 32)
	if err != nil {
		return 0, 0, errors.New("No number in \"amount\"")
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
		return errors.New("No Content-Type in header")
	}
}
