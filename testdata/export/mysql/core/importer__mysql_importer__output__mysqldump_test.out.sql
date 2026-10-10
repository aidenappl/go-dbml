CREATE TABLE `api_key_outer_auth` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT,
  `api_key_outer_id` bigint NOT NULL,
  `flag` varchar(255) NOT NULL,
  `created_at` datetime
);

CREATE TABLE `entity` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT COMMENT 'The database table row id',
  `dateCreated` timestamp DEFAULT (now()) COMMENT 'Date / time created',
  `dateEdited` timestamp DEFAULT (now()) COMMENT 'Date / time last modified',
  `createdBy` varchar(40) NOT NULL DEFAULT '' COMMENT 'The person who created the record',
  `editedBy` varchar(40) NOT NULL DEFAULT '' COMMENT 'The person who last modified the record',
  `constituentId` varchar(40) NOT NULL COMMENT 'The globally unique identifier for this constituent',
  `name` varchar(255) NOT NULL,
  `ein` varchar(16) NOT NULL
);

CREATE TABLE `entity2` (
  `id` bigint PRIMARY KEY NOT NULL AUTO_INCREMENT COMMENT 'The database table row id',
  `dateCreated` timestamp DEFAULT (now()) COMMENT 'Date / time created',
  `dateEdited` timestamp DEFAULT (now()) COMMENT 'Date / time last modified',
  `createdBy` varchar(40) NOT NULL DEFAULT '' COMMENT 'The person who created the record',
  `editedBy` varchar(40) NOT NULL DEFAULT '' COMMENT 'The person who last modified the record',
  `constituentId` varchar(40) NOT NULL COMMENT 'The globally unique identifier for this constituent',
  `name` varchar(255) NOT NULL DEFAULT 'hello',
  `ein` varchar(16) NOT NULL
);

CREATE UNIQUE INDEX `UQE_api_key_outer_auth_api_key_outer_auth_af` ON `api_key_outer_auth` (`api_key_outer_id`, `flag`) USING BTREE;

CREATE INDEX `IDX_api_key_outer_auth_id` ON `api_key_outer_auth` (`id`) USING BTREE;

CREATE INDEX `IDX_api_key_outer_auth_api_key_outer_id` ON `api_key_outer_auth` (`api_key_outer_id`) USING BTREE;

CREATE INDEX `IDX_api_key_outer_auth_flag` ON `api_key_outer_auth` (`flag`) USING BTREE;

CREATE INDEX `IDX_api_key_outer_auth_created_at` ON `api_key_outer_auth` (`created_at`) USING BTREE;

CREATE UNIQUE INDEX `id` ON `entity` (`id`);

CREATE UNIQUE INDEX `constituentId` ON `entity` (`constituentId`);

CREATE UNIQUE INDEX `ein` ON `entity` (`ein`);

CREATE INDEX `entity_index_0` ON `entity` (`constituentId`);

CREATE INDEX `entity_index_1` ON `entity` (`ein`);

CREATE UNIQUE INDEX `id` ON `entity2` (`id`);

CREATE UNIQUE INDEX `constituentId` ON `entity2` (`constituentId`);

CREATE UNIQUE INDEX `ein` ON `entity2` (`ein`);
