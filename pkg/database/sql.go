package database

const (
	checkTables = "select count(*) from `sqlite_master`;"
)

var migrations = []string{
	`create table blacklist_list
		(
			tag        TEXT    not null,
			url        TEXT    not null,
			active     INTEGER default 1 not null,
			updated_at INTEGER not null
		);`,
	`create unique index blacklist_list_tag_uindex on blacklist_list (tag);`,
	`create unique index blacklist_list_url_uindex on blacklist_list (url);`,
	`create index blacklist_list_active_index on blacklist_list (active);`,
	`create table blacklist_domains
		(
			tag        TEXT    not null,
			domain     TEXT    not null,
			active     INTEGER default 1 not null,
			updated_at INTEGER not null
		);`,
	`create index blacklist_domains_tag_index on blacklist_domains (tag);`,
	`create unique index blacklist_domains_domain_uindex on blacklist_domains (domain);`,
	`create index blacklist_domains_active_index on blacklist_domains (active);`,
	`create table rules
		(
			types      TEXT,
			domain     TEXT,
			ips        TEXT,
			active     INTEGER default 1 not null,
			updated_at INTEGER
		);`,
	`create unique index rules_types_domain_uindex on rules (types, domain);`,
	`create index rules_active_index on rules (active);`,
	`insert into rules (types, domain, ips, active, updated_at) 
		values ('ns', 'cloudflare.com', '1.1.1.1, 1.0.0.1', 1, 0);`,
}
