package api

import (
	"github.com/vahriin/MT/db"
	"net/http"
	"errors"
	"github.com/vahriin/MT/model"
	"strconv"
	"encoding/json"
)

func TransactionHandler(cdb *db.CacheDB) (http.Handler) {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		transactionId, err := getTransactionParsedForm(req)
		if err != nil {
			http.Error(rw, http.StatusText(http.StatusBadRequest) + "\n" +
				err.Error(), http.StatusBadRequest)
				return
		}

		/* Method Processing */
		if req.Method == http.MethodGet {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)

			subtransactions, err := cdb.GetSubtransactionsOfTransaction(transactionId)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError) + "\n",
					http.StatusInternalServerError)
			}

			encoder := json.NewEncoder(rw)
			err = encoder.Encode(subtransactions)
		} else if req.Method == http.MethodDelete {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)

			err := cdb.DeleteTransactionById(transactionId)
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError) + "\n",
					http.StatusInternalServerError)
			}
		} else {
			http.Error(rw, http.StatusText(http.StatusBadRequest) + "\n" +
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
		return model.Id(0), errors.New("No number in \"first\"")
	}

	return model.Id(int(id)), nil
}