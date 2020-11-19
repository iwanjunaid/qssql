package where

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractNotIn(q iface.GenericQSSQL, key string, values interface{}) error {
	whereAliases := q.GetWhereAliases()
	aliasValue, ok := whereAliases[key]

	if !ok {
		q.AddWhereValue(key + " NOT IN (?)")
	} else {
		q.AddWhereValue(aliasValue + " NOT IN (?) ")
	}

	q.AddWhereBindValue(values)

	return nil
}
