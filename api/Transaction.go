package api

import (
	"encoding/json"
	"errors"
	"github.com/vahriin/MT/db"
	"github.com/vahriin/MT/model"
	"net/http"
	"strconv"
)

func TransactionHandler(cdb *db.CacheDB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		transactionId, err := getTransactionParsedForm(req)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest)+"\n"+
				err.Error(), http.StatusBadRequest)
			return
		}

		/* Method Processing */
		if req.Method == http.MethodGet {
			subtransactions, err := cdb.GetSubtransactionsOfTransaction(transactionId)
			if err != nil {
				if err == db.ErrNotFound {
					http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				} else {
					http.Error(rw, http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
				return
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)

			encoder := json.NewEncoder(rw)
			err = encoder.Encode(subtransactions)
		} else if req.Method == http.MethodDelete {
			err := cdb.DeleteTransactionById(transactionId)
			if err != nil {
				if err == db.ErrNotFound {
					http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				} else {
					http.Error(rw, http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
				return
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
		} else {
			http.Error(rw, http.StatusText(http.StatusBadRequest)+"\n"+
				"This method is not supported in the current version of api", http.StatusBadRequest)
			return
		}
	})
}

func getTransactionParsedForm(req *http.Request) (model.Id, error) {
	req.ParseForm()

	form1, ok := req.Form["id"]
	if !ok {
		return model.Id(0), errors.New("\"id\" parameter not found")
	}

	id, err := strconv.ParseInt(form1[0], 10, 32)
	if err != nil {
		return model.Id(0), errors.New("No number in \"first\" ")
	}

	return model.Id(int(id)), nil
}
