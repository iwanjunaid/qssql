package having

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractNotIn(q iface.GenericQSSQL, key string, values interface{}) error {
	havingAliases := q.GetHavingAliases()
	aliasValue, ok := havingAliases[key]

	if !ok {
		q.AddHavingValue(key + " NOT IN (?)")
	} else {
		q.AddHavingValue(aliasValue + " NOT IN (?) ")
	}

	q.AddHavingBindValue(values)

	return nil
}
