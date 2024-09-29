CREATE TABLE IF NOT EXISTS `dns`
(
    `id`         int unsigned                                                  NOT NULL AUTO_INCREMENT,
    `zone`       varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    `data`       json                                                          NOT NULL,
    `updated_at` datetime                                                      NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `zone` (`zone`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;