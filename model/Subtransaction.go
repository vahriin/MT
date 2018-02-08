package model

type Subtransaction struct {
	TransactionId Id
	Source        Id
	Target        Id
	Sum           int
	Proportion    int
}
