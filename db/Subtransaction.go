package db

import (
	"database/sql"
	"github.com/vahriin/MT/model"
)

func (adb *AppDB) GetSubtransactionsOfTransactions(transaction *model.Transaction) ([]model.Subtransaction, error) {
	rows, err := adb.db.Query("SELECT target, sum FROM subtransactions WHERE tr_id=$1", transaction.Id)
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
		)
		if err != nil {
			return nil, err
		}

		for _, target := range transaction.Targets {
			if target.Id == targetId {
				subtransaction.Target = target
			}
		}

		subtransaction.TransactionId = transaction.Id
		subtransaction.Source = transaction.Source

		subtransactions = append(subtransactions, subtransaction)
	}
	return subtransactions, nil
}

func (adb *AppDB) Difference(source *model.User, target *model.User) (int, error) {
	row := adb.db.QueryRow("SELECT SUM(sum) FROM subtransactions WHERE source=$1 AND target=$2",
		source.Id, target.Id)

	var sumSource int
	err := row.Scan(&sumSource)
	if err != nil {
		return 0, err
	}

	row = adb.db.QueryRow("SELECT SUM(sum) FROM subtransactions WHERE source=$1 AND target=$2",
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
	sum
	) VALUES(
	$1, $2, $3, $4)`,
		subtransaction.TransactionId,
		subtransaction.Source.Id,
		subtransaction.Target.Id,
		subtransaction.Sum,
	)

	return err
}

func (adb *AppDB) getTargetsOfTransaction(transactionDB *model.TransactionDB) ([]model.User, error) {
	rows, err := adb.db.Query(
		"SELECT target FROM subtransactions WHERE tr_id=$1 AND target != source",
		transactionDB.Id)
	var targets []model.User
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id model.Id

		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		target, err := adb.GetUserById(id)
		if err != nil {
			return nil, err
		}

		targets = append(targets, *target)
	}
	return targets, nil
}

func deleteSubtransactionsPack(tx *sql.Tx, transaction *model.TransactionDB) error {
	_, err := tx.Exec("DELETE FROM subtransactions WHERE tr_id=$1", transaction.Id)
	return err
}
