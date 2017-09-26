PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS `cluster` (
    `id` TEXT PRIMARY KEY NOT NULL,
    `name` TEXT NOT NULL COLLATE NOCASE, 
    `description` TEXT NULL, 
    `state` TEXT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS `UQE_cluster_name` ON `cluster` (`name`);

CREATE TABLE IF NOT EXISTS `cluster_component` (
    `cluster_id` TEXT NOT NULL, 
    `component_id` TEXT NOT NULL, 
    `component_name` TEXT NOT NULL, 
    `component` TEXT NOT NULL, 
    PRIMARY KEY ( `cluster_id`,`component_name` ),
    FOREIGN KEY(cluster_id) REFERENCES cluster(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS `UQE_cluster_component_name` ON `cluster_component` (`cluster_id`, `component_name`);

CREATE TABLE IF NOT EXISTS `cluster_host` (
    `cluster_id` TEXT NOT NULL, 
    `host_id` TEXT NOT NULL, 
    `ip` TEXT NOT NULL, 
    `hostname` TEXT NOT NULL, 
    `host` TEXT NOT NULL, 
    PRIMARY KEY ( `cluster_id`,`host_id` ),
    FOREIGN KEY(cluster_id) REFERENCES cluster(id)
);

CREATE UNIQUE INDEX IF NOT EXISTS `UQE_cluster_host_ip` ON `cluster_host` (`cluster_id`, `ip`);

CREATE UNIQUE INDEX IF NOT EXISTS `UQE_cluster_host_hostname` ON `cluster_host` (`cluster_id`, `hostname`);