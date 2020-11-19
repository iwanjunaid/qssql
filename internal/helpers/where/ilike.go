package where

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractILike(q iface.GenericQSSQL, key string, value interface{}) error {
	whereAliases := q.GetWhereAliases()
	aliasValue, ok := whereAliases[key]

	if !ok {
		q.AddWhereValue(key + " ILIKE ?")
	} else {
		q.AddWhereValue(aliasValue + " ILIKE ?")
	}

	q.AddWhereBindValue(value)

	return nil
}
