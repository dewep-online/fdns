package database_test

import (
	"testing"

	"github.com/dewep-online/fdns/pkg/database"
	"github.com/stretchr/testify/require"
)

func TestBlacklistDomainsModel_ToMap(t *testing.T) {
	list := make(database.BlacklistDomainsModel, 0)
	list = append(list, database.BlacklistDomainModel{
		Tag:    "123",
		Domain: "aaa",
		Active: 0,
	}, database.BlacklistDomainModel{
		Tag:    "456",
		Domain: "bbb",
		Active: 1,
	})
	require.Equal(t, map[string]string{"bbb": ""}, list.ToMap(database.ActiveTrue))
	require.Equal(t, map[string]string{"aaa": ""}, list.ToMap(database.ActiveFalse))
}
