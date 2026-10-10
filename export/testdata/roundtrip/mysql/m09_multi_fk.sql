CREATE TABLE `tenants` (`id` int NOT NULL, `slug` varchar(40) NOT NULL, PRIMARY KEY (`id`));
CREATE TABLE `projects` (
  `tenant_id` int NOT NULL,
  `id` int NOT NULL,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`tenant_id`,`id`),
  CONSTRAINT `fk_proj_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`)
);
CREATE TABLE `tasks` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `tenant_id` int NOT NULL,
  `project_id` int NOT NULL,
  `parent_task_id` bigint DEFAULT NULL,
  `title` varchar(200) NOT NULL,
  `priority` tinyint NOT NULL DEFAULT 3,
  `due` date DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_task_project` (`tenant_id`,`project_id`),
  CONSTRAINT `fk_task_project` FOREIGN KEY (`tenant_id`, `project_id`) REFERENCES `projects` (`tenant_id`, `id`) ON DELETE CASCADE,
  CONSTRAINT `fk_task_parent` FOREIGN KEY (`parent_task_id`) REFERENCES `tasks` (`id`) ON DELETE SET NULL
);
