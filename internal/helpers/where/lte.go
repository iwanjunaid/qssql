package where

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractLowerThanOrEqual(q iface.GenericQSSQL, key string, values interface{}) error {
	whereAliases := q.GetWhereAliases()
	aliasValue, ok := whereAliases[key]

	if !ok {
		q.AddWhereValue(key + " <= ?")
	} else {
		q.AddWhereValue(aliasValue + " <= ?")
	}

	q.AddWhereBindValue(values)

	return nil
}
