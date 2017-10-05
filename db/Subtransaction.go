package db

import (
	"database/sql"
	"github.com/vahriin/MT/model"
	"log"
)

func (adb AppDB) getSubtransactionsById(transactionId model.Id) ([]model.Subtransaction, error) {
	row := adb.db.QueryRow("SELECT DISTINCT source FROM subtransactions WHERE tr_id=$1", transactionId)
	var sourceId model.Id
	err := row.Scan(&sourceId)
	if err != nil {
		return nil, ErrNotFound
	}

	rows, err := adb.db.Query("SELECT target, sum, proportion FROM subtransactions WHERE tr_id=$1", transactionId)
	if err != nil {
		log.Println("getSubtransactionsById returned this message: " + err.Error())
		return nil, ErrInternal
	}
	defer rows.Close()

	var subtransactions []model.Subtransaction

	for rows.Next() {
		var subtransaction model.Subtransaction

		err := rows.Scan(
			&subtransaction.Target,
			&subtransaction.Sum,
			&subtransaction.Proportion,
		)
		if err != nil {
			log.Println("getSubtransactionsById returned this message: " + err.Error())
			return nil, ErrInternal
		}

		if err != nil {
			log.Println("getSubtransactionsById returned this message: " + err.Error())
			return nil, ErrInternal
		}

		subtransaction.TransactionId = transactionId
		subtransaction.Source = sourceId

		subtransactions = append(subtransactions, subtransaction)
	}
	return subtransactions, nil
}

func (cdb CacheDB) Difference(sourceId model.Id, targetId model.Id) (int, error) {
	row := cdb.adb.db.QueryRow("SELECT SUM(sum) FROM subtransactions WHERE source=$1 AND target=$2",
		sourceId, targetId)

	var sumSource int
	err := row.Scan(&sumSource)
	if err != nil {
		log.Println("Difference returned this message: " + err.Error())
		return 0, ErrNotFound
	}

	row = cdb.adb.db.QueryRow("SELECT SUM(sum) FROM subtransactions WHERE source=$1 AND target=$2",
		targetId, sourceId)

	var sumTarget int
	err = row.Scan(&sumTarget)
	if err != nil {
		log.Println("Difference returned this message: " + err.Error())
		return 0, ErrNotFound
	}

	return sumSource - sumTarget, nil
}

func addSubtransaction(tx *sql.Tx, subtransaction *model.Subtransaction) error {
	_, err := tx.Exec(`
	INSERT INTO subtransactions(
	tr_id,
	source,
	target,
	sum,
	proportion
	) VALUES(
	$1, $2, $3, $4, $5)`,
		subtransaction.TransactionId,
		subtransaction.Source,
		subtransaction.Target,
		subtransaction.Sum,
		subtransaction.Proportion,
	)

	return err
}

func deleteSubtransactionsPack(tx *sql.Tx, transactionId model.Id) error {
	result, err := tx.Exec("DELETE FROM subtransactions WHERE tr_id=$1", transactionId)
	if err != nil {
		log.Println("deleteSubtransactionPack returned this message: " + err.Error())
	}
	if affectedRows, _ := result.RowsAffected(); affectedRows == 0 {
		return ErrNotFound
	}
	return err
}
