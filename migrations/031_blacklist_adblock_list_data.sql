# https://filters.adtidy.org/extension/ublock/filters/3.txt
# https://filters.adtidy.org/extension/ublock/filters/20.txt
# https://filters.adtidy.org/extension/ublock/filters/2_without_easylist.txt
# https://filters.adtidy.org/extension/ublock/filters/11.txt
# https://cdn.statically.io/gh/uBlockOrigin/uAssetsCDN/main/thirdparties/easylist.txt

INSERT IGNORE INTO `blacklist_adblock_list` (`data`, `hash`, `type`, `updated_at`, `deleted_at`)
VALUES ('LOCAL',
        '9be34046ca1ba588dc09de6a2c470501a652545c',
        'static',
        now(),
        NULL),
       ('https://cdn.osspkg.com/adblock/ublock3.txt',
        'efd0942b36c013991d5d7a90b8b134d33a999e36',
        'dynamic',
        now(),
        NULL),
       ('https://cdn.osspkg.com/adblock/ublock20.txt',
        '3710bd237f378c7cc7bdf4e29478cec6befdad9a',
        'dynamic',
        now(),
        NULL),
       ('https://cdn.osspkg.com/adblock/adtidy2.txt',
        '4519274a65b57a8387bfbf0bacea1d7b6b965550',
        'dynamic',
        now(),
        NULL),
       ('https://cdn.osspkg.com/adblock/adtidy11.txt',
        '7089c97c84d32bbda302ae2a1d336a67e3a75c8f',
        'dynamic',
        now(),
        NULL),
       ('https://cdn.osspkg.com/adblock/easylist.txt',
        'cbfc00603de237ae610e90d417c78b1d1e10d876',
        'dynamic',
        now(),
        NULL);