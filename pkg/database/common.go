package database

import "fmt"

type (
	Types  string
	Active int64
)

var (
	NS    Types = "ns"
	DNS   Types = "dns"
	Host  Types = "host"
	Regex Types = "regexp"
	Query Types = "query"
	IP    Types = "ip"

	typeList = map[Types]struct{}{
		NS:    {},
		DNS:   {},
		Host:  {},
		Query: {},
	}

	ActiveTrue  Active = 1
	ActiveFalse Active = 0
	ActiveALL   Active = -1
)

type (
	RuleModel struct {
		Types  string
		Domain string
		IPs    string
		Active Active
	}
	RulesModel []RuleModel
)

func (v RulesModel) ToMap() map[string]string {
	result := make(map[string]string, len(v))
	for _, m := range v {
		result[m.Domain] = m.IPs
	}
	return result
}

func ValidateType(v string) (Types, error) {
	vv := Types(v)
	if _, ok := typeList[vv]; ok {
		return vv, nil
	}
	return Types(""), fmt.Errorf("type is not supported")
}
