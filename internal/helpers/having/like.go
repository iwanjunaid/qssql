package having

import (
	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractLike(q iface.GenericQSSQL, key string, value interface{}) error {
	havingAliases := q.GetHavingAliases()
	aliasValue, ok := havingAliases[key]

	if !ok {
		q.AddHavingValue(key + " LIKE ?")
	} else {
		q.AddHavingValue(aliasValue + " LIKE ?")
	}

	q.AddHavingBindValue(value)

	return nil
}
