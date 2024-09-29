INSERT IGNORE INTO `static_regexp_rules` (`rule`, `qtype`, `data`, `updated_at`, `deleted_at`)
VALUES (
        '^(\\d+)-(\\d+)-(\\d+)-(\\d+)\\.local\\.$',
        1,
        '[\"$1.$2.$3.$4\"]',
        now(),
        NULL);