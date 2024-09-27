package tables

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUp_20240927081858(t *testing.T) {
	db := applyUpToPrev(t)

	applyNext(t, db)

	var names []string
	err := db.Select(&names, `SELECT name FROM labels`)
	require.NoError(t, err)

	require.Contains(t, names, "Fedora Linux")
}
