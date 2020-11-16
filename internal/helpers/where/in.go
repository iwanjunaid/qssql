package where

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractIn(q iface.GenericQSSQL, key string, values interface{}) error {
	whereAliases := q.GetWhereAliases()
	aliasValue, ok := whereAliases[key]

	if !ok {
		q.AddWhereValue(key + " IN (?)")
	} else {
		q.AddWhereValue(aliasValue + " IN (?)")
	}

	q.AddWhereBindValue(values)

	return nil
}
