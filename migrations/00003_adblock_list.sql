CREATE TABLE `adblock_list`
(
    `id`         int unsigned                                              NOT NULL AUTO_INCREMENT,
    `data`       text COLLATE utf8mb4_unicode_ci                           NOT NULL,
    `hash`       char(40) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `updated_at` datetime                                                  NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `hash` (`hash`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;