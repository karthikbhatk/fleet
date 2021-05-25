package mysql

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/VividCortex/mysqlerr"
	"github.com/fleetdm/fleet/server/kolide"
	"github.com/go-kit/kit/log"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v3"
)

func TestSanitizeColumn(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input  string
		output string
	}{
		{"foobar-column", "foobar-column"},
		{"foobar_column", "foobar_column"},
		{"foobar;column", "foobarcolumn"},
		{"foobar#", "foobar"},
		{"foobar*baz", "foobarbaz"},
	}

	for _, tt := range testCases {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.output, sanitizeColumn(tt.input))
		})
	}
}

func TestSearchLike(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		inSQL     string
		inParams  []interface{}
		match     string
		columns   []string
		outSQL    string
		outParams []interface{}
	}{
		{
			inSQL:     "SELECT * FROM HOSTS WHERE TRUE",
			inParams:  []interface{}{},
			match:     "foobar",
			columns:   []string{"hostname"},
			outSQL:    "SELECT * FROM HOSTS WHERE TRUE AND (hostname LIKE ?)",
			outParams: []interface{}{"%foobar%"},
		},
		{
			inSQL:     "SELECT * FROM HOSTS WHERE TRUE",
			inParams:  []interface{}{3},
			match:     "foobar",
			columns:   []string{},
			outSQL:    "SELECT * FROM HOSTS WHERE TRUE",
			outParams: []interface{}{3},
		},
		{
			inSQL:     "SELECT * FROM HOSTS WHERE TRUE",
			inParams:  []interface{}{1},
			match:     "foobar",
			columns:   []string{"hostname"},
			outSQL:    "SELECT * FROM HOSTS WHERE TRUE AND (hostname LIKE ?)",
			outParams: []interface{}{1, "%foobar%"},
		},
		{
			inSQL:     "SELECT * FROM HOSTS WHERE TRUE",
			inParams:  []interface{}{1},
			match:     "foobar",
			columns:   []string{"hostname", "uuid"},
			outSQL:    "SELECT * FROM HOSTS WHERE TRUE AND (hostname LIKE ? OR uuid LIKE ?)",
			outParams: []interface{}{1, "%foobar%", "%foobar%"},
		},
		{
			inSQL:     "SELECT * FROM HOSTS WHERE TRUE",
			inParams:  []interface{}{1},
			match:     "foobar",
			columns:   []string{"hostname", "uuid"},
			outSQL:    "SELECT * FROM HOSTS WHERE TRUE AND (hostname LIKE ? OR uuid LIKE ?)",
			outParams: []interface{}{1, "%foobar%", "%foobar%"},
		},
		{
			inSQL:     "SELECT * FROM HOSTS WHERE 1=1",
			inParams:  []interface{}{1},
			match:     "forty_%",
			columns:   []string{"ipv4", "uuid"},
			outSQL:    "SELECT * FROM HOSTS WHERE 1=1 AND (ipv4 LIKE ? OR uuid LIKE ?)",
			outParams: []interface{}{1, "%forty\\_\\%%", "%forty\\_\\%%"},
		},
		{
			inSQL:     "SELECT * FROM HOSTS WHERE 1=1",
			inParams:  []interface{}{1},
			match:     "forty_%",
			columns:   []string{"ipv4", "uuid"},
			outSQL:    "SELECT * FROM HOSTS WHERE 1=1 AND (ipv4 LIKE ? OR uuid LIKE ?)",
			outParams: []interface{}{1, "%forty\\_\\%%", "%forty\\_\\%%"},
		},
	}

	for _, tt := range testCases {
		t.Run("", func(t *testing.T) {
			t.Parallel()

			sql, params := searchLike(tt.inSQL, tt.inParams, tt.match, tt.columns...)
			assert.Equal(t, tt.outSQL, sql)
			assert.Equal(t, tt.outParams, params)
		})
	}
}

func mockDatastore(t *testing.T) (sqlmock.Sqlmock, *Datastore) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	ds := &Datastore{
		db:     sqlx.NewDb(db, "sqlmock"),
		logger: log.NewNopLogger(),
	}

	return mock, ds
}

