package tables

import (
	"database/sql"
	"fmt"
)

func init() {
	MigrationClient.AddMigration(Up_20240925112748, Down_20240925112748)
}

func Up_20240925112748(tx *sql.Tx) error {
	_, err := tx.Exec(`
CREATE TABLE setup_experience_status_results (
	id		INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	host_uuid	VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	-- Type of step
	type		ENUM('bootstrap-package', 'software-install', 'post-install-script') NOT NULL,
	name		VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL,
	-- Status of the step
	status		ENUM('pending', 'running', 'success', 'failure') NOT NULL,
	host_software_installs_id -- FIXME
	nano_commands_command_uuid -- FIXME
	execution_id	VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL, -- FIXME host_script_results, maybe ID in new table?
	error 		VARCHAR(255) COLLATE utf8mb4_unicode_ci NOT NULL,

	PRIMARY KEY (id)
)
`)
	// Service layer state machine like SetupExperienceNestStep()?
	// Called from each of the three endpoints (software install, vpp
	// mdm, scripts) involved in the setup when an eligible installer
	// writes its results
	if err != nil {
		return fmt.Errorf("failed to create setup_experience_status_results table: %w", err)
	}

	return nil
}

func Down_20240925112748(tx *sql.Tx) error {
	return nil
}
