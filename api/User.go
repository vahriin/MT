package api

import (
	"encoding/json"
	"github.com/vahriin/MT/db"
	"github.com/vahriin/MT/model"
	"net/http"
	"strings"
)

func UserHandler(cdb *db.CacheDB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(rw, http.StatusText(http.StatusBadRequest)+"\n"+
				"This method are unsupported", http.StatusBadRequest)
			return
		}

		req.ParseForm()
		//fmt.Fprintln(rw, req.Form["id"])
		//return
		if form1, ok := req.Form["id"]; ok {
			usersId, err := parseUser(strings.Split(form1[0], ","))
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest)+"\n"+
					"\"id\" parameter not found", http.StatusBadRequest)
				return
			}

			var users []model.User
			var usersNotFound []model.Id

			for _, userId := range usersId {
				user, err := cdb.GetUserById(userId)
				if err == nil {
					users = append(users, *user)
				} else {
					usersNotFound = append(usersNotFound, userId)
				}
			}

			encoder := json.NewEncoder(rw)
			if len(usersNotFound) == 0 {
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(http.StatusOK)

				encoder.Encode(users)
			} else {
				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(http.StatusNotFound)

				encoder.Encode(usersNotFound)
			}

		} else {
			http.Error(rw, http.StatusText(http.StatusBadRequest)+"\n"+
				"\"id\" parameter not found", http.StatusBadRequest)
		}
		return
	})
}