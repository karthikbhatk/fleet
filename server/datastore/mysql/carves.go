package mysql

import (
	"database/sql"

	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
)

func (d *Datastore) NewCarve(metadata *kolide.CarveMetadata) (*kolide.CarveMetadata, error) {
	sql := `INSERT INTO carve_metadata (
		host_id,
		name,
		block_count,
		block_size,
		carve_size,
		carve_id,
		request_id,
		session_id
	) VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	)`

	result, err := d.db.Exec(
		sql,
		metadata.HostId,
		metadata.Name,
		metadata.BlockCount,
		metadata.BlockSize,
		metadata.CarveSize,
		metadata.CarveId,
		metadata.RequestId,
		metadata.SessionId,
	)
	if err != nil {
		return nil, errors.Wrap(err, "insert carve metadata")
	}

	id, _ := result.LastInsertId()
	metadata.ID = id

	return metadata, nil
}

func (d *Datastore) CarveBySessionId(sessionId string) (*kolide.CarveMetadata, error) {
	// Selecting max_block should be very efficient because MySQL is able to use
	// the index metadata and optimizes away the SELECT.
	sql := `
		SELECT
			id,
			host_id,
			created_at,
			name,
			block_count,
			block_size,
			carve_size,
			carve_id,
			request_id,
			session_id,
			(SELECT COALESCE(MAX(block_id), -1) FROM carve_blocks WHERE metadata_id = cm.id) AS max_block
		FROM carve_metadata cm
		WHERE session_id = ?
`

	var metadata kolide.CarveMetadata
	if err := d.db.Get(&metadata, sql, sessionId); err != nil {
		return nil, errors.Wrap(err, "get carve by session ID")
	}

	return &metadata, nil
}

func (d *Datastore) CarveByName(name string) (*kolide.CarveMetadata, error) {
	// Selecting max_block should be very efficient because MySQL is able to use
	// the index metadata and optimizes away the SELECT.
	sql := `
		SELECT
			id,
			host_id,
			created_at,
			name,
			block_count,
			block_size,
			carve_size,
			carve_id,
			request_id,
			session_id,
			(SELECT COALESCE(MAX(block_id), -1) FROM carve_blocks WHERE metadata_id = cm.id) AS max_block
		FROM carve_metadata cm
		WHERE name = ?
`

	var metadata kolide.CarveMetadata
	if err := d.db.Get(&metadata, sql, name); err != nil {
		return nil, errors.Wrap(err, "get carve by name")
	}

	return &metadata, nil
}

func (d *Datastore) ListCarves(opt kolide.ListOptions) ([]*kolide.CarveMetadata, error) {
	// Selecting max_block should be very efficient because MySQL is able to use
	// the index metadata and optimizes away the SELECT.
	query := `
		SELECT
			id,
			host_id,
			created_at,
			name,
			block_count,
			block_size,
			carve_size,
			carve_id,
			request_id,
			session_id,
			(SELECT COALESCE(MAX(block_id), -1) FROM carve_blocks WHERE metadata_id = cm.id) AS max_block
		FROM carve_metadata cm
`
	query = appendListOptionsToSQL(query, opt)
	carves := []*kolide.CarveMetadata{}
	if err := d.db.Select(&carves, query); err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "list carves")
	}

	return carves, nil
}

func (d *Datastore) NewBlock(metadataId int64, blockId int64, data []byte) error {
	sql := `
		INSERT INTO carve_blocks (
			metadata_id,
			block_id,
			data
		) VALUES (
			?,
			?,
			?
		)`

	if _, err := d.db.Exec(sql, metadataId, blockId, data); err != nil {
		return errors.Wrap(err, "insert carve block")
	}

	return nil
}

func (d *Datastore) GetBlock(metadataId int64, blockId int64) ([]byte, error) {
	sql := `
		SELECT data
		FROM carve_blocks
		WHERE metadata_id = ? AND block_id = ?
	`
	var data []byte
	if err := d.db.Get(&data, sql, metadataId, blockId); err != nil {
		return nil, errors.Wrap(err, "select data")
	}

	return data, nil
}
