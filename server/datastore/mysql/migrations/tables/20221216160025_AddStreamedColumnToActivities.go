package tables

import (
	"database/sql"

	"github.com/pkg/errors"
)

func init() {
	MigrationClient.AddMigration(Up_20221216160025, Down_20221216160025)
}

func Up_20221216160025(tx *sql.Tx) error {
	if _, err := tx.Exec(
		"ALTER TABLE `activities` ADD COLUMN `streamed` TINYINT(1) NOT NULL DEFAULT FALSE;",
	); err != nil {
		return errors.Wrap(err, "adding streamed column to activities")
	}
	if _, err := tx.Exec(
		"CREATE INDEX activities_streamed_idx ON activities (id, streamed);",
	); err != nil {
		return errors.Wrap(err, "create activities_streamed_idx")
	}
	return nil
}

func Down_20221216160025(tx *sql.Tx) error {
	return nil
}
