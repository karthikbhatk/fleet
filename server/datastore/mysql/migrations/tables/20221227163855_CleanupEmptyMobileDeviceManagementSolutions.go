package tables

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func init() {
	MigrationClient.AddMigration(Up_20221227163855, Down_20221227163855)
}

func Up_20221227163855(tx *sql.Tx) error {
	// Fetch the IDs of the solutions with empty server_url.
	txx := sqlx.Tx{Tx: tx}
	var mdmIDs []uint
	if err := txx.Select(&mdmIDs,
		`SELECT id FROM mobile_device_management_solutions WHERE server_url = '';`,
	); err != nil {
		return errors.Wrap(err, "select mobile_device_management_solutions")
	}

	// Cleanup mobile_device_management_solutions.
	query, args, err := sqlx.In(
		"DELETE FROM mobile_device_management_solutions WHERE id IN (?)",
		mdmIDs,
	)
	if err != nil {
		return errors.Wrap(err, "sqlx.In mobile_device_management_solutions")
	}
	if _, err := txx.Exec(query, args...); err != nil {
		return errors.Wrap(err, "mobile_device_management_solutions clean up")
	}

	// Cleanup host_mdm.
	query, args, err = sqlx.In(
		"DELETE FROM host_mdm WHERE mdm_id IN (?)",
		mdmIDs,
	)
	if err != nil {
		return errors.Wrap(err, "sqlx.In host_mdm")
	}
	if _, err := txx.Exec(query, args...); err != nil {
		return errors.Wrap(err, "host_mdm clean up")
	}
	return nil
}

func Down_20221227163855(tx *sql.Tx) error {
	return nil
}
