package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/fleetdm/fleet/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) NewHost(host *kolide.Host) (*kolide.Host, error) {
	sqlStatement := `
	INSERT INTO hosts (
		osquery_host_id,
		detail_update_time,
		label_update_time,
		node_key,
		host_name,
		uuid,
		platform,
		osquery_version,
		os_version,
		uptime,
		physical_memory,
		seen_time
	)
	VALUES( ?,?,?,?,?,?,?,?,?,?,?,? )
	`
	result, err := d.db.Exec(
		sqlStatement,
		host.OsqueryHostID,
		host.DetailUpdateTime,
		host.LabelUpdateTime,
		host.NodeKey,
		host.HostName,
		host.UUID,
		host.Platform,
		host.OsqueryVersion,
		host.OSVersion,
		host.Uptime,
		host.PhysicalMemory,
		host.SeenTime,
	)
	if err != nil {
		return nil, errors.Wrap(err, "new host")
	}
	id, _ := result.LastInsertId()
	host.ID = uint(id)
	return host, nil
}

// TODO needs test
func (d *Datastore) SaveHost(host *kolide.Host) error {
	sqlStatement := `
		UPDATE hosts SET
			detail_update_time = ?,
			label_update_time = ?,
			node_key = ?,
			host_name = ?,
			uuid = ?,
			platform = ?,
			osquery_version = ?,
			os_version = ?,
			uptime = ?,
			physical_memory = ?,
			cpu_type = ?,
			cpu_subtype = ?,
			cpu_brand = ?,
			cpu_physical_cores = ?,
			hardware_vendor = ?,
			hardware_model = ?,
			hardware_version = ?,
			hardware_serial = ?,
			computer_name = ?,
			build = ?,
			platform_like = ?,
			code_name = ?,
			cpu_logical_cores = ?,
			seen_time = ?,
			distributed_interval = ?,
			config_tls_refresh = ?,
			logger_tls_period = ?,
			additional = COALESCE(?, additional),
			enroll_secret_name = ?,
			primary_ip = ?,
			primary_mac = ?
		WHERE id = ?
	`
	_, err := d.db.Exec(sqlStatement,
		host.DetailUpdateTime,
		host.LabelUpdateTime,
		host.NodeKey,
		host.HostName,
		host.UUID,
		host.Platform,
		host.OsqueryVersion,
		host.OSVersion,
		host.Uptime,
		host.PhysicalMemory,
		host.CPUType,
		host.CPUSubtype,
		host.CPUBrand,
		host.CPUPhysicalCores,
		host.HardwareVendor,
		host.HardwareModel,
		host.HardwareVersion,
		host.HardwareSerial,
		host.ComputerName,
		host.Build,
		host.PlatformLike,
		host.CodeName,
		host.CPULogicalCores,
		host.SeenTime,
		host.DistributedInterval,
		host.ConfigTLSRefresh,
		host.LoggerTLSPeriod,
		host.Additional,
		host.EnrollSecretName,
		host.PrimaryIP,
		host.PrimaryMac,
		host.ID,
	)
	if err != nil {
		return errors.Wrapf(err, "save host with id %d", host.ID)
	}

	return nil
}

func (d *Datastore) DeleteHost(hid uint) error {
	err := d.deleteEntity("hosts", hid)
	if err != nil {
		return errors.Wrapf(err, "deleting host with id %d", hid)
	}
	return nil
}

func (d *Datastore) Host(id uint) (*kolide.Host, error) {
	sqlStatement := `
		SELECT * FROM hosts
		WHERE id = ? LIMIT 1
	`
	host := &kolide.Host{}
	err := d.db.Get(host, sqlStatement, id)
	if err != nil {
		return nil, errors.Wrap(err, "getting host by id")
	}

	return host, nil
}

func (d *Datastore) ListHosts(opt kolide.HostListOptions) ([]*kolide.Host, error) {
	sql := `
		SELECT * FROM hosts
	`
	var params []interface{}
	switch opt.StatusFilter {
	case "new":
		sql += "WHERE DATE_ADD(created_at, INTERVAL 1 DAY) >= ?"
		params = append(params, time.Now())
	case "online":
		sql += fmt.Sprintf("WHERE DATE_ADD(seen_time, INTERVAL LEAST(distributed_interval, config_tls_refresh) + %d SECOND) > ?", kolide.OnlineIntervalBuffer)
		params = append(params, time.Now())
	case "offline":
		sql += fmt.Sprintf("WHERE DATE_ADD(seen_time, INTERVAL LEAST(distributed_interval, config_tls_refresh) + %d SECOND) <= ? AND DATE_ADD(seen_time, INTERVAL 30 DAY) >= ?", kolide.OnlineIntervalBuffer)
		params = append(params, time.Now(), time.Now())
	case "mia":
		sql += "WHERE DATE_ADD(seen_time, INTERVAL 30 DAY) <= ?"
		params = append(params, time.Now())
	}
	sql = appendListOptionsToSQL(sql, opt.ListOptions)

	hosts := []*kolide.Host{}
	if err := d.db.Select(&hosts, sql, params...); err != nil {
		return nil, errors.Wrap(err, "list hosts")
	}

	return hosts, nil
}

