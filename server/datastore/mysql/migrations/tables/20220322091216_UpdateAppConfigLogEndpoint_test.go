package tables

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUp_20220322091216(t *testing.T) {
	db := applyUpToPrev(t) // must be done in top-level test as the migration comes from the test name
	t.Run("no entry", func(t *testing.T) {
		_, err := db.Exec(`DELETE FROM app_config_json`)
		require.NoError(t, err)

		// Apply current migration.
		applyNext(t, db)

		var count int
		err = db.Get(&count, `SELECT 1 FROM app_config_json`)
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	db = applyUpToPrev(t) // must be done in top-level test as the migration comes from the test name
	t.Run("required update", func(t *testing.T) {
		var raw string
		err := db.Get(&raw, `SELECT json_value FROM app_config_json`)
		require.NoError(t, err)
		require.Contains(t, raw, "/api/v1/osquery/log")
		require.NotContains(t, raw, "/api/latest/osquery/log")

		// Apply current migration.
		applyNext(t, db)

		err = db.Get(&raw, `SELECT json_value FROM app_config_json`)
		require.NoError(t, err)
		require.NotContains(t, raw, "/api/v1/osquery/log")
		require.Contains(t, raw, "/api/latest/osquery/log")
	})

	db = applyUpToPrev(t) // must be done in top-level test as the migration comes from the test name
	t.Run("no update required", func(t *testing.T) {
		var raw string
		err := db.Get(&raw, `SELECT json_value FROM app_config_json`)
		require.NoError(t, err)
		raw = strings.ReplaceAll(raw, "/api/v1/osquery/log", "/api/v2/osquery/log")
		_, err = db.Exec(`UPDATE app_config_json SET json_value = ? WHERE id = 1`, raw)
		require.NoError(t, err)

		// Apply current migration.
		applyNext(t, db)

		err = db.Get(&raw, `SELECT json_value FROM app_config_json`)
		require.NoError(t, err)
		require.NotContains(t, raw, "/api/latest/osquery/log")
		require.Contains(t, raw, "/api/v2/osquery/log")
	})
}
