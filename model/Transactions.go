package model

import (
	"time"
	"fmt"
)

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
