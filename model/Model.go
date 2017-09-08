package model

import (
	"time"
)

type Id int

type User struct {
	Id Id
	Nick string
}


type PassUser struct {
	Id Id
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
	TransactionID Id
	Source User
	Target User
	Sum int
}


type Transaction struct {
	Id Id
	Datetime time.Time
	Source User
	Targets []User
	Sum int
	Matter string
	Comment string
}

