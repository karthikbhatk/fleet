package datastore

import (
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (orm gormDB) EnrollHost(uuid, hostname, ip, platform string, nodeKeySize int) (*kolide.Host, error) {
	if uuid == "" {
		return nil, errors.New("missing uuid for host enrollment", "programmer error?")
	}
	host := kolide.Host{UUID: uuid}
	err := orm.DB.Where(&host).First(&host).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			// Create new Host
			host = kolide.Host{
				UUID:             uuid,
				HostName:         hostname,
				PrimaryIP:        ip,
				Platform:         platform,
				DetailUpdateTime: time.Unix(0, 0).Add(24 * time.Hour),
			}

		default:
			return nil, err
		}
	}

	// Generate a new key each enrollment
	host.NodeKey, err = generateRandomText(nodeKeySize)
	if err != nil {
		return nil, err
	}

	// Update these fields if provided
	if hostname != "" {
		host.HostName = hostname
	}
	if ip != "" {
		host.PrimaryIP = ip
	}
	if platform != "" {
		host.Platform = platform
	}

	if err := orm.DB.Save(&host).Error; err != nil {
		return nil, err
	}

	return &host, nil
}

func (orm gormDB) AuthenticateHost(nodeKey string) (*kolide.Host, error) {
	host := kolide.Host{NodeKey: nodeKey}
	err := orm.DB.Where("node_key = ?", host.NodeKey).First(&host).Error
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			e := errors.NewFromError(
				err,
				http.StatusUnauthorized,
				"invalid node key",
			)
			// osqueryd expects the literal string "true" here
			e.Extra = map[string]interface{}{"node_invalid": "true"}
			return nil, e
		default:
			return nil, errors.DatabaseError(err)
		}
	}

	return &host, nil
}

func (orm gormDB) SaveHost(host *kolide.Host) error {
	if err := orm.DB.Save(host).Error; err != nil {
		return errors.DatabaseError(err)
	}
	return nil
}

func (orm gormDB) DeleteHost(host *kolide.Host) error {
	return orm.DB.Delete(host).Error
}

func (orm gormDB) Host(id uint) (*kolide.Host, error) {
	host := &kolide.Host{
		ID: id,
	}
	err := orm.DB.Where(host).First(host).Error
	if err != nil {
		return nil, err
	}
	return host, nil
}

func (orm gormDB) ListHosts(opt kolide.ListOptions) ([]*kolide.Host, error) {
	var hosts []*kolide.Host
	err := orm.applyListOptions(opt).Find(&hosts).Error
	if err != nil {
		return nil, err
	}
	return hosts, nil
}

func (orm gormDB) NewHost(host *kolide.Host) (*kolide.Host, error) {
	if host == nil {
		return nil, errors.New(
			"error creating host",
			"nil pointer passed to NewHost",
		)
	}
	err := orm.DB.Create(host).Error
	if err != nil {
		return nil, err
	}
	return host, err
}

func (orm gormDB) MarkHostSeen(host *kolide.Host, t time.Time) error {
	err := orm.DB.Exec("UPDATE hosts SET updated_at=? WHERE node_key=?", t, host.NodeKey).Error
	if err != nil {
		return errors.DatabaseError(err)
	}
	host.UpdatedAt = t
	return nil
}

func (orm gormDB) SearchHosts(query string, omit []uint) ([]kolide.Host, error) {
	sql := `
SELECT *
FROM hosts
WHERE MATCH(host_name, primary_ip)
AGAINST(? IN BOOLEAN MODE)
`
	results := []kolide.Host{}

	var db *gorm.DB
	if len(omit) > 0 {
		sql += "AND id NOT IN (?) LIMIT 10;"
		db = orm.DB.Raw(sql, query+"*", omit)
	} else {
		sql += "LIMIT 10;"
		db = orm.DB.Raw(sql, query+"*")
	}

	err := db.Scan(&results).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.DatabaseError(err)
	}
	return results, nil
}

func (orm gormDB) DistributedQueriesForHost(host *kolide.Host) (map[uint]string, error) {
	sql := `
SELECT DISTINCT dqc.id, q.query
FROM distributed_query_campaigns dqc
JOIN distributed_query_campaign_targets dqct
    ON (dqc.id = dqct.distributed_query_campaign_id)
LEFT JOIN label_query_executions lqe
    ON (dqct.type = ? AND dqct.target_id = lqe.label_id AND lqe.matches)
LEFT JOIN hosts h
    ON ((dqct.type = ? AND lqe.host_id = h.id) OR (dqct.type = ? AND dqct.target_id = h.id))
LEFT JOIN distributed_query_executions dqe
    ON (h.id = dqe.host_id AND dqc.id = dqe.distributed_query_id)
JOIN queries q
    ON (dqc.query_id = q.id)
WHERE dqe.status IS NULL AND dqc.status = ? AND h.id = ?;
`
	rows, err := orm.DB.Raw(sql, kolide.TargetLabel, kolide.TargetLabel, kolide.TargetHost, kolide.QueryRunning, host.ID).Rows()
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.DatabaseError(err)
	}
	defer rows.Close()

	results := map[uint]string{}
	for rows.Next() {
		var id uint
		var query string
		err = rows.Scan(&id, &query)
		if err != nil {
			return nil, errors.DatabaseError(err)
		}
		results[id] = query
	}

	return results, nil

}
