CREATE TABLE IF NOT EXISTS `static_rules`
(
    `id`         int               NOT NULL AUTO_INCREMENT,
    `rule`       varchar(255)      NOT NULL,
    `qtype`      smallint unsigned NOT NULL,
    `data`       json              NOT NULL,
    `updated_at` datetime          NOT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `qtype_rule` (`qtype`,`rule`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;