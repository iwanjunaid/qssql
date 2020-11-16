package having

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractILike(q iface.GenericQSSQL, key string, value interface{}) error {
	havingAliases := q.GetHavingAliases()
	aliasValue, ok := havingAliases[key]

	if !ok {
		q.AddHavingValue(key + " ILIKE ?")
	} else {
		q.AddHavingValue(aliasValue + " ILIKE ?")
	}

	q.AddHavingBindValue(value)

	return nil
}
