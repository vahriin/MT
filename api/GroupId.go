package api

import (
	"github.com/vahriin/MT/db"
	"net/http"
	"encoding/json"
	"github.com/vahriin/MT/model"
)

type UserGroup struct {
	User model.Id
	Group model.Id
}

func GroupIdHandler(cdb *db.CacheDB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			groupId, err := getGroupIdForm(req)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}

			encoder := json.NewEncoder(rw)

			groups, err := cdb.GetGroupMember(groupId)
			if err != nil {
				if err == db.ErrNotFound {
					http.Error(rw, err.Error(), http.StatusNotFound)
				} else {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
				}
				return
			}

			err = encoder.Encode(groups)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)

		} else if req.Method == http.MethodPost {
			userGroup := new(UserGroup)

			decoder := json.NewDecoder(req.Body)

			if err := decoder.Decode(&userGroup); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := cdb.AddUserToGroup(userGroup.User, userGroup.Group); err != nil {
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
