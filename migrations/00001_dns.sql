CREATE TABLE `dns`
(
    `id`         int unsigned                                                  NOT NULL AUTO_INCREMENT,
    `zone`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `data`       json                                                          NOT NULL,
    `updated_at` datetime                                                      NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb3
  COLLATE = utf8mb3_unicode_ci;