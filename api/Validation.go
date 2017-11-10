package api

import (
	"github.com/vahriin/MT/model"
	"errors"
)

var errWrongSource = errors.New("wrong source")
var errWrongSum = errors.New("wrong sum")
var errWrongMatter = errors.New("wrong matter")
var errEmptyTargets = errors.New("targets is empty")
var errDifferentLengths = errors.New("lengths of targets array and proportions array must be len(Targets) + 1 = len(Proportions)")


/*prevent validation*/
func inputTransactionValidation(it *model.InputTransaction) error {
	/*warning: shitcode*/
	if it.Source == 0 {
		return errWrongSource
	}
	if it.Sum <= 0 {
		return errWrongSum
	}
	if it.Matter == "" {
	 	return errWrongMatter
	}
	if len(it.Targets) == 0 {
		return errEmptyTargets
	}
	if len(it.Targets) + 1 != len(it.Proportions) {
		return errDifferentLengths
	}
	return nil
}