package model

import "time"

type Subtransaction struct {
	TransactionId Id
	Source        User
	Target        User
	Sum           int
	Proportion    int
}

type Transaction struct {
	Id      Id
	Date    time.Time
	Source  User
	Sum     int
	Matter  string
	Comment string
}

type InputTransaction struct {
	Source      User
	Sum         int
	Matter      string
	Comment     string
	Targets     []User
	Proportions []int
}

func (it InputTransaction) SumProportions() int {
	if it.Proportions != nil {
		var sum int
		for _, proportion := range it.Proportions {
			sum += proportion
		}
		return sum
	} else {
		return len(it.Targets)
	}
}
