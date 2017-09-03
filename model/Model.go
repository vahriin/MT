package model

import (
	"time"
)

type User struct {
	id int64
	nick string
}

type Subtransaction struct {
	TransactionID int64
	Source User
	Target User
	Sum int
}

type Transaction struct {
	Id int64
	Datetime time.Time
	Source User
	Targets []User
	Sum int
	Matter string
	Comment string
}