func TestWithRetryTxxSuccess(t *testing.T) {
	mock, ds := mockDatastore(t)
	defer ds.Close()

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	require.NoError(t, ds.withRetryTxx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("SELECT 1")
		return err
	}))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWithRetryTxxRollbackSuccess(t *testing.T) {
	mock, ds := mockDatastore(t)
	defer ds.Close()

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnError(errors.New("fail"))
	mock.ExpectRollback()

	require.Error(t, ds.withRetryTxx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("SELECT 1")
		return err
	}))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWithRetryTxxRollbackError(t *testing.T) {
	mock, ds := mockDatastore(t)
	defer ds.Close()

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnError(errors.New("fail"))
	mock.ExpectRollback().WillReturnError(errors.New("rollback failed"))

	require.Error(t, ds.withRetryTxx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("SELECT 1")
		return err
	}))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWithRetryTxxRetrySuccess(t *testing.T) {
	mock, ds := mockDatastore(t)
	defer ds.Close()

	mock.ExpectBegin()
	// Return a retryable error
	mock.ExpectExec("SELECT 1").WillReturnError(&mysql.MySQLError{Number: mysqlerr.ER_LOCK_DEADLOCK})
	mock.ExpectRollback()
	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	assert.NoError(t, ds.withRetryTxx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("SELECT 1")
		return err
	}))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWithRetryTxxCommitRetrySuccess(t *testing.T) {
	mock, ds := mockDatastore(t)
	defer ds.Close()

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	// Return a retryable error
	mock.ExpectCommit().WillReturnError(&mysql.MySQLError{Number: mysqlerr.ER_LOCK_DEADLOCK})
	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	assert.NoError(t, ds.withRetryTxx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("SELECT 1")
		return err
	}))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestWithRetryTxxCommitError(t *testing.T) {
	mock, ds := mockDatastore(t)
	defer ds.Close()

	mock.ExpectBegin()
	mock.ExpectExec("SELECT 1").WillReturnResult(sqlmock.NewResult(1, 1))
	// Return a retryable error
	mock.ExpectCommit().WillReturnError(errors.New("fail"))

	assert.Error(t, ds.withRetryTxx(func(tx *sqlx.Tx) error {
		_, err := tx.Exec("SELECT 1")
		return err
	}))

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAppendListOptionsToSQL(t *testing.T) {
	sql := "SELECT * FROM app_configs"
	opts := kolide.ListOptions{
		OrderKey: "name",
	}

	actual := appendListOptionsToSQL(sql, opts)
	expected := "SELECT * FROM app_configs ORDER BY name ASC LIMIT 1000000"
	if actual != expected {
		t.Error("Expected", expected, "Actual", actual)
	}

	sql = "SELECT * FROM app_configs"
	opts.OrderDirection = kolide.OrderDescending
	actual = appendListOptionsToSQL(sql, opts)
	expected = "SELECT * FROM app_configs ORDER BY name DESC LIMIT 1000000"
	if actual != expected {
		t.Error("Expected", expected, "Actual", actual)
	}

	opts = kolide.ListOptions{
		PerPage: 10,
	}

	sql = "SELECT * FROM app_configs"
	actual = appendListOptionsToSQL(sql, opts)
	expected = "SELECT * FROM app_configs LIMIT 10"
	if actual != expected {
		t.Error("Expected", expected, "Actual", actual)
	}

	sql = "SELECT * FROM app_configs"
	opts.Page = 2
	actual = appendListOptionsToSQL(sql, opts)
	expected = "SELECT * FROM app_configs LIMIT 10 OFFSET 20"
	if actual != expected {
		t.Error("Expected", expected, "Actual", actual)
	}

	opts = kolide.ListOptions{}
	sql = "SELECT * FROM app_configs"
	actual = appendListOptionsToSQL(sql, opts)
	expected = "SELECT * FROM app_configs LIMIT 1000000"

	if actual != expected {
		t.Error("Expected", expected, "Actual", actual)
	}

}

func TestWhereFilterHostsByTeams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		filter   kolide.TeamFilter
		expected string
	}{
		// No teams or global role
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{},
			},
			expected: "FALSE",
		},
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{Teams: []kolide.UserTeam{}},
			},
			expected: "FALSE",
		},

		// Global role
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{GlobalRole: null.StringFrom(kolide.RoleAdmin)},
			},
			expected: "TRUE",
		},
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{GlobalRole: null.StringFrom(kolide.RoleMaintainer)},
			},
			expected: "TRUE",
		},
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{GlobalRole: null.StringFrom(kolide.RoleObserver)},
			},
			expected: "FALSE",
		},
		{
			filter: kolide.TeamFilter{
				User:            &kolide.User{GlobalRole: null.StringFrom(kolide.RoleObserver)},
				IncludeObserver: true,
			},
			expected: "TRUE",
		},

		// Team roles
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{
					Teams: []kolide.UserTeam{
						{Role: kolide.RoleObserver, Team: kolide.Team{ID: 1}},
					},
				},
			},
			expected: "FALSE",
		},
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{
					Teams: []kolide.UserTeam{
						{Role: kolide.RoleObserver, Team: kolide.Team{ID: 1}},
					},
				},
				IncludeObserver: true,
			},
			expected: "hosts.team_id IN (1)",
		},
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{
					Teams: []kolide.UserTeam{
						{Role: kolide.RoleObserver, Team: kolide.Team{ID: 1}},
						{Role: kolide.RoleObserver, Team: kolide.Team{ID: 2}},
					},
				},
			},
			expected: "FALSE",
		},
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{
					Teams: []kolide.UserTeam{
						{Role: kolide.RoleObserver, Team: kolide.Team{ID: 1}},
						{Role: kolide.RoleMaintainer, Team: kolide.Team{ID: 2}},
					},
				},
			},
			expected: "hosts.team_id IN (2)",
		},
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{
					Teams: []kolide.UserTeam{
						{Role: kolide.RoleObserver, Team: kolide.Team{ID: 1}},
						{Role: kolide.RoleMaintainer, Team: kolide.Team{ID: 2}},
					},
				},
				IncludeObserver: true,
			},
			expected: "hosts.team_id IN (1,2)",
		},
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{
					Teams: []kolide.UserTeam{
						{Role: kolide.RoleObserver, Team: kolide.Team{ID: 1}},
						{Role: kolide.RoleMaintainer, Team: kolide.Team{ID: 2}},
						// Invalid role should be ignored
						{Role: "bad", Team: kolide.Team{ID: 37}},
					},
				},
			},
			expected: "hosts.team_id IN (2)",
		},
		{
			filter: kolide.TeamFilter{
				User: &kolide.User{
					Teams: []kolide.UserTeam{
						{Role: kolide.RoleObserver, Team: kolide.Team{ID: 1}},
						{Role: kolide.RoleMaintainer, Team: kolide.Team{ID: 2}},
						{Role: kolide.RoleAdmin, Team: kolide.Team{ID: 3}},
						// Invalid role should be ignored
					},
				},
			},
			expected: "hosts.team_id IN (2,3)",
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			ds := &Datastore{logger: log.NewNopLogger()}
			sql := ds.whereFilterHostsByTeams(tt.filter, "hosts")
			assert.Equal(t, tt.expected, sql)
		})
	}
}
