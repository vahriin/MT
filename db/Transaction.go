package db

import (
	"database/sql"
	"github.com/vahriin/MT/model"
)

func (cdb CacheDB) AddTransaction(transaction *model.InputTransaction) error {
	err := cdb.adb.addTransaction(transaction)
	if err != nil {
		return err
	}

	/* Add AddToCache here*/

	return nil
}

func (cdb CacheDB) GetTransactions(limit int, offset int) ([]model.Transaction, error) {
	/* Add cache search */

	return cdb.adb.getTransactionsPack(limit, offset)
}

func (cdb CacheDB) GetTransactionsByUser(user *model.User, number int) ([]model.Transaction, error) {
	/* Add search transaction from cache */

	return cdb.adb.getTransactionsByUserId(user.Id)
}

func (cdb CacheDB) GetSubtransactionsOfTransaction(transactionId model.Id) ([]model.Subtransaction, error) {
	/* Add search subtransaction pack from cache */

	return cdb.adb.getSubtransactionsByIdFromDB(transactionId) //temp
}

func (cdb CacheDB) DeleteTransactionById(transactionId model.Id) error {
	err := cdb.adb.deleteTransactionByIdFromDB(transactionId)
	if err != nil {
		return err
	}

	/* add call of cache collector */
	return nil
}

func (adb AppDB) addTransaction(inputTransaction *model.InputTransaction) error {
	sumProps := inputTransaction.SumProportions()

	addTX, err := adb.db.Begin()
	if err != nil {
		return err
	}

	transactionId, err := addMainTransaction(addTX, inputTransaction)
	if err != nil {
		addTX.Rollback()
		return err
	}

	subtransaction := new(model.Subtransaction)
	var sumSubs int

	for i, target := range inputTransaction.Targets {
		subtransaction.Proportion = inputTransaction.Proportions[i]
		subtransaction.Sum = roundMoney(
			float64(subtransaction.Proportion*inputTransaction.Sum) / float64(sumProps))
		subtransaction.Target = target
		subtransaction.Source = inputTransaction.Source
		subtransaction.TransactionId = transactionId

		err := addSubtransaction(addTX, subtransaction)
		if err != nil {
			addTX.Rollback()
			return err
		}

		sumSubs += subtransaction.Sum
	}

	subtransaction.Proportion = inputTransaction.Proportions[len(inputTransaction.Proportions)-1]
	subtransaction.Sum = inputTransaction.Sum - sumSubs
	subtransaction.Source = inputTransaction.Source
	subtransaction.Target = inputTransaction.Source //transactions "to oneself"
	subtransaction.TransactionId = transactionId

	err = addSubtransaction(addTX, subtransaction)
	if err != nil {
		addTX.Rollback()
		return err
	}

	addTX.Commit()
	return nil
}

func (adb AppDB) getTransactionsPack(limit int, offset int) ([]model.Transaction, error) {
	rows, err := adb.db.Query(`SELECT
		tr_id, date, source, sum, matter, comment
		FROM transactions ORDER BY date
		LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, err
	}

	var transactions []model.Transaction

	for rows.Next() {
		var transaction model.Transaction
		var source model.Id

		err := rows.Scan(&transaction.Id,
			&transaction.Date,
			&source,
			&transaction.Sum,
			&transaction.Matter,
			&transaction.Comment)

		if err != nil {
			return nil, err
		}

		transaction.Source, err = adb.getUserById(source)

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (adb AppDB) deleteTransactionByIdFromDB(transactionId model.Id) error {
	delTX, err := adb.db.Begin()
	if err != nil {
		return err
	}

	err = deleteTransactionById(delTX, transactionId)
	if err != nil {
		delTX.Rollback()
		return err
	}

	err = deleteSubtransactionsPack(delTX, transactionId)
	if err != nil {
		delTX.Rollback()
		return err
	}

	delTX.Commit()
	return nil
}

func (adb AppDB) deleteTransaction(transaction *model.Transaction) error {
	delTX, err := adb.db.Begin()
	if err != nil {
		return err
	}

	err = deleteTransactionById(delTX, transaction.Id)
	if err != nil {
		delTX.Rollback()
		return err
	}

	err = deleteSubtransactionsPack(delTX, transaction.Id)
	if err != nil {
		delTX.Rollback()
		return err
	}

	delTX.Commit()
	return nil
}

func addMainTransaction(tx *sql.Tx, inputTransaction *model.InputTransaction) (model.Id, error) {
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
		inputTransaction.Source.Id,
		inputTransaction.Sum,
		inputTransaction.Matter,
		inputTransaction.Comment)
	if err != nil {
		return 0, err
	}

	row := tx.QueryRow("SELECT MAX(tr_id) FROM transactions")

	var trId model.Id
	err = row.Scan(&trId)
	if err != nil {
		return 0, err
	}

	return trId, nil
}


/* temp */
/* Add count of transaction */
func (adb AppDB) getTransactionsByUserId(sourceId model.Id) ([]model.Transaction, error) {
	var transactionsOfUser []model.Transaction

	rows, err := adb.db.Query(`SELECT tr_id, date, sum, matter, comment
		FROM transactions WHERE source=$1`, sourceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var currentTransaction model.Transaction

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

		currentTransaction.Source, err = adb.getUserById(sourceId)

		transactionsOfUser = append(transactionsOfUser, currentTransaction)
	}
	return transactionsOfUser, nil
}

func deleteTransactionById(tx *sql.Tx, id model.Id) error {
	_, err := tx.Exec("DELETE FROM transactions WHERE tr_id=$1", id)
	return err
}

func roundMoney(val float64) int {
	return int(val + 0.5)
}
