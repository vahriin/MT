package db

import (
	"database/sql"
	"github.com/vahriin/MT/model"
	"log"
)

func (adb AppDB) getSubtransactionsById(transactionId model.Id) (*[]model.Subtransaction, error) {
	row := adb.db.QueryRow("SELECT DISTINCT source FROM app_subtransaction WHERE tr_id=$1", transactionId)
	var sourceId model.Id
	err := row.Scan(&sourceId)
	if err != nil {
		return nil, ErrNotFound
	}

	rows, err := adb.db.Query("SELECT target, sum, proportion FROM app_subtransaction WHERE tr_id=$1", transactionId)
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
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			} else {
				log.Println("getSubtransactionsById returned this message: " + err.Error())
				return nil, ErrInternal
			}
		}

		subtransaction.TransactionId = transactionId
		subtransaction.Source = sourceId

		subtransactions = append(subtransactions, subtransaction)
	}
	return &subtransactions, nil
}

func (cdb CacheDB) Difference(sourceId model.Id, targetId model.Id, groupId model.Id) (int, error) {
	row := cdb.adb.db.QueryRow(`
		SELECT SUM(app_subtransaction.sum)
		FROM app_subtransaction
		JOIN app_transaction
		ON app_subtransaction.tr_id = app_transaction.tr_id
		WHERE app_subtransaction.source=$1
		AND app_subtransaction.target=$2 AND gr_id=$3`,
		sourceId, targetId, groupId)

	var sumSource int
	err := row.Scan(&sumSource)
	if err != nil {
		log.Println("Difference returned this message: " + err.Error())
		return 0, ErrNotFound
	}

	row = cdb.adb.db.QueryRow("SELECT SUM(sum) FROM app_subtransaction WHERE source=$1 AND target=$2",
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
	INSERT INTO app_subtransaction(
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
	result, err := tx.Exec("DELETE FROM app_subtransaction WHERE tr_id=$1", transactionId)
	if err != nil {
		log.Println("deleteSubtransactionPack returned this message: " + err.Error())
		return ErrInternal
	}
	if affectedRows, _ := result.RowsAffected(); affectedRows == 0 {
		return ErrNotFound
	}
	return err
}
