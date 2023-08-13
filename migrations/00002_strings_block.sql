CREATE TABLE `strings_block`
(
    `id`         int unsigned                            NOT NULL AUTO_INCREMENT,
    `data`       varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    `updated_at` datetime                                NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `data` (`data`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
