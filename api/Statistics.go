package api

import (
	"github.com/vahriin/MT/db"
	"net/http"
	"encoding/json"
)

type Sum struct {
	Sum int
}

func StatisticsHandler(cdb *db.CacheDB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			sourceId, targetId, groupId, err := getStatisticsForm(req)

			sum, err := cdb.Difference(sourceId, targetId, groupId)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}

			sum0 := Sum{Sum:sum}

			encoder := json.NewEncoder(rw)

			if err := encoder.Encode(sum0); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			rw.WriteHeader(http.StatusOK)
		} else {
			http.Error(rw, http.StatusText(http.StatusBadRequest)+
				"\nSupported only POST and GET", http.StatusBadRequest)
		}
		return
	})
}