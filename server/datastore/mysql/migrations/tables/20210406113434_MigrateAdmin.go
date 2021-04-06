package tables

import (
	"database/sql"

	"github.com/pkg/errors"
)

func init() {
	MigrationClient.AddMigration(Up_20210406113434, Down_20210406113434)
}

func Up_20210406113434(tx *sql.Tx) error {
	// Old admins become global admins
	query := `
		UPDATE users
			SET global_role = 'admin'
			WHERE admin = TRUE
	`
	if _, err := tx.Exec(query); err != nil {
		return errors.Wrap(err, "update admins")
	}

	// Old non-admins become global maintainers
	query = `
		UPDATE users
			SET global_role = 'maintainer'
			WHERE admin = FALSE
	`
	if _, err := tx.Exec(query); err != nil {
		return errors.Wrap(err, "update maintainers")
	}

	// Drop the old admin column
	query = `
		ALTER TABLE users
			DROP COLUMN admin
	`
	if _, err := tx.Exec(query); err != nil {
		return errors.Wrap(err, "drop old admin column")
	}

	return nil
}

func Down_20210406113434(tx *sql.Tx) error {
	return nil
}
