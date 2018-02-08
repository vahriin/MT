package api

import (
	"github.com/vahriin/MT/db"
	"net/http"
	"bufio"
	"encoding/json"
	"github.com/vahriin/MT/model"
	"io"
)

type GoogleAuth struct {
	Iss string
	Sub string
	Azp string
	Aud string
	Iat string
	Exp string

	Email string
	EmailVerified bool //probably error
	Name string
	Picture string
	Given_Name string
	Family_Name string
	Locale string
}

func VerifyTokenHandler(cdb *db.CacheDB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "only POST supported", http.StatusBadRequest)
			return
		}

		reqReader := bufio.NewReader(req.Body)

		token, err := reqReader.ReadString('\n')
		if err != nil && err != io.EOF {
			http.Error(rw, "token cannot be read: " + err.Error() + token, http.StatusInternalServerError)
			return
		}

		resp, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + token)
		if err != nil {
			http.Error(rw, "failed to connect with google.com: " + err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(resp.Body)

		googleAuth := new(GoogleAuth)

		if err := decoder.Decode(googleAuth); err != nil {
			http.Error(rw, "failed to parse JSON response: " + err.Error(), http.StatusInternalServerError)
			return
		}

		encoder := json.NewEncoder(rw)

		/*probably shitcode*/
		user, err := cdb.GetUserByGoogleId([]byte(googleAuth.Sub))

		if err == nil {
			encoder.Encode(user)
			rw.WriteHeader(http.StatusOK)
			return

		} else if err == db.ErrNotFound {
			gUser := new(model.GoogleUser)
			gUser.Id = 0
			gUser.GoogleId = []byte(googleAuth.Sub)
			gUser.Nick = googleAuth.Given_Name

			if err := cdb.AddGoogleUser(gUser); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			user, _ = cdb.GetUserByGoogleId([]byte(googleAuth.Sub))

			encoder.Encode(user)
			rw.WriteHeader(http.StatusCreated)
			return

		} else {
			http.Error(rw, "Unexpected error" + err.Error(), http.StatusInternalServerError)
			return
		}

		return
	})
}