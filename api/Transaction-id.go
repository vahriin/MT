package api

import (
	"encoding/json"
	"github.com/vahriin/MT/db"
	"net/http"
)

func TransactionIdHandler(cdb *db.CacheDB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		transactionId, err := getTransactionIdForm(req)
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

			encoder := json.NewEncoder(rw)
			if err = encoder.Encode(subtransactions); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
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
		}
		return
	})
}
