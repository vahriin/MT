package db

import (
	"database/sql"
	"github.com/vahriin/MT/model"
	"log"
)

/* TODO: make to return "Not Found" error for all public functions, that take model.Id */

func (cdb CacheDB) AddTransaction(transaction *model.InputTransaction) error {
	if err := cdb.adb.addTransaction(transaction); err != nil {
		return err
	}

	/* Add AddToCache here*/

	return nil
}

func (cdb CacheDB) GetTransactions(limit int, offset int) (*[]model.MainTransaction, error) {
	/* Add cache search */

	return cdb.adb.getMainTransactionsPack(limit, offset)
}

func (cdb CacheDB) GetTransactionsByUser(user *model.User, number int) (*[]model.MainTransaction, error) {
	/* Add search transaction from cache */

	return cdb.adb.getMainTransactionsByUserId(user.Id)
}

func (cdb CacheDB) GetSubtransactionsOfTransaction(transactionId model.Id) (*[]model.Subtransaction, error) {
	/* Add search subtransaction pack from cache */

	return cdb.adb.getSubtransactionsById(transactionId) //temp
}

func (cdb CacheDB) DeleteTransactionById(transactionId model.Id) error {
	err := cdb.adb.deleteTransactionById(transactionId)
	if err != nil {
		return err
	}

	/* add call of cache collector */
	return nil
}

/* ------------------------------------------------------------------------------*/

func (adb AppDB) addTransaction(inputTransaction *model.InputTransaction) error {
	sumProps := inputTransaction.SumProportions()

	addTX, err := adb.db.Begin()
	if err != nil {
		log.Println("addTransaction returned this message: " + err.Error())
		return ErrInternal
	}

	defer addTX.Commit()

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

	return nil
}

func (adb AppDB) deleteTransactionById(transactionId model.Id) error {
	delTX, err := adb.db.Begin()
	if err != nil {
		log.Println("deleteTransactionById returned this message: " + err.Error())
		return ErrInternal
	}

	defer delTX.Commit()

	err = deleteMainTransaction(delTX, transactionId)
	if err != nil {
		delTX.Rollback()
		return err
	}

	err = deleteSubtransactionsPack(delTX, transactionId)
	if err != nil {
		delTX.Rollback()
		return err
	}

	return nil
}

func (adb AppDB) getMainTransactionsPack(limit int, offset int) (*[]model.MainTransaction, error) {
	rows, err := adb.db.Query(`SELECT
		tr_id, date, source, sum, matter, comment
		FROM app_transaction ORDER BY date
		LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		log.Println("getMainTransactionsPack returned this message: " + err.Error())
		return nil, ErrInternal
	}
	defer rows.Close()

	var transactions []model.MainTransaction

	for rows.Next() {
		var transaction model.MainTransaction

		if err := rows.Scan(&transaction.Id,
			&transaction.Date,
			&transaction.Source,
			&transaction.Sum,
			&transaction.Matter,
			&transaction.Comment); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			} else {
				log.Println("getMainTransactionsPack returned this message: " + err.Error())
				return nil, ErrInternal
			}
		}

		transactions = append(transactions, transaction)
	}

	return &transactions, nil
}

func addMainTransaction(tx *sql.Tx, inputTransaction *model.InputTransaction) (model.Id, error) {
	_, err := tx.Exec(`
		INSERT INTO app_transaction(
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
		inputTransaction.Source,
		inputTransaction.Sum,
		inputTransaction.Matter,
		inputTransaction.Comment)
	if err != nil {
		log.Println("addMainTransaction returned this message: " + err.Error())
		return 0, ErrInternal
	}

	row := tx.QueryRow("SELECT MAX(tr_id) FROM app_transaction")

	var trId model.Id
	err = row.Scan(&trId)
	if err != nil {
		log.Println("addMainTransaction returned this message: " + err.Error())
		return 0, ErrInternal
	}

	return trId, nil
}

/* temp */
/* Add count of transaction */
func (adb AppDB) getMainTransactionsByUserId(sourceId model.Id) (*[]model.MainTransaction, error) {
	var transactionsOfUser []model.MainTransaction

	rows, err := adb.db.Query(`SELECT tr_id, date, sum, matter, comment
		FROM app_transaction WHERE source=$1`, sourceId)
	if err != nil {
		log.Println("getMainTransactionsByUserId returned this message: " + err.Error())
		return nil, ErrInternal
	}
	defer rows.Close()

	for rows.Next() {
		var currentTransaction model.MainTransaction

		err = rows.Scan(
			&currentTransaction.Id,
			&currentTransaction.Date,
			&currentTransaction.Sum,
			&currentTransaction.Matter,
			&currentTransaction.Comment,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			} else {
				log.Println("getMainTransactionsByUserId returned this message: " + err.Error())
				return nil, ErrInternal
			}
		}

		currentTransaction.Source = sourceId

		transactionsOfUser = append(transactionsOfUser, currentTransaction)
	}
	return &transactionsOfUser, nil
}

func deleteMainTransaction(tx *sql.Tx, id model.Id) error {
	result, err := tx.Exec("DELETE FROM app_transaction WHERE tr_id=$1", id)
	if err != nil {
		log.Println("deleteMainTransaction returned this message: " + err.Error())
		return ErrInternal
	}
	if affectedRows, _ := result.RowsAffected(); affectedRows == 0 {
		return ErrNotFound
	}
	return err
}

func roundMoney(val float64) int {
	return int(val + 0.5)
}
