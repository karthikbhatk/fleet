//go:build darwin
// +build darwin

package tcc_access

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/osquery/osquery-go/plugin/table"
)

var (
	testContext   = false
	tccPathPrefix = ""
	tccPathSuffix = "/Library/Application Support/com.apple.TCC/TCC.db"
	dbQuery       = "SELECT service, client, client_type, auth_value, auth_reason, last_modified, policy_id, indirect_object_identifier, indirect_object_identifier_type FROM access;"
	sqlite3Path   = "/usr/bin/sqlite3"
	dbColNames    = []string{"service", "client", "client_type", "auth_value", "auth_reason", "last_modified", "policy_id", "indirect_object_identifier", "indirect_object_identifier_type"}
)

// Columns is the schema of the table.
func Columns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		// added here
		table.TextColumn("source"),
		table.IntegerColumn("uid"),
		// derived from a TCC.db
		table.TextColumn("service"),
		table.TextColumn("client"),
		table.IntegerColumn("client_type"),
		table.IntegerColumn("auth_value"),
		table.IntegerColumn("auth_reason"),
		table.BigIntColumn("last_modified"),
		table.IntegerColumn("policy_id"),
		table.TextColumn("indirect_object_identifier"),
		table.IntegerColumn("indirect_object_identifier_type"),
	}
}

// Generate is called to return the results for the table at query time.
// Constraints for generating can be retrieved from the queryContext.

func Generate(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
	// get all usernames on the mac
	cmd := exec.Command("dscl", ".", "list", "/Users")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	allUsernames := strings.Split(string(out[:]), "\n")
	var usernames []string
	for _, username := range allUsernames {
		if !strings.HasPrefix(username, "_") && username != "nobody" && username != "root" && username != "daemon" {
			usernames = append(usernames, username)
		}
	}

	if testContext {
		usernames = []string{"testUser1", "testUser2"}
	}

	// build rows for every user-level TCC.db
	var rows []map[string]string
	for _, username := range usernames {
		uRs, err := getTCCAccessRows(username)
		if err != nil {
			return nil, err
		}
		rows = append(rows, uRs...)
	}

	// and for the system-level TCC.db
	sRs, err := getTCCAccessRows("")
	if err != nil {
		return nil, err
	}
	rows = append(rows, sRs...)

	return rows, nil
}

func getTCCAccessRows(username string) ([]map[string]string, error) {
	// avoids additional C compilation requirements that would be introduced by using
	// https://github.com/mattn/go-sqlite3

	tccPath := tccPathPrefix
	if username != "" {
		tccPath += "/Users/" + username
	}
	tccPath += tccPathSuffix
	cmd := exec.Command(sqlite3Path, tccPath, dbQuery)
	var dbOut bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &dbOut
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Generate failed at `cmd.Run()`%s: %w", stderr.String(), err)
	}

	parsedRows := parseTCCDbReadOutput(dbOut.Bytes())

	rows, err := buildTableRows(username, parsedRows)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func parseTCCDbReadOutput(dbOut []byte) [][]string {
	// split by newLine for rows, then by "|" for columns
	rawRows := strings.Split(string(dbOut[:]), "\n")
	n := len(rawRows)
	// the end of the db response is "\n", making the final row "", which we want to omit
	rawRows = rawRows[:n-1]

	parsedRows := make([][]string, 0, len(rawRows))
	for _, rawRow := range rawRows {
		parsedRows = append(parsedRows, strings.Split(rawRow, "|"))
	}
	return parsedRows
}

func buildTableRows(username string, parsedRows [][]string) ([]map[string]string, error) {
	// root's uid, for "system" rows by default
	uid := "0"
	source := "system"
	if username != "" {
		// a user-scoped table, so get its uid
		parsedUid, err := getUidFromUsername(username)
		if err != nil {
			return nil, fmt.Errorf("generate failed: %w", err)
		}
		uid = parsedUid
		source = "user"

	}

	var rows []map[string]string
	for _, parsedRow := range parsedRows {
		row := make(map[string]string)
		row["source"] = source
		row["uid"] = uid
		for i, rowColVal := range parsedRow {
			row[dbColNames[i]] = rowColVal
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func getUidFromUsername(username string) (string, error) {
	if testContext {
		testUId, ok := map[string]string{"testUser1": "1", "testUser2": "2"}[username]
		if ok {
			return testUId, nil
		}
		return "", errors.New("Invalid test username")
	}
	cmd := exec.Command("id", username)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("generate failed: %w", err)
	}

	sOut := string(out[:])
	unI := strings.Index(sOut, username)
	return sOut[4 : unI-1], nil
}
