package where

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractString(q iface.GenericQSSQL, key string, value string) error {
	whereAliases := q.GetWhereAliases()
	aliasValue, ok := whereAliases[key]

	if !ok {
		q.AddWhereValue(key + " = ?")
	} else {
		q.AddWhereValue(aliasValue + " = ?")
	}

	q.AddWhereBindValue(value)

	return nil
}
