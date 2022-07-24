package database

import (
	"context"
	"fmt"
	"time"

	"github.com/deweppro/go-logger"

	"github.com/dewep-online/fdns/pkg/utils"

	"github.com/deweppro/go-orm"
)

func (v *Database) GetRulesMap(ctx context.Context, types Types, ts int64, call func(map[string]string) error) error {
	list, err := v.GetRules(ctx, types, ts, ActiveTrue)
	if err != nil {
		return fmt.Errorf("load rules [%s]: %w", string(types), err)
	}
	logger.Infof("update rules [%s]: %d", string(types), len(list))
	return call(list.ToMap())
}

func (v *Database) GetAllRules(ctx context.Context, ts int64) (RulesModel, error) {
	result := make(RulesModel, 0, 1000)
	err := v.pool.QueryContext("get_all_rules", ctx, func(q orm.Querier) {
		q.SQL("select `types`, `domain`, `ips`, `active` from `rules` where `updated_at` >= ?;", ts)
		q.Bind(func(bind orm.Scanner) error {
			m := RuleModel{}
			if err := bind.Scan(&m.Types, &m.Domain, &m.IPs, &m.Active); err != nil {
				return err
			}
			result = append(result, m)
			return nil
		})
	})
	return result, err
}

func (v *Database) GetRules(ctx context.Context, t Types, ts int64, active Active) (RulesModel, error) {
	result := make(RulesModel, 0, 1000)
	err := v.pool.QueryContext("get_rules_type", ctx, func(q orm.Querier) {
		q.SQL("select `domain`, `ips` from `rules` where `types` = ? and `active` = ? and updated_at >= ?;",
			string(t), int64(active), ts)
		q.Bind(func(bind orm.Scanner) error {
			m := RuleModel{}
			if err := bind.Scan(&m.Domain, &m.IPs); err != nil {
				return err
			}
			result = append(result, m)
			return nil
		})
	})
	return result, err
}

func (v *Database) SetRules(ctx context.Context, t Types, domain, ips string, active Active) error {
	return v.pool.ExecContext("set_rule", ctx, func(q orm.Executor) {
		q.SQL("replace into rules (types, domain, ips, active, updated_at) values (?, ?, ?, ?, ?);")
		q.Params(string(t), domain, ips, int64(active), time.Now().Unix())
	})
}

func (v *Database) DelRules(ctx context.Context, t Types, domain string) error {
	return v.pool.ExecContext("del_rule", ctx, func(q orm.Executor) {
		q.SQL("delete from rules where types = ? and domain = ?;")
		q.Params(string(t), domain)
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type (
	BlacklistDomainModel struct {
		Tag    string
		Domain string
		Active Active
	}
	BlacklistDomainsModel []BlacklistDomainModel
)

func (v BlacklistDomainsModel) ToMap(a Active) map[string]string {
	result := make(map[string]string, len(v))
	for _, m := range v {
		if a != ActiveALL && m.Active != a {
			continue
		}
		result[m.Domain] = ""
	}
	return result
}

func (v *Database) GetBlacklistDomain(ctx context.Context, ts int64) (BlacklistDomainsModel, error) {
	result := make(BlacklistDomainsModel, 0, 1000)
	err := v.pool.QueryContext("get_blacklist_domain", ctx, func(q orm.Querier) {
		q.SQL("select `tag`, `domain`, `active` from `blacklist_domains` where updated_at >= ?;", ts)
		q.Bind(func(bind orm.Scanner) error {
			m := BlacklistDomainModel{}
			if err := bind.Scan(&m.Tag, &m.Domain, &m.Active); err != nil {
				return err
			}
			result = append(result, m)
			return nil
		})
	})
	return result, err
}

func (v *Database) SetBlacklistDomain(ctx context.Context, tag string, domain []string) error {
	return v.pool.TransactionContext("set_blacklist_domain", ctx, func(v orm.Tx) {
		v.Exec(func(e orm.Executor) {
			e.SQL(`insert or ignore into blacklist_domains (tag, domain, updated_at) values (?, ?, ?);`)
			ts := time.Now().Unix()
			for _, s := range domain {
				e.Params(tag, s, ts)
			}
		})
	})
}

func (v *Database) SetBlacklistDomainActive(ctx context.Context, domain string, active Active) error {
	return v.pool.TransactionContext("set_blacklist_domain", ctx, func(v orm.Tx) {
		v.Exec(func(e orm.Executor) {
			e.SQL(`update blacklist_domains set active = ?, updated_at = ? where domain = ?;`)
			e.Params(int64(active), time.Now().Unix(), domain)
			e.Bind(func(result orm.Result) error {
				if result.RowsAffected == 0 {
					return fmt.Errorf("cant update active: %s", domain)
				}
				return nil
			})
		})
	})
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type (
	BlacklistURIModel struct {
		Tag    string
		URI    string
		Active Active
	}
	BlacklistURIsModel []BlacklistURIModel
)

func (v BlacklistURIsModel) ToMap(a Active) map[string]string {
	result := make(map[string]string, len(v))
	for _, m := range v {
		if a != ActiveALL && m.Active != a {
			continue
		}
		result[m.Tag] = m.URI
	}
	return result
}

func (v *Database) GetBlacklistURI(ctx context.Context, ts int64) (BlacklistURIsModel, error) {
	result := make(BlacklistURIsModel, 0, 1000)
	err := v.pool.QueryContext("get_blacklist_list", ctx, func(q orm.Querier) {
		q.SQL("select `tag`, `url`, `active` from `blacklist_list` where updated_at >= ?;", ts)
		q.Bind(func(bind orm.Scanner) error {
			m := BlacklistURIModel{}
			if err := bind.Scan(&m.Tag, &m.URI, &m.Active); err != nil {
				return err
			}
			result = append(result, m)
			return nil
		})
	})
	return result, err
}

func (v *Database) SetBlacklistURI(ctx context.Context, uri string) (string, error) {
	tag := utils.Tag(uri)
	err := v.pool.ExecContext("set_blacklist_list", ctx, func(q orm.Executor) {
		q.SQL(`insert or ignore into blacklist_list (tag, url, updated_at) values (?, ?, ?);`)
		q.Params(tag, uri, time.Now().Unix())
	})
	return tag, err
}
