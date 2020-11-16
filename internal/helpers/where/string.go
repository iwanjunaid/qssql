package where

import (
	"errors"
	"fmt"

	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func ExtractString(q iface.GenericQSSQL, key string, value interface{}) error {
	assertedValue, ok := value.(string)

	if !ok {
		errorMessage := fmt.Sprintf("Can't assert to string for key %s", key)

		return errors.New(errorMessage)
	}

	whereAliases := q.GetWhereAliases()
	aliasValue, ok := whereAliases[key]

	if !ok {
		q.AddWhereValue(key + " = ?")
	} else {
		q.AddWhereValue(aliasValue + " = ?")
	}

	q.AddWhereBindValue(assertedValue)

	return nil
}
