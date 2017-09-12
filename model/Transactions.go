package model

import "time"

type Subtransaction struct {
	TransactionId Id
	Source        User
	Target        User
	Sum           int
}

type TransactionDB struct {
	Id      Id
	Date    time.Time
	Source  User
	Sum     int
	Matter  string
	Comment string
}

type Transaction struct {
	TransactionDB
	Targets []User //contains only targets, without a source
}