func (d *Datastore) CleanupIncomingHosts(now time.Time) error {
	sqlStatement := `
		DELETE FROM hosts
		WHERE host_name = '' AND osquery_version = ''
		AND created_at < (? - INTERVAL 5 MINUTE)
	`
	if _, err := d.db.Exec(sqlStatement, now); err != nil {
		return errors.Wrap(err, "cleanup incoming hosts")
	}

	return nil
}

func (d *Datastore) GenerateHostStatusStatistics(now time.Time) (online, offline, mia, new uint, e error) {
	// The logic in this function should remain synchronized with
	// host.Status and CountHostsInTargets

	sqlStatement := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(CASE WHEN DATE_ADD(seen_time, INTERVAL 30 DAY) <= ? THEN 1 ELSE 0 END), 0) mia,
			COALESCE(SUM(CASE WHEN DATE_ADD(seen_time, INTERVAL LEAST(distributed_interval, config_tls_refresh) + %d SECOND) <= ? AND DATE_ADD(seen_time, INTERVAL 30 DAY) >= ? THEN 1 ELSE 0 END), 0) offline,
			COALESCE(SUM(CASE WHEN DATE_ADD(seen_time, INTERVAL LEAST(distributed_interval, config_tls_refresh) + %d SECOND) > ? THEN 1 ELSE 0 END), 0) online,
			COALESCE(SUM(CASE WHEN DATE_ADD(created_at, INTERVAL 1 DAY) >= ? THEN 1 ELSE 0 END), 0) new
		FROM hosts
		LIMIT 1;
	`, kolide.OnlineIntervalBuffer, kolide.OnlineIntervalBuffer)

	counts := struct {
		MIA     uint `db:"mia"`
		Offline uint `db:"offline"`
		Online  uint `db:"online"`
		New     uint `db:"new"`
	}{}
	err := d.db.Get(&counts, sqlStatement, now, now, now, now, now)
	if err != nil && err != sql.ErrNoRows {
		e = errors.Wrap(err, "generating host statistics")
		return
	}

	mia = counts.MIA
	offline = counts.Offline
	online = counts.Online
	new = counts.New
	return online, offline, mia, new, nil
}

// EnrollHost enrolls a host
func (d *Datastore) EnrollHost(osqueryHostID, nodeKey, secretName string) (*kolide.Host, error) {
	if osqueryHostID == "" {
		return nil, fmt.Errorf("missing osquery host identifier")
	}

	zeroTime := time.Unix(0, 0).Add(24 * time.Hour)
	sqlInsert := `
		INSERT INTO hosts (
			detail_update_time,
			label_update_time,
			osquery_host_id,
			seen_time,
			node_key,
			enroll_secret_name
		) VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			node_key = VALUES(node_key)
	`

	var result sql.Result
	result, err := d.db.Exec(sqlInsert, zeroTime, zeroTime, osqueryHostID, time.Now().UTC(), nodeKey, secretName)

	if err != nil {
		return nil, errors.Wrap(err, "inserting")
	}

	id, _ := result.LastInsertId()
	sqlSelect := `
		SELECT * FROM hosts WHERE id = ? LIMIT 1
	`
	host := &kolide.Host{}
	err = d.db.Get(host, sqlSelect, id)
	if err != nil {
		return nil, errors.Wrap(err, "getting the host to return")
	}

	_, err = d.db.Exec(`INSERT IGNORE INTO label_membership (host_id, label_id) VALUES (?, (SELECT id FROM labels WHERE name = 'All Hosts' AND label_type = 1))`, id)
	if err != nil {
		return nil, errors.Wrap(err, "insert new host into all hosts label")
	}

	return host, nil

}

func (d *Datastore) AuthenticateHost(nodeKey string) (*kolide.Host, error) {
	sqlStatement := `
		SELECT
			id,
			osquery_host_id,
			created_at,
			updated_at,
			detail_update_time,
			label_update_time,
			node_key,
			host_name,
			uuid,
			platform,
			osquery_version,
			os_version,
			build,
			platform_like,
			code_name,
			uptime,
			physical_memory,
			cpu_type,
			cpu_subtype,
			cpu_brand,
			cpu_physical_cores,
			cpu_logical_cores,
			hardware_vendor,
			hardware_model,
			hardware_version,
			hardware_serial,
			computer_name,
			primary_ip_id,
			seen_time,
			distributed_interval,
			logger_tls_period,
			config_tls_refresh,
			enroll_secret_name
		FROM hosts
		WHERE node_key = ?
		LIMIT 1
	`

	host := &kolide.Host{}
	if err := d.db.Get(host, sqlStatement, nodeKey); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, notFound("Host")
		default:
			return nil, errors.New("finding host")
		}
	}

	return host, nil
}

func (d *Datastore) MarkHostSeen(host *kolide.Host, t time.Time) error {
	sqlStatement := `
		UPDATE hosts SET
			seen_time = ?
		WHERE node_key=?
	`

	_, err := d.db.Exec(sqlStatement, t, host.NodeKey)
	if err != nil {
		return errors.Wrap(err, "marking host seen")
	}

	host.UpdatedAt = t
	return nil
}

func (d *Datastore) searchHostsWithOmits(query string, omit ...uint) ([]*kolide.Host, error) {
	hostQuery := transformQuery(query)
	ipQuery := `"` + query + `"`

	sqlStatement :=
		`
		SELECT DISTINCT *
		FROM hosts
		WHERE
		(
			MATCH (host_name, uuid) AGAINST (? IN BOOLEAN MODE)
			OR MATCH (primary_ip, primary_mac) AGAINST (? IN BOOLEAN MODE)
		)
		AND id NOT IN (?)
		LIMIT 10
	`

	sql, args, err := sqlx.In(sqlStatement, hostQuery, ipQuery, omit)
	if err != nil {
		return nil, errors.Wrap(err, "searching hosts")
	}
	sql = d.db.Rebind(sql)

	hosts := []*kolide.Host{}

	err = d.db.Select(&hosts, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "searching hosts rebound")
	}

	return hosts, nil
}

func (d *Datastore) searchHostsDefault(omit ...uint) ([]*kolide.Host, error) {
	sqlStatement := `
	SELECT * FROM hosts
	WHERE id NOT in (?)
	ORDER BY seen_time DESC
	LIMIT 5
	`

	var in interface{}
	{
		// use -1 if there are no values to omit.
		// Avoids empty args error for `sqlx.In`
		in = omit
		if len(omit) == 0 {
			in = -1
		}
	}

	var hosts []*kolide.Host
	sql, args, err := sqlx.In(sqlStatement, in)
	if err != nil {
		return nil, errors.Wrap(err, "searching default hosts")
	}
	sql = d.db.Rebind(sql)
	err = d.db.Select(&hosts, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "searching default hosts rebound")
	}
	return hosts, nil
}

// SearchHosts find hosts by query containing an IP address, a host name or UUID.
// Optionally pass a list of IDs to omit from the search
func (d *Datastore) SearchHosts(query string, omit ...uint) ([]*kolide.Host, error) {
	hostQuery := transformQuery(query)
	if !queryMinLength(hostQuery) {
		return d.searchHostsDefault(omit...)
	}
	if len(omit) > 0 {
		return d.searchHostsWithOmits(query, omit...)
	}

	// Needs quotes to avoid each . marking a word boundary
	ipQuery := `"` + query + `"`

	sqlStatement :=
		`
		SELECT DISTINCT *
		FROM hosts
		WHERE
		(
			MATCH (host_name, uuid) AGAINST (? IN BOOLEAN MODE)
			OR MATCH (primary_ip, primary_mac) AGAINST (? IN BOOLEAN MODE)
		)
		LIMIT 10
	`
	hosts := []*kolide.Host{}

	if err := d.db.Select(&hosts, sqlStatement, hostQuery, ipQuery); err != nil {
		return nil, errors.Wrap(err, "searching hosts")
	}

	return hosts, nil

}

func (d *Datastore) HostIDsByName(hostnames []string) ([]uint, error) {
	if len(hostnames) == 0 {
		return []uint{}, nil
	}

	sqlStatement := `
		SELECT id FROM hosts
		WHERE host_name IN (?)
	`

	sql, args, err := sqlx.In(sqlStatement, hostnames)
	if err != nil {
		return nil, errors.Wrap(err, "building query to get host IDs")
	}

	var hostIDs []uint
	if err := d.db.Select(&hostIDs, sql, args...); err != nil {
		return nil, errors.Wrap(err, "get host IDs")
	}

	return hostIDs, nil

}

func (d *Datastore) HostByIdentifier(identifier string) (*kolide.Host, error) {
	sql := `
		SELECT * FROM hosts
		WHERE ? IN (host_name, osquery_host_id, node_key, uuid)
		LIMIT 1
	`
	host := &kolide.Host{}
	err := d.db.Get(host, sql, identifier)
	if err != nil {
		return nil, errors.Wrap(err, "get host by identifier")
	}

	return host, nil
}
