package model

import (
	"time"
)

type MainTransaction struct {
	Id      Id
	Group   Id
	Date    time.Time
	Source  Id
	Sum     int
	Matter  string
	Comment string
}

type InputTransaction struct {
	Group       Id
	Source      Id
	Sum         int
	Matter      string
	Comment     string
	Targets     []Id
	Proportions []int
}

func (it *InputTransaction) SumProportions() int {
	var sum int
	for _, proportion := range it.Proportions {
		sum += proportion
	}
	return sum
}
