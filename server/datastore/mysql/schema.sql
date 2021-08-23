/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `activities` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_id` int(10) unsigned DEFAULT NULL,
  `user_name` varchar(255) DEFAULT NULL,
  `activity_type` varchar(255) NOT NULL,
  `details` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_activities_user_id` (`user_id`),
  CONSTRAINT `activities_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_config_json` (
  `id` int(10) unsigned NOT NULL DEFAULT '1',
  `json_value` json NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
INSERT INTO `app_config_json` VALUES (1,'{\"org_info\": {\"org_name\": \"\", \"org_logo_url\": \"\"}, \"sso_settings\": {\"idp_name\": \"\", \"metadata\": \"\", \"entity_id\": \"\", \"enable_sso\": false, \"issuer_uri\": \"\", \"metadata_url\": \"\", \"idp_image_url\": \"\", \"enable_sso_idp_login\": false}, \"agent_options\": {\"config\": {\"options\": {\"logger_plugin\": \"tls\", \"pack_delimiter\": \"/\", \"logger_tls_period\": 10, \"distributed_plugin\": \"tls\", \"disable_distributed\": false, \"logger_tls_endpoint\": \"/api/v1/osquery/log\", \"distributed_interval\": 10, \"distributed_tls_max_attempts\": 3}, \"decorators\": {\"load\": [\"SELECT uuid AS host_uuid FROM system_info;\", \"SELECT hostname AS hostname FROM system_info;\"]}}, \"overrides\": {}}, \"host_settings\": {\"enable_host_users\": true, \"enable_software_inventory\": false}, \"smtp_settings\": {\"port\": 587, \"domain\": \"\", \"server\": \"\", \"password\": \"\", \"user_name\": \"\", \"configured\": false, \"enable_smtp\": false, \"enable_ssl_tls\": true, \"sender_address\": \"\", \"enable_start_tls\": true, \"verify_ssl_certs\": true, \"authentication_type\": \"0\", \"authentication_method\": \"0\"}, \"server_settings\": {\"server_url\": \"\", \"enable_analytics\": false, \"live_query_disabled\": false}, \"host_expiry_settings\": {\"host_expiry_window\": 0, \"host_expiry_enabled\": false}, \"vulnerability_settings\": {\"databases_path\": \"\"}}','2021-08-23 19:19:48','2021-08-23 19:19:48');
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_configs` (
  `id` int(10) unsigned NOT NULL DEFAULT '1',
  `org_name` varchar(255) NOT NULL DEFAULT '',
  `org_logo_url` varchar(255) NOT NULL DEFAULT '',
  `server_url` varchar(255) NOT NULL DEFAULT '',
  `smtp_configured` tinyint(1) NOT NULL DEFAULT '0',
  `smtp_sender_address` varchar(255) NOT NULL DEFAULT '',
  `smtp_server` varchar(255) NOT NULL DEFAULT '',
  `smtp_port` int(10) unsigned NOT NULL DEFAULT '587',
  `smtp_authentication_type` int(11) NOT NULL DEFAULT '0',
  `smtp_enable_ssl_tls` tinyint(1) NOT NULL DEFAULT '1',
  `smtp_authentication_method` int(11) NOT NULL DEFAULT '0',
  `smtp_domain` varchar(255) NOT NULL DEFAULT '',
  `smtp_user_name` varchar(255) NOT NULL DEFAULT '',
  `smtp_password` varchar(255) NOT NULL DEFAULT '',
  `smtp_verify_ssl_certs` tinyint(1) NOT NULL DEFAULT '1',
  `smtp_enable_start_tls` tinyint(1) NOT NULL DEFAULT '1',
  `entity_id` varchar(255) NOT NULL DEFAULT '',
  `issuer_uri` varchar(255) NOT NULL DEFAULT '',
  `idp_image_url` varchar(512) NOT NULL DEFAULT '',
  `metadata` text NOT NULL,
  `metadata_url` varchar(512) NOT NULL DEFAULT '',
  `idp_name` varchar(255) NOT NULL DEFAULT '',
  `enable_sso` tinyint(1) NOT NULL DEFAULT '0',
  `fim_interval` int(11) NOT NULL DEFAULT '300',
  `fim_file_accesses` varchar(255) NOT NULL DEFAULT '',
  `host_expiry_enabled` tinyint(1) NOT NULL DEFAULT '0',
  `host_expiry_window` int(11) DEFAULT '0',
  `live_query_disabled` tinyint(1) NOT NULL DEFAULT '0',
  `additional_queries` json DEFAULT NULL,
  `enable_sso_idp_login` tinyint(1) NOT NULL DEFAULT '0',
  `agent_options` json DEFAULT NULL,
  `enable_analytics` tinyint(1) NOT NULL DEFAULT '0',
  `vulnerability_databases_path` text,
  `enable_host_users` tinyint(1) DEFAULT '1',
  `enable_software_inventory` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
INSERT INTO `app_configs` VALUES (1,'','','',0,'','',587,0,1,0,'','','',1,1,'','','','','','',0,300,'',0,0,0,NULL,0,'{\"config\": {\"options\": {\"logger_plugin\": \"tls\", \"pack_delimiter\": \"/\", \"logger_tls_period\": 10, \"distributed_plugin\": \"tls\", \"disable_distributed\": false, \"logger_tls_endpoint\": \"/api/v1/osquery/log\", \"distributed_interval\": 10, \"distributed_tls_max_attempts\": 3}, \"decorators\": {\"load\": [\"SELECT uuid AS host_uuid FROM system_info;\", \"SELECT hostname AS hostname FROM system_info;\"]}}, \"overrides\": {}}',0,NULL,1,0);
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `carve_blocks` (
  `metadata_id` int(10) unsigned NOT NULL,
  `block_id` int(11) NOT NULL,
  `data` longblob,
  PRIMARY KEY (`metadata_id`,`block_id`),
  CONSTRAINT `carve_blocks_ibfk_1` FOREIGN KEY (`metadata_id`) REFERENCES `carve_metadata` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `carve_metadata` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` int(10) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `name` varchar(255) DEFAULT NULL,
  `block_count` int(10) unsigned NOT NULL,
  `block_size` int(10) unsigned NOT NULL,
  `carve_size` bigint(20) unsigned NOT NULL,
  `carve_id` varchar(64) NOT NULL,
  `request_id` varchar(64) NOT NULL,
  `session_id` varchar(255) NOT NULL,
  `expired` tinyint(4) DEFAULT '0',
  `max_block` int(11) DEFAULT '-1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_session_id` (`session_id`),
  UNIQUE KEY `idx_name` (`name`),
  KEY `host_id` (`host_id`),
  CONSTRAINT `carve_metadata_ibfk_1` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `distributed_query_campaign_targets` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `type` int(11) DEFAULT NULL,
  `distributed_query_campaign_id` int(10) unsigned DEFAULT NULL,
  `target_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `distributed_query_campaigns` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `query_id` int(10) unsigned DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `user_id` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `email_changes` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `token` varchar(128) NOT NULL,
  `new_email` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_email_changes_token` (`token`) USING BTREE,
  KEY `fk_email_changes_users` (`user_id`),
  CONSTRAINT `fk_email_changes_users` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `enroll_secrets` (
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `secret` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `team_id` int(10) unsigned DEFAULT NULL,
  UNIQUE KEY `secret` (`secret`),
  KEY `fk_enroll_secrets_team_id` (`team_id`),
  CONSTRAINT `enroll_secrets_ibfk_1` FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `host_additional` (
  `host_id` int(10) unsigned NOT NULL,
  `additional` json DEFAULT NULL,
  PRIMARY KEY (`host_id`),
  CONSTRAINT `host_additional_ibfk_1` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `host_software` (
  `host_id` int(10) unsigned NOT NULL,
  `software_id` bigint(20) unsigned NOT NULL,
  PRIMARY KEY (`host_id`,`software_id`),
  KEY `host_software_software_fk` (`software_id`),
  CONSTRAINT `host_software_ibfk_1` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`) ON DELETE CASCADE,
  CONSTRAINT `host_software_ibfk_2` FOREIGN KEY (`software_id`) REFERENCES `software` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `host_users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` int(10) unsigned NOT NULL,
  `uid` int(10) unsigned NOT NULL,
  `username` varchar(255) DEFAULT NULL,
  `groupname` varchar(255) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `removed_at` timestamp NULL DEFAULT NULL,
  `user_type` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_uid_username` (`host_id`,`uid`,`username`),
  CONSTRAINT `host_users_ibfk_1` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `hosts` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `osquery_host_id` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `detail_updated_at` timestamp NULL DEFAULT NULL,
  `node_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT NULL,
  `hostname` varchar(255) NOT NULL DEFAULT '',
  `uuid` varchar(255) NOT NULL DEFAULT '',
  `platform` varchar(255) NOT NULL DEFAULT '',
  `osquery_version` varchar(255) NOT NULL DEFAULT '',
  `os_version` varchar(255) NOT NULL DEFAULT '',
  `build` varchar(255) NOT NULL DEFAULT '',
  `platform_like` varchar(255) NOT NULL DEFAULT '',
  `code_name` varchar(255) NOT NULL DEFAULT '',
  `uptime` bigint(20) NOT NULL DEFAULT '0',
  `memory` bigint(20) NOT NULL DEFAULT '0',
  `cpu_type` varchar(255) NOT NULL DEFAULT '',
  `cpu_subtype` varchar(255) NOT NULL DEFAULT '',
  `cpu_brand` varchar(255) NOT NULL DEFAULT '',
  `cpu_physical_cores` int(11) NOT NULL DEFAULT '0',
  `cpu_logical_cores` int(11) NOT NULL DEFAULT '0',
  `hardware_vendor` varchar(255) NOT NULL DEFAULT '',
  `hardware_model` varchar(255) NOT NULL DEFAULT '',
  `hardware_version` varchar(255) NOT NULL DEFAULT '',
  `hardware_serial` varchar(255) NOT NULL DEFAULT '',
  `computer_name` varchar(255) NOT NULL DEFAULT '',
  `primary_ip_id` int(10) unsigned DEFAULT NULL,
  `seen_time` timestamp NULL DEFAULT NULL,
  `distributed_interval` int(11) DEFAULT '0',
  `logger_tls_period` int(11) DEFAULT '0',
  `config_tls_refresh` int(11) DEFAULT '0',
  `primary_ip` varchar(45) NOT NULL DEFAULT '',
  `primary_mac` varchar(17) NOT NULL DEFAULT '',
  `label_updated_at` timestamp NOT NULL DEFAULT '2000-01-01 00:00:00',
  `last_enrolled_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `refetch_requested` tinyint(1) NOT NULL DEFAULT '0',
  `team_id` int(10) unsigned DEFAULT NULL,
  `gigs_disk_space_available` float NOT NULL DEFAULT '0',
  `percent_disk_space_available` float NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_osquery_host_id` (`osquery_host_id`),
  UNIQUE KEY `idx_host_unique_nodekey` (`node_key`),
  KEY `fk_hosts_team_id` (`team_id`),
  FULLTEXT KEY `hosts_search` (`hostname`,`uuid`),
  FULLTEXT KEY `host_ip_mac_search` (`primary_ip`,`primary_mac`),
  CONSTRAINT `hosts_ibfk_1` FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invite_teams` (
  `invite_id` int(10) unsigned NOT NULL,
  `team_id` int(10) unsigned NOT NULL,
  `role` varchar(64) NOT NULL,
  PRIMARY KEY (`invite_id`,`team_id`),
  KEY `fk_team_id` (`team_id`),
  CONSTRAINT `invite_teams_ibfk_1` FOREIGN KEY (`invite_id`) REFERENCES `invites` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `invite_teams_ibfk_2` FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invites` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `invited_by` int(10) unsigned NOT NULL,
  `email` varchar(255) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `position` varchar(255) DEFAULT NULL,
  `token` varchar(255) NOT NULL,
  `sso_enabled` tinyint(1) NOT NULL DEFAULT '0',
  `global_role` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_invite_unique_email` (`email`),
  UNIQUE KEY `idx_invite_unique_key` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `label_membership` (
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `label_id` int(10) unsigned NOT NULL,
  `host_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`host_id`,`label_id`),
  KEY `idx_lm_label_id` (`label_id`),
  CONSTRAINT `fk_lm_host_id` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `fk_lm_label_id` FOREIGN KEY (`label_id`) REFERENCES `labels` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `labels` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `query` mediumtext NOT NULL,
  `platform` varchar(255) DEFAULT NULL,
  `label_type` int(10) unsigned NOT NULL DEFAULT '1',
  `label_membership_type` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_label_unique_name` (`name`),
  FULLTEXT KEY `labels_search` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `locks` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `owner` varchar(255) DEFAULT NULL,
  `expires_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `migration_status_tables` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `version_id` bigint(20) NOT NULL,
  `is_applied` tinyint(1) NOT NULL,
  `tstamp` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
INSERT INTO `migration_status_tables` VALUES (1,0,1,'2021-08-23 19:19:44'),(2,20161118193812,1,'2021-08-23 19:19:44'),(3,20161118211713,1,'2021-08-23 19:19:44'),(4,20161118212436,1,'2021-08-23 19:19:44'),(5,20161118212515,1,'2021-08-23 19:19:44'),(6,20161118212528,1,'2021-08-23 19:19:44'),(7,20161118212538,1,'2021-08-23 19:19:45'),(8,20161118212549,1,'2021-08-23 19:19:45'),(9,20161118212557,1,'2021-08-23 19:19:45'),(10,20161118212604,1,'2021-08-23 19:19:45'),(11,20161118212613,1,'2021-08-23 19:19:45'),(12,20161118212621,1,'2021-08-23 19:19:45'),(13,20161118212630,1,'2021-08-23 19:19:45'),(14,20161118212641,1,'2021-08-23 19:19:45'),(15,20161118212649,1,'2021-08-23 19:19:45'),(16,20161118212656,1,'2021-08-23 19:19:45'),(17,20161118212758,1,'2021-08-23 19:19:45'),(18,20161128234849,1,'2021-08-23 19:19:45'),(19,20161230162221,1,'2021-08-23 19:19:45'),(20,20170104113816,1,'2021-08-23 19:19:45'),(21,20170105151732,1,'2021-08-23 19:19:45'),(22,20170108191242,1,'2021-08-23 19:19:45'),(23,20170109094020,1,'2021-08-23 19:19:45'),(24,20170109130438,1,'2021-08-23 19:19:45'),(25,20170110202752,1,'2021-08-23 19:19:45'),(26,20170111133013,1,'2021-08-23 19:19:45'),(27,20170117025759,1,'2021-08-23 19:19:45'),(28,20170118191001,1,'2021-08-23 19:19:45'),(29,20170119234632,1,'2021-08-23 19:19:45'),(30,20170124230432,1,'2021-08-23 19:19:45'),(31,20170127014618,1,'2021-08-23 19:19:45'),(32,20170131232841,1,'2021-08-23 19:19:45'),(33,20170223094154,1,'2021-08-23 19:19:45'),(34,20170306075207,1,'2021-08-23 19:19:46'),(35,20170309100733,1,'2021-08-23 19:19:46'),(36,20170331111922,1,'2021-08-23 19:19:46'),(37,20170502143928,1,'2021-08-23 19:19:46'),(38,20170504130602,1,'2021-08-23 19:19:46'),(39,20170509132100,1,'2021-08-23 19:19:46'),(40,20170519105647,1,'2021-08-23 19:19:46'),(41,20170519105648,1,'2021-08-23 19:19:46'),(42,20170831234300,1,'2021-08-23 19:19:46'),(43,20170831234301,1,'2021-08-23 19:19:46'),(44,20170831234303,1,'2021-08-23 19:19:46'),(45,20171116163618,1,'2021-08-23 19:19:46'),(46,20171219164727,1,'2021-08-23 19:19:46'),(47,20180620164811,1,'2021-08-23 19:19:46'),(48,20180620175054,1,'2021-08-23 19:19:46'),(49,20180620175055,1,'2021-08-23 19:19:46'),(50,20191010101639,1,'2021-08-23 19:19:46'),(51,20191010155147,1,'2021-08-23 19:19:46'),(52,20191220130734,1,'2021-08-23 19:19:46'),(53,20200311140000,1,'2021-08-23 19:19:46'),(54,20200405120000,1,'2021-08-23 19:19:46'),(55,20200407120000,1,'2021-08-23 19:19:46'),(56,20200420120000,1,'2021-08-23 19:19:47'),(57,20200504120000,1,'2021-08-23 19:19:47'),(58,20200512120000,1,'2021-08-23 19:19:47'),(59,20200707120000,1,'2021-08-23 19:19:47'),(60,20201011162341,1,'2021-08-23 19:19:47'),(61,20201021104586,1,'2021-08-23 19:19:47'),(62,20201102112520,1,'2021-08-23 19:19:47'),(63,20201208121729,1,'2021-08-23 19:19:47'),(64,20201215091637,1,'2021-08-23 19:19:47'),(65,20210119174155,1,'2021-08-23 19:19:47'),(66,20210326182902,1,'2021-08-23 19:19:47'),(67,20210421112652,1,'2021-08-23 19:19:47'),(68,20210506095025,1,'2021-08-23 19:19:47'),(69,20210513115729,1,'2021-08-23 19:19:48'),(70,20210526113559,1,'2021-08-23 19:19:48'),(71,20210601000001,1,'2021-08-23 19:19:48'),(72,20210601000002,1,'2021-08-23 19:19:48'),(73,20210601000003,1,'2021-08-23 19:19:48'),(74,20210601000004,1,'2021-08-23 19:19:48'),(75,20210601000005,1,'2021-08-23 19:19:48'),(76,20210601000006,1,'2021-08-23 19:19:48'),(77,20210601000007,1,'2021-08-23 19:19:48'),(78,20210601000008,1,'2021-08-23 19:19:48'),(79,20210606151329,1,'2021-08-23 19:19:48'),(80,20210616163757,1,'2021-08-23 19:19:48'),(81,20210617174723,1,'2021-08-23 19:19:48'),(82,20210622160235,1,'2021-08-23 19:19:48'),(83,20210623100031,1,'2021-08-23 19:19:48'),(84,20210623133615,1,'2021-08-23 19:19:48'),(85,20210708143152,1,'2021-08-23 19:19:48'),(86,20210709124443,1,'2021-08-23 19:19:48'),(87,20210712155608,1,'2021-08-23 19:19:48'),(88,20210714102108,1,'2021-08-23 19:19:48'),(89,20210719153709,1,'2021-08-23 19:19:48'),(90,20210721171531,1,'2021-08-23 19:19:48'),(91,20210723135713,1,'2021-08-23 19:19:48'),(92,20210802135933,1,'2021-08-23 19:19:48'),(93,20210806112844,1,'2021-08-23 19:19:48'),(94,20210810095603,1,'2021-08-23 19:19:48'),(95,20210811150223,1,'2021-08-23 19:19:48'),(96,20210816141251,1,'2021-08-23 19:19:48'),(97,20210818151827,1,'2021-08-23 19:19:48'),(98,20210818182258,1,'2021-08-23 19:19:49'),(99,20210819131107,1,'2021-08-23 19:19:49'),(100,20210819143446,1,'2021-08-23 19:19:49');
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `network_interfaces` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `host_id` int(10) unsigned NOT NULL,
  `mac` varchar(255) NOT NULL DEFAULT '',
  `ip_address` varchar(255) NOT NULL DEFAULT '',
  `broadcast` varchar(255) NOT NULL DEFAULT '',
  `ibytes` bigint(20) NOT NULL DEFAULT '0',
  `interface` varchar(255) NOT NULL DEFAULT '',
  `ipackets` bigint(20) NOT NULL DEFAULT '0',
  `last_change` bigint(20) NOT NULL DEFAULT '0',
  `mask` varchar(255) NOT NULL DEFAULT '',
  `metric` int(11) NOT NULL DEFAULT '0',
  `mtu` int(11) NOT NULL DEFAULT '0',
  `obytes` bigint(20) NOT NULL DEFAULT '0',
  `ierrors` bigint(20) NOT NULL DEFAULT '0',
  `oerrors` bigint(20) NOT NULL DEFAULT '0',
  `opackets` bigint(20) NOT NULL DEFAULT '0',
  `point_to_point` varchar(255) NOT NULL DEFAULT '',
  `type` int(11) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_network_interfaces_unique_ip_host_intf` (`ip_address`,`host_id`,`interface`),
  KEY `idx_network_interfaces_hosts_fk` (`host_id`),
  FULLTEXT KEY `ip_address_search` (`ip_address`),
  CONSTRAINT `network_interfaces_ibfk_1` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `osquery_options` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `override_type` int(1) NOT NULL,
  `override_identifier` varchar(255) NOT NULL DEFAULT '',
  `options` json NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
INSERT INTO `osquery_options` VALUES (1,0,'','{\"options\": {\"logger_plugin\": \"tls\", \"pack_delimiter\": \"/\", \"logger_tls_period\": 10, \"distributed_plugin\": \"tls\", \"disable_distributed\": false, \"logger_tls_endpoint\": \"/api/v1/osquery/log\", \"distributed_interval\": 10, \"distributed_tls_max_attempts\": 3}, \"decorators\": {\"load\": [\"SELECT uuid AS host_uuid FROM system_info;\", \"SELECT hostname AS hostname FROM system_info;\"]}}');
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pack_targets` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `pack_id` int(10) unsigned DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `target_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `constraint_pack_target_unique` (`pack_id`,`target_id`,`type`),
  CONSTRAINT `pack_targets_ibfk_1` FOREIGN KEY (`pack_id`) REFERENCES `packs` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `packs` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `disabled` tinyint(1) NOT NULL DEFAULT '0',
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `platform` varchar(255) DEFAULT NULL,
  `pack_type` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_pack_unique_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `password_reset_requests` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `expires_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `user_id` int(10) unsigned NOT NULL,
  `token` varchar(1024) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `policies` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `query_id` int(10) unsigned NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `fk_policies_query_id` (`query_id`),
  CONSTRAINT `policies_ibfk_1` FOREIGN KEY (`query_id`) REFERENCES `queries` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
SET @saved_cs_client     = @@character_set_client;
SET character_set_client = utf8;
/*!50001 CREATE VIEW `policy_membership` AS SELECT 
 1 AS `id`,
 1 AS `policy_id`,
 1 AS `host_id`,
 1 AS `passes`,
 1 AS `created_at`,
 1 AS `updated_at`*/;
SET character_set_client = @saved_cs_client;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `policy_membership_history` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `policy_id` int(10) unsigned DEFAULT NULL,
  `host_id` int(10) unsigned NOT NULL,
  `passes` tinyint(1) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_policy_membership_passes` (`passes`),
  KEY `idx_policy_membership_policy_id` (`policy_id`),
  KEY `idx_policy_membership_host_id_passes` (`host_id`,`passes`),
  CONSTRAINT `policy_membership_history_ibfk_1` FOREIGN KEY (`policy_id`) REFERENCES `policies` (`id`) ON DELETE CASCADE,
  CONSTRAINT `policy_membership_history_ibfk_2` FOREIGN KEY (`host_id`) REFERENCES `hosts` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `queries` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `saved` tinyint(1) NOT NULL DEFAULT '0',
  `name` varchar(255) NOT NULL,
  `description` mediumtext NOT NULL,
  `query` mediumtext NOT NULL,
  `author_id` int(10) unsigned DEFAULT NULL,
  `observer_can_run` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_query_unique_name` (`name`),
  UNIQUE KEY `constraint_query_name_unique` (`name`),
  KEY `author_id` (`author_id`),
  CONSTRAINT `queries_ibfk_1` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scheduled_queries` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `pack_id` int(10) unsigned DEFAULT NULL,
  `query_id` int(10) unsigned DEFAULT NULL,
  `interval` int(10) unsigned DEFAULT NULL,
  `snapshot` tinyint(1) DEFAULT NULL,
  `removed` tinyint(1) DEFAULT NULL,
  `platform` varchar(255) DEFAULT '',
  `version` varchar(255) DEFAULT '',
  `shard` int(10) unsigned DEFAULT NULL,
  `query_name` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(1023) DEFAULT '',
  `denylist` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_names_in_packs` (`name`,`pack_id`),
  KEY `scheduled_queries_pack_id` (`pack_id`),
  KEY `scheduled_queries_query_name` (`query_name`),
  CONSTRAINT `scheduled_queries_pack_id` FOREIGN KEY (`pack_id`) REFERENCES `packs` (`id`) ON DELETE CASCADE,
  CONSTRAINT `scheduled_queries_query_name` FOREIGN KEY (`query_name`) REFERENCES `queries` (`name`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `scheduled_query_stats` (
  `host_id` int(10) unsigned NOT NULL,
  `scheduled_query_id` int(10) unsigned NOT NULL,
  `average_memory` int(11) DEFAULT NULL,
  `denylisted` tinyint(1) DEFAULT NULL,
  `executions` int(11) DEFAULT NULL,
  `schedule_interval` int(11) DEFAULT NULL,
  `last_executed` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `output_size` int(11) DEFAULT NULL,
  `system_time` int(11) DEFAULT NULL,
  `user_time` int(11) DEFAULT NULL,
  `wall_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`host_id`,`scheduled_query_id`),
  KEY `scheduled_query_id` (`scheduled_query_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `sessions` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `accessed_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `user_id` int(10) unsigned NOT NULL,
  `key` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_session_unique_key` (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `software` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `version` varchar(255) NOT NULL DEFAULT '',
  `source` varchar(64) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name_version` (`name`,`version`,`source`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `software_cpe` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `software_id` bigint(20) unsigned DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `cpe` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_software_cpe_software_id` (`software_id`),
  CONSTRAINT `software_cpe_ibfk_1` FOREIGN KEY (`software_id`) REFERENCES `software` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `software_cve` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `cpe_id` int(10) unsigned DEFAULT NULL,
  `cve` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_cpe_cve` (`cpe_id`,`cve`),
  CONSTRAINT `software_cve_ibfk_1` FOREIGN KEY (`cpe_id`) REFERENCES `software_cpe` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `statistics` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `anonymous_identifier` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `teams` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `name` varchar(255) NOT NULL,
  `description` varchar(1023) NOT NULL DEFAULT '',
  `agent_options` json DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_teams` (
  `user_id` int(10) unsigned NOT NULL,
  `team_id` int(10) unsigned NOT NULL,
  `role` varchar(64) NOT NULL,
  PRIMARY KEY (`user_id`,`team_id`),
  KEY `fk_user_teams_team_id` (`team_id`),
  CONSTRAINT `user_teams_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `user_teams_ibfk_2` FOREIGN KEY (`team_id`) REFERENCES `teams` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `password` varbinary(255) NOT NULL,
  `salt` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL DEFAULT '',
  `email` varchar(255) NOT NULL,
  `admin_forced_password_reset` tinyint(1) NOT NULL DEFAULT '0',
  `gravatar_url` varchar(255) NOT NULL DEFAULT '',
  `position` varchar(255) NOT NULL DEFAULT '',
  `sso_enabled` tinyint(4) NOT NULL DEFAULT '0',
  `global_role` varchar(64) DEFAULT NULL,
  `api_only` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_unique_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!50001 DROP VIEW IF EXISTS `policy_membership`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `policy_membership` AS select `policy_membership_history`.`id` AS `id`,`policy_membership_history`.`policy_id` AS `policy_id`,`policy_membership_history`.`host_id` AS `host_id`,`policy_membership_history`.`passes` AS `passes`,`policy_membership_history`.`created_at` AS `created_at`,`policy_membership_history`.`updated_at` AS `updated_at` from `policy_membership_history` where `policy_membership_history`.`id` in (select max(`policy_membership_history`.`id`) AS `max_id` from `policy_membership_history` group by `policy_membership_history`.`host_id`) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
