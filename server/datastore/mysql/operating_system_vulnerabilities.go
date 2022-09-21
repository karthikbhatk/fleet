package mysql

import (
	"context"
	"fmt"
	"strings"

	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/jmoiron/sqlx"
)

func (ds *Datastore) ListOSVulnerabilitiesByHostID(ctx context.Context, hostID uint) ([]fleet.OSVulnerability, error) {
	stmt := `
        SELECT host_id, operating_system_id, cve
        FROM operating_system_vulnerabilities
        WHERE host_id = ?
	`
	r := []fleet.OSVulnerability{}

	if err := sqlx.SelectContext(ctx, ds.reader, &r, stmt, hostID); err != nil {
		return nil, ctxerr.Wrap(ctx, err, "listing OS vulnerabilities")
	}

	return r, nil
}

func (ds *Datastore) InsertOSVulnerabilities(ctx context.Context, vulnerabilities []fleet.OSVulnerability, source fleet.VulnerabilitySource) (int64, error) {
	var args []interface{}

	if len(vulnerabilities) == 0 {
		return 0, nil
	}

	values := strings.TrimSuffix(strings.Repeat("(?,?,?,?),", len(vulnerabilities)), ",")
	sql := fmt.Sprintf(`INSERT IGNORE INTO operating_system_vulnerabilities (host_id, operating_system_id, cve, source) VALUES %s`, values)

	for _, v := range vulnerabilities {
		args = append(args, v.HostID, v.OSID, v.CVE, source)
	}
	res, err := ds.writer.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, ctxerr.Wrap(ctx, err, "insert operating system vulnerabilities")
	}
	count, _ := res.RowsAffected()

	return count, nil
}

func (ds *Datastore) DeleteOSVulnerabilities(ctx context.Context, vulnerabilities []fleet.OSVulnerability) error {
	if len(vulnerabilities) == 0 {
		return nil
	}

	sql := fmt.Sprintf(
		`DELETE FROM operating_system_vulnerabilities WHERE (host_id, cve) IN (%s)`,
		strings.TrimSuffix(strings.Repeat("(?,?),", len(vulnerabilities)), ","),
	)

	var args []interface{}
	for _, v := range vulnerabilities {
		args = append(args, v.HostID, v.CVE)
	}
	if _, err := ds.writer.ExecContext(ctx, sql, args...); err != nil {
		return ctxerr.Wrapf(ctx, err, "deleting operating system vulnerabilities")
	}
	return nil
}
