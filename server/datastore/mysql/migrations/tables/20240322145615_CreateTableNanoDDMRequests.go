package tables

import (
	"database/sql"
	"fmt"
)

func init() {
	MigrationClient.AddMigration(Up_20240322145615, Down_20240322145615)
}

func Up_20240322145615(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE mdm_apple_declarative_requests (
  id BIGINT NOT NULL AUTO_INCREMENT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  enrollment_id VARCHAR(255) NOT NULL,
  -- Should be one of "tokens", "declaration-items", "status", or "declaration/…/…" where the ellipses reference a declaration on the server
  message_type VARCHAR(255),
  -- json payload
  raw_json TEXT,
  PRIMARY KEY (id),
  CONSTRAINT mdm_apple_declarative_requests_enrollment_id FOREIGN KEY (enrollment_id) REFERENCES nano_enrollments (id) ON DELETE CASCADE
)
`)

	if err != nil {
		return fmt.Errorf("creating mdm_apple_declarative_requsts: %w", err)
	}

	return nil
}

func Down_20240322145615(tx *sql.Tx) error {
	_, err := tx.Exec(`DROP TABLE mdm_apple_declarative_requests`)
	if err != nil {
		return fmt.Errorf("dropping mdm_apple_declarative_requests: %w", err)
	}

	return nil
}
