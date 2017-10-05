package model

import (
	"fmt"
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
	if len(it.Proportions) != 0 {
		fmt.Println(it.Proportions)
		var sum int
		for _, proportion := range it.Proportions {
			sum += proportion
		}
		return sum
	} else {
		for i := 0; i < len(it.Targets); i++ {
			it.Proportions = append(it.Proportions, 1)
		}
		fmt.Println(it.Proportions)
		return len(it.Targets)
	}
}
