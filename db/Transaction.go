package db

import (
	"database/sql"
	"github.com/vahriin/MT/model"
)

func (adb *AppDB) AddTransaction(transaction *model.Transaction, proportion []int) error {
	sumProps := sumProportions(&proportion)

	var sumSubs int
	subtransaction := new(model.Subtransaction)

	addTX, err := adb.db.Begin()
	if err != nil {
		return err
	}

	for i, target := range transaction.Targets {
		subtransaction.Sum = roundMoney(
			float64(proportion[i]*transaction.Sum) / float64(sumProps))
		subtransaction.Target = target
		subtransaction.Source = transaction.Source
		subtransaction.TransactionId = transaction.Id

		err := addSubtransaction(addTX, subtransaction)
		if err != nil {
			addTX.Rollback()
			return err
		}

		sumSubs += subtransaction.Sum
	}

	subtransaction.Sum = transaction.Sum - sumSubs
	subtransaction.Source = transaction.Source
	subtransaction.Target = transaction.Source //transactions "to oneself"
	subtransaction.TransactionId = transaction.Id

	err = addSubtransaction(addTX, subtransaction)
	if err != nil {
		addTX.Rollback()
		return err
	}

	err = addTransactionDB(addTX, &transaction.TransactionDB)
	if err != nil {
		addTX.Rollback()
		return err
	}

	addTX.Commit()
	return nil
}

func (adb *AppDB) GetTransactionsBySource(source *model.User) ([]model.Transaction, error) {
	var transactions []model.Transaction
	transactionsDB, err := adb.getTransactionsDB(source)
	if err != nil {
		return nil, err
	}
	for _, transactionDB := range transactionsDB {
		var transaction model.Transaction
		transaction.TransactionDB = transactionDB
		transaction.Targets, err = adb.getTargetsOfTransaction(&transactionDB)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}
	return transactions, err
}

func (adb *AppDB) DeleteTransaction(transaction *model.Transaction) error {
	delTX, err := adb.db.Begin()
	if err != nil {
		return err
	}

	err = deleteTransactionsDB(delTX, &transaction.TransactionDB)
	if err != nil {
		delTX.Rollback()
		return err
	}

	err = deleteSubtransactionsPack(delTX, &transaction.TransactionDB)
	if err != nil {
		delTX.Rollback()
		return err
	}

	delTX.Commit()
	return nil
}

func addTransactionDB(tx *sql.Tx, transaction *model.TransactionDB) error {
	_, err := tx.Exec(`
	INSERT INTO transactions(
	date,
	source,
	sum,
	matter,
	comment
	) VALUES(
	LOCALTIMESTAMP(0),
	$1,
	$2,
	$3,
	$4
	);`,
		transaction.Source.Id,
		transaction.Sum,
		transaction.Matter,
		transaction.Comment)

	return err
}

func (adb *AppDB) getTransactionsDB(source *model.User) ([]model.TransactionDB, error) {
	var transactionsOfUser []model.TransactionDB

	rows, err := adb.db.Query(`SELECT tr_id, date, sum, matter, comment
		FROM transactions WHERE source=$1`, source.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var currentTransaction model.TransactionDB

		err = rows.Scan(
			&currentTransaction.Id,
			&currentTransaction.Date,
			&currentTransaction.Sum,
			&currentTransaction.Matter,
			&currentTransaction.Comment,
		)
		if err != nil {
			return nil, err
		}

		currentTransaction.Source = *source

		transactionsOfUser = append(transactionsOfUser, currentTransaction)
	}
	return transactionsOfUser, nil
}

func deleteTransactionsDB(tx *sql.Tx, transaction *model.TransactionDB) error {
	_, err := tx.Exec("DELETE FROM transactions WHERE tr_id=$1", transaction.Id)
	return err
}

func sumProportions(proportions *[]int) int {
	var sum int
	for _, proportion := range *proportions {
		sum += proportion
	}
	return sum
}

func roundMoney(val float64) int {
	return int(val + 0.5)
}
