CREATE TABLE `adblock_rules`
(
    `id`         int unsigned                        NOT NULL AUTO_INCREMENT,
    `list_id`    int unsigned                        NOT NULL,
    `data`       text COLLATE utf8mb4_unicode_ci     NOT NULL,
    `hash`       char(40) COLLATE utf8mb4_unicode_ci NOT NULL,
    `updated_at` datetime                            NOT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `hash` (`hash`),
    KEY `list_id` (`list_id`),
    CONSTRAINT `adblock_rules_ibfk_1` FOREIGN KEY (`list_id`) REFERENCES `adblock_list` (`id`) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
