package api

import (
	"github.com/vahriin/MT/model"
	"errors"
	"strconv"
	"net/http"
)

func getTransactionIdForm(req *http.Request) (model.Id, error) {
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

func getTransactionsForm(req *http.Request) (int, int, model.Id, error) {
	req.ParseForm()

	form1, ok := req.Form["first"]
	if !ok {
		return 0, 0, 0, errors.New("\"first\" parameter not found")
	}

	form2, ok := req.Form["amount"]
	if !ok {
		return 0, 0, 0, errors.New("\"amount\" parameter not found")
	}

	form3, ok := req.Form["group"]
	if !ok {
		return 0, 0, 0, errors.New("\"group\" parameter not found")
	}

	first, err := strconv.ParseInt(form1[0], 10, 32)
	if err != nil {
		return 0, 0, 0, errors.New("No number in \"first\" ")
	}

	amount, err := strconv.ParseInt(form2[0], 10, 32)
	if err != nil {
		return 0, 0, 0, errors.New("No number in \"amount\" ")
	}

	group, err := strconv.ParseInt(form3[0], 10, 32)
	if err != nil {
		return 0, 0, 0, errors.New("No number in \"amount\" ")
	}

	return int(first), int(amount), model.Id(int(group)), nil
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
		return errors.New("no Content-Type in header")
	}
}

func parseUser(strUser []string) ([]model.Id, error) {
	var usersId []model.Id
	for _, strUserId := range strUser {
		userId, err := strconv.ParseInt(strUserId, 10, 32)
		if err != nil {
			return nil, errors.New("number not found")
		}
		usersId = append(usersId, model.Id(userId))
	}
	return usersId, nil
}

func getGroupsForm(req *http.Request) (model.Id, bool, error) {
	req.ParseForm()
	form1, ok := req.Form["id"]
	if !ok {
		return 0, false, errors.New("\"id\" parameter not found")
	}

	id, err := strconv.ParseInt(form1[0], 10, 32)
	if err != nil {
		return 0, false, errors.New("no number in \"id\"")
	}

	form2, ok := req.Form["creator"]
	if !ok {
		return model.Id(id), false, nil
	}

	creator, err := strconv.ParseBool(form2[0])
	if err != nil {
		return 0, false, errors.New("value of \"creator\" is not bool")
	}

	return model.Id(id), creator, nil
}

func getGroupIdForm(req *http.Request) (model.Id, error) {
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