package where

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractLike(q iface.GenericQSSQL, key string, value interface{}) error {
	whereAliases := q.GetWhereAliases()
	aliasValue, ok := whereAliases[key]

	if !ok {
		q.AddWhereValue(key + " LIKE ?")
	} else {
		q.AddWhereValue(aliasValue + " LIKE ?")
	}

	q.AddWhereBindValue(value)

	return nil
}
