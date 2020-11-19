package having

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractIn(q iface.GenericQSSQL, key string, values interface{}) error {
	havingAliases := q.GetHavingAliases()
	aliasValue, ok := havingAliases[key]

	if !ok {
		q.AddHavingValue(key + " IN (?)")
	} else {
		q.AddHavingValue(aliasValue + " IN (?)")
	}

	q.AddHavingBindValue(values)

	return nil
}
