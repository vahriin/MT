package db

import (
	"database/sql"
	"github.com/vahriin/MT/model"
)

func (adb AppDB) getSubtransactionsByIdFromDB(transactionId model.Id) ([]model.Subtransaction, error) {
	row := adb.db.QueryRow("SELECT DISTINCT source FROM subtransactions WHERE tr_id=$1", transactionId)
	var sourceId model.Id
	err := row.Scan(&sourceId)
	if err != nil {
		return nil, err
	}

	source, err := adb.getUserById(sourceId) //temp
	if err != nil {
		return nil, err
	}


	rows, err := adb.db.Query("SELECT target, sum, proportion FROM subtransactions WHERE tr_id=$1", transactionId)
	if err != nil {
		return nil, err
	}


	var subtransactions []model.Subtransaction

	for rows.Next() {
		var subtransaction model.Subtransaction
		var targetId model.Id

		err := rows.Scan(
			&targetId,
			&subtransaction.Sum,
			&subtransaction.Proportion,
		)
		if err != nil {
			return nil, err
		}

		subtransaction.Target, err = adb.getUserById(targetId)
		if err != nil {
			return nil, err
		}

		subtransaction.TransactionId = transactionId
		subtransaction.Source = source

		subtransactions = append(subtransactions, subtransaction)
	}
	return subtransactions, nil
}

func (cdb CacheDB) Difference(source *model.User, target *model.User) (int, error) {
	row := cdb.adb.db.QueryRow("SELECT SUM(sum) FROM subtransactions WHERE source=$1 AND target=$2",
		source.Id, target.Id)

	var sumSource int
	err := row.Scan(&sumSource)
	if err != nil {
		return 0, err
	}

	row = cdb.adb.db.QueryRow("SELECT SUM(sum) FROM subtransactions WHERE source=$1 AND target=$2",
		target.Id, source.Id)

	var sumTarget int
	err = row.Scan(&sumTarget)
	if err != nil {
		return 0, err
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
		subtransaction.Source.Id,
		subtransaction.Target.Id,
		subtransaction.Sum,
		subtransaction.Proportion,
	)

	return err
}

func deleteSubtransactionsPack(tx *sql.Tx, transactionId model.Id) error {
	_, err := tx.Exec("DELETE FROM subtransactions WHERE tr_id=$1", transactionId)
	return err
}
