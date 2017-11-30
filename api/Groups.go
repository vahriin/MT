package api

import (
	"github.com/vahriin/MT/db"
	"net/http"
	"encoding/json"
	"github.com/vahriin/MT/model"
)

func GroupsHandler(cdb *db.CacheDB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			id, creator, err := getGroupsForm(req)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusBadRequest)
				return
			}

			encoder := json.NewEncoder(rw)

			if !creator {
				userGroups, err := cdb.GetGroupsByUser(id)
				if err != nil {
					if err == db.ErrNotFound {
						http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					} else {
						http.Error(rw, err.Error(), http.StatusInternalServerError)
					}
					return
				}

				if err := encoder.Encode(userGroups); err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				creatorGroups, err := cdb.GetGroupsByCreator(id)
				if err != nil {
					if err == db.ErrNotFound {
						http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					} else {
						http.Error(rw, err.Error(), http.StatusInternalServerError)
					}
					return
				}
				if err := encoder.Encode(creatorGroups); err != nil {
					http.Error(rw, err.Error(), http.StatusInternalServerError)
				}
			}

		} else if req.Method == http.MethodPost {
			group := new(model.Group)

			decoder := json.NewDecoder(req.Body)

			if err := decoder.Decode(group); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := cdb.AddGroup(group); err != nil {
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
