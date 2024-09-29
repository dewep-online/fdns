CREATE TABLE IF NOT EXISTS `blacklist_adblock_list`
(
    `id`         int unsigned                                              NOT NULL AUTO_INCREMENT,
    `data`       text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci     NOT NULL,
    `hash`       char(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `type`       enum ('dynamic','static')                                 NOT NULL,
    `updated_at` datetime                                                  NOT NULL,
    `deleted_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `hash` (`hash`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;