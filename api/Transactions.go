package api

import "net/http"

func Transactions(rw http.ResponseWriter, req *http.Request)  {
	//temp
	if req.Method != "GET" {
		http.Error(rw, http.StatusText(405), 405)
	}
}

func getTransactions
