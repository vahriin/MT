package model

import (
	"time"
)

type Subtransaction struct {
	TransactionId Id
	Source        Id
	Target        Id
	Sum           int
	Proportion    int
}

type MainTransaction struct {
	Id      Id
	Date    time.Time
	Source  Id
	Sum     int
	Matter  string
	Comment string
}

type InputTransaction struct {
	Source      Id
	Sum         int
	Matter      string
	Comment     string
	Targets     []Id
	Proportions []int
}

// shit
func (it *InputTransaction) SumProportions() int {
	var sum int
	for _, proportion := range it.Proportions {
		sum += proportion
	}
	return sum
}
