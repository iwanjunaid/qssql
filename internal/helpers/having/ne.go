package having

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractNotEqual(q iface.GenericQSSQL, key string, values interface{}) error {
	havingAliases := q.GetHavingAliases()
	aliasValue, ok := havingAliases[key]

	if !ok {
		q.AddHavingValue(key + " <> ?")
	} else {
		q.AddHavingValue(aliasValue + " <> ?")
	}

	q.AddHavingBindValue(values)

	return nil
}
