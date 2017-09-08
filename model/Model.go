package model

import (
	"time"
)

type User struct {
	Id int
	Nick string
}


type PassUser struct {
	Id int
	Nick string
	PassHash []byte
}

func (pu PassUser) ToUser() (User) {
	return User{
		Id: pu.Id,
		Nick: pu.Nick,
	}
}


type Subtransaction struct {
	TransactionID int
	Source User
	Target User
	Sum int
}


type Transaction struct {
	Id int
	Datetime time.Time
	Source User
	Targets []User
	Sum int
	Matter string
	Comment string
}

